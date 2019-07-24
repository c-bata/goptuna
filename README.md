# Goptuna

![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/c-bata/goptuna?status.svg)](https://godoc.org/github.com/c-bata/goptuna) 


Experimental Black-box optimization library written in pure Go, inspired by [optuna](https://github.com/pfnet/optuna).
This library helps the parameter tuning of the systems built with Go.

## Example

Goptuna supports define-by-run style user API.

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
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()

    study, _ := goptuna.CreateStudy(
        "goptuna-example",
        goptuna.StudyOptionSampler(tpe.NewTPESampler()),
        goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
        goptuna.StudyOptionSetLogger(logger),
    )
    _ = study.Optimize(objective, 100)

    v, _ := study.GetBestValue()
    params, _ := study.GetBestParams()
    fmt.Println("result:", v, params)
}
```

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.
