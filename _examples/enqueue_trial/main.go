package main

import (
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb.v2"
	"github.com/c-bata/goptuna/tpe"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestFloat("x1", -10, 10)
	x2, _ := trial.SuggestFloat("x2", -10, 10)
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	err = rdb.RunAutoMigrate(db)
	if err != nil {
		log.Fatal("failed to run auto migrate:", err)
	}
	storage := rdb.NewStorage(db)

	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionLoadIfExists(true),
	)
	if err != nil {
		log.Fatal("failed to create study:", err)
	}

	for i := 0; i < 5; i++ {
		err = study.EnqueueTrial(map[string]float64{
			"x1": float64(i), "x2": float64(i),
		})
		if err != nil {
			log.Fatal("failed to enqueue trial: ", err)
		}
	}

	if err = study.Optimize(objective, 50); err != nil {
		log.Fatal("failed to optimize:", err)
	}

	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	log.Printf("Best evaluation=%f (x1=%f, x2=%f)",
		v, params["x1"].(float64), params["x2"].(float64))
}
