package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcVniaListInstanceDiagnosisRecordApi
/* 获取实例诊断记录列表
 */type CtvpcVniaListInstanceDiagnosisRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaListInstanceDiagnosisRecordApi(client *core.CtyunClient) *CtvpcVniaListInstanceDiagnosisRecordApi {
	return &CtvpcVniaListInstanceDiagnosisRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/list-instance-diagnosis-record",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaListInstanceDiagnosisRecordApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaListInstanceDiagnosisRecordRequest) (*CtvpcVniaListInstanceDiagnosisRecordResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ResourceID != nil {
		ctReq.AddParam("resourceID", *req.ResourceID)
	}
	if req.ResourceType != nil {
		ctReq.AddParam("resourceType", *req.ResourceType)
	}
	if req.DiagnosisRecordID != nil {
		ctReq.AddParam("diagnosisRecordID", *req.DiagnosisRecordID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcVniaListInstanceDiagnosisRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaListInstanceDiagnosisRecordRequest struct {
	RegionID          string  /*  资源池 ID  */
	ResourceID        *string /*  资源 ID  */
	ResourceType      *string /*  资源类型, 仅支持 eip / natgw / elb  */
	DiagnosisRecordID *string /*  实例诊断记录 ID  */
	PageNumber        int32   /*  列表的页码，默认值为 1。  */
	PageSize          int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcVniaListInstanceDiagnosisRecordResponse struct {
	StatusCode  int32                                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcVniaListInstanceDiagnosisRecordReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaListInstanceDiagnosisRecordReturnObjResponse struct {
	Results      []*CtvpcVniaListInstanceDiagnosisRecordReturnObjResultsResponse `json:"results"`      /*  网卡列表  */
	TotalCount   int32                                                           `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                           `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                           `json:"totalPage"`    /*  总页数  */
}

type CtvpcVniaListInstanceDiagnosisRecordReturnObjResultsResponse struct {
	CreatedAt         *string `json:"createdAt,omitempty"`         /*  创建时间  */
	UpdatedAt         *string `json:"updatedAt,omitempty"`         /*  更新时间  */
	DiagnosisRecordID *string `json:"diagnosisRecordID,omitempty"` /*  实例诊断记录 ID  */
	ResourceID        *string `json:"resourceID,omitempty"`        /*  资源 ID  */
	ResourceType      *string `json:"resourceType,omitempty"`      /*  资源类型, 仅支持 eip / natgw / elb  */
	ErrMsg            *string `json:"errMsg,omitempty"`            /*  错误信息  */
	DiagnosisStatus   *string `json:"diagnosisStatus,omitempty"`   /*  诊断状态  */
	HealthStatus      int32   `json:"healthStatus"`                /*  健康检查状态  */
}
