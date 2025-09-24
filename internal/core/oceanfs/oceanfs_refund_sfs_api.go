package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsRefundSfsApi
/* 退订文件系统
 */type OceanfsRefundSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsRefundSfsApi(client *core.CtyunClient) *OceanfsRefundSfsApi {
	return &OceanfsRefundSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oceanfs/refund-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsRefundSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsRefundSfsRequest) (*OceanfsRefundSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*OceanfsRefundSfsRequest
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
	var resp OceanfsRefundSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsRefundSfsRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	SfsUID      string `json:"sfsUID,omitempty"`      /*  海量文件功能系统唯一 ID  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
}

type OceanfsRefundSfsResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800为成功，900为失败/订单处理中)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
