package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsModifyPolicyStatusEbsSnapApi
/* 您可设置自动快照策略状态为启用/停用，启用状态的自动快照策略才可以执行创建自动快照任务。
 */type EbsModifyPolicyStatusEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsModifyPolicyStatusEbsSnapApi(client *core.CtyunClient) *EbsModifyPolicyStatusEbsSnapApi {
	return &EbsModifyPolicyStatusEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/modify-policy-status-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsModifyPolicyStatusEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsModifyPolicyStatusEbsSnapRequest) (*EbsModifyPolicyStatusEbsSnapResponse, error) {
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
	var resp EbsModifyPolicyStatusEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsModifyPolicyStatusEbsSnapRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID。  */
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty"` /*  快照策略ID，多个策略用英文逗号隔开。  */
	TargetStatus     string `json:"targetStatus,omitempty"`     /*  目标状态，取值范围：
	●activated：启用。
	●nonactivated：停用。  */
}

type EbsModifyPolicyStatusEbsSnapResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                        `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                        `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsModifyPolicyStatusEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                        `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                        `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsModifyPolicyStatusEbsSnapReturnObjResponse struct {
	SnapshotPolicyJobResult *string   `json:"snapshotPolicyJobResult,omitempty"` /*  任务执行状态。  */
	FailedList              []*string `json:"failedList"`                        /*  执行“启用/关闭快照策略”操作失败的快照策略列表。  */
}
