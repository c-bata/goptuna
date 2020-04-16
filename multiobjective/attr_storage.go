package multiobjective

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c-bata/goptuna"
)

const (
	prefixSubMetricValue = "sub_metric_value"
)

func keySubMetricsValue(metricName string) string {
	return fmt.Sprintf("%s_%s", prefixSubMetricValue, metricName)
}

func parseMetricNameFromKey(metricName string) string {
	if !strings.HasPrefix(metricName, prefixSubMetricValue) {
		return ""
	}
	return metricName[len(prefixSubMetricValue+"_"):]
}

// ReportSubMetrics save the metrics for multi-objective optimization.
// You can retrieve Pareto-optimal solutions by `study.GetParetoOptimalTrials()`.
func ReportSubMetrics(trial goptuna.Trial, metrics map[string]float64) error {
	for name := range metrics {
		key := keySubMetricsValue(name)
		value := strconv.FormatFloat(metrics[name], 'f', -1, 64)
		err := trial.Study.Storage.SetStudySystemAttr(trial.Study.ID, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSubMetrics(trialSystemAttrs map[string]string) (map[string]float64, error) {
	metrics := make(map[string]float64, 4)
	for key := range trialSystemAttrs {
		if !strings.HasPrefix(key, prefixSubMetricValue) {
			continue
		}

		name := parseMetricNameFromKey(key)
		value, err := strconv.ParseFloat(trialSystemAttrs[key], 64)
		if err != nil {
			return nil, err
		}
		metrics[name] = value
	}
	return metrics, nil
}
