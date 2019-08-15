package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func objective(trial goptuna.Trial) (float64, error) {
	ctx := trial.GetContext()

	x1, _ := trial.SuggestUniform("x1", -10, 10)
	x2, _ := trial.SuggestUniform("x2", -10, 10)

	cmd := exec.CommandContext(ctx, "sleep", "1")
	err := cmd.Run()
	if err != nil {
		return -1, err
	}
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		logger.Fatal("failed to open db", zap.Error(err))
	}
	defer db.Close()
	db.DB().SetMaxOpenConns(1)
	rdb.RunAutoMigrate(db)

	// create a study
	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionStorage(rdb.NewStorage(db)),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSetLogger(logger),
	)
	if err != nil {
		logger.Fatal("Failed to create a study", zap.Error(err))
	}

	// create a context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	study.WithContext(ctx)

	// set signal handler
	sigch := make(chan os.Signal, 1)
	defer close(sigch)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig, ok := <-sigch
		if !ok {
			return
		}
		cancel()
		logger.Error("Catch a kill signal", zap.String("signal", sig.String()))
	}()

	// run optimize with multiple goroutine workers
	concurrency := runtime.NumCPU() - 1
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = study.Optimize(objective, 100/concurrency)
			if err != nil {
				logger.Error("Optimize error", zap.Error(err))
			}
		}()
	}
	wg.Wait()

	// print best hyper-parameters and the result
	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	logger.Info("result", zap.Float64("value", v),
		zap.String("params", fmt.Sprintf("%#v", params)))
}
