package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetTriggerApi
/* 查询触发器详情 */
type CfGetTriggerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetTriggerApi(client *core.CtyunClient) *CfGetTriggerApi {
	return &CfGetTriggerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/triggers/{triggerName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetTriggerApi) Do(ctx context.Context, credential core.Credential, req *CfGetTriggerRequest) (*CfGetTriggerResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("triggerName", req.TriggerName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetTriggerResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetTriggerRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称，函数必须已存在  */
	TriggerName  string `json:"triggerName"`  /*  触发器名称，触发器必须已存在  */
	RegionId     string `json:"regionId"`     /*  资源池id，标识不同的地区，如：华东1、西南1  */
}

type CfGetTriggerResponse struct {
	StatusCode int32                         `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      string                        `json:"error"`      /*  错误码，CF_0表示成功，其他表示出错  */
	Message    string                        `json:"message"`    /*  错误提示信息  */
	ReturnObj  *CfGetTriggerReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetTriggerReturnObjResponse struct {
	CreatedAt     *string `json:"createdAt"`     /*  创建时间  */
	UpdatedAt     *string `json:"updatedAt"`     /*  更新时间  */
	Creator       *int32  `json:"creator"`       /*  创建者ID  */
	Editor        *int32  `json:"editor"`        /*  编辑者ID  */
	TriggerName   *string `json:"triggerName"`   /*  触发器名称  */
	TriggerConfig *string `json:"triggerConfig"` /*  触发器配置，JSON字符串 ，不同触发器类型配置不同，详情请查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=53&api=16023&data=42&isNormal=1&vid=40">创建触发器</a>  */
	TriggerType   *string `json:"triggerType"`   /*  触发器类型。
	schedule：定时触发器
	http：Http触发器
	kafka：Kafka触发器
	rocketmq：RocketMQ触发器
	rabbitmq： RabbitMQ触发器
	mqtt：MQTT触发器
	als：日志触发器
	apigateway：云原生网关触发器
	zos：对象存储触发器  */
	Status      *int32  `json:"status"`      /*  触发器状态。1：启用；2：禁用；3：系统禁用  */
	Version     *string `json:"version"`     /*  版本或别名，版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
	Region      *string `json:"region"`      /*  资源池id，标识不同的地区，如：华东1、西南1  */
	FunctionId  *int32  `json:"functionId"`  /*  函数ID  */
	UrlInternet *string `json:"urlInternet"` /*  外网URL ，非http触发器无值  */
	UrlIntranet *string `json:"urlIntranet"` /*  内网URL，非http触发器无值  */
	IsVersion   *bool   `json:"isVersion"`   /*  是否是版本。版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
}
