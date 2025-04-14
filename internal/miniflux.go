package internal

import (
	"bytes"
	"io"

	"github.com/radulucut/cleed/internal/utils"
	"miniflux.app/v2/client"
)

func (f *TerminalFeed) minifluxInitClient(token string) error {
	f.miniflux = client.NewClient("https://reader.miniflux.app/v1/", token)
	return nil
}

func (f *TerminalFeed) MinifluxPush() error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	if config.MinifluxToken == "" {
		return utils.NewInternalError("miniflux token is not set. Use `cleed config --miniflux-token=\"your_token_here\"` to set it.")
	}
	if f.miniflux == nil {
		err := f.minifluxInitClient(config.MinifluxToken)
		if err != nil {
			return err
		}
	}
	buf := &bytes.Buffer{}
	f.writeOPML(buf, "", false)
	return f.miniflux.Import(io.NopCloser(buf))
}

func (f *TerminalFeed) MinifluxPull() error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	if config.MinifluxToken == "" {
		return utils.NewInternalError("miniflux token is not set. Use `cleed config --miniflux-token=\"your_token_here\"` to set it.")
	}
	if f.miniflux == nil {
		err := f.minifluxInitClient(config.MinifluxToken)
		if err != nil {
			return err
		}
	}
	b, err := f.miniflux.Export()
	if err != nil {
		return err
	}
	opml, err := utils.ParseOPMLBytes(b)
	if err != nil {
		return utils.NewInternalError("failed to parse OPML: " + err.Error())
	}
	return f.importOPML(opml, "")
}
