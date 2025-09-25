package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcVniaListInstanceDiagnosisReportApi
/* 获取实例诊断报告
 */type CtvpcVniaListInstanceDiagnosisReportApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaListInstanceDiagnosisReportApi(client *core.CtyunClient) *CtvpcVniaListInstanceDiagnosisReportApi {
	return &CtvpcVniaListInstanceDiagnosisReportApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/list-instance-diagnosis-report",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaListInstanceDiagnosisReportApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaListInstanceDiagnosisReportRequest) (*CtvpcVniaListInstanceDiagnosisReportResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("resourceID", req.ResourceID)
	ctReq.AddParam("resourceType", req.ResourceType)
	if req.DiagnosisReportID != nil {
		ctReq.AddParam("diagnosisReportID", *req.DiagnosisReportID)
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
	var resp CtvpcVniaListInstanceDiagnosisReportResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaListInstanceDiagnosisReportRequest struct {
	RegionID          string  /*  资源池 ID  */
	ResourceID        string  /*  资源 ID  */
	ResourceType      string  /*  资源类型, 仅支持 eip / natgw / elb  */
	DiagnosisReportID *string /*  实例诊断报告 ID  */
	DiagnosisRecordID *string /*  实例诊断记录 ID  */
	PageNumber        int32   /*  列表的页码，默认值为 1。  */
	PageSize          int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcVniaListInstanceDiagnosisReportResponse struct {
	StatusCode  int32                                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcVniaListInstanceDiagnosisReportReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaListInstanceDiagnosisReportReturnObjResponse struct {
	Results      []*CtvpcVniaListInstanceDiagnosisReportReturnObjResultsResponse `json:"results"`      /*  网卡列表  */
	TotalCount   int32                                                           `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                           `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                           `json:"totalPage"`    /*  总页数  */
}

type CtvpcVniaListInstanceDiagnosisReportReturnObjResultsResponse struct {
	DiagnosisReportID *string `json:"diagnosisReportID,omitempty"` /*  实例诊断报告 ID  */
	CreatedAt         *string `json:"createdAt,omitempty"`         /*  创建时间  */
	UpdatedAt         *string `json:"updatedAt,omitempty"`         /*  更新时间  */
	ItemCategory      *string `json:"itemCategory,omitempty"`      /*  诊断类型: config_diagnosis配置诊断 heath_check健康检测  */
	DiagnosisRecordID *string `json:"diagnosisRecordID,omitempty"` /*  实例诊断记录 ID  */
	ResourceID        *string `json:"resourceID,omitempty"`        /*  资源 ID  */
	ResourceType      *string `json:"resourceType,omitempty"`      /*  资源类型, 仅支持 eip / natgw / elb  */
	ErrMsg            *string `json:"errMsg,omitempty"`            /*  错误信息  */
	DiagnosisMsg      *string `json:"diagnosisMsg,omitempty"`      /*  诊断信息  */
	ItemCode          *string `json:"itemCode,omitempty"`          /*  诊断项，具体含义请参考文档：https://docs.qq.com/doc/DZWlNYU9RVlRMWEx4  */
	ItemStatus        int32   `json:"itemStatus"`                  /*  诊断状态：0未通过。1通过，2诊断系统异常  */
}
