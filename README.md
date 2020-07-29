# Goptuna

![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/c-bata/goptuna?status.svg)](https://godoc.org/github.com/c-bata/goptuna)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-bata/goptuna)](https://goreportcard.com/report/github.com/c-bata/goptuna)

Distributed hyperparameter optimization framework, inspired by [Optuna](https://github.com/optuna/optuna) [1].
This library is particularly designed for machine learning, but everything will be able to optimize if you can define the objective function
(e.g. Optimizing the number of goroutines of your server and the memory buffer size of the caching systems).

**Supported algorithms:**

Goptuna supports following state-of-the-art algorithms.
These algorithms are written in pure Go with a few dependencies and continuously benchmarked on GitHub Actions.

* Random search
* TPE: Tree-structured Parzen Estimators [2]
* CMA-ES: Covariance Matrix Adaptation Evolution Strategy [3]
* IPOP-CMA-ES: CMA-ES with increasing population size [4]
* BIPOP-CMA-ES: BI-population CMA-ES [5]
* Median Stopping Rule [6]
* ASHA: Asynchronous Successive Halving Algorithm (Optuna flavored version) [1,7,8]

**Other key features:**

| Easy to use | Optuna compatibility |
| ----------- | -------------------- |
| <img width="750" alt="easy-to-use" src="https://user-images.githubusercontent.com/5564044/88830312-c2eb6700-d208-11ea-982a-fd034116280a.png"> | <img width="750" alt="optuna-compatibility" src="https://user-images.githubusercontent.com/5564044/88843168-a3aa0500-d21b-11ea-8fc1-d1cdca890a3f.png"> |

**Projects using Goptuna:**

* [Kubeflow/Katib: Kubernetes-based system for hyperparameter tuning and neural architecture search.](https://github.com/kubeflow/katib)
* (If you have a project which uses Goptuna and want your own project to be listed here, please submit a GitHub issue.)


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
    "log"
    "math"

    "github.com/c-bata/goptuna"
    "github.com/c-bata/goptuna/tpe"
)

// Define an objective function we want to minimize.
func objective(trial goptuna.Trial) (float64, error) {
    // Define a search space of the input values.
    x1, _ := trial.SuggestFloat("x1", -10, 10)
    x2, _ := trial.SuggestFloat("x2", -10, 10)

    // Here is a two-dimensional quadratic function.
    // F(x1, x2) = (x1 - 2)^2 + (x2 + 5)^2
    return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
    study, err := goptuna.CreateStudy(
        "goptuna-example",
        goptuna.StudyOptionSampler(tpe.NewSampler()),
    )
    if err != nil { ... }

    // Run an objective function 100 times to find a global minimum.
    err = study.Optimize(objective, 100)
    if err != nil { ... }

    // Print the best evaluation value and the parameters.
    // Mathematically, argmin F(x1, x2) is (x1, x2) = (+2, -5).
    v, _ := study.GetBestValue()
    p, _ := study.GetBestParams()
    log.Printf("Best evaluation value=%f (x1=%f, x2=%f)",
        v, p["x1"].(float64), p["x2"].(float64))
}
```

```console
$ go run main.go
...
2019/08/18 00:54:51 Best evaluation=0.038327 (x1=2.181604, x2=-4.926880)
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

<summary>Distributed optimization using RDB storage backend with MySQL</summary>

There is no complicated setup for distributed optimization but all Goptuna workers need to use the same RDB storage backend.
First, setup MySQL server like following to share the optimization result.

```console
$ cat mysql/my.cnf
[mysqld]
bind-address = 0.0.0.0
default_authentication_plugin=mysql_native_password

$ docker pull mysql:8.0
$ docker run \
  -d \
  --rm \
  -p 3306:3306 \
  --mount type=volume,src=mysql,dst=/etc/mysql/conf.d \
  -e MYSQL_USER=goptuna \
  -e MYSQL_DATABASE=goptuna \
  -e MYSQL_PASSWORD=password \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  --name goptuna-mysql \
  mysql:8.0
```

Then, create a study object using goptuna CLI

```console
$ goptuna create-study --storage mysql://goptuna:password@localhost:3306/yourdb --study yourstudy
yourstudy
```

```mysql
$ mysql --host 127.0.0.1 --port 3306 --user goptuna -ppassword -e "SELECT * FROM studies;"
+----------+------------+-----------+
| study_id | study_name | direction |
+----------+------------+-----------+
|        1 | yourstudy  | MINIMIZE  |
+----------+------------+-----------+
1 row in set (0.00 sec)
```

Finally, run the Goptuna workers which contains following code.

```go
package main

import ...

func main() {
    db, _ := gorm.Open("mysql", "goptuna:password@tcp(localhost:3306)/yourdb?parseTime=true")
    storage := rdb.NewStorage(db)
    defer db.Close()

    study, _ := goptuna.LoadStudy(
        "yourstudy",
        goptuna.StudyOptionStorage(storage),
        ...,
    )
    _ = study.Optimize(objective, 50)
    ...
}
```

The schema of Goptuna RDB storage backend is compatible with Optuna's one.
So you can check optimization result with Optuna's dashboard like following:

```console
$ pip install optuna bokeh mysqlclient
$ optuna dashboard --storage mysql+mysqldb://goptuna:password@127.0.0.1:3306/yourdb --study yourstudy
...
```

[shell script to reproduce this](./_examples/simple_rdb/check_mysql.sh)

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

* [1] [Takuya Akiba, Shotaro Sano, Toshihiko Yanase, Takeru Ohta, Masanori Koyama. 2019. Optuna: A Next-generation Hyperparameter Optimization Framework. In The 25th ACM SIGKDD Conference on Knowledge Discovery and Data Mining (KDD ’19), August 4–8, 2019.](https://dl.acm.org/citation.cfm?id=3330701)
* [2] [James S. Bergstra, Remi Bardenet, Yoshua Bengio, and Balázs Kégl. Algorithms for hyper-parameter optimization. In Advances in Neural Information Processing Systems 25. 2011.](https://papers.nips.cc/paper/4443-algorithms-for-hyper-parameter-optimization.pdf)
* [3] [N. Hansen, The CMA Evolution Strategy: A Tutorial. arXiv:1604.00772, 2016.](https://arxiv.org/abs/1604.00772)
* [4] [Auger, A., Hansen, N.: A restart CMA evolution strategy with increasing population size. In: Proceedings of the 2005 IEEE Congress on Evolutionary Computation (CEC’2005), pp. 1769–1776 (2005a)](https://sci2s.ugr.es/sites/default/files/files/TematicWebSites/EAMHCO/contributionsCEC05/auger05ARCMA.pdf)
* [5] [Hansen N. Benchmarking a BI-Population CMA-ES on the BBOB-2009 Function Testbed. In the workshop Proceedings of the Genetic and Evolutionary Computation Conference, GECCO, pages 2389–2395. ACM, 2009.](https://hal.inria.fr/inria-00382093/document)
* [6] [Golovin, B. Sonik, S. Moitra, G. Kochanski, J. Karro, and D.Sculley. Google Vizier: A service for black-box optimization. In Knowledge Discovery and Data Mining (KDD), 2017.](http://www.kdd.org/kdd2017/papers/view/google-vizier-a-service-for-black-box-optimization)
* [7] [Kevin G. Jamieson and Ameet S. Talwalkar. Non-stochastic best arm identification and hyperparameter optimization. In AISTATS, 2016.](http://proceedings.mlr.press/v51/jamieson16.pdf)
* [8] [Liam Li, Kevin Jamieson, Afshin Rostamizadeh, Ekaterina Gonina, Moritz Hardt, Benjamin Recht, and Ameet Talwalkar. Massively parallel hyperparameter tuning. arXiv preprint arXiv:1810.05934, 2018.](https://arxiv.org/abs/1810.05934)

Blog posts:

* [Practical bayesian optimization using Goptuna](https://medium.com/@c_bata_/practical-bayesian-optimization-in-go-using-goptuna-edf97195fcb5).

Status:

* [godoc.org](http://godoc.org/github.com/c-bata/goptuna)
* [gocover.io](https://gocover.io/github.com/c-bata/goptuna)
* [goreportcard.com](https://goreportcard.com/report/github.com/c-bata/goptuna)

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.

