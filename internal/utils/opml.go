package utils

import (
	"encoding/xml"
	"os"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OPMLHead `xml:"head"`
	Body    OPMLBody `xml:"body"`
}

type OPMLHead struct {
	Title       string `xml:"title,omitempty"`
	DateCreated string `xml:"dateCreated,omitempty"`
}

type OPMLBody struct {
	Text      string         `xml:"text,attr,omitempty"`
	Outltines []*OPMLOutline `xml:"outline"`
}

type OPMLOutline struct {
	Text        string         `xml:"text,attr,omitempty"`
	Description string         `xml:"description,attr,omitempty"`
	XMLURL      string         `xml:"xmlUrl,attr,omitempty"`
	Outlines    []*OPMLOutline `xml:"outline"`
}

func ParseOPMLFile(path string) (*OPML, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseOPMLBytes(b)
}

func ParseOPMLBytes(data []byte) (*OPML, error) {
	opml := &OPML{}
	err := xml.Unmarshal(data, opml)
	if err != nil {
		return nil, err
	}
	return opml, nil
}
