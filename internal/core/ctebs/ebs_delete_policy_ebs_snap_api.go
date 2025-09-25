package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsDeletePolicyEbsSnapApi
/* 如果您不再需要自动快照策略，可直接删除策略，与该策略关联的云硬盘将同时自动解绑。
 */type EbsDeletePolicyEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsDeletePolicyEbsSnapApi(client *core.CtyunClient) *EbsDeletePolicyEbsSnapApi {
	return &EbsDeletePolicyEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/delete-policy-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsDeletePolicyEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsDeletePolicyEbsSnapRequest) (*EbsDeletePolicyEbsSnapResponse, error) {
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
	var resp EbsDeletePolicyEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsDeletePolicyEbsSnapRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID。  */
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty"` /*  快照策略ID，只支持单个删除。  */
}

type EbsDeletePolicyEbsSnapResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                  `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                  `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsDeletePolicyEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                  `json:"details,omitempty"`     /*  详情描述。  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                  `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsDeletePolicyEbsSnapReturnObjResponse struct {
	SnapshotPolicyJobResult *string `json:"snapshotPolicyJobResult,omitempty"` /*  任务执行状态。  */
}
