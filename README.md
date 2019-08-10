# Goptuna

![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/c-bata/goptuna?status.svg)](https://godoc.org/github.com/c-bata/goptuna)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-bata/goptuna)](https://goreportcard.com/report/github.com/c-bata/goptuna)


Bayesian optimization framework for black-box functions, inspired by [Optuna](https://github.com/pfnet/optuna).
This library is not only for hyperparameter tuning of machine learning models.
Everything will be able to optimized if you can define the objective function (e.g. Optimizing the number of goroutines of your server and the memory buffer size of the caching systems).

Currently following algorithms are implemented:

* Random Search
* Tree-structured Parzen Estimators (TPE)
* Median Stopping Rule (Google Vizier)

See the blog post for more details: [Practical bayesian optimization in Go using Goptuna](https://medium.com/@c_bata_/practical-bayesian-optimization-in-go-using-goptuna-edf97195fcb5).

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

<summary>Parallel optimization with multiple goroutine workers</summary>

``Optimize`` method of ``goptuna.Study`` object is designed as the goroutine safe.
So you can easily optimize your objective function using multiple goroutine workers.

```go
package main

import ...

func main() {
    study, _ := goptuna.CreateStudy(...)

    eg, ctx := errgroup.WithContext(context.Background())
    study.WithContext(ctx)
    for i := 0; i < 5; i++ {
        eg.Go(func() error {
            return study.Optimize(objective, 100)
        })
    }
    if err := eg.Wait(); err != nil { ... }
    ...
}
```

[full source code](./_examples/concurrency/main.go)

</details>

<details>

<summary>Receive notifications of each trials</summary>

You can receive notifications of each trials via channel.
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
    if err != nil { ... }
    ...
}
```

[full source code](./_examples/trialnotify/main.go)

</details>

## Links

References:

* TPE: [James S. Bergstra, Remi Bardenet, Yoshua Bengio, and Balázs Kégl. Algorithms for hyper-parameter optimization. In Advances in Neural Information Processing Systems 25. 2011.](https://papers.nips.cc/paper/4443-algorithms-for-hyper-parameter-optimization.pdf)
* Optuna (Define-by-Run interface): [Takuya Akiba, Shotaro Sano, Toshihiko Yanase, Takeru Ohta, Masanori Koyama. 2019. Optuna: A Next-generation Hyperparameter Optimization Framework. In The 25th ACM SIGKDD Conference on Knowledge Discovery and Data Mining (KDD ’19), August 4–8, 2019.](https://dl.acm.org/citation.cfm?id=3330701)
* Median stopping rule: [Golovin, B. Sonik, S. Moitra, G. Kochanski, J. Karro, and D.Sculley. Google Vizier: A service for black-box optimization. In Knowledge Discovery and Data Mining (KDD), 2017.](http://www.kdd.org/kdd2017/papers/view/google-vizier-a-service-for-black-box-optimization)

Status:

* [godoc.org](http://godoc.org/github.com/c-bata/goptuna)
* [gocover.io](https://gocover.io/github.com/c-bata/goptuna)
* [goreportcard.com](https://goreportcard.com/report/github.com/c-bata/goptuna)

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.
