package htmlgen

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Page contains the configuration of a page.
type Page struct {
	NoDocType bool
	Title     string
}

// Render renders a page.
func Render(w io.Writer, page *Page, body *Node) error {
	if !page.NoDocType {
		if _, err := io.WriteString(w, "<!doctype html>\n"); err != nil {
			return err
		}
	}

	doc := Element(atom.Html, Attrs{"lang": "en"})

	head := Head(NewMeta("charset", "UTF-8"))
	if page.Title != "" {
		head.Add(Title(page.Title))
	}
	doc.Add(head, body)

	if err := html.Render(w, doc.Node); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}
