package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/c-bata/goptuna/successivehalving"
	"github.com/c-bata/goptuna/tpe"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	pb "gopkg.in/cheggaaa/pb.v1"
	"gorgonia.org/gorgonia"
	"gorgonia.org/gorgonia/examples/mnist"
	"gorgonia.org/tensor"
)

var (
	epochs    int
	dataset   string
	dtype     string
	batchsize int
	seed      int
	datadir   string

	numExamples int
	batchNum    int

	inputs  tensor.Tensor
	targets tensor.Tensor
	dt      tensor.Dtype
)

func init() {
	flag.IntVar(&epochs, "epochs", 100, "Number of epochs to train for")
	flag.StringVar(&dataset, "dataset", "train", "Which dataset to train on? Valid options are \"train\" or \"test\"")
	flag.StringVar(&dtype, "dtype", "float64", "Which dtype to use")
	flag.IntVar(&batchsize, "batchsize", 100, "Batch size")
	flag.IntVar(&seed, "seed", 1337, "Seed number")
	flag.StringVar(&datadir, "datadir", ".", "Directory path to mnist data")
	flag.Parse()
	rand.Seed(int64(seed))

	// parse dtype
	switch dtype {
	case "float64":
		dt = tensor.Float64
	case "float32":
		dt = tensor.Float32
	default:
		log.Fatalf("Unknown dtype: %v", dtype)
	}

	// check batch size
	if batchsize == 0 {
		log.Fatal("batch size should be larger than 0")
	}

	// Load dataset
	var err error
	inputs, targets, err = mnist.Load(dataset, datadir, dt)
	if err != nil {
		log.Fatal(err)
	}

	batchNum = numExamples / batchsize
	log.Printf("Batches %d", batchNum)

	// the data is in (numExamples, 784).
	// In order to use a convnet, we need to massage the data
	// into this format (batchsize, numberOfChannels, height, width).
	//
	// This translates into (numExamples, 1, 28, 28).
	// This is because the convolution operators actually understand height and width.
	// The 1 indicates that there is only one channel (MNIST data is black and white).
	numExamples = inputs.Shape()[0]

	err = inputs.Reshape(numExamples, 1, 28, 28)
	if err != nil {
		log.Fatal(err)
	}
}

type sli struct {
	start, end int
}

func (s sli) Start() int { return s.start }
func (s sli) End() int   { return s.end }
func (s sli) Step() int  { return 1 }

type convnet struct {
	g                  *gorgonia.ExprGraph
	w0, w1, w2, w3, w4 *gorgonia.Node // weights. the number at the back indicates which layer it's used for
	d0, d1, d2, d3     float64        // dropout probabilities
	out                *gorgonia.Node
}

func newConvNet(g *gorgonia.ExprGraph) *convnet {
	w0 := gorgonia.NewTensor(g, dt, 4, gorgonia.WithShape(32, 1, 3, 3), gorgonia.WithName("w0"), gorgonia.WithInit(gorgonia.GlorotN(1.0)))
	w1 := gorgonia.NewTensor(g, dt, 4, gorgonia.WithShape(64, 32, 3, 3), gorgonia.WithName("w1"), gorgonia.WithInit(gorgonia.GlorotN(1.0)))
	w2 := gorgonia.NewTensor(g, dt, 4, gorgonia.WithShape(128, 64, 3, 3), gorgonia.WithName("w2"), gorgonia.WithInit(gorgonia.GlorotN(1.0)))
	w3 := gorgonia.NewMatrix(g, dt, gorgonia.WithShape(128*3*3, 625), gorgonia.WithName("w3"), gorgonia.WithInit(gorgonia.GlorotN(1.0)))
	w4 := gorgonia.NewMatrix(g, dt, gorgonia.WithShape(625, 10), gorgonia.WithName("w4"), gorgonia.WithInit(gorgonia.GlorotN(1.0)))
	return &convnet{
		g:  g,
		w0: w0,
		w1: w1,
		w2: w2,
		w3: w3,
		w4: w4,

		d0: 0.2,
		d1: 0.2,
		d2: 0.2,
		d3: 0.55,
	}
}

func (m *convnet) learnables() gorgonia.Nodes {
	return gorgonia.Nodes{m.w0, m.w1, m.w2, m.w3, m.w4}
}

// This function is particularly verbose for educational reasons. In reality, you'd wrap up the layers within a layer struct type and perform per-layer activations
func (m *convnet) fwd(x *gorgonia.Node) (err error) {
	var c0, c1, c2, fc *gorgonia.Node
	var a0, a1, a2, a3 *gorgonia.Node
	var p0, p1, p2 *gorgonia.Node
	var l0, l1, l2, l3 *gorgonia.Node

	// LAYER 0
	// here we convolve with stride = (1, 1) and padding = (1, 1),
	// which is your bog standard convolution for convnet
	if c0, err = gorgonia.Conv2d(x, m.w0, tensor.Shape{3, 3}, []int{1, 1}, []int{1, 1}, []int{1, 1}); err != nil {
		return errors.Wrap(err, "Layer 0 Convolution failed")
	}
	if a0, err = gorgonia.Rectify(c0); err != nil {
		return errors.Wrap(err, "Layer 0 activation failed")
	}
	if p0, err = gorgonia.MaxPool2D(a0, tensor.Shape{2, 2}, []int{0, 0}, []int{2, 2}); err != nil {
		return errors.Wrap(err, "Layer 0 Maxpooling failed")
	}
	log.Printf("p0 shape %v", p0.Shape())
	if l0, err = gorgonia.Dropout(p0, m.d0); err != nil {
		return errors.Wrap(err, "Unable to apply a dropout")
	}

	// Layer 1
	if c1, err = gorgonia.Conv2d(l0, m.w1, tensor.Shape{3, 3}, []int{1, 1}, []int{1, 1}, []int{1, 1}); err != nil {
		return errors.Wrap(err, "Layer 1 Convolution failed")
	}
	if a1, err = gorgonia.Rectify(c1); err != nil {
		return errors.Wrap(err, "Layer 1 activation failed")
	}
	if p1, err = gorgonia.MaxPool2D(a1, tensor.Shape{2, 2}, []int{0, 0}, []int{2, 2}); err != nil {
		return errors.Wrap(err, "Layer 1 Maxpooling failed")
	}
	if l1, err = gorgonia.Dropout(p1, m.d1); err != nil {
		return errors.Wrap(err, "Unable to apply a dropout to layer 1")
	}

	// Layer 2
	if c2, err = gorgonia.Conv2d(l1, m.w2, tensor.Shape{3, 3}, []int{1, 1}, []int{1, 1}, []int{1, 1}); err != nil {
		return errors.Wrap(err, "Layer 2 Convolution failed")
	}
	if a2, err = gorgonia.Rectify(c2); err != nil {
		return errors.Wrap(err, "Layer 2 activation failed")
	}
	if p2, err = gorgonia.MaxPool2D(a2, tensor.Shape{2, 2}, []int{0, 0}, []int{2, 2}); err != nil {
		return errors.Wrap(err, "Layer 2 Maxpooling failed")
	}
	log.Printf("p2 shape %v", p2.Shape())

	var r2 *gorgonia.Node
	b, c, h, w := p2.Shape()[0], p2.Shape()[1], p2.Shape()[2], p2.Shape()[3]
	if r2, err = gorgonia.Reshape(p2, tensor.Shape{b, c * h * w}); err != nil {
		return errors.Wrap(err, "Unable to reshape layer 2")
	}
	log.Printf("r2 shape %v", r2.Shape())
	if l2, err = gorgonia.Dropout(r2, m.d2); err != nil {
		return errors.Wrap(err, "Unable to apply a dropout on layer 2")
	}

	ioutil.WriteFile("tmp.dot", []byte(m.g.ToDot()), 0644)

	// Layer 3
	if fc, err = gorgonia.Mul(l2, m.w3); err != nil {
		return errors.Wrapf(err, "Unable to multiply l2 and w3")
	}
	if a3, err = gorgonia.Rectify(fc); err != nil {
		return errors.Wrapf(err, "Unable to activate fc")
	}
	if l3, err = gorgonia.Dropout(a3, m.d3); err != nil {
		return errors.Wrapf(err, "Unable to apply a dropout on layer 3")
	}

	// output decode
	var out *gorgonia.Node
	if out, err = gorgonia.Mul(l3, m.w4); err != nil {
		return errors.Wrapf(err, "Unable to multiply l3 and w4")
	}
	m.out, err = gorgonia.SoftMax(out)
	return
}

func objective(trial goptuna.Trial) (float64, error) {
	learningRate, err := trial.SuggestLogFloat("learning_rate", 1e-5, 1e-1)
	if err != nil {
		return 0, err
	}

	g := gorgonia.NewGraph()
	x := gorgonia.NewTensor(g, dt, 4, gorgonia.WithShape(batchsize, 1, 28, 28), gorgonia.WithName("x"))
	y := gorgonia.NewMatrix(g, dt, gorgonia.WithShape(batchsize, 10), gorgonia.WithName("y"))
	m := newConvNet(g)
	if err := m.fwd(x); err != nil {
		return 0, err
	}
	losses := gorgonia.Must(gorgonia.Log(gorgonia.Must(gorgonia.HadamardProd(m.out, y))))
	cost := gorgonia.Must(gorgonia.Mean(losses))
	cost = gorgonia.Must(gorgonia.Neg(cost))

	// we wanna track costs
	var costVal gorgonia.Value
	gorgonia.Read(cost, &costVal)

	prog, locMap, _ := gorgonia.Compile(g)
	log.Printf("%v", prog)

	vm := gorgonia.NewTapeMachine(g, gorgonia.WithPrecompiled(prog, locMap), gorgonia.BindDualValues(m.learnables()...))
	solver := gorgonia.NewRMSPropSolver(
		gorgonia.WithBatchSize(float64(batchsize)),
		gorgonia.WithLearnRate(learningRate),
	)
	defer vm.Close()

	bar := pb.New(batchNum)
	bar.SetRefreshRate(time.Second)
	bar.SetMaxWidth(80)

	for i := 0; i < epochs; i++ {
		bar.Prefix(fmt.Sprintf("Epoch %d", i))
		bar.Set(0)
		bar.Start()
		for b := 0; b < batchNum; b++ {
			start := b * batchsize
			end := start + batchsize
			if start >= numExamples {
				break
			}
			if end > numExamples {
				end = numExamples
			}

			var xVal, yVal tensor.Tensor
			if xVal, err = inputs.Slice(sli{start, end}); err != nil {
				return 0, err
			}

			if yVal, err = targets.Slice(sli{start, end}); err != nil {
				return 0, err
			}
			if err = xVal.(*tensor.Dense).Reshape(batchsize, 1, 28, 28); err != nil {
				return 0, err
			}

			gorgonia.Let(x, xVal)
			gorgonia.Let(y, yVal)
			if err = vm.RunAll(); err != nil {
				return 0, err
			}
			solver.Step(gorgonia.NodesToValueGrads(m.learnables()))
			vm.Reset()
			bar.Increment()
		}

		evaluation := costVal.Data().(float64)
		if err := trial.ShouldPrune(i, evaluation); err != nil {
			return 0, err
		}
		log.Printf("Epoch %d | cost %v", i, costVal)
	}
	return costVal.Data().(float64), nil
}

func main() {
	pruner, _ := successivehalving.NewPruner()

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	rdb.RunAutoMigrate(db)
	storage := rdb.NewStorage(db)

	study, err := goptuna.CreateStudy(
		"gorgonia-convnet",
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionPruner(pruner),
	)
	if err != nil {
		log.Fatal("failed to create study: ", err)
	}
	err = study.Optimize(objective, 50)
	if err != nil {
		log.Fatal("failed to create study: ", err)
	}
}
