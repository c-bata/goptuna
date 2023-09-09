package dashboard_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/dashboard"
)

func TestCachedExtraStudyProperty_UnionUserAttrs(t *testing.T) {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(goptuna.NewRandomSampler()),
	)
	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestFloat("x1", -10, 10)
		x2, _ := trial.SuggestFloat("x2", -10, 10)

		trial.SetUserAttr("attr1", fmt.Sprintf("%f", x1))
		trial.SetUserAttr("attr2", fmt.Sprintf("x1=%f", x1))
		if trial.ID%2 == 0 {
			trial.SetUserAttr("attr3", "foo")
		}
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}
	study.Optimize(objective, 100)

	p := dashboard.NewCachedExtraStudyProperty()

	trials, _ := study.GetTrials()
	p.Update(trials)

	unionUserAttrs := p.GetUnionUserAttrs()
	if len(unionUserAttrs) != 3 {
		t.Errorf("The length of GetUnionUserAttrs must be 3, but got %d", len(unionUserAttrs))
		return
	}
	for _, attr := range unionUserAttrs {
		if attr.Key == "attr1" && !attr.Sortable {
			t.Error("attr1 must be sortable")
			return
		}
		if attr.Key == "attr2" && attr.Sortable {
			t.Error("attr1 must not be sortable")
			return
		}
		if attr.Key == "attr3" && attr.Sortable {
			t.Error("attr1 must not be sortable")
			return
		}
	}
}
