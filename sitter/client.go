package sitter

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"shanhu.io/misc/tempfile"
	"shanhu.io/misc/ziputil"
)

type client struct {
	scheme string
	server string
	space  string
}

func newSpaceClient(server string) (*client, error) {
	c := &client{
		scheme: "http",
		server: server,
	}

	name, err := c.call("new", nil, nil)
	if err != nil {
		return nil, err
	}

	c.space = strings.TrimSpace(name)
	return c, nil
}

func (c *client) httpCall(m string, r io.Reader) (*http.Response, error) {
	if r == nil {
		return http.Get(m)
	}
	return http.Post(m, "", r)
}

func (c *client) call(m string, q url.Values, r io.Reader) (string, error) {
	if q == nil {
		q = url.Values{}
	}
	if c.space != "" {
		q.Add("space", c.space)
	}
	path := &url.URL{
		Scheme: c.scheme,
		Host:   c.server,
		Path:   m,
	}
	if len(q) > 0 {
		path.RawQuery = q.Encode()
	}

	resp, err := c.httpCall(path.String(), r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	return string(body), nil
}

func (c *client) put(r io.Reader) error {
	_, err := c.call("put", nil, r)
	return err
}

func (c *client) putZipFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.put(f); err != nil {
		return err
	}
	return err
}

func (c *client) putDir(dir string) error {
	tmp, err := tempfile.NewFile("", "sitter")
	if err != nil {
		return err
	}
	defer tmp.CleanUp()

	if err := ziputil.ZipDir(dir, tmp); err != nil {
		return err
	}

	if err := tmp.Reset(); err != nil {
		return err
	}
	if err := c.put(tmp); err != nil {
		return err
	}
	return tmp.CleanUp()
}

func (c *client) putFile(file string) error {
	tmp, err := tempfile.NewFile("", "sitter")
	if err != nil {
		return err
	}
	defer tmp.CleanUp()

	if err := ziputil.ZipFile(file, tmp); err != nil {
		return err
	}

	if err := tmp.Reset(); err != nil {
		return err
	}
	if err := c.put(tmp); err != nil {
		return err
	}
	return tmp.CleanUp()
}

func (c *client) run(bin string) error {
	q := url.Values{}
	q.Add("bin", bin)
	_, err := c.call("run", q, nil)
	return err
}

// Push pushes the current directory into a sitter hosted docker
// container.
func Push(dir, server string, out io.Writer) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", dir)
	}

	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	infos, err := f.Readdir(0)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "server: ", server)
	c, err := newSpaceClient(server)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "space: ", c.space)

	bin := ""
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		name := info.Name()
		fullPath := filepath.Join(dir, name)
		if strings.HasSuffix(name, ".zip") {
			fmt.Fprintln(out, "zip: ", name)
			c.putZipFile(fullPath)
		} else {
			fmt.Fprintln(out, "file: ", name)
			c.putFile(fullPath)

			perm := info.Mode() & os.ModePerm
			if (perm & 0100) != 0 {
				// an executable
				if bin == "" || name < bin {
					bin = name
				}
			}
		}
	}

	if err := f.Close(); err != nil {
		return err
	}
	fmt.Fprintln(out, "run: ", bin)
	return c.run(bin)
}

// ClientMain is the entrance of a client.
func ClientMain() {
	var (
		server = flag.String("s", "", "server to connect")
		dir    = flag.String("dir", ".", "project directory")
	)
	flag.Parse()
	if err := Push(*dir, *server, os.Stdout); err != nil {
		log.Println(err)
	}
}
