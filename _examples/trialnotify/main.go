package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	"go.uber.org/zap"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, err := trial.SuggestUniform("x1", -10, 10)
	if err != nil {
		return 0.0, err
	}
	x2, err := trial.SuggestUniform("x2", -10, 10)
	if err != nil {
		return 0.0, err
	}
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	trialchan := make(chan goptuna.FrozenTrial, 8)
	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionIgnoreObjectiveErr(true),
		goptuna.StudyOptionSetTrialNotifyChannel(trialchan),
	)
	if err != nil {
		log.Fatal("failed to create study", zap.Error(err))
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := study.Optimize(objective, 100)
		if err != nil {
			log.Println("error", err)
		}
		close(trialchan)
	}()
	go func() {
		defer wg.Done()
		for t := range trialchan {
			log.Println("trial", t)
		}
	}()

	wg.Wait()
	v, err := study.GetBestValue()
	if err != nil {
		log.Fatal("failed to get best value", zap.Error(err))
	}
	params, err := study.GetBestParams()
	if err != nil {
		log.Fatal("failed to get best params", zap.Error(err))
	}
	fmt.Println("Result:")
	fmt.Println("- best value", v)
	fmt.Println("- best param", params)
}
