package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatQueryPrivatenatIPApi
/* 查询中转IP
 */type CtnatQueryPrivatenatIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatQueryPrivatenatIPApi(client *core.CtyunClient) *CtnatQueryPrivatenatIPApi {
	return &CtnatQueryPrivatenatIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/privatenat/list-ips",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatQueryPrivatenatIPApi) Do(ctx context.Context, credential core.Credential, req *CtnatQueryPrivatenatIPRequest) (*CtnatQueryPrivatenatIPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("natGatewayID", req.NatGatewayID)
	if req.Address != "" {
		ctReq.AddParam("address", req.Address)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtnatQueryPrivatenatIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatQueryPrivatenatIPRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要查询的私网NAT的ID。  */
	Address      string `json:"address,omitempty"`      /*  要搜索的地址。  */
}

type CtnatQueryPrivatenatIPResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                     `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                     `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtnatQueryPrivatenatIPReturnObjResponse `json:"returnObj"`   /*  返回结果  */
}

type CtnatQueryPrivatenatIPReturnObjResponse struct {
	Address   string `json:"address"`   /*  中转IP的地址  */
	Status    string `json:"status"`    /*  中转IP状态: running代表运行中, freeze代表已冻结, expired代表已到期  */
	IsDefault *bool  `json:"isDefault"` /*  是否为默认中转地址  */
	SnatCnt   int32  `json:"snatCnt"`   /*  在使用此中转IP的snat数量  */
	DnarCnt   int32  `json:"dnarCnt"`   /*  在使用此中转IP的dnat数量  */
}
