package cleed

import (
	"github.com/radulucut/cleed/internal"
	"github.com/spf13/cobra"
)

func (r *Root) initExplore() {
	cmd := &cobra.Command{
		Use:   "explore",
		Short: "Explore feeds",
		Long: `Explore feeds from a repository

Examples:
  # Explore feeds from the default repository
  cleed explore

  # Fetch the latest changes and explore feeds from the default repository
  cleed explore --update

  # Explore feeds from a repository
  cleed explore https://github.com/radulucut/cleed-explore.git

  # Limit the number of items to display from each list
  cleed explore --limit 5

  # Search for items (title, description)
  cleed explore --search "news"

  # Import all feeds into my feeds
  cleed explore --import --limit 0

  # Import feeds from search results
  cleed explore --import --search "news"

  # Remove a repository
  cleed explore https://github.com/radulucut/cleed-explore.git --remove
`,

		RunE: r.RunExplore,
	}

	flags := cmd.Flags()
	flags.Bool("import", false, "import feeds")
	flags.Uint("limit", 10, "limit the number of items to display from each list")
	flags.String("search", "", "search for items (title, description)")
	flags.BoolP("update", "u", false, "fetch the latest changes")
	flags.Bool("remove", false, "remove the repository")

	r.Cmd.AddCommand(cmd)
}

func (r *Root) RunExplore(cmd *cobra.Command, args []string) error {
	if cmd.Flag("remove").Changed {
		url := ""
		if len(args) > 0 {
			url = args[0]
		}
		return r.feed.ExploreRemove(url)
	}
	limit, err := cmd.Flags().GetUint("limit")
	if err != nil {
		return err
	}
	opts := &internal.ExploreOptions{
		Limit:  int(limit),
		Update: cmd.Flag("update").Changed,
		Query:  cmd.Flag("search").Value.String(),
	}
	if len(args) > 0 {
		opts.Url = args[0]
	}
	if opts.Query != "" {
		if !cmd.Flag("limit").Changed {
			opts.Limit = 25
		}
		opts.Import, err = cmd.Flags().GetBool("import")
		if err != nil {
			return err
		}
		return r.feed.ExploreSearch(opts)
	}
	if cmd.Flag("import").Changed {
		return r.feed.ExploreImport(opts)
	}
	return r.feed.Explore(opts)
}
