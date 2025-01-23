package utils

import (
	"encoding/xml"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OPMLHead `xml:"head"`
	Body    OPMLBody `xml:"body"`
}

type OPMLHead struct {
	Title string `xml:"title,omitempty"`
}

type OPMLBody struct {
	Text     string         `xml:"text,attr,omitempty"`
	Oultines []*OPMLOutline `xml:"outline"`
}

type OPMLOutline struct {
	Text        string         `xml:"text,attr,omitempty"`
	Description string         `xml:"description,attr,omitempty"`
	XMLURL      string         `xml:"xmlUrl,attr,omitempty"`
	Outlines    []*OPMLOutline `xml:"outline"`
}
