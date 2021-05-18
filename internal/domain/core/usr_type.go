package core

import (
	"github.com/rendau/gms_temp/internal/cns"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

type UsrType struct {
	r *St
}

func NewUsrType(r *St) *UsrType {
	return &UsrType{r: r}
}

func (c *UsrType) List() []*entities.UsrTypeSt {
	return []*entities.UsrTypeSt{
		{
			Id:   cns.UsrTypeUndefined,
			Name: "не определен",
		},
		{
			Id:   cns.UsrTypeAdmin,
			Name: "Админ",
		},
	}
}

func (c *UsrType) GetName(id int) string {
	for _, rec := range c.List() {
		if rec.Id == id {
			return rec.Name
		}
	}

	return ""
}
