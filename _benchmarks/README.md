# Continuous benchmarking using kurobako and GitHub Actions

Benchmark scripts are built on [kurobako](https://github.com/sile/kurobako) and [kurobako-go](https://github.com/sile/kurobako-go).
See [Introduction to Kurobako: A Benchmark Tool for Hyperparameter Optimization Algorithms](https://medium.com/optuna/kurobako-a2e3f7b760c7) for more details.

## How to run benchmark scripts

GitHub Actions continuously run the benchmark scripts and comment on your pull request using [github-actions-kurobako](https://github.com/c-bata/github-actions-kurobako).
If you want to run on your local machines, please execute following after installed kurobako.

```console
$ mkdir -p tmp
$ ./_benchmark/runner.sh rosenbrock ./tmp/kurobako.json
$ cat ./tmp/kurobako.json | kurobako plot curve --errorbar -o ./tmp
```

`kurobako plot curve` requires gnuplot. If you want to run on Docker container, please execute following:

```
$ docker pull sile/kurobako
$ ./_benchmarks/runner.sh -h
runner.sh is an entrypoint to run benchmarkers.
Usage:
    $ runner.sh <problem> <json-path>
Problem:
    rosenbrock     : https://www.sfu.ca/~ssurjano/rosen.html
    himmelblau     : https://en.wikipedia.org/wiki/Himmelblau%27s_function
    ackley         : https://www.sfu.ca/~ssurjano/ackley.html
    rastrigin      : https://www.sfu.ca/~ssurjano/rastr.html
    weierstrass    : Weierstrass function in https://github.com/sigopt/evalset
    schwefel20     : https://www.sfu.ca/~ssurjano/schwef.html
    schwefel36     : https://www.sfu.ca/~ssurjano/schwef.html
Options:
    --help, -h         print this
Example:
    $ runner.sh rosenbrock ./tmp/kurobako.json
    $ cat ./tmp/kurobako.json | kurobako plot curve --errorbar -o ./tmp
$ ./_benchmark/runner.sh rosenbrock ./tmp/kurobako.json
$ cat ./tmp/kurobako.json | docker run -v $PWD/tmp/images/:/images/ --rm -i sile/kurobako plot curve
```

If you got something error, please investigate using:

```
$ docker run -it --rm -v $PWD/tmp:/volume --entrypoint sh sile/kurobako
```
