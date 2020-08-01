package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/c-bata/goptuna/successivehalving"
	"github.com/c-bata/goptuna/tpe"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/jinzhu/gorm"
	"gonum.org/v1/gonum/mat"
	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	dataset string
)

func init() {
	flag.StringVar(&dataset, "dataset", "iris.csv", "File path to iris dataset")
	flag.Parse()
}

func main() {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	rdb.RunAutoMigrate(db)
	storage := rdb.NewStorage(db)

	pruner, _ := successivehalving.NewPruner(
		successivehalving.OptionSetReductionFactor(3))
	study, err := goptuna.CreateStudy(
		"gorgonia-iris",
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionPruner(pruner),
		goptuna.StudyOptionDirection(goptuna.StudyDirectionMaximize),
	)
	if err != nil {
		log.Fatal("failed to create study: ", err)
	}
	err = study.Optimize(objective, 200)
	if err != nil {
		log.Fatal("failed to optimize: ", err)
	}

	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	log.Printf("Best evaluation=%f", v)
	log.Printf("Solver: %s", params["solver"].(string))
	if params["solver"].(string) == "Vanilla" {
		log.Printf("Learning rate (vanilla): %f", params["vanilla_learning_rate"].(float64))
	}
}

// https://www.kaggle.com/amarpandey/implementing-linear-regression-on-iris-dataset/notebook
func objective(trial goptuna.Trial) (float64, error) {
	g := gorgonia.NewGraph()
	x, y := getXYMat()
	xT := tensor.FromMat64(mat.DenseCopyOf(x))
	yT := tensor.FromMat64(mat.DenseCopyOf(y))

	s := yT.Shape()
	yT.Reshape(s[0])

	X := gorgonia.NodeFromAny(g, xT, gorgonia.WithName("x"))
	Y := gorgonia.NodeFromAny(g, yT, gorgonia.WithName("y"))
	theta := gorgonia.NewVector(
		g,
		gorgonia.Float64,
		gorgonia.WithName("theta"),
		gorgonia.WithShape(xT.Shape()[1]),
		gorgonia.WithInit(gorgonia.Gaussian(0, 1)))

	pred, err := gorgonia.Mul(X, theta)
	if err != nil {
		return 0, err
	}

	// Gorgonia might delete values from nodes so we are going to save it
	// and print it out later
	var predicted gorgonia.Value
	gorgonia.Read(pred, &predicted)

	predError, err := gorgonia.Sub(pred, Y)
	if err != nil {
		return 0, err
	}
	squaredError, err := gorgonia.Square(predError)
	if err != nil {
		return 0, err
	}
	cost, err := gorgonia.Mean(squaredError)
	if err != nil {
		return 0, err
	}

	if _, err := gorgonia.Grad(cost, theta); err != nil {
		return 0, err
	}

	machine := gorgonia.NewTapeMachine(g, gorgonia.BindDualValues(theta))
	defer machine.Close()

	var solver gorgonia.Solver
	solverName, _ := trial.SuggestCategorical("solver", []string{"Adam", "Vanilla"})
	if solverName == "Adam" {
		solver = gorgonia.NewAdamSolver()
	} else if solverName == "Vanilla" {
		learnRate, _ := trial.SuggestLogFloat("vanilla_learning_rate", 1e-5, 1e-1)
		solver = gorgonia.NewVanillaSolver(gorgonia.WithLearnRate(learnRate))
	}
	model := []gorgonia.ValueGrad{theta}

	iter := 10000
	var acc float64
	for i := 1; i <= iter; i++ {
		if err = machine.RunAll(); err != nil {
			fmt.Printf("Error during iteration: %v: %v\n", i, err)
			return 0, err
		}

		if err = solver.Step(model); err != nil {
			return 0, err
		}
		acc = accuracy(predicted.Data().([]float64), Y.Value().Data().([]float64))
		machine.Reset() // Reset is necessary in a loop like this

		if i%100 == 0 {
			if err := trial.ShouldPrune(i, acc); err != nil {
				return 0, err
			}
		}
	}
	return acc, nil
}

func accuracy(prediction, y []float64) float64 {
	var ok float64
	for i := 0; i < len(prediction); i++ {
		if math.Round(prediction[i]-y[i]) == 0 {
			ok += 1.0
		}
	}
	return ok / float64(len(y))
}

func getXYMat() (*matrix, *matrix) {
	f, err := os.Open(dataset)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	df := dataframe.ReadCSV(f)

	toValue := func(s series.Series) series.Series {
		records := s.Records()
		floats := make([]float64, len(records))
		m := map[string]int{}
		for i, r := range records {
			if _, ok := m[r]; !ok {
				m[r] = len(m) + 1
			}
			floats[i] = float64(m[r])
		}
		return series.Floats(floats)
	}

	xDF := df.Drop("species")
	yDF := df.Select("species").Capply(toValue)
	numRows, _ := xDF.Dims()
	xDF = xDF.Mutate(series.New(one(numRows), series.Float, "bias"))
	return &matrix{xDF}, &matrix{yDF}
}

type matrix struct {
	dataframe.DataFrame
}

func (m matrix) At(i, j int) float64 {
	return m.Elem(i, j).Float()
}

func (m matrix) T() mat.Matrix {
	return mat.Transpose{Matrix: m}
}

func one(size int) []float64 {
	one := make([]float64, size)
	for i := 0; i < size; i++ {
		one[i] = 1.0
	}
	return one
}
