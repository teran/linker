package v1

import "github.com/teran/linker/models"

func NewLink(l models.Link) (*Link, error) {
	parameters := map[string]*Parameter{}
	for k, v := range l.Parameters {
		parameters[k] = &Parameter{Value: v}
	}

	return &Link{
		DestinationUrl:          l.DestinationURL,
		Parameters:              parameters,
		AllowParametersOverride: l.AllowParametersOverride,
	}, nil
}

func (l *Link) ToSvc() (models.Link, error) {
	parameters := map[string][]string{}
	for k, v := range l.GetParameters() {
		parameters[k] = v.GetValue()
	}

	return models.Link{
		DestinationURL:          l.GetDestinationUrl(),
		Parameters:              parameters,
		AllowParametersOverride: l.GetAllowParametersOverride(),
	}, nil
}
