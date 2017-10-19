package sitter

// Config contains the configuration to start a babysitter.
type Config struct {
	StateFile string

	Home string

	Service   string
	External1 string
	External2 string
	Internal1 string
	Internal2 string
	Control   string

	Line1 string
	Line2 string
}

// Default return the default sitter setting.
func Default(local bool) *Config {
	c := &Config{
		StateFile: "sitter.json",
		Home:      ".",

		Service:   ":8100",
		External1: "localhost:8101",
		External2: "localhost:8102",
		Internal1: "localhost:8103",
		Internal2: "localhost:8104",
		Control:   ":8105",

		Line1: "Look to the left.",
		Line2: "Look to the right.",
	}
	if local {
		c.Service = "localhost:8100"
		c.Control = "localhost:8105"
	}
	return c
}
