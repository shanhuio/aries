package gometa

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"shanhu.io/misc/pathutil"
)

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "ascii":
		return r, nil
	}
	return nil, fmt.Errorf("charset %q not supported", charset)
}

func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}

// ParseGoImport takes an HTML page and parses for the go-import meta tag.
func ParseGoImport(r io.Reader, pkg string) (*Repo, error) {
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charsetReader
	dec.Strict = false

	for {
		t, err := dec.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if e, ok := t.(xml.StartElement); ok {
			if strings.EqualFold(e.Name.Local, "body") {
				break
			}
		}

		if e, ok := t.(xml.EndElement); ok {
			if strings.EqualFold(e.Name.Local, "head") {
				break
			}
		}

		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}

		if attrValue(e.Attr, "name") != "go-import" {
			continue
		}

		fields := strings.Fields(attrValue(e.Attr, "content"))
		if len(fields) != 3 {
			continue
		}

		repoPath := fields[0]
		if !pathutil.IsParent(repoPath, pkg) {
			continue
		}

		// we found it
		vcs := fields[1]
		url := fields[2]
		return &Repo{
			ImportRoot: repoPath,
			VCS:        vcs,
			VCSRoot:    url,
		}, nil
	}

	return nil, fmt.Errorf("go meta not found for %q", pkg)
}