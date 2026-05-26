package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListTriggersV2024Api
/* 查询触发器列表 */
type CfListTriggersV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListTriggersV2024Api(client *core.CtyunClient) *CfListTriggersV2024Api {
	return &CfListTriggersV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/triggers",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListTriggersV2024Api) Do(ctx context.Context, credential core.Credential, req *CfListTriggersV2024Request) (*CfListTriggersV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.TriggerName != nil && *req.TriggerName != "" {
		ctReq.AddParam("triggerName", *req.TriggerName)
	}
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.Version != nil && *req.Version != "" {
		ctReq.AddParam("version", *req.Version)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListTriggersV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListTriggersV2024Request struct {
	FunctionName string  `json:"functionName"`          /*  函数名称，函数必须已存在。  */
	RegionId     string  `json:"regionId"`              /*  资源池id，标识不同的地区，如：华东1、西南1  */
	TriggerName  *string `json:"triggerName,omitempty"` /*  支持模糊搜索  */
	PageIndex    *int32  `json:"pageIndex,omitempty"`   /*  页码。如果pageIndex和pageSize都不传，则不分页，pageIndex<=0时，默认pageIndex=1  */
	PageSize     *int32  `json:"pageSize,omitempty"`    /*  每页大小。如果pageIndex和pageSize都不传，则不分页，pageSize<=0时，默认pageSize=10  */
	Version      *string `json:"version,omitempty"`     /*  触发器对应的版本或别名，版本包括特殊版本LATEST和普通版本1,2,...。  */
}

type CfListTriggersV2024Response struct {
	StatusCode *int32                                `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Code       *string                               `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string                               `json:"message"`    /*  错误提示信息  */
	ReturnObj  *CfListTriggersV2024ReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListTriggersV2024ReturnObjResponse struct {
	Data       []*CfListTriggersV2024ReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListTriggersV2024ReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListTriggersV2024ReturnObjDataResponse struct {
	CreatedAt     *string `json:"createdAt"`     /*  创建时间  */
	UpdatedAt     *string `json:"updatedAt"`     /*  更新时间  */
	Creator       *int32  `json:"creator"`       /*  创建者ID  */
	Editor        *int32  `json:"editor"`        /*  编辑者ID  */
	TriggerName   *string `json:"triggerName"`   /*  触发器名称  */
	TriggerConfig *string `json:"triggerConfig"` /*  触发器配置，JSON字符串 ，不同触发器类型配置不同，详情请查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=53&api=16023&data=42&isNormal=1&vid=40">创建触发器</a>  */
	TriggerType   *string `json:"triggerType"`   /*  触发器类型。
	schedule: 定时触发器
	http: Http触发器
	kafka：Kafka触发器
	rocketmq：RocketMQ触发器
	rabbitmq： RabbitMQ触发器
	mqtt: MQTT触发器
	als： 日志触发器
	apigateway: 云原生网关触发器
	zos: 对象存储触发器  */
	Status      *int32  `json:"status"`      /*  触发器状态。1：启用；2：禁用；3：系统禁用  */
	Version     *string `json:"version"`     /*  别名或版本  */
	Region      *string `json:"region"`      /*  资源池id，标识不同的地区，如：华东1、西南1  */
	FunctionId  *int32  `json:"functionId"`  /*  函数ID  */
	UrlInternet *string `json:"urlInternet"` /*  外网URL ，非http触发器无值  */
	UrlIntranet *string `json:"urlIntranet"` /*  内网URL，非http触发器无值  */
	IsVersion   *bool   `json:"isVersion"`   /*  是否是版本，用于和别名区分开来，LATEST是特殊的版本  */
}

type CfListTriggersV2024ReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}
