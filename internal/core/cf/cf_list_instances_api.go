package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfListInstancesApi
/* 查询实例列表 */
type CfListInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListInstancesApi(client *core.CtyunClient) *CfListInstancesApi {
	return &CfListInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/observability/listInstances",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListInstancesApi) Do(ctx context.Context, credential core.Credential, req *CfListInstancesRequest) (*CfListInstancesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfListInstancesRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListInstancesRequest struct {
	RegionId     string  `json:"regionId"`              /*  资源池id  */
	InstanceID   *string `json:"instanceID,omitempty"`  /*  实例id  */
	FunctionName string  `json:"functionName"`          /*  函数名称  */
	Qualifier    *string `json:"qualifier,omitempty"`   /*  函数版本  */
	PageSize     *int32  `json:"pageSize,omitempty"`    /*  每页个数  */
	CurrentPage  *int32  `json:"currentPage,omitempty"` /*  页数  */
}

type CfListInstancesResponse struct {
	StatusCode *int32                            `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                           `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                           `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListInstancesReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListInstancesReturnObjResponse struct {
	Data       []*CfListInstancesReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListInstancesReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListInstancesReturnObjDataResponse struct {
	DurationMS     *int32  `json:"durationMS"`     /*  实例运行时长(ms)  */
	FunctionName   *string `json:"functionName"`   /*  函数名称  */
	IpAddress      *string `json:"ipAddress"`      /*  podIP  */
	Status         *string `json:"status"`         /*  实例状态  */
	InstanceID     *string `json:"instanceID"`     /*  实例id  */
	Revision       *string `json:"revision"`       /*  版本revision  */
	LatestRevision *bool   `json:"latestRevision"` /*  是否最新版本  */
	Qualifier      *string `json:"qualifier"`      /*  函数别名  */
	Version        *string `json:"version"`        /*  函数版本  */
	StartTime      *int32  `json:"startTime"`      /*  实例起始时间  */
}

type CfListInstancesReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}
