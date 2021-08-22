package basic

import (
	"github.com/fanovilla/terraform-unit-testing/tut"
	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestS3(t *testing.T) {
	plan := tut.PlanWithConfig(t, tut.PlanFixtureConfig{
		Vars: map[string]string{"bucket_name": "my-tf-test-bucket"},
	})

	results, _ := jmespath.Search("resource_changes[?type=='aws_s3_bucket'].change.after.bucket", plan.Plan)
	assert.Len(t, results, 1)

	bucketResult, _ := jmespath.Search("[0]", results)
	assert.Equal(t, "my-tf-test-bucket", bucketResult)
}
