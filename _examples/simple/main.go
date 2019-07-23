package main

import (
	"fmt"
	"math"
	"os"

	"github.com/c-bata/goptuna"
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
	return math.Pow(x1-1, 2) + math.Pow(x2-2, 2) + 1, nil
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.NewInMemoryStorage(),
		&goptuna.RandomSearchSampler{},
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSetLogger(logger),
	)
	if err != nil {
		logger.Fatal("failed to create study", zap.Error(err))
	}

	if err = study.Optimize(objective, 20); err != nil {
		logger.Fatal("failed to optimize", zap.Error(err))
	}

	v, err := study.GetBestValue()
	if err != nil {
		logger.Fatal("failed to get best value", zap.Error(err))
	}
	params, err := study.GetBestParams()
	if err != nil {
		logger.Fatal("failed to get best value", zap.Error(err))
	}
	logger.Info("best value",
		zap.Float64("value", v),
		zap.String("suggested", fmt.Sprintf("%#v", params)))
}
