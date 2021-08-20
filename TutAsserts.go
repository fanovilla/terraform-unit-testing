package terraform_unit_testing

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/assert"
	"log"
	"sort"
	"testing"
)

type ResourceCount struct {
	Resource string
	Count    int
}

func (r *ResourceCount) Equal(s ResourceCount) bool {
	return r.Resource == s.Resource && r.Count == s.Count
}

type ResourceCounts struct {
	Counts []ResourceCount
}

func makeCounts(counts []ResourceCount) ResourceCounts {
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Resource < counts[j].Resource
	})
	return ResourceCounts{Counts: counts}
}

func (f *PlanFixture) AssertResourceCounts(t *testing.T, expected map[string]int) {
	results, _ := jmespath.Search("resource_changes[*].type", f.Plan)
	types := convertToStringArray(results)
	actual := make(map[string]int)
	for _, num := range types {
		actual[num] = actual[num] + 1
	}

	expectedCounts := makeCounts(bag(expected))
	actualCounts := makeCounts(bag(actual))
	log.Println(cmp.Diff(expectedCounts, actualCounts))
	assert.True(t, cmp.Equal(expectedCounts, actualCounts))
}

func bag(dict map[string]int) []ResourceCount {
	var counts []ResourceCount // == nil
	for key, val := range dict {
		counts = append(counts, ResourceCount{Resource: key, Count: val})
	}
	return counts
}

func convertToStringArray(results interface{}) []string {
	types := results.([]interface{})
	b := make([]string, len(types))
	for i := range types {
		b[i] = fmt.Sprintf("%v", types[i])
	}
	return b
}
