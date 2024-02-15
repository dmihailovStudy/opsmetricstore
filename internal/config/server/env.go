package server

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Envs struct {
	Address string `env:"ADDRESS"`
}

func (e *Envs) Load() error {
	if err := env.Parse(e); err != nil {
		return errors.WithMessage(err, "Load(): parse env config")
	}
	return nil
}
