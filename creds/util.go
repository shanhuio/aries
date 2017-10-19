package creds

// Poke pokes a server with the default creds.
func Poke(server string, path string) error {
	c, err := Dial(server)
	if err != nil {
		return err
	}

	return c.Poke("/api/rebuild")
}
