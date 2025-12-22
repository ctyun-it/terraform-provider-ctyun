package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatModifySpecApi
/* 变配私网NAT
 */type CtnatModifySpecApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatModifySpecApi(client *core.CtyunClient) *CtnatModifySpecApi {
	return &CtnatModifySpecApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/modify-spec",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatModifySpecApi) Do(ctx context.Context, credential core.Credential, req *CtnatModifySpecRequest) (*CtnatModifySpecResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatModifySpecRequest
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
	var resp CtnatModifySpecResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatModifySpecRequest struct {
	RegionID        string `json:"regionID,omitempty"`        /*  私网NAT网关所在的地域 ID  */
	NatGatewayID    string `json:"natGatewayID,omitempty"`    /*  私网NAT的ID  */
	Spec            string `json:"spec,omitempty"`            /*  规格(可传值：small, medium, large, xlarge)  */
	ClientToken     string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 6  */
	PayVoucherPrice string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位，仅包周期支持代金券  */
}

type CtnatModifySpecResponse struct {
	StatusCode  int32                             `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                            `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                            `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                            `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatModifySpecReturnObjResponse `json:"returnObj"`   /*  返回结果  */
}

type CtnatModifySpecReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID"` /*  订单 id  */
	MasterOrderNO string `json:"masterOrderNO"` /*  订单编号，可以为null  */
	RegionID      string `json:"regionID"`      /*  资源池ID  */
}
