package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanEdgeBranchApi
/* 创建edge互联关系 */
type SdwanCreateSdwanEdgeBranchApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanEdgeBranchApi(client *core.CtyunClient) *SdwanCreateSdwanEdgeBranchApi {
	return &SdwanCreateSdwanEdgeBranchApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-branch/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanEdgeBranchApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanEdgeBranchRequest) (*SdwanCreateSdwanEdgeBranchResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanEdgeBranchRequest
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
	var resp SdwanCreateSdwanEdgeBranchResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanEdgeBranchRequest struct {
	SrcEdgeID string `json:"srcEdgeID"` /*  源 edge id  */
	DstEdgeID string `json:"dstEdgeID"` /*  目的 edge id  */
}

type SdwanCreateSdwanEdgeBranchResponse struct {
	StatusCode   int32     `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode    *string   `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string   `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string   `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	OperationID  *string   `json:"operationID"`  /*  操作id  */
	EdgeBranchID []*string `json:"edgeBranchID"` /*  branch id  ，值类型为string  */
	Error        *string   `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}
