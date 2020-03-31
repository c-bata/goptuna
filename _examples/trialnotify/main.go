package main

import (
	"log"
	"math"
	"os"
	"sync"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestFloat("x1", -10, 10)
	x2, _ := trial.SuggestFloat("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	trialchan := make(chan goptuna.FrozenTrial, 8)
	study, _ := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionIgnoreError(true),
		goptuna.StudyOptionSetTrialNotifyChannel(trialchan),
	)

	var wg sync.WaitGroup
	wg.Add(2)
	var err error
	go func() {
		defer wg.Done()
		err = study.Optimize(objective, 100)
		close(trialchan)
	}()
	go func() {
		defer wg.Done()
		for t := range trialchan {
			log.Println("trial", t)
		}
	}()
	wg.Wait()
	if err != nil {
		os.Exit(1)
	}

	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	log.Printf("Best evaluation=%f (x1=%f, x2=%f)",
		v, params["x1"].(float64), params["x2"].(float64))
}
