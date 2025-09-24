package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsModifyPolicyEbsSnapApi
/* 如果您需要修改自动快照的定时创建时间、保留规则等信息，可以通过接口编辑修改自动快照策略。
 */ /* 修改自动快照策略后，之前已应用该策略的云硬盘随即执行更新后的自动快照策略。
 */type EbsModifyPolicyEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsModifyPolicyEbsSnapApi(client *core.CtyunClient) *EbsModifyPolicyEbsSnapApi {
	return &EbsModifyPolicyEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/modify-policy-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsModifyPolicyEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsModifyPolicyEbsSnapRequest) (*EbsModifyPolicyEbsSnapResponse, error) {
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
	var resp EbsModifyPolicyEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsModifyPolicyEbsSnapRequest struct {
	RegionID           string  `json:"regionID,omitempty"`           /*  资源池ID。  */
	SnapshotPolicyID   string  `json:"snapshotPolicyID,omitempty"`   /*  要修改的快照策略ID。  */
	SnapshotPolicyName *string `json:"snapshotPolicyName,omitempty"` /*  修改后的快照策略名称，只能由英文字母、数字、下划线、中划线组成，不能以特殊字符、数字开头，长度2-63字符。  */
	RepeatWeekdays     *string `json:"repeatWeekdays,omitempty"`     /*  修改后的创建快照的重复日期，0-6分别代表周日-周六，多个日期用英文逗号隔开。  */
	RepeatTimes        *string `json:"repeatTimes,omitempty"`        /*  修改后的创建快照的重复时间，0-23分别代表零点-23点，多个时间用英文逗号隔开。  */
	RetentionTime      int32   `json:"retentionTime"`                /*  修改后的创建快照的保留时间，输入范围为[ -1，1-65535], -1代表永久保留。单位为天。  */
}

type EbsModifyPolicyEbsSnapResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                  `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                  `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsModifyPolicyEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                  `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                  `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsModifyPolicyEbsSnapReturnObjResponse struct {
	SnapshotPolicyJobResult *string `json:"snapshotPolicyJobResult,omitempty"` /*  任务执行状态。  */
}
