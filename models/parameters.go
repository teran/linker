package models

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

type Parameters map[string][]string

func (p Parameters) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Parameters) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}

func (p Parameters) String() string {
	return url.Values(p).Encode()
}
