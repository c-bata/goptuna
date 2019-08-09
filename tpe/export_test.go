package tpe

var (
	ExportGetObservationPairs   = getObservationPairs
	ExportSplitObservationPairs = (*Sampler).splitObservationPairs
	ExportSampleCategorical     = (*Sampler).sampleCategorical
)
