package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsSetDeletePolicyEbsApi
/* 云硬盘的快照默认不随云硬盘自动释放，如果您不再需要某些云硬盘快照，可以将其手动释放。您还可以通过设置云硬盘释放策略，选择在释放云硬盘的同时释放该盘的全部快照。
 */type EbsSetDeletePolicyEbsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsSetDeletePolicyEbsApi(client *core.CtyunClient) *EbsSetDeletePolicyEbsApi {
	return &EbsSetDeletePolicyEbsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/set-delete-policy-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsSetDeletePolicyEbsApi) Do(ctx context.Context, credential core.Credential, req *EbsSetDeletePolicyEbsRequest) (*EbsSetDeletePolicyEbsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsSetDeletePolicyEbsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsSetDeletePolicyEbsRequest struct {
	RegionID          string `json:"regionID,omitempty"` /*  资源池ID。  */
	DiskID            string `json:"diskID,omitempty"`   /*  云硬盘ID。  */
	DeleteSnapWithEbs bool   `json:"deleteSnapWithEbs"`  /*  设置快照是否随盘删除。  */
}

type EbsSetDeletePolicyEbsResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                 `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                 `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsSetDeletePolicyEbsReturnObjResponse `json:"returnObj"`             /*  返回数据信息。  */
	Details     *string                                 `json:"details,omitempty"`     /*  详情描述。  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                 `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsSetDeletePolicyEbsReturnObjResponse struct {
	SetDeletePolicyJobResult *string `json:"setDeletePolicyJobResult,omitempty"` /*  任务执行状态。  */
}
