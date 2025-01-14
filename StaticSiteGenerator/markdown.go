package main

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func parseMarkdown(markdown string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
