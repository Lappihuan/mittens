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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var (
	// Set by CI.
	version = "dev"
	date    = "not_set" //nolint: gochecknoglobals
	commit  = "not_set" //nolint: gochecknoglobals
)

const (
	annotationOriginalTargetPort = "mittens.io/original-port"
	annotationConfigMap          = "mittens.io/proxy-config"
	annotationIsTapped           = "mittens.io/tapped"

	defaultImageHTTP = "ghcr.io/lappihuan/mittens-mitmproxy:latest"
	defaultImageRaw  = "ghcr.io/lappihuan/mittens-raw:latest"
	defaultImageGRPC = "ghcr.io/lappihuan/mittens-grpc:latest"
)

// die exit the program, printing the error.
func die(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func main() {
	exiter := &Exit{}
	rootCmd := NewRootCmd(exiter)

	kubernetesConfigFlags := genericclioptions.NewConfigFlags(false)
	kubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())

	config, err := kubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		die(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		die(err)
	}
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		die(err)
	}

	versionCmd := NewVersionCmd()
	tapCmd := &cobra.Command{
		Use:     "mittens <service>",
		Short:   "Enable mittens for a Service",
		Long:    "Proxy a Kubernetes Service with mitmproxy for interactive debugging",
		Example: "kubectl mittens -n my-namespace -p443 --https my-sample-service",
		PreRunE: bindTapFlags,
		RunE:    NewTapCommand(client, config, viper.GetViper()),
		Args:    cobra.ExactArgs(1),
	}

	tapCmd.Flags().StringP("port", "p", "", "target Service port")
	tapCmd.Flags().StringP("image", "i", defaultImageHTTP, "image to run in proxy container")
	tapCmd.Flags().Bool("https", false, "enable if target listener uses HTTPS")
	tapCmd.Flags().String("command-args", "mitmproxy", "specify command arguments for the proxy sidecar container")
	tapCmd.Flags().String("protocol", "http", "specify a protocol. Supported protocols: [ http ]")

	rootCmd.AddCommand(versionCmd, tapCmd)

	if err := rootCmd.Execute(); err != nil {
		exiter.Exit(1)
	}
}

// bindTapFlags is a workaround for https://github.com/spf13/viper/issues/233
func bindTapFlags(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlag("proxyPort", cmd.Flags().Lookup("port")); err != nil {
		return err
	}
	if err := viper.BindPFlag("proxyImage", cmd.Flags().Lookup("image")); err != nil {
		return err
	}
	if err := viper.BindPFlag("https", cmd.Flags().Lookup("https")); err != nil {
		return err
	}
	if err := viper.BindPFlag("commandArgs", cmd.Flags().Lookup("command-args")); err != nil {
		return err
	}
	if err := viper.BindPFlag("protocol", cmd.Flags().Lookup("protocol")); err != nil {
		return err
	}
	return nil
}

func NewRootCmd(e Exiter) *cobra.Command {
	return &cobra.Command{
		// HACK: there is a "bug" in cobra's handling of Use strings with spaces, so the space
		// below in the Use field isn't a true space -- it's a unicode Em space.
		// Also, if you can't see the space below prominently, you need to change your editor settings
		// to reveal Unicode characters. Otherwise, you're likely to PR some malicious code with a unicode
		// domain at some point.
		Use:   "kubectl mittens",
		Short: "mittens",
		Example: ` Proxy a Service with mitmproxy:
   kubectl mittens -n demo -p443 --https sample-service

 Show mittens version:
   kubectl mittens version`,
		Long: `mittens - proxy Services in Kubernetes with mitmproxy TUI.

 Mittens is a fork of kubetap by Soluble, redesigned around mitmproxy's interactive terminal UI.
 
 More information is available at the project website:
	 https://github.com/Lappihuan/mittens

 Original project (kubetap):
   https://github.com/soluble-ai/kubetap
`,
		Run: func(cmd *cobra.Command, _ []string) {
			// NOTE: explicitly print out usage here, but overriding it for subcommands by way
			// of SilenceUsage: true
			_ = cmd.Usage()
			e.Exit(64) // EX_USAGE
		},
		SilenceUsage: true,
	}
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version of mittens",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "version: %s, commit: %s, built at: %s\n", version, commit, date)
		},
	}
}

// Exiter exits the program, calling os.Exit(code), nothing more.
type Exiter interface {
	Exit(code int)
}

type Exit struct{}

func (e *Exit) Exit(code int) {
	os.Exit(code)
}
