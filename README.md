# Goptuna

![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/c-bata/goptuna?status.svg)](https://godoc.org/github.com/c-bata/goptuna)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-bata/goptuna)](https://goreportcard.com/report/github.com/c-bata/goptuna)

Distributed hyperparameter optimization framework, inspired by [Optuna](https://github.com/optuna/optuna) [1].
This library is particularly designed for machine learning, but everything will be able to optimize if you can define the objective function
(e.g. Optimizing the number of goroutines of your server and the memory buffer size of the caching systems).

**Key features:**

| State-of-the-art algorithms | Optuna compatibility |
| --------------------------- | -------------------- |
| <img width="750" alt="state-of-the-art-algorithms" src="https://user-images.githubusercontent.com/5564044/88860180-2d66cc00-d236-11ea-9a2f-de731c54a870.png"> | <img width="750" alt="optuna-compatibility" src="https://user-images.githubusercontent.com/5564044/88843168-a3aa0500-d21b-11ea-8fc1-d1cdca890a3f.png"> |

**Supported algorithms:**

Goptuna supports various state-of-the-art algorithms.
These algorithms are written in pure Go with a few dependencies and continuously benchmarked on GitHub Actions.

* Random search
* TPE: Tree-structured Parzen Estimators [2]
* CMA-ES: Covariance Matrix Adaptation Evolution Strategy [3]
* IPOP-CMA-ES: CMA-ES with increasing population size [4]
* BIPOP-CMA-ES: BI-population CMA-ES [5]
* Median Stopping Rule [6]
* ASHA: Asynchronous Successive Halving Algorithm (Optuna flavored version) [1,7,8]

**Projects using Goptuna:**

* [Kubeflow/Katib: Kubernetes-based system for hyperparameter tuning and neural architecture search.](https://github.com/kubeflow/katib)
* [c-bata/goptuna-bayesopt: Goptuna sampler for Gaussian Process based bayesian optimization using d4l3k/go-bayesopt.](https://github.com/c-bata/goptuna-bayesopt) [9]
* [c-bata/goptuna-isucon9q: Applying bayesian optimization for the parameters of MySQL, Nginx and Go web applications.](https://github.com/c-bata/goptuna-isucon9q)
* (If you have a project which uses Goptuna and want your own project to be listed here, please submit a GitHub issue.)


## Installation

You can integrate Goptuna in wide variety of Go projects because of its portability of pure Go.

```console
$ go get -u github.com/c-bata/goptuna
```

## Usage

Goptuna supports Define-by-Run style API like Optuna.
You can dynamically construct the search spaces.

<table><tr><td valign="top" width="50%">

### 5 steps to use Goptuna.

1. Define an objective function which returns a value you want to minimize.
1. Define the search space via Suggest APIs.
1. Create a study which manages each experiment.
1. Evaluate your objective function.
1. Print the best evaluation parameters.

Furthermore, I recommend you to use RDB storage backend for following purposes.

* Continue from where we stopped in the previous optimizations.
* Scale studies to tens of workers that connecting to the same RDB storage.
* Visualize parameters on Jupyter notebook using Optuna.

</td><td valign="top" width="50%">

```go
package main

import (
    "log"
    "math"

    "github.com/c-bata/goptuna"
    "github.com/c-bata/goptuna/tpe"
)

func objective(trial goptuna.Trial) (float64, error) {
    x1, _ := trial.SuggestFloat("x1", -10, 10)
    x2, _ := trial.SuggestFloat("x2", -10, 10)
    return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
    study, err := goptuna.CreateStudy(
        "goptuna-example",
        goptuna.StudyOptionSampler(tpe.NewSampler()))
    if err != nil { ... }

    err = study.Optimize(objective, 100)
    if err != nil { ... }

    v, _ := study.GetBestValue()
    p, _ := study.GetBestParams()
    log.Printf("Best value=%f (x1=%f, x2=%f)",
        v, p["x1"].(float64), p["x2"].(float64))
}
```

</td></tr></table>

**Advanced usages**

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
* [9] [J. Snoek, H. Larochelle, and R. Adams. Practical Bayesian optimization of machine learning algorithms. In Advances in Neural Information Processing Systems 25, pages 2960–2968, 2012.](https://arxiv.org/abs/1206.2944)

Blog posts:

* [Practical bayesian optimization using Goptuna](https://medium.com/@c_bata_/practical-bayesian-optimization-in-go-using-goptuna-edf97195fcb5).

Status:

* [godoc.org](http://godoc.org/github.com/c-bata/goptuna)
* [gocover.io](https://gocover.io/github.com/c-bata/goptuna)
* [goreportcard.com](https://goreportcard.com/report/github.com/c-bata/goptuna)

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.

