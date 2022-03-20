package exec

import (
	"os/exec"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

func Command(name string, argv ...string) Cmd {
	cmd := exec.Command(name, argv...)

	return Cmd{
		Cmd: *cmd,
		Env: EnvFrom(cmd.Env),
	}
}

type Cmd struct {
	exec.Cmd
	Env Env
}

func (c Cmd) Run() error {
	log.Debug().
		Interface("env", c.Env).
		Str("dir", c.Dir).
		Msgf("%s %s", c.Cmd.Path, strings.Join(c.Cmd.Args, " "))

	c.Cmd.Env = c.Env.Strings()
	return c.Cmd.Run()
}

type Env map[string]string

func (e Env) Strings() []string {
	s := make([]string, 0, len(e))
	for k, v := range e {
		s = append(s, k+"="+v)
	}
	sort.Strings(s)
	return s
}

func EnvFrom(strs []string) Env {
	e := make(Env)
	for _, s := range e {
		k, v, ok := strings.Cut(s, "=")
		if !ok {
			panic(e)
		}
		e[k] = v
	}
	return e
}
