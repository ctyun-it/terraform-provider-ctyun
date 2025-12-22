package cda

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "cda"

type Apis struct {
	CdaCdaGatewayListApi                 *CdaCdaGatewayListApi
	CdaCdaPhysicalLineUpdateApi          *CdaCdaPhysicalLineUpdateApi
	CdaCdaPhysicalLineListApi            *CdaCdaPhysicalLineListApi
	CdaCdaSharedPhysicalLineListApi      *CdaCdaSharedPhysicalLineListApi
	CdaCdaPhysicalLineBindApi            *CdaCdaPhysicalLineBindApi
	CdaCdaPhysicalLineUnbindApi          *CdaCdaPhysicalLineUnbindApi
	CdaCdaGatewayPhysicalLineListApi     *CdaCdaGatewayPhysicalLineListApi
	CdaCdaVpcAddApi                      *CdaCdaVpcAddApi
	CdaCdaVpcDeleteApi                   *CdaCdaVpcDeleteApi
	CdaCdaVpcUpdateApi                   *CdaCdaVpcUpdateApi
	CdaCdaStaticRouteAddApi              *CdaCdaStaticRouteAddApi
	CdaCdaStaticRouteDeleteApi           *CdaCdaStaticRouteDeleteApi
	CdaCdaStaticRouteUpdateApi           *CdaCdaStaticRouteUpdateApi
	CdaCdaStaticRouteListApi             *CdaCdaStaticRouteListApi
	CdaCdaBgpRouteAddApi                 *CdaCdaBgpRouteAddApi
	CdaCdaBgpRouteDeleteApi              *CdaCdaBgpRouteDeleteApi
	CdaCdaBgpRouteUpdateApi              *CdaCdaBgpRouteUpdateApi
	CdaCdaBgpRouteListApi                *CdaCdaBgpRouteListApi
	CdaCdaPhysicalLineAccessPointListApi *CdaCdaPhysicalLineAccessPointListApi
	CdaCdaPhysicalLineAddApi             *CdaCdaPhysicalLineAddApi
	CdaCdaPhysicalLineDeleteApi          *CdaCdaPhysicalLineDeleteApi
	CdaCdaPhysicalLineCountApi           *CdaCdaPhysicalLineCountApi
	CdaCdaVpcListApi                     *CdaCdaVpcListApi
	CdaCdaVpcCountApi                    *CdaCdaVpcCountApi
	CdaCdaGatewayCountApi                *CdaCdaGatewayCountApi
	CdaCdaGatewayDeleteApi               *CdaCdaGatewayDeleteApi
	CdaCdaSwitchListApi                  *CdaCdaSwitchListApi
	CdaCdaGatewayAddApi                  *CdaCdaGatewayAddApi
	CdaCdaHealthCheckQueryApi            *CdaCdaHealthCheckQueryApi
	CdaCdaHealthCheckStatusQueryApi      *CdaCdaHealthCheckStatusQueryApi
	CdaCdaHealthCheckAddApi              *CdaCdaHealthCheckAddApi
	CdaCdaHealthCheckUpdateApi           *CdaCdaHealthCheckUpdateApi
	CdaCdaHealthCheckDeleteApi           *CdaCdaHealthCheckDeleteApi
	CdaCdaLinkProbeAddApi                *CdaCdaLinkProbeAddApi
	CdaCdaLinkProbeDeleteApi             *CdaCdaLinkProbeDeleteApi
	CdaCdaLinkProbeQueryApi              *CdaCdaLinkProbeQueryApi
	CdaCdaVPCQueryApi                    *CdaCdaVPCQueryApi
	CdaCdaECQueryApi                     *CdaCdaECQueryApi
	CdaCdaListAccountAuthApi             *CdaCdaListAccountAuthApi
	CdaCdaCreateAccountAuthApi           *CdaCdaCreateAccountAuthApi
	CdaCdaDeleteAccountAuthApi           *CdaCdaDeleteAccountAuthApi
	CdaCdaStaticsAccountAuthApi          *CdaCdaStaticsAccountAuthApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		CdaCdaGatewayListApi:                 NewCdaCdaGatewayListApi(client),
		CdaCdaPhysicalLineUpdateApi:          NewCdaCdaPhysicalLineUpdateApi(client),
		CdaCdaPhysicalLineListApi:            NewCdaCdaPhysicalLineListApi(client),
		CdaCdaSharedPhysicalLineListApi:      NewCdaCdaSharedPhysicalLineListApi(client),
		CdaCdaPhysicalLineBindApi:            NewCdaCdaPhysicalLineBindApi(client),
		CdaCdaPhysicalLineUnbindApi:          NewCdaCdaPhysicalLineUnbindApi(client),
		CdaCdaGatewayPhysicalLineListApi:     NewCdaCdaGatewayPhysicalLineListApi(client),
		CdaCdaVpcAddApi:                      NewCdaCdaVpcAddApi(client),
		CdaCdaVpcDeleteApi:                   NewCdaCdaVpcDeleteApi(client),
		CdaCdaVpcUpdateApi:                   NewCdaCdaVpcUpdateApi(client),
		CdaCdaStaticRouteAddApi:              NewCdaCdaStaticRouteAddApi(client),
		CdaCdaStaticRouteDeleteApi:           NewCdaCdaStaticRouteDeleteApi(client),
		CdaCdaStaticRouteUpdateApi:           NewCdaCdaStaticRouteUpdateApi(client),
		CdaCdaStaticRouteListApi:             NewCdaCdaStaticRouteListApi(client),
		CdaCdaBgpRouteAddApi:                 NewCdaCdaBgpRouteAddApi(client),
		CdaCdaBgpRouteDeleteApi:              NewCdaCdaBgpRouteDeleteApi(client),
		CdaCdaBgpRouteUpdateApi:              NewCdaCdaBgpRouteUpdateApi(client),
		CdaCdaBgpRouteListApi:                NewCdaCdaBgpRouteListApi(client),
		CdaCdaPhysicalLineAccessPointListApi: NewCdaCdaPhysicalLineAccessPointListApi(client),
		CdaCdaPhysicalLineAddApi:             NewCdaCdaPhysicalLineAddApi(client),
		CdaCdaPhysicalLineDeleteApi:          NewCdaCdaPhysicalLineDeleteApi(client),
		CdaCdaPhysicalLineCountApi:           NewCdaCdaPhysicalLineCountApi(client),
		CdaCdaVpcListApi:                     NewCdaCdaVpcListApi(client),
		CdaCdaVpcCountApi:                    NewCdaCdaVpcCountApi(client),
		CdaCdaGatewayCountApi:                NewCdaCdaGatewayCountApi(client),
		CdaCdaGatewayDeleteApi:               NewCdaCdaGatewayDeleteApi(client),
		CdaCdaSwitchListApi:                  NewCdaCdaSwitchListApi(client),
		CdaCdaGatewayAddApi:                  NewCdaCdaGatewayAddApi(client),
		CdaCdaHealthCheckQueryApi:            NewCdaCdaHealthCheckQueryApi(client),
		CdaCdaHealthCheckStatusQueryApi:      NewCdaCdaHealthCheckStatusQueryApi(client),
		CdaCdaHealthCheckAddApi:              NewCdaCdaHealthCheckAddApi(client),
		CdaCdaHealthCheckUpdateApi:           NewCdaCdaHealthCheckUpdateApi(client),
		CdaCdaHealthCheckDeleteApi:           NewCdaCdaHealthCheckDeleteApi(client),
		CdaCdaLinkProbeAddApi:                NewCdaCdaLinkProbeAddApi(client),
		CdaCdaLinkProbeDeleteApi:             NewCdaCdaLinkProbeDeleteApi(client),
		CdaCdaLinkProbeQueryApi:              NewCdaCdaLinkProbeQueryApi(client),
		CdaCdaVPCQueryApi:                    NewCdaCdaVPCQueryApi(client),
		CdaCdaECQueryApi:                     NewCdaCdaECQueryApi(client),
		CdaCdaListAccountAuthApi:             NewCdaCdaListAccountAuthApi(client),
		CdaCdaCreateAccountAuthApi:           NewCdaCdaCreateAccountAuthApi(client),
		CdaCdaDeleteAccountAuthApi:           NewCdaCdaDeleteAccountAuthApi(client),
		CdaCdaStaticsAccountAuthApi:          NewCdaCdaStaticsAccountAuthApi(client),
	}
}
