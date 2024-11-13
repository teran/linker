package v1

import (
	"time"

	"github.com/teran/linker/models"
)

func NewRequest(r models.Request) (*Request, error) {
	parameters := map[string]*Parameter{}
	for k, v := range r.Parameters {
		parameters[k] = &Parameter{Value: v}
	}

	return &Request{
		Timestamp:  uint32(r.Timestamp.Unix()),
		LinkId:     r.LinkID,
		ClientIp:   r.ClientIP,
		CookieId:   r.CookieID,
		UserAgent:  r.UserAgent,
		Parameters: parameters,
	}, nil
}

func (r *Request) ToSvc() (models.Request, error) {
	parameters := map[string][]string{}
	for k, v := range r.GetParameters() {
		parameters[k] = v.GetValue()
	}

	return models.Request{
		Timestamp:  time.Unix(int64(r.Timestamp), 0).UTC(),
		LinkID:     r.LinkId,
		ClientIP:   r.ClientIp,
		CookieID:   r.CookieId,
		UserAgent:  r.UserAgent,
		Parameters: parameters,
	}, nil
}
