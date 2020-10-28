package main

import (
	"flag"
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb.v2"
	"github.com/c-bata/goptuna/tpe"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestFloat("x1", -10, 10)
	x2, _ := trial.SuggestFloat("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("please pass dialect and dsn")
	}
	dialect := flag.Arg(0)
	dsn := flag.Arg(1)

	var db *gorm.DB
	var err error
	if dialect == "sqlite3" {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else if dialect == "mysql" {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		log.Fatal("unsupported dialect")
	}
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	err = rdb.RunAutoMigrate(db)
	if err != nil {
		log.Fatal("failed to run auto migrate:", err)
	}

	study, err := goptuna.CreateStudy(
		"rdb",
		goptuna.StudyOptionStorage(rdb.NewStorage(db)),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionDirection(goptuna.StudyDirectionMinimize),
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
