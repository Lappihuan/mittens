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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSpinner(t *testing.T) {
	// Disable pterm output during tests to avoid race conditions
	oldEnv := os.Getenv("CI")
	defer func() {
		if oldEnv != "" {
			os.Setenv("CI", oldEnv)
		} else {
			os.Unsetenv("CI")
		}
	}()
	os.Setenv("CI", "true")

	t.Run("spinner_creation", func(t *testing.T) {
		done := make(chan bool, 1)
		go func() {
			spinner := NewSpinner("Testing...")
			require.NotNil(t, spinner)
			require.NotNil(t, spinner.spinner)
			require.NotNil(t, spinner.done)

			// Stop should not panic
			spinner.Stop("Done!")
			time.Sleep(50 * time.Millisecond) // Allow goroutine to settle
			done <- true
		}()

		select {
		case <-done:
			// Success
		case <-time.After(3 * time.Second):
			t.Fatal("Test timeout")
		}
	})

	t.Run("spinner_fail", func(t *testing.T) {
		done := make(chan bool, 1)
		go func() {
			spinner := NewSpinner("Testing...")
			require.NotNil(t, spinner)

			// Fail should not panic
			spinner.Fail("Failed!")
			time.Sleep(50 * time.Millisecond) // Allow goroutine to settle
			done <- true
		}()

		select {
		case <-done:
			// Success
		case <-time.After(3 * time.Second):
			t.Fatal("Test timeout")
		}
	})

	t.Run("spinner_update", func(t *testing.T) {
		done := make(chan bool, 1)
		go func() {
			spinner := NewSpinner("Testing...")
			require.NotNil(t, spinner)

			// Update should not panic
			spinner.Update("Updated message")
			time.Sleep(20 * time.Millisecond)

			spinner.Stop("Done!")
			time.Sleep(50 * time.Millisecond) // Allow goroutine to settle
			done <- true
		}()

		select {
		case <-done:
			// Success
		case <-time.After(3 * time.Second):
			t.Fatal("Test timeout")
		}
	})
}
