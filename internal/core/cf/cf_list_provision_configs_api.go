package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListProvisionConfigsApi
/* 查询函数版本预留实例配置列表 */
type CfListProvisionConfigsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListProvisionConfigsApi(client *core.CtyunClient) *CfListProvisionConfigsApi {
	return &CfListProvisionConfigsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/resources/provision-configs",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListProvisionConfigsApi) Do(ctx context.Context, credential core.Credential, req *CfListProvisionConfigsRequest) (*CfListProvisionConfigsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.FunctionName != nil && *req.FunctionName != "" {
		ctReq.AddParam("functionName", *req.FunctionName)
	}
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.FunctionId != nil && *req.FunctionId != "" {
		ctReq.AddParam("functionId", *req.FunctionId)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListProvisionConfigsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListProvisionConfigsRequest struct {
	RegionId     string  `json:"regionId"`               /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	FunctionName *string `json:"functionName,omitempty"` /*  函数名称、版本、别名，支持模糊匹配。若不指定则列出所有函数的预留配置  */
	PageIndex    *int32  `json:"pageIndex,omitempty"`    /*  页码，取值 >= 1，默认值为 1  */
	PageSize     *int32  `json:"pageSize,omitempty"`     /*  每页大小，取值范围[1, 100]，默认值为 50  */
	FunctionId   *string `json:"functionId,omitempty"`   /*  函数id。若不指定则列出所有函数的预留配置  */
}

type CfListProvisionConfigsResponse struct {
	StatusCode *int32                                   `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                  `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                  `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListProvisionConfigsReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListProvisionConfigsReturnObjResponse struct {
	Data       []*CfListProvisionConfigsReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListProvisionConfigsReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListProvisionConfigsReturnObjDataResponse struct {
	Current          *int32                                                         `json:"current"`          /*  实际资源个数  */
	CurrentError     *string                                                        `json:"currentError"`     /*  预留实例创建失败时的错误信息  */
	FunctionId       *string                                                        `json:"functionId"`       /*  函数id  */
	FunctionName     *string                                                        `json:"functionName"`     /*  函数名称  */
	KsvcName         *string                                                        `json:"ksvcName"`         /*  ksvc名称  */
	Version          *string                                                        `json:"version"`          /*  函数版本  */
	Target           *int32                                                         `json:"target"`           /*  目标资源个数  */
	IdleMode         *bool                                                          `json:"idleMode"`         /*  是否为闲置模式  */
	Qualifier        *string                                                        `json:"qualifier"`        /*  函数别名  */
	ScheduledActions []*CfListProvisionConfigsReturnObjDataScheduledActionsResponse `json:"scheduledActions"` /*  定时策略配置  */
}

type CfListProvisionConfigsReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}

type CfListProvisionConfigsReturnObjDataScheduledActionsResponse struct {
	ScheduleExpression *string `json:"scheduleExpression"` /*  定时配置表达式  */
	Name               *string `json:"name"`               /*  定时策略名称  */
	StartTime          *string `json:"startTime"`          /*  策略生效时间  */
	EndTime            *string `json:"endTime"`            /*  策略失效时间  */
	Target             *int32  `json:"target"`             /*  预留目标实例数  */
	Timezone           *string `json:"timezone"`           /*  时区  */
}
