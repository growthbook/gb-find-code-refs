package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/growthbook/gb-find-code-refs/coderefs"
	"github.com/growthbook/gb-find-code-refs/internal/log"
	"github.com/growthbook/gb-find-code-refs/internal/version"
	o "github.com/growthbook/gb-find-code-refs/options"
)

var extinctions = &cobra.Command{
	Use:     "extinctions",
	Example: "gb-find-code-refs extinctions",
	Short:   "Find and Post extinctions for branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := o.InitYAML()
		if err != nil {
			return err
		}

		opts, err := o.GetOptions()
		if err != nil {
			return err
		}
		err = opts.ValidateRequired()
		if err != nil {
			return err
		}

		log.Init(opts.Debug)
		coderefs.Run(opts, false)
		return nil
	},
}

var cmd = &cobra.Command{
	Use: "gb-find-code-refs",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := o.InitYAML()
		if err != nil {
			return err
		}

		opts, err := o.GetOptions()
		if err != nil {
			return err
		}
		err = opts.Validate()
		if err != nil {
			return err
		}

		log.Init(opts.Debug)
		coderefs.Run(opts, true)
		return nil
	},
	Version: version.Version,
}

func main() {
	if err := o.Init(cmd.PersistentFlags()); err != nil {
		panic(err)
	}

	cmd.AddCommand(extinctions)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
