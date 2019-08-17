package main

import (
	"context"
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	"golang.org/x/sync/errgroup"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	study, _ := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
	)

	eg, ctx := errgroup.WithContext(context.Background())
	study.WithContext(ctx)
	for i := 0; i < 5; i++ {
		eg.Go(func() error {
			return study.Optimize(objective, 20)
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal("Optimize error", err)
	}

	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	log.Printf("Best evaluation=%f (x1=%f, x2=%f)",
		v, params["x1"].(float64), params["x2"].(float64))
}
