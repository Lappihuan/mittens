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
	"context"
	"sync"
	"time"

	"github.com/pterm/pterm"
)

// Spinner wraps pterm spinner for pod readiness waiting with nice terminal output.
type Spinner struct {
	spinner *pterm.SpinnerPrinter
	done    chan bool
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewSpinner creates a new spinner with the given message.
func NewSpinner(message string) *Spinner {
	ctx, cancel := context.WithCancel(context.Background())
	spinner, _ := pterm.DefaultSpinner.Start(message)
	// Give the spinner goroutine time to start
	time.Sleep(10 * time.Millisecond)
	return &Spinner{
		spinner: spinner,
		done:    make(chan bool, 1),
		mu:      sync.RWMutex{},
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Stop stops the spinner and displays a success message.
func (s *Spinner) Stop(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.spinner != nil {
		s.cancel()
		_ = s.spinner.Stop()
		// Give the spinner goroutine time to fully stop
		time.Sleep(10 * time.Millisecond)
		pterm.Success.Println(message)
	}
}

// Fail stops the spinner and displays a failure message.
func (s *Spinner) Fail(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.spinner != nil {
		s.cancel()
		_ = s.spinner.Stop()
		// Give the spinner goroutine time to fully stop
		time.Sleep(10 * time.Millisecond)
		pterm.Error.Println(message)
	}
}

// Update updates the spinner text.
func (s *Spinner) Update(message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.spinner != nil {
		s.spinner.UpdateText(message)
	}
}

// StartProgressWithTimeout shows a spinner for up to the given timeout duration,
// updating with progress information.
func StartProgressWithTimeout(message string, timeout time.Duration, onUpdate func(elapsed time.Duration)) {
	spinner := NewSpinner(message)
	defer spinner.Stop("Ready!")

	start := time.Now()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		elapsed := time.Since(start)
		if elapsed > timeout {
			return
		}
		if onUpdate != nil {
			onUpdate(elapsed)
		}
	}
}
