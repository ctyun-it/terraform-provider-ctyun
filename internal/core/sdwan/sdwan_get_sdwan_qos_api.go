package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanQosApi
/* 查询qos */
type SdwanGetSdwanQosApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanQosApi(client *core.CtyunClient) *SdwanGetSdwanQosApi {
	return &SdwanGetSdwanQosApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/qos/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanQosApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanQosRequest) (*SdwanGetSdwanQosResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ProjectID != nil && *req.ProjectID != "" {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	if req.QosID != nil && *req.QosID != "" {
		ctReq.AddParam("qosID", *req.QosID)
	}
	if req.QosName != nil && *req.QosName != "" {
		ctReq.AddParam("qosName", *req.QosName)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanQosResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanQosRequest struct {
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目名称  */
	QosID     *string `json:"qosID,omitempty"`     /*  qos策略ID  */
	QosName   *string `json:"qosName,omitempty"`   /*  qos策略名称  */
	PageNo    int32   `json:"pageNo"`              /*  页数  */
	PageSize  int32   `json:"pageSize"`            /*  页大小  */
	Search    *string `json:"search,omitempty"`    /*  模糊查询  */
}

type SdwanGetSdwanQosResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanGetSdwanQosReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanQosReturnObjResponse struct {
	Result       []*SdwanGetSdwanQosReturnObjResultResponse `json:"result"`       /*  查询qos  */
	TotalCount   int32                                      `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                      `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetSdwanQosReturnObjResultResponse struct {
	QosID         *string `json:"qosID"`         /*  qos策略ID  */
	QosName       *string `json:"qosName"`       /*  qos策略名称  */
	Description   *string `json:"description"`   /*  描述  */
	ProjectID     *string `json:"projectID"`     /*  企业项目名称  */
	Bandwidth     int32   `json:"bandwidth"`     /*  带宽峰值  */
	BandwidthType *string `json:"bandwidthType"` /*  本参数表示带宽类型<br/><br/>取值范围:<br/>internet:互联网带宽<br/>sdwan:SD-WAN带宽  */
	RuleCount     int32   `json:"ruleCount"`     /*  规则数  */
	EdgeCount     int32   `json:"edgeCount"`     /*  关联实例数  */
}
