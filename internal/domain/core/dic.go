package core

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/rendau/gms_temp/internal/domain/entities"
)

type Dic struct {
	r *St

	dataMu sync.RWMutex
	dataHs string
	data   []byte
}

func NewDic(r *St) *Dic {
	return &Dic{r: r}
}

func (c *Dic) GetJson(ctx context.Context, reqHs string) (string, []byte, error) {
	var err error

	c.dataMu.RLock()
	if c.dataHs != "" {
		defer c.dataMu.RUnlock()

		if reqHs == c.dataHs {
			return c.dataHs, nil, nil
		}

		return c.dataHs, c.data, nil
	}
	c.dataMu.RUnlock()

	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if c.dataHs != "" {
		if reqHs == c.dataHs {
			return c.dataHs, nil, nil
		}

		return c.dataHs, c.data, nil
	}

	data := &entities.DicDataSt{}

	data.UsrTypes = c.r.UsrType.List()

	data.Config, err = c.r.Config.Get(ctx)
	if err != nil {
		return "", nil, err
	}

	c.data, err = json.Marshal(data)
	if err != nil {
		c.data = nil
		return "", nil, err
	}

	c.dataHs = fmt.Sprintf("%x", sha256.Sum256(c.data))

	// c.r.lg.Info(c.dataHs)

	return c.dataHs, c.data, nil
}

func (c *Dic) Refresh() {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	c.dataHs = ""
	c.data = nil
}
