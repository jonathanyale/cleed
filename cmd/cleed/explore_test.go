package cleed

import (
	"bytes"
	"os"
	"path"
	"testing"
	"time"

	"github.com/radulucut/cleed/internal"
	_storage "github.com/radulucut/cleed/internal/storage"
	"github.com/radulucut/cleed/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Explore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(tempRepository, 0755)

	err = os.WriteFile(path.Join(tempRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `RSS Feed
RSS Feed description
https://rss-feed.com/rss

Atom Feed
Atom Feed description
https://atom-feed.com/atom

test 2 (2/2)
--------------------------------
RSS Feed
RSS Feed description
https://rss-feed.com/rss

Atom Feed
Atom Feed description
https://atom-feed.com/atom

https://test.com

test (3/3)
--------------------------------
Displayed 5 out of 5 feeds from 2 lists
`, out.String())
}

func Test_Explore_Custom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(tempRepository, 0755)

	customRepository, err := storage.JoinExploreDir("custom")
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(customRepository, 0755)

	err = os.WriteFile(path.Join(customRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore", "custom"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `RSS Feed
RSS Feed description
https://rss-feed.com/rss

Atom Feed
Atom Feed description
https://atom-feed.com/atom

test 2 (2/2)
--------------------------------
RSS Feed
RSS Feed description
https://rss-feed.com/rss

Atom Feed
Atom Feed description
https://atom-feed.com/atom

https://test.com

test (3/3)
--------------------------------
Displayed 5 out of 5 feeds from 2 lists
`, out.String())
}

func Test_Explore_Limit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(tempRepository, 0755)

	err = os.WriteFile(path.Join(tempRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore", "--limit", "1"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `RSS Feed
RSS Feed description
https://rss-feed.com/rss

test 2 (1/2)
--------------------------------
RSS Feed
RSS Feed description
https://rss-feed.com/rss

test (1/3)
--------------------------------
Displayed 2 out of 5 feeds from 2 lists
`, out.String())
}

func Test_Explore_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(tempRepository, 0755)

	err = os.WriteFile(path.Join(tempRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom 1 Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom 2 Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore", "--search", "Atom"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `Atom 2 Feed
Atom Feed description
https://atom-feed.com/atom

Atom 1 Feed
Atom Feed description
https://atom-feed.com/atom

Displayed 2 out of 2 items
`, out.String())
}

func Test_Explore_Import(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(tempRepository, 0755)

	err = os.WriteFile(path.Join(tempRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore", "--import"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "Imported 5 feeds from 2 lists\n", out.String())

	lists, err := storage.LoadLists()
	assert.NoError(t, err)

	expectedLists := []string{"test", "test 2"}
	assert.Equal(t, expectedLists, lists)

	feeds, err := storage.GetFeedsFromList("test")
	assert.NoError(t, err)
	expectedFeeds := []*_storage.ListItem{
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://rss-feed.com/rss"},
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://atom-feed.com/atom"},
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://test.com"},
	}
	assert.Equal(t, expectedFeeds, feeds)

	feeds, err = storage.GetFeedsFromList("test 2")
	assert.NoError(t, err)
	expectedFeeds = []*_storage.ListItem{
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://rss-feed.com/rss"},
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://atom-feed.com/atom"},
	}
	assert.Equal(t, expectedFeeds, feeds)
}

func Test_Explore_Import_With_Limit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(tempRepository, 0755)

	err = os.WriteFile(path.Join(tempRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore", "--import", "--limit", "1"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "Imported 2 feeds from 2 lists\n", out.String())

	lists, err := storage.LoadLists()
	assert.NoError(t, err)

	expectedLists := []string{"test", "test 2"}
	assert.Equal(t, expectedLists, lists)

	feeds, err := storage.GetFeedsFromList("test")
	assert.NoError(t, err)
	expectedFeeds := []*_storage.ListItem{
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://rss-feed.com/rss"},
	}
	assert.Equal(t, expectedFeeds, feeds)

	feeds, err = storage.GetFeedsFromList("test 2")
	assert.NoError(t, err)
	expectedFeeds = []*_storage.ListItem{
		{AddedAt: time.Unix(defaultCurrentTime.Unix(), 0), Address: "https://rss-feed.com/rss"},
	}
	assert.Equal(t, expectedFeeds, feeds)
}

func Test_Explore_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	tempRepository, err := storage.JoinExploreDir("repo")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(tempRepository, 0755)

	err = os.WriteFile(path.Join(tempRepository, "feeds.opml"), []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Export from cleed/test</title>
    <dateCreated>Mon, 01 Jan 2024 00:00:00 +0000</dateCreated>
  </head>
  <body>
    <outline text="test">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
      <outline xmlUrl="https://test.com" />
    </outline>
    <outline text="test 2">
      <outline text="RSS Feed" description="RSS Feed description" xmlUrl="https://rss-feed.com/rss" />
      <outline text="Atom Feed" description="Atom Feed description" xmlUrl="https://atom-feed.com/atom" />
    </outline>
  </body>
</opml>`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)
	feed.SetDefaultExploreRepository("repo")

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "explore", "repo", "--remove"}
	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "repo was removed\n", out.String())

	_, err = os.Stat(tempRepository)
	assert.True(t, os.IsNotExist(err))
}
