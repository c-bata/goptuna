package main

import (
	"flag"
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/c-bata/goptuna/tpe"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("please pass dialect and dsn")
	}
	dialect := flag.Arg(0)
	dsn := flag.Arg(1)

	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}

	rdb.RunAutoMigrate(db)
	storage := rdb.NewStorage(db)
	defer db.Close()

	study, err := goptuna.CreateStudy(
		"rdb",
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionLoadIfExists(true),
	)
	if err != nil {
		log.Fatal("failed to create study", err)
	}

	if err = study.Optimize(objective, 50); err != nil {
		log.Fatal("failed to optimize", err)
	}

	v, err := study.GetBestValue()
	if err != nil {
		log.Fatal("failed to get best value", err)
	}
	params, err := study.GetBestParams()
	if err != nil {
		log.Fatal("failed to get best params:", err)
	}
	log.Printf("Best evaluation=%f (x1=%f, x2=%f)",
		v, params["x1"].(float64), params["x2"].(float64))
}
