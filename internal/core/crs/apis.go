package crs

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "crs"

type Apis struct {
	CrsGetValuesApi    *CrsGetValuesApi
	CrsListPluginApi   *CrsListPluginApi
	CrsListTagApi      *CrsListTagApi
	CrsListTemplateApi *CrsListTemplateApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		CrsGetValuesApi:    NewCrsGetValuesApi(client),
		CrsListPluginApi:   NewCrsListPluginApi(client),
		CrsListTagApi:      NewCrsListTagApi(client),
		CrsListTemplateApi: NewCrsListTemplateApi(client),
	}
}
