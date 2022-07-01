package core

import (
	"context"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

type Config struct {
	r *St
}

func NewConfig(r *St) *Config {
	return &Config{r: r}
}

func (c *Config) Get(ctx context.Context) (*entities.ConfigSt, error) {
	return c.r.repo.ConfigGet(ctx)
}

func (c *Config) Set(ctx context.Context, config *entities.ConfigSt) error {
	return c.r.repo.ConfigSet(ctx, config)
}
