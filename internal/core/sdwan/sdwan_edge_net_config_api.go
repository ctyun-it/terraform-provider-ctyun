package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanEdgeNetConfigApi
/* edge私网配置 */
type SdwanEdgeNetConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanEdgeNetConfigApi(client *core.CtyunClient) *SdwanEdgeNetConfigApi {
	return &SdwanEdgeNetConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-net-config/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanEdgeNetConfigApi) Do(ctx context.Context, credential core.Credential, req *SdwanEdgeNetConfigRequest) (*SdwanEdgeNetConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanEdgeNetConfigRequest
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
	var resp SdwanEdgeNetConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanEdgeNetConfigRequest struct {
	EdgeID        string   `json:"edgeID"`                  /*  edge的ID  */
	EdgeCIDR      []string `json:"edgeCIDR"`                /*  edge的子网  ，值类型为string,单个字符串为两部分，前面是Ipv4,后面是Ipv6  */
	LanControlIP  *string  `json:"lanControlIP,omitempty"`  /*  lan侧管理IP  */
	LanBusinessIP *string  `json:"lanBusinessIP,omitempty"` /*  lan侧业务IP  */
}

type SdwanEdgeNetConfigResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
