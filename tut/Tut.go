package tut

import (
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type PlanFixture struct {
	Plan interface{}
	Json string
}

type PlanFixtureConfig struct {
	Vars               map[string]string
	ModuleReplacements map[string]string
}

// PlanWithConfig runs a terraform plan.
// See plan elements in https://www.terraform.io/docs/internals/json-format.html.
func PlanWithConfig(t *testing.T, config PlanFixtureConfig) PlanFixture {

	dir, err := prepareAndPlan(t, config)
	failTestOnError(t, err)
	defer os.RemoveAll(dir)

	planBytes, err := ioutil.ReadFile(".tut/plan.json")
	failTestOnError(t, err)

	var planData interface{}
	err = json.Unmarshal(planBytes, &planData)

	log.Println(planData)
	return PlanFixture{Json: string(planBytes), Plan: planData}
}

func Plan(t *testing.T) PlanFixture {
	return PlanWithConfig(t, PlanFixtureConfig{})
}

func prepareAndPlan(t *testing.T, config PlanFixtureConfig) (string, error) {
	testDir, err := os.Getwd()
	failTestOnError(t, err)
	err = os.Chdir("../..") // change to root module directory for copying
	failTestOnError(t, err)

	dir, err := copyRootModuleForProcessing()
	failTestOnError(t, err)
	err = os.Chdir(testDir) // change back to test directory
	failTestOnError(t, err)

	err = replaceModuleCalls(dir, config.ModuleReplacements)
	failTestOnError(t, err)

	err = copyTestOverrides(dir)
	failTestOnError(t, err)

	err = terraformInit(dir)
	failTestOnError(t, err)

	err = terraformPlan(dir, config.Vars)
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

func replaceModuleCalls(dir string, replacements map[string]string) error {
	err := filepath.Walk(dir, getWalkFunc(replacements))
	return err
}

// Copies .tf files from current test suite folder to the root module for processing.
// This is useful for overriding provider and data entries to disable remote queries.
func copyTestOverrides(dir string) error {
	err := copy.Copy(".", dir, copy.Options{
		Skip: func(src string) (bool, error) {
			return src == ".tut", nil
		},
	},
	)
	return err
}

// Copies the whole terraform root module to a directory in os tempdir - used as a base for processing.
// See https://www.terraform.io/docs/language/modules/index.html#the-root-module
// Assumes a directory layout where the test is in a 'tests/<suite_name>/' directory under the root module as
// described in https://www.terraform.io/docs/language/modules/testing-experiment.html#writing-tests-for-a-module.
func copyRootModuleForProcessing() (string, error) {
	dir, err := ioutil.TempDir("", "tut-")
	if err != nil {
		log.Fatal(err)
	}

	err = copy.Copy("./", dir, copy.Options{ // assume already in terraform root module dir or copying
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, "/tests") || strings.HasSuffix(src, ".terraform"), nil
		},
		OnSymlink: func(src string) copy.SymlinkAction {
			return copy.Deep
		},
	},
	)
	return dir, err
}

func getWalkFunc(replacements map[string]string) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !!fi.IsDir() {
			return nil //
		}

		matched, err := filepath.Match("*.tf", fi.Name())

		if err != nil {
			panic(err)
			return err
		}

		if matched {
			read, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			log.Printf("Processing module call replacements in %s\n", path)

			newContents := string(read)
			for old, new := range replacements {
				newContents = strings.Replace(newContents, old, new, -1)
			}

			err = ioutil.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}
}
