package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	study, _ := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionSetLogger(logger),
	)

	var eg errgroup.Group
	for i := 0; i < 5; i++ {
		eg.Go(func() error {
			return study.Optimize(objective, 100)
		})
	}
	if err := eg.Wait(); err != nil {
		log.Println("error", err)
		os.Exit(1)
	}

	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	fmt.Println("best value:", v)
	fmt.Println("best params:", params)
}
