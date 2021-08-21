package basic

import (
	tut "github.com/fanovilla/terraform-unit-testing"
	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJmespath(t *testing.T) {
	plan := tut.PlanWithConfig(t, tut.PlanFixtureConfig{
		ModuleReplacements: map[string]string{"terraform-aws-modules/s3-bucket/aws": "./modules/module_s3"},
	})

	results, _ := jmespath.Search("resource_changes[?type=='aws_ssm_parameter'].change.after.value", plan.Plan)
	assert.Len(t, results, 1)

	bucketArnResult, _ := jmespath.Search("[0]", results)
	assert.Equal(t, "test_s3_bucket_arn", bucketArnResult)
}
