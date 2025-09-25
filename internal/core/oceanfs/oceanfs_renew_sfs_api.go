package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsRenewSfsApi
/* 文件系统续订
 */type OceanfsRenewSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsRenewSfsApi(client *core.CtyunClient) *OceanfsRenewSfsApi {
	return &OceanfsRenewSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oceanfs/renew-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsRenewSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsRenewSfsRequest) (*OceanfsRenewSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*OceanfsRenewSfsRequest
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
	var resp OceanfsRenewSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsRenewSfsRequest struct {
	SfsUID      string `json:"sfsUID,omitempty"`      /*  海量文件功能系统唯一 ID  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	CycleType   string `json:"cycleType,omitempty"`   /*  包周期类型，year/month  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  包周期数。周期最大长度不能超过 3 年  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
}

type OceanfsRenewSfsResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800为成功，900为失败/订单处理中)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
