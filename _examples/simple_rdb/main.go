package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/c-bata/goptuna/tpe"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	flag.Parse()
	if len(flag.Args()) == 0 {
		logger.Fatal("please pass dialect and dsn")
	}
	dialect := flag.Arg(0)
	dsn := flag.Arg(1)

	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		logger.Fatal("failed to open db", zap.Error(err))
	}

	storage := rdb.NewStorage(db)
	defer db.Close()

	study, err := goptuna.LoadStudy(
		"rdb",
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSetLogger(logger),
	)
	if err != nil {
		logger.Fatal("failed to create study", zap.Error(err))
	}

	if err = study.Optimize(objective, 50); err != nil {
		logger.Fatal("failed to optimize", zap.Error(err))
	}

	v, err := study.GetBestValue()
	if err != nil {
		logger.Fatal("failed to get best value", zap.Error(err))
	}
	params, err := study.GetBestParams()
	if err != nil {
		logger.Fatal("failed to get best params", zap.Error(err))
	}
	fmt.Println("Result:")
	fmt.Println("- best value", v)
	fmt.Println("- best param", params)
}
