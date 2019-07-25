# Goptuna

![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/c-bata/goptuna?status.svg)](https://godoc.org/github.com/c-bata/goptuna) 


Experimental Black-box optimization library, inspired by [optuna](https://github.com/pfnet/optuna).
This library is not only for machine learning but also we can use the parameter tuning of the systems built with Go.
Currently two algorithms are implemented:

* Random Search
* Tree of Parzen Estimators (TPE)

## Installation

You can integrate Goptuna in wide variety of Go projects because of its portability of pure Go.

```console
$ go get -u github.com/c-bata/goptuna
```

## Usage

Goptuna supports Define-By-Run style user API like Optuna.
It makes the modularity high, and the user can dynamically construct the search spaces.

```go
package main

import (
    "fmt"
    "math"

    "github.com/c-bata/goptuna"
    "github.com/c-bata/goptuna/tpe"
    "go.uber.org/zap"
)

func objective(trial goptuna.Trial) (float64, error) {
    x1, _ := trial.SuggestUniform("x1", -10, 10)
    x2, _ := trial.SuggestUniform("x2", -10, 10)
    return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()

    study, _ := goptuna.CreateStudy(
        "goptuna-example",
        goptuna.StudyOptionSampler(tpe.NewSampler()),
        goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
        goptuna.StudyOptionSetLogger(logger),
    )
    _ = study.Optimize(objective, 100)

    v, _ := study.GetBestValue()
    params, _ := study.GetBestParams()
    fmt.Println("result:", v, params)
}
```

**Advanced usages**

<details>

<summary>Parallel optimization</summary>

[full source code](./_examples/concurrency/main.go).

Goptuna can easily implement parallel optimization using goroutine.

```go
package main

import ...

func main() {
	study, _ := goptuna.CreateStudy(...)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := study.Optimize(objective, 100)
			if err != nil {
				log.Println("error", err)
			}
		}()
	}
	wg.Wait()

    v, _ := study.GetBestValue()
    fmt.Println("best evaluation value:", v)
}
```

</details>

<details>

<summary>Notification system</summary>

[full source code](./_examples/trialnotify/main.go).

You can receive notification of each trials via channel.
It can be used for logging and any notification systems.

```go
package main

import ...

func main() {
	trialchan := make(chan goptuna.FrozenTrial, 8)
	study, _ := goptuna.CreateStudy(
		...
		goptuna.StudyOptionIgnoreObjectiveErr(true),
		goptuna.StudyOptionSetTrialNotifyChannel(trialchan),
	)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		study.Optimize(objective, 100)
		close(trialchan)
	}()
	go func() {
		defer wg.Done()
		for t := range trialchan {
			log.Println("trial", t)
		}
	}()
	wg.Wait()

    v, _ := study.GetBestValue()
    fmt.Println("best evaluation value:", v)
}
```

</details>

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.
