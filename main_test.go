package main

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/waler4ik/kk/cucumber"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"cucumber/features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^kk tool installed$`, cucumber.InstallKK)
	sc.Step(`^a directory without folder (.*)$`, cucumber.CleanFolder)
	sc.Step(`^I create a (.*) project with uri (.*)$`, cucumber.CreateProject)
	sc.Step(`^(.*) contents are same as in (.*)$`, cucumber.CompareFolder)
	sc.Step(`^I create a websocket with path (.*) in folder (.*)$`, cucumber.CreateWebSocket)
	sc.Step(`^I create a resource with path (.*) in folder (.*)$`, cucumber.CreateRESTResource)
}
