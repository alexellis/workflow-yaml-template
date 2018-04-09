package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

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
func handle(req []byte) []byte {
	dir, _ := os.Getwd()

	workflowBytes, readErr := ioutil.ReadFile(path.Join(dir, "./function/workflow.yml"))
	if readErr != nil {
		return []byte(readErr.Error())
	}

	workflow := WorkflowTop{}

	err := yaml.Unmarshal(workflowBytes, &workflow)
	if err != nil {
		return []byte(err.Error())
	}

	previousInput := req
	for i, step := range workflow.Workflow.Steps {
		st := time.Now()
		result, statusCode, resErr := runStep(step, workflow.Workflow.GatewayURL, &previousInput)
		log.Printf("[%d] %s %d byte(s) HTTP: %d - %fs\n",
			i,
			step.Name,
			len(result),
			statusCode,
			time.Since(st).Seconds())

		if resErr != nil {
			return []byte(resErr.Error())
		}
		previousInput = result
	}

	return previousInput
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
