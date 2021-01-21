package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/vetyy/kubetools/internal/cli/genericflags"
	log "github.com/vetyy/kubetools/internal/logging"
)

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kubetools",
		Short: "kubetools is CLI tool for managing Kubernetes clusters",
	}

	var cmdCompletion = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run:

source <(kubetools completion)

To configure your bash shell to load completions for each session add to your bashrc:

# ~/.bashrc or ~/.bash_profile
source <(kubetools completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = rootCmd.GenBashCompletion(os.Stdout)
		},
	}


	genericOptions, err := genericflags.NewGenericFlags()
	if err != nil {
		log.Fatalf("failed to create generic options: %v", err)
	}
	rootCmd.AddCommand(cmdCompletion)
	rootCmd.AddCommand(NewCmdUpdateContext(genericOptions))
	return rootCmd
}

// Execute executes CLI application
func Execute() {
	if err := NewCmdRoot().Execute(); err != nil {
		os.Exit(1)
	}
}

