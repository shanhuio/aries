package static

import (
	"shanhu.io/aries"
)

type config struct {
	Dir string // home directory
}

func build(c interface{}, _ *aries.Logger) (aries.Service, error) {
	static := aries.NewStaticFiles(c.(*config).Dir)
	return static, nil
}

// Main is the main entrance for smlstatic binary
func Main() {
	aries.Main(build, new(config), "localhost:8000")
}
