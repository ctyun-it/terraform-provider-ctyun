package crs

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "crs"

type Apis struct {
	CrsGetValuesApi                      *CrsGetValuesApi
	CrsListPluginApi                     *CrsListPluginApi
	CrsListTagApi                        *CrsListTagApi
	CrsListTemplateApi                   *CrsListTemplateApi
	CrsCreateInstanceVpceLinkedVpcsV2Api *CrsCreateInstanceVpceLinkedVpcsV2Api
	CrsDeleteInstanceVpceLinkedVpcsV2Api *CrsDeleteInstanceVpceLinkedVpcsV2Api
	CrsGetInstanceVpceLinkedVpcsV2Api    *CrsGetInstanceVpceLinkedVpcsV2Api
	CrsListOpenSourceRepositoryV2Api     *CrsListOpenSourceRepositoryV2Api
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		CrsGetValuesApi:                      NewCrsGetValuesApi(client),
		CrsListPluginApi:                     NewCrsListPluginApi(client),
		CrsListTagApi:                        NewCrsListTagApi(client),
		CrsListTemplateApi:                   NewCrsListTemplateApi(client),
		CrsCreateInstanceVpceLinkedVpcsV2Api: NewCrsCreateInstanceVpceLinkedVpcsV2Api(client),
		CrsDeleteInstanceVpceLinkedVpcsV2Api: NewCrsDeleteInstanceVpceLinkedVpcsV2Api(client),
		CrsGetInstanceVpceLinkedVpcsV2Api:    NewCrsGetInstanceVpceLinkedVpcsV2Api(client),
		CrsListOpenSourceRepositoryV2Api:     NewCrsListOpenSourceRepositoryV2Api(client),
	}
}
