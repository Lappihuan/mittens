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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestDetectServicePort(t *testing.T) {
	tests := []struct {
		name          string
		service       *v1.Service
		expectedPort  int32
		expectedZero  bool
		expectError   bool
		expectedError string
	}{
		{
			name: "single_port",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test-svc"},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{Name: "http", Port: 8080, TargetPort: intstr.FromInt(8080)},
					},
				},
			},
			expectedPort: 8080,
			expectError:  false,
		},
		{
			name: "multiple_ports_returns_zero",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test-svc"},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{Name: "http", Port: 8080, TargetPort: intstr.FromInt(8080)},
						{Name: "https", Port: 8443, TargetPort: intstr.FromInt(8443)},
					},
				},
			},
			expectedZero: true,
			expectError:  false,
		},
		{
			name: "no_ports_error",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test-svc"},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{},
				},
			},
			expectError:   true,
			expectedError: "has no ports defined",
		},
		{
			name: "named_single_port",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test-svc"},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{Name: "api", Port: 3000, TargetPort: intstr.FromString("api-port")},
					},
				},
			},
			expectedPort: 3000,
			expectError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			port, err := DetectServicePort(tc.service)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
				if tc.expectedZero {
					assert.Equal(t, int32(0), port)
				} else {
					assert.Equal(t, tc.expectedPort, port)
				}
			}
		})
	}
}
