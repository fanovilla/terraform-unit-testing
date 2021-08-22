package basic

import (
	"github.com/fanovilla/terraform-unit-testing/tut"
	"testing"
)

func TestJmespath(t *testing.T) {
	plan := tut.Plan(t)

	plan.AssertResourceCounts(t, map[string]int{
		"aws_iam_group": 1,
		"aws_iam_user":  2,
		"aws_s3_bucket": 3,
	})
}
