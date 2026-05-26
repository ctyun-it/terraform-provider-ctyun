package opensearch

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ctcsx"

type Apis struct {
	OpensearchNewClusterApi          *OpensearchNewClusterApi
	OpensearchUnsubscribeInstanceApi *OpensearchUnsubscribeInstanceApi
	OpensearchListInstancesApi       *OpensearchListInstancesApi
	OpensearchGetInstanceApi         *OpensearchGetInstanceApi
	OpensearchScaleInstanceApi       *OpensearchScaleInstanceApi
	OpensearchRestartInstanceApi     *OpensearchRestartInstanceApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		OpensearchNewClusterApi:          NewOpensearchNewClusterApi(client),
		OpensearchUnsubscribeInstanceApi: NewOpensearchUnsubscribeInstanceApi(client),
		OpensearchListInstancesApi:       NewOpensearchListInstancesApi(client),
		OpensearchGetInstanceApi:         NewOpensearchGetInstanceApi(client),
		OpensearchScaleInstanceApi:       NewOpensearchScaleInstanceApi(client),
		OpensearchRestartInstanceApi:     NewOpensearchRestartInstanceApi(client),
	}
}
