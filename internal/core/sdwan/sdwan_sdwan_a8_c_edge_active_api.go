package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanA8CEdgeActiveApi
/* A8C 智能网关激活 */
type SdwanSdwanA8CEdgeActiveApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanA8CEdgeActiveApi(client *core.CtyunClient) *SdwanSdwanA8CEdgeActiveApi {
	return &SdwanSdwanA8CEdgeActiveApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/a8c-edge-active/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanA8CEdgeActiveApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanA8CEdgeActiveRequest) (*SdwanSdwanA8CEdgeActiveResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanA8CEdgeActiveRequest
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
	var resp SdwanSdwanA8CEdgeActiveResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanA8CEdgeActiveRequest struct {
	EdgeID     string   `json:"edgeID"`     /*  edge的ID  */
	DcID       string   `json:"dcID"`       /*  资源池ID  */
	SubnetList []string `json:"subnetList"` /*  子网网段，值类型为string  */
}

type SdwanSdwanA8CEdgeActiveResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
