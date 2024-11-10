package models

type Link struct {
	DestinationURL          string
	Parameters              Parameters
	AllowParametersOverride bool
}
