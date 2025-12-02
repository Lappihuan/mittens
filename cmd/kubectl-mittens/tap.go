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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	k8sappsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	mittensContainerName   = "mittens"
	mittensServicePortName = "mittens-web"
	mittensPortName        = "mittens-listen"
	mittensWebPortName     = "mittens-web"
	mittensProxyListenPort = 7777
	mittensConfigMapPrefix = "mittens-target-"

	interactiveTimeoutSeconds = 90
	configMapAnnotationPrefix = "target-"

	// Flux drift detection annotation to prevent automatic rollback of Service mutations.
	fluxDriftDetectionAnnotation = "helm.toolkit.fluxcd.io/driftDetection"
	fluxDriftDetectionDisabled   = "disabled"

	protocolHTTP Protocol = "http"
	protocolTCP  Protocol = "tcp"
	protocolUDP  Protocol = "udp"
	protocolGRPC Protocol = "grpc"
)

var (
	ErrNamespaceNotExist          = errors.New("the provided Namespace does not exist")
	ErrServiceMissingPort         = errors.New("the target Service does not have the provided port")
	ErrServiceSelectorNoMatch     = errors.New("the Service selector did not match any Deployments")
	ErrServiceSelectorMultiMatch  = errors.New("the Service selector matched multiple Deployments")
	ErrDeploymentOutsideNamespace = errors.New("the Service selector matched Deployment outside the specified Namespace")
	ErrSelectorsMissing           = errors.New("no selectors are set for the target Service")
	ErrConfigMapNoMatch           = errors.New("the ConfigMap list did not match any ConfigMaps")
	ErrMittensPodNoMatch          = errors.New("a Mittens Pod was not found")
	ErrCreateResourceMismatch     = errors.New("the created resource did not match the desired state")
	ErrDeploymentMissingPorts     = errors.New("error resolving Service port number by name from Deployment")
)

// Protocol is a supported tap method, and ultimately determines what container
// is injected as a sidecar.
type Protocol string

// Tap is a method of implementing a "Tap" for a Kubernetes cluster.
type Tap interface {
	// Sidecar produces a sidecar container to be added to a
	// Deployment.
	Sidecar(string) v1.Container

	// PatchDeployment tweaks a Deployment after a Sidecar is added
	// during the tap process.
	// Example: mitmproxy calls this function to configure the ConfigMap volume refs.
	PatchDeployment(*k8sappsv1.Deployment)

	// ReadyEnv and UnreadyEnv are used to prepare the environment
	// with resources that will be necessary for the sidecar, but do
	// not exist within a given Deployment.
	// Example: mitmproxy calls this function to apply and remove ConfigMaps for mitmproxy.
	ReadyEnv() error
	UnreadyEnv() error

	// String prints the tap method, be it mitmproxy, tcpdump, etc.
	String() string

	// Protocols returns a slice of protocols supported by the tap.
	Protocols() []Protocol
}

// ProxyOptions are options used to configure the Tap implementation.
type ProxyOptions struct {
	// Target is the target Service
	Target string `json:"target"`
	// Protocol is the protocol type, one of [http, https]
	Protocol Protocol `json:"protocol"`
	// UpstreamHTTPS should be set to true if the target is using HTTPS
	UpstreamHTTPS bool `json:"upstreamHttps"`
	// UpstreamPort is the listening port for the target Service
	UpstreamPort string `json:"upstreamPort"`
	// Mode is the proxy mode. Only "reverse" is currently supported.
	Mode string `json:"mode"`
	// Namespace is the namespace that the Service and Deployment are in
	Namespace string `json:"namespace"`
	// Image is the proxy image to deploy as a sidecar
	Image string `json:"image"`

	// dplName tracks the current deployment target
	dplName string
}

// NewTapCommand identifies a target employment through service selectors and modifies that
// deployment to add a proxy sidecar.
func NewTapCommand(client kubernetes.Interface, _ *rest.Config, viper *viper.Viper) func(*cobra.Command, []string) error { //nolint: gocyclo
	return func(cmd *cobra.Command, args []string) error {
		targetSvcName := args[0]

		protocol := viper.GetString("protocol")
		targetSvcPort := viper.GetInt32("proxyPort")
		namespace := viper.GetString("namespace")
		image := viper.GetString("proxyImage")
		https := viper.GetBool("https")

		commandArgs := strings.Fields(viper.GetString("commandArgs"))
		if namespace == "" {
			// TODO: There is probably a way to get the default namespace from the
			// client context, but I'm not sure what that API is. Will dig
			// for that at some point.
			// BUG: "default" is not always the "correct default".
			viper.Set("namespace", "default")
			namespace = "default"
		}
		exists, err := hasNamespace(client, namespace)
		if err != nil {
			return fmt.Errorf("error fetching namespaces: %w", err)
		}
		if !exists {
			return ErrNamespaceNotExist
		}

		// Get the service first for port auto-detection
		targetService, err := client.CoreV1().Services(namespace).Get(context.TODO(), targetSvcName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Auto-detect or prompt for port if not provided
		if targetSvcPort == 0 {
			var detectionErr error
			targetSvcPort, detectionErr = DetectServicePort(targetService)
			if detectionErr != nil {
				return detectionErr
			}

			// If multiple ports, prompt user to select
			if targetSvcPort == 0 {
				selectedPort, err := InteractivePortSelection(targetService)
				if err != nil {
					return err
				}
				targetSvcPort = selectedPort
			}

			viper.Set("proxyPort", targetSvcPort)
		}

		proxyOpts := ProxyOptions{
			Target:        targetSvcName,
			UpstreamHTTPS: https,
			Mode:          "reverse", // eventually this may be configurable
			Namespace:     namespace,
		}
		// Adjust default image by protocol if not manually set
		if image == defaultImageHTTP {
			switch Protocol(protocol) { //nolint: exhaustive
			case protocolTCP, protocolUDP:
				// TODO: make this container and remove error
				image = defaultImageRaw
				return fmt.Errorf("mode %q is currently not supported", image)
			case protocolGRPC:
				// TODO: make this container and remove error
				image = defaultImageGRPC
				return fmt.Errorf("mode %q is currently not supported", image)
			}
			viper.Set("proxyImage", image)
		}

		deploymentsClient := client.AppsV1().Deployments(namespace)
		servicesClient := client.CoreV1().Services(namespace)
		podsClient := client.CoreV1().Pods(namespace)

		// Check if this service is already tapped
		anns := targetService.GetAnnotations()
		alreadyTapped := anns[annotationOriginalTargetPort] != ""

		if !alreadyTapped {
			if err := performTap(cmd, client, deploymentsClient, servicesClient, targetService, targetSvcName, targetSvcPort, image, commandArgs, protocol, proxyOpts, viper); err != nil {
				return err
			}
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout())
		if !alreadyTapped {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Port %d of Service %q has been tapped!\n", targetSvcPort, targetSvcName)
		} else {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Service already tapped. Attaching to existing mitmproxy session...")
		}

		// Only wait for pod and exec when explicitly requested by running from a terminal
		// Check if stdout is going to a terminal (not pipes/redirects)
		outFile, isTerminal := cmd.OutOrStdout().(*os.File)
		if !isTerminal || outFile == nil {
			// Output is redirected or in a test, skip the waiting/exec
			return nil
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nWaiting for pod to start...\n\n")
		ic := make(chan os.Signal, 1)
		signal.Notify(ic, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-ic
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Stopping mittens...")
			_ = NewUntapCommand(client, viper)(cmd, args)
			die()
		}()

		s := make(chan struct{})
		defer close(s)
		go func() {
			// Skip the first few checks to give pods time to come up.
			// Race: If the first few cycles are not skipped, the condition status may be "Ready".
			time.Sleep(5 * time.Second)
			s <- struct{}{}
		}()

		// Use spinner instead of progress bar
		spinner := NewSpinner("Waiting for Pod containers to become ready...")

		var ready bool
		for range interactiveTimeoutSeconds {
			if ready {
				break
			}
			// Check if context was cancelled (e.g., in tests)
			if cmd.Context() != nil {
				select {
				case <-cmd.Context().Done():
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), "")
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Context cancelled. Stopping mittens...")
					spinner.Fail("Cancelled")
					_ = NewUntapCommand(client, viper)(cmd, args)
					return cmd.Context().Err()
				default:
				}
			}
			time.Sleep(1 * time.Second)
			select {
			case <-time.After(1 * time.Nanosecond):
				// if not ready this cycle, abort
				continue
			case <-s:
				dp, err := deploymentFromSelectors(deploymentsClient, targetService.Spec.Selector)
				if err != nil {
					spinner.Fail("Error getting deployment")
					return err
				}
				pod, err := mittensPod(podsClient, dp.Name)
				if err != nil {
					spinner.Fail("Error getting pod")
					return err
				}
				for _, cond := range pod.Status.Conditions {
					if cond.Type == "ContainersReady" {
						if cond.Status == "True" {
							ready = true
						}
					}
				}
				go func() {
					s <- struct{}{}
				}()
			}
		}
		if !ready {
			spinner.Fail("Pod not running after 90 seconds. Cancelling.")
			die()
		}
		spinner.Stop("Pod ready!")
		dp, err := deploymentFromSelectors(deploymentsClient, targetService.Spec.Selector)
		if err != nil {
			return err
		}
		pod, err := mittensPod(podsClient, dp.Name)
		if err != nil {
			return err
		}

		// Spawn kubectl exec to attach to mitmproxy tmux session
		execCmd := exec.CommandContext(cmd.Context(), "kubectl", "exec", "-it", pod.Name, "-n", namespace, "-c", "mittens", "--", "tmux", "attach-session", "-t", "mitmproxy")
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err = execCmd.Run()

		// User has exited the tmux session, clean up the tap
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "")
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Cleaning up litter...")
		untapErr := NewUntapCommand(client, viper)(cmd, args)
		if untapErr != nil {
			return untapErr
		}
		return err
	}
}

// performTap handles the actual tapping logic for a service.
func performTap(cmd *cobra.Command, client kubernetes.Interface, deploymentsClient appsv1.DeploymentInterface, servicesClient corev1.ServiceInterface, targetService *v1.Service, targetSvcName string, targetSvcPort int32, image string, commandArgs []string, protocol string, proxyOpts ProxyOptions, v *viper.Viper) error {
	// set the upstream port so the proxy knows where to forward traffic
	for _, ports := range targetService.Spec.Ports {
		if ports.Port != targetSvcPort {
			continue
		}
		if ports.TargetPort.Type == intstr.Int {
			proxyOpts.UpstreamPort = ports.TargetPort.String()
		}
		// if named, must determine port from deployment spec
		if ports.TargetPort.Type == intstr.String {
			var err error
			targetDpl, err := deploymentFromSelectors(deploymentsClient, targetService.Spec.Selector)
			if err != nil {
				return fmt.Errorf("error resolving Deployment from Service selectors while setting proxy ports: %w", err)
			}
			for _, c := range targetDpl.Spec.Template.Spec.Containers {
				for _, p := range c.Ports {
					if p.Name == ports.TargetPort.String() {
						// Set the upstream (target) Service port
						proxyOpts.UpstreamPort = strconv.Itoa(int(p.ContainerPort))
					}
				}
			}
			if proxyOpts.UpstreamPort == "" {
				return ErrDeploymentMissingPorts
			}
		}
	}

	targetDpl, err := deploymentFromSelectors(deploymentsClient, targetService.Spec.Selector)
	if err != nil {
		return fmt.Errorf("error resolving Deployment from Service selectors: %w", err)
	}
	proxyOpts.dplName = targetDpl.Name

	// Get a proxy based on the protocol type
	var proxy Tap
	switch Protocol(protocol) { //nolint: exhaustive
	case protocolTCP, protocolUDP:
	case protocolGRPC:
	default:
		// AKA, case protocolHTTP:
		proxy = NewMitmproxy(client, proxyOpts)
	}

	// Prepare the environment (configmaps, secrets, volumes, etc).
	if err := proxy.ReadyEnv(); err != nil {
		return err
	}

	// Setup the sidecar
	dpl, err := deploymentFromSelectors(deploymentsClient, targetService.Spec.Selector)
	if err != nil {
		return err
	}

	sidecar := proxy.Sidecar(dpl.Name)
	sidecar.Image = image
	sidecar.Args = commandArgs

	// Apply the Deployment configuration
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		dpl.Spec.Template.Spec.Containers = append(dpl.Spec.Template.Spec.Containers, sidecar)
		proxy.PatchDeployment(&dpl)
		// set annotation on pod to know what pods are tapped
		anns := dpl.Spec.Template.GetAnnotations()
		if anns == nil {
			anns = map[string]string{}
		}
		anns[annotationIsTapped] = dpl.Name
		dpl.Spec.Template.SetAnnotations(anns)
		_, updateErr := deploymentsClient.Update(context.TODO(), &dpl, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Error modifying Deployment, reverting tap...")
		args := []string{targetSvcName}
		_ = NewUntapCommand(client, v)(cmd, args)
		return fmt.Errorf("failed to add sidecars to Deployment: %w", retryErr)
	}

	// Tap the Service to redirect the incoming traffic to our proxy
	if err := tapSvc(servicesClient, targetSvcName, targetSvcPort); err != nil {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Error modifying Service, reverting tap...")
		args := []string{targetSvcName}
		_ = NewUntapCommand(client, v)(cmd, args)
		return err
	}

	return nil
}

// NewUntapCommand unconditionally removes all proxies, taps, and artifacts. This is
// the inverse of NewTapCommand.
func NewUntapCommand(client kubernetes.Interface, viper *viper.Viper) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		targetSvcName := args[0]
		namespace := viper.GetString("namespace")
		if namespace == "" {
			namespace = "default"
		}
		exists, err := hasNamespace(client, namespace)
		if err != nil {
			return fmt.Errorf("error fetching namespaces: %w", err)
		}
		if !exists {
			return ErrNamespaceNotExist
		}

		deploymentsClient := client.AppsV1().Deployments(namespace)
		servicesClient := client.CoreV1().Services(namespace)

		targetService, err := client.CoreV1().Services(namespace).Get(context.TODO(), targetSvcName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		dpl, err := deploymentFromSelectors(deploymentsClient, targetService.Spec.Selector)
		if err != nil {
			return err
		}
		if dpl.Namespace != namespace {
			panic(ErrDeploymentOutsideNamespace)
		}

		proxy := NewMitmproxy(client, ProxyOptions{
			Namespace: namespace,
			Target:    targetSvcName,
			dplName:   dpl.Name,
		})

		if err := proxy.UnreadyEnv(); err != nil {
			// both error types below can be thrown
			if !errors.Is(ErrConfigMapNoMatch, err) {
				return err
			}
		}

		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Explicitly re-fetch the deployment to reduce the chance of having a race
			deployment, getErr := deploymentsClient.Get(context.TODO(), dpl.Name, metav1.GetOptions{})
			if getErr != nil {
				return getErr
			}
			var containersNoProxy []v1.Container
			for _, c := range deployment.Spec.Template.Spec.Containers {
				if c.Name != mittensContainerName {
					containersNoProxy = append(containersNoProxy, c)
				}
			}
			deployment.Spec.Template.Spec.Containers = containersNoProxy
			var volumes []v1.Volume
			for _, v := range deployment.Spec.Template.Spec.Volumes {
				if !strings.HasPrefix(v.Name, "mittens") {
					volumes = append(volumes, v)
				}
			}
			deployment.Spec.Template.Spec.Volumes = volumes
			anns := deployment.Spec.Template.GetAnnotations()
			if anns != nil {
				delete(anns, annotationIsTapped)
				deployment.Spec.Template.SetAnnotations(anns)
			}
			_, updateErr := deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return fmt.Errorf("failed to remove sidecars from Deployment: %w", retryErr)
		}
		if err := untapSvc(servicesClient, targetSvcName); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Untapped Service %q\n", targetSvcName)
		return nil
	}
}

// deploymentFromSelectors returns a deployment given selector labels.
func deploymentFromSelectors(deploymentsClient appsv1.DeploymentInterface, selectors map[string]string) (k8sappsv1.Deployment, error) {
	var sel string
	switch len(selectors) {
	case 0:
		return k8sappsv1.Deployment{}, ErrSelectorsMissing
	case 1:
		for k, v := range selectors {
			sel = k + "=" + v
		}
	default:
		for k, v := range selectors {
			sel = strings.Join([]string{sel, k + "=" + v}, ",")
		}
		sel = strings.TrimLeft(sel, ",")
	}
	dpls, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: sel,
	})
	if err != nil {
		return k8sappsv1.Deployment{}, err
	}
	switch len(dpls.Items) {
	case 0:
		return k8sappsv1.Deployment{}, ErrServiceSelectorNoMatch
	case 1:
		return dpls.Items[0], nil
	default:
		return k8sappsv1.Deployment{}, ErrServiceSelectorMultiMatch
	}
}

// mittensPod returns a mittens pod matching a given Deployment name and Namespace.
func mittensPod(podClient corev1.PodInterface, deploymentName string) (v1.Pod, error) {
	pods, err := podClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return v1.Pod{}, err
	}
	var p v1.Pod
	for _, pod := range pods.Items {
		anns := pod.GetAnnotations()
		if anns == nil {
			continue
		}
		for k, v := range anns {
			if k == annotationIsTapped && v == deploymentName {
				return pod, nil
			}
		}
	}
	return p, ErrMittensPodNoMatch
}

// tapSvc modifies a target port to point to a new proxy service.
func tapSvc(svcClient corev1.ServiceInterface, svcName string, targetPort int32) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		svc, getErr := svcClient.Get(context.TODO(), svcName, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		anns := svc.GetAnnotations()
		// If anns is nil, it means that the target Service had no annotations.
		// BUG: this is probably safe if the comment above is true, but if it isn't
		// true this could wipe the annotations associated with a Service, which
		// would be highly undesirable....
		if anns == nil {
			anns = make(map[string]string)
		}

		var targetSvcPort v1.ServicePort
		var hasPort bool
		for _, sp := range svc.Spec.Ports {
			if sp.Port == targetPort {
				hasPort = true
				targetSvcPort = sp
			}
		}
		if !hasPort {
			return ErrServiceMissingPort
		}

		anns[annotationOriginalTargetPort] = targetSvcPort.TargetPort.String()
		// Add Flux drift detection annotation to prevent automatic rollback
		anns[fluxDriftDetectionAnnotation] = fluxDriftDetectionDisabled
		svc.SetAnnotations(anns)

		// then do the swap and build a new ports list
		var servicePorts []v1.ServicePort
		for _, sp := range svc.Spec.Ports {
			if sp.Port == targetSvcPort.Port {
				if sp.Name == "" {
					sp.Name = mittensPortName
				}
				sp.TargetPort = intstr.FromInt(mittensProxyListenPort)
			}
			servicePorts = append(servicePorts, sp)
		}
		svc.Spec.Ports = servicePorts

		_, updateErr := svcClient.Update(context.TODO(), svc, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return fmt.Errorf("failed to tap Service: %w", retryErr)
	}
	return nil
}

// untapSvc modifies a target port point to the original service, not our proxy sidecar.
func untapSvc(svcClient corev1.ServiceInterface, svcName string) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		svc, getErr := svcClient.Get(context.TODO(), svcName, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		// NOTE: it is critical to Parse here (vs FromString)
		origSvcTargetPort := intstr.Parse(svc.GetAnnotations()[annotationOriginalTargetPort])
		var servicePorts []v1.ServicePort
		for _, sp := range svc.Spec.Ports {
			if sp.Name == mittensServicePortName {
				continue
			}
			if sp.TargetPort.IntValue() == mittensProxyListenPort {
				if sp.Name == mittensPortName {
					sp.Name = ""
				}
				sp.TargetPort = origSvcTargetPort
			}
			servicePorts = append(servicePorts, sp)
		}
		svc.Spec.Ports = servicePorts
		anns := svc.GetAnnotations()
		newAnns := make(map[string]string)
		for k, v := range anns {
			// Remove mittens and Flux annotations added during tap
			if k != annotationOriginalTargetPort && k != fluxDriftDetectionAnnotation {
				newAnns[k] = v
			}
		}
		svc.SetAnnotations(newAnns)
		_, updateErr := svcClient.Update(context.TODO(), svc, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return fmt.Errorf("failed to untap Service: %w", retryErr)
	}
	return nil
}

// hasNamespace checks if a given Namespace exists.
func hasNamespace(client kubernetes.Interface, namespace string) (bool, error) {
	if namespace == "" {
		return false, os.ErrInvalid
	}
	ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, n := range ns.Items {
		if n.Name == namespace {
			return true, nil
		}
	}
	return false, nil
}
