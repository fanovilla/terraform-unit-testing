package terraform_unit_testing

import (
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type PlanFixture struct {
	Plan interface{}
	Json string
}

// Plan runs a terraform plan.
// See plan elements in https://www.terraform.io/docs/internals/json-format.html.
func Plan(t *testing.T, vars map[string]string) PlanFixture {

	dir, err := prepareAndPlan(t, vars)
	failTestOnError(t, err)
	defer os.RemoveAll(dir)

	planBytes, err := ioutil.ReadFile(".tut/plan.json")
	failTestOnError(t, err)

	var planData interface{}
	err = json.Unmarshal(planBytes, &planData)

	log.Println(planData)
	return PlanFixture{Json: string(planBytes), Plan: planData}
}

func prepareAndPlan(t *testing.T, vars map[string]string) (string, error) {
	dir, err := copyRootModuleForProcessing()
	failTestOnError(t, err)

	err = copyTestOverrides(dir)
	failTestOnError(t, err)

	err = terraformInit(dir)
	failTestOnError(t, err)

	err = terraformPlan(dir, vars)
	failTestOnError(t, err)

	err = terraformShow(dir)
	return dir, err
}

func failTestOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Failed with %s", err)
	}
}

func terraformShow(dir string) error {
	cmd := exec.Command("terraform", "show", "-json", "tf.plan")

	_ = os.Mkdir("./.tut", 0755)
	outfile, err := os.Create("./.tut/plan.json")
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmd.Dir = dir
	cmd.Stdout = outfile
	cmd.Stderr = os.Stdout

	err = cmd.Run()
	return err
}

func terraformPlan(dir string, vars map[string]string) error {

	args := []string{"plan", "-refresh=false", "-out=tf.plan"}
	for key, val := range vars {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, val))
	}
	log.Println(fmt.Sprintf("Running terraform %v", args))
	cmd := exec.Command("terraform", args...)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	return err
}

func terraformInit(dir string) error {
	cmd := exec.Command("terraform", "init", "-input=false")

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	return err
}

// Copies .tf files from current test suite folder to the root module for processing.
// This is useful for overriding provider and data entries to disable remote queries.
func copyTestOverrides(dir string) error {
	err := copy.Copy(".", dir, copy.Options{
		Skip: func(src string) (bool, error) {
			return !strings.HasSuffix(src, ".tf"), nil
		},
	},
	)
	return err
}

// Copies the whole terraform root module to a sibling temp directory
// so any relative module references outside the root module directory is resolved as is.
// See https://www.terraform.io/docs/language/modules/index.html#the-root-module
// Assumes a directory layout where the test is in a 'tests/<suite_name>/' directory under the root module as
// described in https://www.terraform.io/docs/language/modules/testing-experiment.html#writing-tests-for-a-module.
func copyRootModuleForProcessing() (string, error) {
	dir, err := ioutil.TempDir("../../..", "tut-")
	if err != nil {
		log.Fatal(err)
	}

	err = copy.Copy("../..", dir, copy.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, "/tests"), nil
		},
	},
	)
	return dir, err
}
