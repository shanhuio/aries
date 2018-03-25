package static

import (
	"shanhu.io/aries"
)

type config struct {
	Dir string // home directory
}

func main(env *aries.Env) (aries.Service, error) {
	static := aries.NewStaticFiles(env.Config.(*config).Dir)
	return static, nil
}

// Main is the main entrance for smlstatic binary
func Main() {
	aries.Main(main, new(config), "localhost:8000")
}
