package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsCreatePolicyEbsSnapApi
/* 您可以通过此接口在指定地域下创建一条自动快照策略。自动快照策略可以周期性地为云硬盘创建自动快照。
 */type EbsCreatePolicyEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsCreatePolicyEbsSnapApi(client *core.CtyunClient) *EbsCreatePolicyEbsSnapApi {
	return &EbsCreatePolicyEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/create-policy-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsCreatePolicyEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsCreatePolicyEbsSnapRequest) (*EbsCreatePolicyEbsSnapResponse, error) {
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
	var resp EbsCreatePolicyEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsCreatePolicyEbsSnapRequest struct {
	RegionID           string `json:"regionID,omitempty"`           /*  资源池ID。  */
	SnapshotPolicyName string `json:"snapshotPolicyName,omitempty"` /*  快照策略名称，只能由英文字母、数字、下划线、中划线组成，只能以英文字母开头，长度2-63字符。  */
	RepeatWeekdays     string `json:"repeatWeekdays,omitempty"`     /*  创建快照的重复日期，0-6分别代表周日-周六，多个日期用英文逗号隔开。  */
	RepeatTimes        string `json:"repeatTimes,omitempty"`        /*  创建快照的重复时间，0-23分别代表零点-23点，多个时间用英文逗号隔开。  */
	RetentionTime      int32  `json:"retentionTime"`                /*  创建快照的保留时间，输入范围为[-1，1-65535]，-1代表永久保留。单位为天。  */
	IsEnabled          *bool  `json:"isEnabled"`                    /*  是否启用策略，取值范围：
	●true：启用。
	●false：不启用。
	默认为true。  */
	ProjectID string `json:"projectID,omitempty"` /*  企业项目，默认为“0”。  */
}

type EbsCreatePolicyEbsSnapResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                  `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                  `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsCreatePolicyEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                  `json:"details,omitempty"`     /*  详情描述。  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                  `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsCreatePolicyEbsSnapReturnObjResponse struct {
	SnapshotPolicyID        *string `json:"snapshotPolicyID,omitempty"`        /*  策略ID。  */
	SnapshotPolicyJobResult *string `json:"snapshotPolicyJobResult,omitempty"` /*  任务执行状态。  */
}
