package internal

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/radulucut/cleed/internal/utils"
)

type ExploreOptions struct {
	Url    string
	Update bool
	Limit  int
	Query  string
}

func (f *TerminalFeed) ExploreRemove(url string) error {
	if url == "" {
		url = f.defaultExploreRepository
	}
	err := f.storage.RemoveExploreRepository(url)
	if err != nil {
		return utils.NewInternalError("failed to remove repository: " + err.Error())
	}
	f.printer.Printf("%s was removed\n", url)
	return nil
}

type ExploreSearchItem struct {
	ListColor uint8
	Outline   *utils.OPMLOutline
	Score     int
}

func (f *TerminalFeed) ExploreSearch(opts *ExploreOptions) error {
	queryTokens := utils.Tokenize(opts.Query, nil)
	if len(queryTokens) == 0 {
		return utils.NewInternalError("query is empty")
	}
	lists, err := f.getExploreFeeds(opts.Url, opts.Update)
	if err != nil {
		return err
	}
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	items := make([]*ExploreSearchItem, 0)
	for i, list := range lists {
		for _, outline := range list.Outlines {
			if outline.XMLURL == "" {
				continue
			}
			itemTokens := utils.Tokenize(outline.Text, nil)
			itemTokens = utils.Tokenize(outline.Description, itemTokens)
			score := utils.Score(queryTokens, itemTokens)
			if score == -1 {
				continue
			}
			items = append(items, &ExploreSearchItem{
				ListColor: mapColor(uint8(i+1%256), config),
				Outline:   outline,
				Score:     score,
			})
		}
	}
	if len(items) == 0 {
		f.printer.Printf("no results found for %s\n", opts.Query)
	}
	slices.SortFunc(items, func(a, b *ExploreSearchItem) int {
		if a.Score == b.Score {
			return strings.Compare(a.Outline.Text, b.Outline.Text)
		} else if a.Score < b.Score {
			return -1
		}
		return 1
	})
	secondaryTextColor := mapColor(7, config)
	totalDisplayed := 0
	l := min(opts.Limit, len(items))
	if l <= 0 {
		l = len(items)
	}
	for i := l - 1; i >= 0; i-- {
		item := items[i]
		if item.Outline.XMLURL == "" {
			continue
		}
		if item.Outline.Text != "" {
			f.printer.Print(f.printer.ColorForeground(item.Outline.Text, item.ListColor), "\n")
		}
		if item.Outline.Description != "" {
			f.printer.Print(item.Outline.Description, "\n")
		}
		f.printer.Print(f.printer.ColorForeground(item.Outline.XMLURL, secondaryTextColor), "\n\n")
		totalDisplayed++
	}
	f.printer.Printf("Displayed %d out of %s\n", totalDisplayed, utils.Pluralize(int64(len(items)), "item"))
	return nil
}

func (f *TerminalFeed) Explore(opts *ExploreOptions) error {
	lists, err := f.getExploreFeeds(opts.Url, opts.Update)
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		f.printer.Printf("No feeds found\n")
		return nil
	}
	slices.SortFunc(lists, func(a, b *utils.OPMLOutline) int {
		return strings.Compare(a.Text, b.Text)
	})
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	secondaryTextColor := mapColor(7, config)
	total := 0
	totalDisplayed := 0
	for i := len(lists) - 1; i >= 0; i-- {
		list := lists[i]
		listColor := mapColor(uint8(i+1%256), config)
		for j, outline := range list.Outlines {
			if outline.XMLURL == "" {
				continue
			}
			if opts.Limit > 0 && j >= opts.Limit {
				break
			}
			if outline.Text != "" {
				f.printer.Print(f.printer.ColorForeground(outline.Text, listColor), "\n")
			}
			if outline.Description != "" {
				f.printer.Print(outline.Description, "\n")
			}
			f.printer.Print(f.printer.ColorForeground(outline.XMLURL, secondaryTextColor), "\n\n")
		}
		displayed := len(list.Outlines)
		if opts.Limit > 0 && opts.Limit < displayed {
			displayed = opts.Limit
		}
		totalDisplayed += displayed
		total += len(list.Outlines)
		f.printer.Printf("%s (%d/%d)\n--------------------------------\n",
			list.Text,
			displayed,
			len(list.Outlines),
		)
	}
	f.printer.Printf("Displayed %d out of %s from %s\n", totalDisplayed, utils.Pluralize(int64(total), "feed"), utils.Pluralize(int64(len(lists)), "list"))
	return nil
}

func (f *TerminalFeed) ExploreImport(opts *ExploreOptions) error {
	lists, err := f.getExploreFeeds(opts.Url, opts.Update)
	if err != nil {
		return err
	}
	urls := make([]string, 0)
	totalImported := 0
	for _, list := range lists {
		for _, outline := range list.Outlines {
			if outline.XMLURL == "" {
				continue
			}
			if opts.Limit > 0 && len(urls) >= opts.Limit {
				break
			}
			urls = append(urls, outline.XMLURL)
		}
		totalImported += len(urls)
		err = f.storage.AddToList(urls, list.Text)
		if err != nil {
			f.printer.Printf("failed to add feeds to list %s: %v\n", list.Text, err)
		}
		urls = urls[:0]
	}
	f.printer.Printf("Imported %s from %s\n", utils.Pluralize(int64(totalImported), "feed"), utils.Pluralize(int64(len(lists)), "list"))
	return nil
}

func (f *TerminalFeed) getExploreFeeds(url string, update bool) ([]*utils.OPMLOutline, error) {
	if url == "" {
		url = f.defaultExploreRepository
	}
	root, err := f.storage.GetExploreRepositoryPath(url, update)
	if err != nil {
		return nil, err
	}
	lists := make([]*utils.OPMLOutline, 0)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".opml" {
			opml, err := utils.ParseOPMLFile(path)
			if err != nil {
				p := strings.TrimPrefix(path, root)
				f.printer.Printf("failed to parse OPML file %s: %v\n", p, err)
				return nil
			}
			lists = append(lists, opml.Body.Outltines...)
			return nil
		}
		return nil
	})
	return lists, err
}
