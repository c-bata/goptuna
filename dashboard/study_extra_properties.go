package dashboard

import (
	"strconv"
	"sync"

	"github.com/c-bata/goptuna"
)

type CachedExtraStudyProperty struct {
	cursor                int
	unionUserAttrs        map[string]bool
	hasIntermediateValues bool

	mu sync.RWMutex
}

func NewCachedExtraStudyProperty() *CachedExtraStudyProperty {
	return &CachedExtraStudyProperty{
		cursor:                0,
		unionUserAttrs:        make(map[string]bool, 8),
		hasIntermediateValues: false,
	}
}

func (c *CachedExtraStudyProperty) GetUnionUserAttrs() []AttributeSpec {
	c.mu.RLock()
	defer c.mu.RUnlock()

	attrs := make([]AttributeSpec, 0, len(c.unionUserAttrs))
	for k := range c.unionUserAttrs {
		attrs = append(attrs, AttributeSpec{
			Key:      k,
			Sortable: c.unionUserAttrs[k],
		})
	}
	return attrs
}

func (c *CachedExtraStudyProperty) Update(trials []goptuna.FrozenTrial) {
	nextCursor := c.cursor
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := len(trials) - 1; i > 0; i-- {
		if c.cursor > trials[i].Number {
			break
		}
		if !trials[i].State.IsFinished() {
			nextCursor = trials[i].Number
		}

		c.updateUserAttrs(trials[i])
		if trials[i].State != goptuna.TrialStateFail {
			c.updateIntermediateValues(trials[i])
		}
	}
	c.cursor = nextCursor
}

func (c *CachedExtraStudyProperty) updateUserAttrs(trial goptuna.FrozenTrial) {
	for k := range trial.UserAttrs {
		if sortable, ok := c.unionUserAttrs[k]; ok && sortable {
			c.unionUserAttrs[k] = isNumberString(trial.UserAttrs[k])
		} else {
			c.unionUserAttrs[k] = isNumberString(trial.UserAttrs[k])
		}
	}
}

func (c *CachedExtraStudyProperty) updateIntermediateValues(trial goptuna.FrozenTrial) {
	if !c.hasIntermediateValues && trial.IntermediateValues != nil && len(trial.IntermediateValues) > 0 {
		c.hasIntermediateValues = true
	}
}

func isNumberString(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
