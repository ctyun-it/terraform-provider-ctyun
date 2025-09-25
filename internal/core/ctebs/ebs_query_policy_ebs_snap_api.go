package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsQueryPolicyEbsSnapApi
/* 您可以调用此接口查询某资源池下自动快照策略详情或列表。
 */type EbsQueryPolicyEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsQueryPolicyEbsSnapApi(client *core.CtyunClient) *EbsQueryPolicyEbsSnapApi {
	return &EbsQueryPolicyEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs_snapshot/query-policy-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsQueryPolicyEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsQueryPolicyEbsSnapRequest) (*EbsQueryPolicyEbsSnapResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.SnapshotPolicyID != nil {
		ctReq.AddParam("snapshotPolicyID", *req.SnapshotPolicyID)
	}
	if req.SnapshotPolicyName != nil {
		ctReq.AddParam("snapshotPolicyName", *req.SnapshotPolicyName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsQueryPolicyEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsQueryPolicyEbsSnapRequest struct {
	RegionID           string  /*  资源池ID。  */
	SnapshotPolicyID   *string /*  快照策略ID。  */
	SnapshotPolicyName *string /*  快照策略名称。  */
}

type EbsQueryPolicyEbsSnapResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                 `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                 `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsQueryPolicyEbsSnapReturnObjResponse `json:"returnObj"`             /*  参考响应示例。  */
	Details     *string                                 `json:"details,omitempty"`     /*  详情描述。  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                 `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsQueryPolicyEbsSnapReturnObjResponse struct {
	SnapshotPolicyTotalCount int32                                                       `json:"snapshotPolicyTotalCount"` /*  自动快照策略总数量。  */
	SnapshotPolicyList       []*EbsQueryPolicyEbsSnapReturnObjSnapshotPolicyListResponse `json:"snapshotPolicyList"`       /*  自动快照策略列表。  */
}

type EbsQueryPolicyEbsSnapReturnObjSnapshotPolicyListResponse struct {
	SnapshotPolicyID     *string `json:"snapshotPolicyID,omitempty"`     /*  自动快照策略ID。  */
	SnapshotPolicyName   *string `json:"snapshotPolicyName,omitempty"`   /*  自动快照策略名称。  */
	SnapshotPolicyStatus *string `json:"snapshotPolicyStatus,omitempty"` /*  自动快照策略状态，取值范围：
	●activated:启用。
	●deactivated：停用。  */
	RepeatWeekdays           *string `json:"repeatWeekdays,omitempty"`           /*  创建快照的重复日期，0-6分别代表周日-周六，多个日期用英文逗号隔开。  */
	RepeatTimes              *string `json:"repeatTimes,omitempty"`              /*  创建快照的重复时间，0-23分别代表零点-23点，多个时间用英文逗号隔开。  */
	RetentionTime            int32   `json:"retentionTime"`                      /*  创建快照的保留时间，-1代表永久保留。  */
	ProjectID                *string `json:"projectID,omitempty"`                /*  企业项目。  */
	BoundDiskNum             int32   `json:"boundDiskNum"`                       /*  关联云硬盘的数量。  */
	SnapshotPolicyCreateTime *string `json:"snapshotPolicyCreateTime,omitempty"` /*  策略创建时间。  */
}
