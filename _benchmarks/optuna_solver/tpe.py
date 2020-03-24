import argparse
import optuna

from kurobako import solver
from kurobako.solver.optuna import OptunaSolverFactory

parser = argparse.ArgumentParser()
parser.add_argument('--loglevel', choices=['debug', 'info', 'warning', 'error'], default='warning')
args = parser.parse_args()

if args.loglevel == 'debug':
    optuna.logging.set_verbosity(optuna.logging.DEBUG)
elif args.loglevel == 'info':
    optuna.logging.set_verbosity(optuna.logging.INFO)
elif args.loglevel == 'warning':
    optuna.logging.set_verbosity(optuna.logging.WARNING)
elif args.loglevel == 'error':
    optuna.logging.set_verbosity(optuna.logging.ERROR)


class CustomOptunaSolverFactory(OptunaSolverFactory):
    def specification(self):
        spec = super().specification()
        spec.name = "Optuna (TPE)"
        return spec


def create_study(seed):
    sampler = optuna.samplers.TPESampler(seed=seed)
    return optuna.create_study(sampler=sampler)


if __name__ == '__main__':
    factory = CustomOptunaSolverFactory(create_study)
    runner = solver.SolverRunner(factory)
    runner.run()
