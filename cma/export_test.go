package cma

var (
	ExportOptimizerIsFeasible = (*Optimizer).isFeasible
)

func ExportDim(optimizer *Optimizer) int {
	return optimizer.dim
}
