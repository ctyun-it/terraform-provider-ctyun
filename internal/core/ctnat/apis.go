package ctnat

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ctnat"

type Apis struct {
	CtnatCreatePrivatenatApi     *CtnatCreatePrivatenatApi
	CtnatModifySpecApi           *CtnatModifySpecApi
	CtnatQueryRenewPriceApi      *CtnatQueryRenewPriceApi
	CtnatRenewPrivatenatApi      *CtnatRenewPrivatenatApi
	CtnatListPrivatenatApi       *CtnatListPrivatenatApi
	CtnatUpdatePrivatenatApi     *CtnatUpdatePrivatenatApi
	CtnatListPrivatenatCidrsApi  *CtnatListPrivatenatCidrsApi
	CtnatCreatePrivatenatIPApi   *CtnatCreatePrivatenatIPApi
	CtnatQueryPrivatenatIPApi    *CtnatQueryPrivatenatIPApi
	CtnatDeletePrivatenatIPApi   *CtnatDeletePrivatenatIPApi
	CtnatCreatePrivatenatSnatApi *CtnatCreatePrivatenatSnatApi
	CtnatQueryPrivatenatSnatApi  *CtnatQueryPrivatenatSnatApi
	CtnatModifyPrivatenatSnatApi *CtnatModifyPrivatenatSnatApi
	CtnatDeletePrivatenatSnatApi *CtnatDeletePrivatenatSnatApi
	CtnatCreatePrivatenatDnatApi *CtnatCreatePrivatenatDnatApi
	CtnatQueryPrivatenatDnatApi  *CtnatQueryPrivatenatDnatApi
	CtnatModifyPrivatenatDnatApi *CtnatModifyPrivatenatDnatApi
	CtnatDeletePrivatenatDnatApi *CtnatDeletePrivatenatDnatApi
	CtnatDeletePrivatenatApi     *CtnatDeletePrivatenatApi
	CtnatQueryCreatePriceApi     *CtnatQueryCreatePriceApi
	CtnatQueryModifySpecPriceApi *CtnatQueryModifySpecPriceApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		CtnatCreatePrivatenatApi:     NewCtnatCreatePrivatenatApi(client),
		CtnatModifySpecApi:           NewCtnatModifySpecApi(client),
		CtnatQueryRenewPriceApi:      NewCtnatQueryRenewPriceApi(client),
		CtnatRenewPrivatenatApi:      NewCtnatRenewPrivatenatApi(client),
		CtnatListPrivatenatApi:       NewCtnatListPrivatenatApi(client),
		CtnatUpdatePrivatenatApi:     NewCtnatUpdatePrivatenatApi(client),
		CtnatListPrivatenatCidrsApi:  NewCtnatListPrivatenatCidrsApi(client),
		CtnatCreatePrivatenatIPApi:   NewCtnatCreatePrivatenatIPApi(client),
		CtnatQueryPrivatenatIPApi:    NewCtnatQueryPrivatenatIPApi(client),
		CtnatDeletePrivatenatIPApi:   NewCtnatDeletePrivatenatIPApi(client),
		CtnatCreatePrivatenatSnatApi: NewCtnatCreatePrivatenatSnatApi(client),
		CtnatQueryPrivatenatSnatApi:  NewCtnatQueryPrivatenatSnatApi(client),
		CtnatModifyPrivatenatSnatApi: NewCtnatModifyPrivatenatSnatApi(client),
		CtnatDeletePrivatenatSnatApi: NewCtnatDeletePrivatenatSnatApi(client),
		CtnatCreatePrivatenatDnatApi: NewCtnatCreatePrivatenatDnatApi(client),
		CtnatQueryPrivatenatDnatApi:  NewCtnatQueryPrivatenatDnatApi(client),
		CtnatModifyPrivatenatDnatApi: NewCtnatModifyPrivatenatDnatApi(client),
		CtnatDeletePrivatenatDnatApi: NewCtnatDeletePrivatenatDnatApi(client),
		CtnatDeletePrivatenatApi:     NewCtnatDeletePrivatenatApi(client),
		CtnatQueryCreatePriceApi:     NewCtnatQueryCreatePriceApi(client),
		CtnatQueryModifySpecPriceApi: NewCtnatQueryModifySpecPriceApi(client),
	}
}
