package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/b4b4r07/afx/pkg/errors"
	"github.com/b4b4r07/afx/pkg/helpers/templates"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type stateCmd struct {
	metaCmd

	opt stateOpt
}

type stateOpt struct {
	force bool
}

var (
	// stateLong is long description of state command
	stateLong = templates.LongDesc(``)

	// stateExample is examples for state command
	stateExample = templates.Examples(``)
)

// newStateCmd creates a new state command
func (m metaCmd) newStateCmd() *cobra.Command {
	c := &stateCmd{metaCmd: m}

	stateCmd := &cobra.Command{
		Use:                   "state [list|refresh|remove]",
		Short:                 "Advanced state management",
		Long:                  stateLong,
		Example:               stateExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(1),
		Hidden:                true,
	}

	stateCmd.AddCommand(
		c.newStateListCmd(),
		c.newStateRefreshCmd(),
		c.newStateRemoveCmd(),
	)

	return stateCmd
}

func (c stateCmd) newStateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "List your state items",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := c.State.List()
			if err != nil {
				return err
			}
			for _, item := range items {
				fmt.Println(item)
			}
			return nil
		},
	}
}

func (c stateCmd) newStateRefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "refresh",
		Short:                 "Refresh your state file",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.opt.force {
				return c.State.New()
			}
			if err := c.State.Refresh(); err != nil {
				return errors.Wrap(err, "failed to refresh state")
			}
			fmt.Println(color.WhiteString("Successfully refreshed"))
			return nil
		},
	}
	cmd.Flags().BoolVarP(&c.opt.force, "force", "", false, "force update")
	return cmd
}

func (c stateCmd) newStateRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "remove",
		Short:                 "Remove selected packages from state file",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Aliases:               []string{"rm"},
		Args:                  cobra.MinimumNArgs(0),
		ValidArgs:             getNameInPackages(c.State.NoChanges),
		RunE: func(cmd *cobra.Command, args []string) error {
			var resources []string
			switch len(cmd.Flags().Args()) {
			case 0:
				list, err := c.State.List()
				if err != nil {
					return errors.Wrap(err, "failed to list state items")
				}
				var selected string
				if err := survey.AskOne(&survey.Select{
					Message: "Choose a package:",
					Options: list,
				}, &selected); err != nil {
					return errors.Wrap(err, "failed to get input from console")
				}
				resources = append(resources, selected)
			default:
				// TODO: check valid or invalid
				resources = cmd.Flags().Args()
			}
			for _, resource := range resources {
				id := c.State.ToID(resource)
				c.State.Remove(id)
			}
			return nil
		},
	}
}
