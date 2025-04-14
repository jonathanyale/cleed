package cleed

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (r *Root) initMiniflux() {
	cmd := &cobra.Command{
		Use:   "miniflux",
		Short: "Push and pull feeds from Miniflux",
		Long: `Push and pull feeds from Miniflux

Examples:
  # Push feeds to Miniflux
  cleed miniflux push

  # Pull feeds from Miniflux
  cleed miniflux pull
`,

		RunE: r.RunMiniflux,
	}

	flags := cmd.Flags()
	flags.Bool("push", false, "push feeds to Miniflux")
	flags.Bool("pull", false, "pull feeds from Miniflux")

	r.Cmd.AddCommand(cmd)
}

func (r *Root) RunMiniflux(cmd *cobra.Command, args []string) error {
	if cmd.Flag("push").Changed {
		return r.feed.MinifluxPush()
	}
	if cmd.Flag("pull").Changed {
		return r.feed.MinifluxPull()
	}
	return fmt.Errorf("invalid command")
}
