package core

import (
	"github.com/rendau/gms_temp/internal/interfaces"
)

type St struct {
	lg    interfaces.Logger
	cache interfaces.Cache
}

func New(lg interfaces.Logger, cache interfaces.Cache) *St {
	return &St{
		lg:    lg,
		cache: cache,
	}
}
