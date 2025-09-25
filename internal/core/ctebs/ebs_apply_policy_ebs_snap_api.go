package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsApplyPolicyEbsSnapApi
/* 您可以调用此接口为一块或者多块云硬盘关联自动快照策略。关联策略后该云盘会按照策略自动创建自动快照。若云硬盘已经绑定策略，将替换为新的策略。
 */type EbsApplyPolicyEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsApplyPolicyEbsSnapApi(client *core.CtyunClient) *EbsApplyPolicyEbsSnapApi {
	return &EbsApplyPolicyEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/apply-policy-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsApplyPolicyEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsApplyPolicyEbsSnapRequest) (*EbsApplyPolicyEbsSnapResponse, error) {
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
	var resp EbsApplyPolicyEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsApplyPolicyEbsSnapRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID。  */
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty"` /*  要关联的快照策略ID。  */
	TargetDiskIDs    string `json:"targetDiskIDs,omitempty"`    /*  要关联的云硬盘ID，多个云硬盘用英文逗号隔开。  */
}

type EbsApplyPolicyEbsSnapResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                 `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                 `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsApplyPolicyEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                 `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                 `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsApplyPolicyEbsSnapReturnObjResponse struct {
	SnapshotPolicyJobResult *string `json:"snapshotPolicyJobResult,omitempty"` /*  任务执行状态。  */
}
