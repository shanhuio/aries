package tasks

import (
	"fmt"
	"log"
	"path"

	"shanhu.io/aries/creds"
)

// Run issues a list of tasks to a particular server.
func Run(server, prefix string, tasks []string) error {
	c, err := creds.Dial(server)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		// list available tasks
		var names []string
		p := path.Join(prefix, "help")
		if err := c.JSONCall(p, nil, &names); err != nil {
			return err
		}

		for _, name := range names {
			fmt.Println(name)
		}
		return nil
	}

	for _, t := range tasks {
		log.Println(t)
		if err := c.Poke(path.Join(prefix, t)); err != nil {
			return err
		}
	}

	return nil
}
