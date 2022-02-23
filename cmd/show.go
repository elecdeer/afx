package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/b4b4r07/afx/pkg/helpers/templates"
	"github.com/b4b4r07/afx/pkg/printers"
	"github.com/spf13/cobra"
)

type showCmd struct {
	metaCmd
}

var (
	// showLong is long description of show command
	showLong = templates.LongDesc(``)

	// showExample is examples for show command
	showExample = templates.Examples(`
		afx show
	`)
)

// newShowCmd creates a new show command
func (m metaCmd) newShowCmd() *cobra.Command {
	c := &showCmd{m}

	showCmd := &cobra.Command{
		Use:                   "show",
		Short:                 "Show packages managed by afx",
		Long:                  showLong,
		Example:               showExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run(args)
		},
	}

	return showCmd
}

func (c *showCmd) run(args []string) error {
	w := printers.GetNewTabWriter(os.Stdout)
	headers := []string{"NAME", "TYPE", "STATUS"}

	type Item struct {
		Name   string
		Type   string
		Status string
	}

	var items []Item
	for _, pkg := range c.state.Additions {
		items = append(items, Item{
			Name:   pkg.Name,
			Type:   pkg.Type,
			Status: "WaitingInstall",
		})
	}
	for _, pkg := range c.state.Changes {
		items = append(items, Item{
			Name:   pkg.Name,
			Type:   pkg.Type,
			Status: "WaitingUpdate",
		})
	}
	for _, pkg := range c.state.Deletions {
		items = append(items, Item{
			Name:   pkg.Name,
			Type:   pkg.Type,
			Status: "WaitingUninstall",
		})
	}
	for _, pkg := range c.state.NoChanges {
		items = append(items, Item{
			Name:   pkg.Name,
			Type:   pkg.Type,
			Status: "Installed",
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	fmt.Fprintf(w, strings.Join(headers, "\t")+"\n")
	for _, item := range items {
		fields := []string{
			item.Name, item.Type, item.Status,
		}
		fmt.Fprintf(w, strings.Join(fields, "\t")+"\n")
	}

	return w.Flush()
}
