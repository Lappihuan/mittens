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

	kubernetesConfigFlags := genericclioptions.NewConfigFlags(false)

	config, err := kubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		die(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		die(err)
	}

	rootCmd := &cobra.Command{
		Use:   "kubectl mittens [SERVICE] [OPTIONS]",
		Short: "mittens",
		Long: `mittens - proxy Services in Kubernetes with mitmproxy TUI.

Mittens is a fork of kubetap by Soluble, redesigned around mitmproxy's interactive terminal UI.

More information is available at the project website:
	https://github.com/Lappihuan/mittens

Original project (kubetap):
	https://github.com/soluble-ai/kubetap
`,
		Example: ` Proxy a Service with mitmproxy:
   kubectl mittens -n demo -p443 --https sample-service

 Show mittens version:
   kubectl mittens version`,
		SilenceUsage: true,
	}

	kubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		die(err)
	}

	// Add version subcommand
	versionCmd := NewVersionCmd()
	rootCmd.AddCommand(versionCmd)

	// Add flags to root command for direct usage
	rootCmd.Flags().StringP("port", "p", "", "target Service port (auto-detected if not provided)")
	rootCmd.Flags().StringP("image", "i", defaultImageHTTP, "image to run in proxy container")
	rootCmd.Flags().Bool("https", false, "enable if target listener uses HTTPS")
	rootCmd.Flags().String("command-args", "mitmproxy", "specify command arguments for the proxy sidecar container")
	rootCmd.Flags().String("protocol", "http", "specify a protocol. Supported protocols: [ http ]")

	// Handle root command with service as positional arg (kubectl mittens <service>)
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = cmd.Usage()
			exiter.Exit(64) // EX_USAGE
			return nil
		}
		// Bind flags and execute tap logic
		if err := bindTapFlags(cmd, args); err != nil {
			return err
		}
		return NewTapCommand(client, config, viper.GetViper())(cmd, args)
	}
	rootCmd.Args = cobra.ArbitraryArgs

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
