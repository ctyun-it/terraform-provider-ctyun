package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsCancelPolicyEbsSnapApi
/* 当云硬盘不再需要创建自动快照时，您可以为云硬盘取消关联自动快照策略。云硬盘取消关联自动快照策略后，将不会再按照自动快照策略继续创建自动快照。
 */type EbsCancelPolicyEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsCancelPolicyEbsSnapApi(client *core.CtyunClient) *EbsCancelPolicyEbsSnapApi {
	return &EbsCancelPolicyEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/cancel-policy-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsCancelPolicyEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsCancelPolicyEbsSnapRequest) (*EbsCancelPolicyEbsSnapResponse, error) {
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
	var resp EbsCancelPolicyEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsCancelPolicyEbsSnapRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID。  */
	TargetDiskIDs string `json:"targetDiskIDs,omitempty"` /*  要取消关联的云硬盘ID，多个云硬盘用英文逗号隔开。  */
}

type EbsCancelPolicyEbsSnapResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                  `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                  `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsCancelPolicyEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                  `json:"details,omitempty"`     /*  详情描述。  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                  `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsCancelPolicyEbsSnapReturnObjResponse struct {
	SnapshotPolicyJobResult *string `json:"snapshotPolicyJobResult,omitempty"` /*  任务执行状态。  */
}
