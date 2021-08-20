package basic

import (
	tut "github.com/fanovilla/terraform-unit-testing"
	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJmespath(t *testing.T) {
	plan := tut.Plan(t, nil)

	results, _ := jmespath.Search("resource_changes[?type=='aws_s3_bucket'].change.after.bucket", plan.Plan)

	assert.Len(t, results, 1)
	assert.ElementsMatch(t, []string{"my-tf-test-bucket"}, results)
}
