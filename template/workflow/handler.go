package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

type WorkflowStep struct {
	Name     string `yaml:"name"`
	Function string `yaml:"function"`
	Method   string `yaml:"method"`
}

type WorkflowTop struct {
	Workflow WorkflowInner `yaml:"workflow"`
}

type WorkflowInner struct {
	Name       string         `yaml:"name"`
	GatewayURL string         `yaml:"gateway_url"`
	Steps      []WorkflowStep `yaml:"steps"`
}

// `yaml:"workflow"`

// handle a serverless request
func handle(req []byte) string {
	dir, _ := os.Getwd()

	workflowBytes, readErr := ioutil.ReadFile(path.Join(dir, "./function/workflow.yml"))
	if readErr != nil {
		return readErr.Error()
	}

	workflow := WorkflowTop{}

	err := yaml.Unmarshal(workflowBytes, &workflow)
	if err != nil {
		return err.Error()
	}

	previousInput := req
	for _, step := range workflow.Workflow.Steps {
		result, _, resErr := runStep(step, workflow.Workflow.GatewayURL, &previousInput)
		// fmt.Println(result, resCode, resErr)
		if resErr != nil {
			return resErr.Error()
		}
		previousInput = result
	}

	return string(previousInput)
}

func runStep(step WorkflowStep, gatewayURL string, input *[]byte) ([]byte, int, error) {

	reader := bytes.NewReader(*input)
	req, _ := http.NewRequest(step.Method, gatewayURL+"function/"+step.Function, reader)
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	var nextInput []byte
	if res.Body != nil {
		defer res.Body.Close()
		nextInput, _ = ioutil.ReadAll(res.Body)
	}

	return nextInput, res.StatusCode, nil
}
