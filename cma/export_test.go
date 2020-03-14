package cma

var (
	ExportOptimizerIsFeasible             = (*Optimizer).isFeasible
	ExportOptimizerRepairInfeasibleParams = (*Optimizer).repairInfeasibleParams
)

func ExportDim(optimizer *Optimizer) int {
	return optimizer.dim
}
