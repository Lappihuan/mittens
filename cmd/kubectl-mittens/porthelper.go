// Copyright 2020 Soluble Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	v1 "k8s.io/api/core/v1"
)

// DetectServicePort attempts to auto-detect the port from a Service.
// Returns the port number if exactly one port is found, 0 if multiple ports exist,
// or an error if no ports are found.
func DetectServicePort(service *v1.Service) (int32, error) {
	if len(service.Spec.Ports) == 0 {
		return 0, fmt.Errorf("service %q has no ports defined", service.Name)
	}

	if len(service.Spec.Ports) == 1 {
		return service.Spec.Ports[0].Port, nil
	}

	// Multiple ports - return 0 to indicate user should select
	return 0, nil
}

// InteractivePortSelection prompts the user to select a port from available options.
// It returns the selected port number.
func InteractivePortSelection(service *v1.Service) (int32, error) {
	if len(service.Spec.Ports) == 0 {
		return 0, fmt.Errorf("service %q has no ports defined", service.Name)
	}

	// Build options for selection
	type portOption struct {
		name string
		port int32
	}

	portCount := len(service.Spec.Ports)
	options := make([]portOption, 0, portCount)
	optionStrings := make([]string, 0, portCount)

	for _, port := range service.Spec.Ports {
		portName := port.Name
		if portName == "" {
			portName = fmt.Sprintf("unnamed (port %d)", port.Port)
		}
		displayStr := fmt.Sprintf("%s (port %d)", portName, port.Port)
		options = append(options, portOption{name: displayStr, port: port.Port})
		optionStrings = append(optionStrings, displayStr)
	}

	// Create prompt
	selectedIndex := 0
	prompt := &survey.Select{
		Message: "Multiple ports found. Which port would you like to tap?",
		Options: optionStrings,
	}

	err := survey.AskOne(prompt, &selectedIndex)
	if err != nil {
		return 0, fmt.Errorf("port selection cancelled: %w", err)
	}

	return options[selectedIndex].port, nil
}
