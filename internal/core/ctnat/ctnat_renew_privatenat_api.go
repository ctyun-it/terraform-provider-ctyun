package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatRenewPrivatenatApi
/* 续订 私网NAT
 */type CtnatRenewPrivatenatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatRenewPrivatenatApi(client *core.CtyunClient) *CtnatRenewPrivatenatApi {
	return &CtnatRenewPrivatenatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/renew-privatenat",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatRenewPrivatenatApi) Do(ctx context.Context, credential core.Credential, req *CtnatRenewPrivatenatRequest) (*CtnatRenewPrivatenatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatRenewPrivatenatRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtnatRenewPrivatenatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatRenewPrivatenatRequest struct {
	RegionID        string `json:"regionID,omitempty"`        /*  私网NAT所在区域的区域id。  */
	NatGatewayID    string `json:"natGatewayID,omitempty"`    /*  私网NAT的ID  */
	ClientToken     string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType       string `json:"cycleType,omitempty"`       /*  订购类型：month / year  */
	PayVoucherPrice string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
}

type CtnatRenewPrivatenatResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                 `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                 `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                 `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatRenewPrivatenatReturnObjResponse `json:"returnObj"`   /*  接口业务数据  */
}

type CtnatRenewPrivatenatReturnObjResponse struct{}
