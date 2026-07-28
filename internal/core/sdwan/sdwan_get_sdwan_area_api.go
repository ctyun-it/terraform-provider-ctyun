package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanAreaApi
/* 查询地区 */
type SdwanGetSdwanAreaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanAreaApi(client *core.CtyunClient) *SdwanGetSdwanAreaApi {
	return &SdwanGetSdwanAreaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/area/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanAreaApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanAreaRequest) (*SdwanGetSdwanAreaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.Province != nil && *req.Province != "" {
		ctReq.AddParam("province", *req.Province)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanAreaResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanAreaRequest struct {
	Province *string `json:"province,omitempty"` /*  省份，若不填则返回所有的省份，填了返回该身份下的地区  */
}

type SdwanGetSdwanAreaResponse struct {
	StatusCode   int32     `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode    *string   `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string   `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string   `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	TotalCount   int32     `json:"totalCount"`   /*  总数  */
	CurrentCount int32     `json:"currentCount"` /*  当前页数量  */
	Area         []*string `json:"area"`         /*  地区，值类型为string  */
	Error        *string   `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}
