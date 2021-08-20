package basic

import (
	tut "github.com/fanovilla/terraform-unit-testing"
	"testing"
)

func TestJmespath(t *testing.T) {
	plan := tut.Plan(t, nil)

	plan.AssertResourceCounts(t, map[string]int{
		"aws_iam_group": 1,
		"aws_iam_user":  2,
		"aws_s3_bucket": 3,
	})
}
