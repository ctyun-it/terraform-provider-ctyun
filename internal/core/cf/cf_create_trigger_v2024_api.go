package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfCreateTriggerV2024Api
/* 创建触发器 */
type CfCreateTriggerV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfCreateTriggerV2024Api(client *core.CtyunClient) *CfCreateTriggerV2024Api {
	return &CfCreateTriggerV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/functions/{functionName}/triggers",
			ContentType:  "application/json",
		},
	}
}

func (a *CfCreateTriggerV2024Api) Do(ctx context.Context, credential core.Credential, req *CfCreateTriggerV2024Request) (*CfCreateTriggerV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfCreateTriggerV2024Request
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfCreateTriggerV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfCreateTriggerV2024Request struct {
	FunctionName  string `json:"functionName"`     /*  函数名称，函数必须存在。  */
	RegionId      string `json:"regionId"`         /*  资源池id，标识不同的地区，如：华东1、西南1  */
	TriggerName   string `json:"triggerName"`      /*  触发器名称。只能包含字母、数字和中划线。只能字母开头,字母数字结尾。长度在 3-63之间。相同资源池、相同函数下的触发器名称不能重复。  */
	TriggerConfig string `json:"triggerConfig"`    /*  触发器配置，JSON字符串 ，不同触发器类型配置不同。<br/>定时触发器：参见ScheduleTriggerConfig<br/>http触发器：参见HttpTriggerConfig<br/>Kafka触发器：参见KafkaTriggerConfig<br/>RocketMq触发器：参见RocketMqTriggerConfig<br/>RabbitMq触发器：参见RabbitMqTriggerConfig<br/>Mqtt触发器：参见MqttTriggerConfig<br/>对象存储触发器：参见ZosTriggerConfig<br/>云原生网关触发器：参见ApiGatewayTriggerConfig<br/>日志触发器：参见AlsTriggerConfig  */
	TriggerType   string `json:"triggerType"`      /*  触发器类型。<br/>schedule: 定时触发器<br/>http: Http触发器<br/>kafka：Kafka触发器<br/>rocketmq：RocketMQ触发器<br/>rabbitmq： RabbitMQ触发器<br/>mqtt: MQTT触发器<br/>als： 日志触发器<br>apigateway: 云原生网关触发器<br/>zos: 对象存储触发器  */
	Version       string `json:"version"`          /*  版本或别名，版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
	Enable        *bool  `json:"enable,omitempty"` /*  触发器是否立即启用，创建场景下固定为启用  */
}

type CfCreateTriggerV2024Response struct {
	StatusCode int32                                 `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      string                                `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    string                                 `json:"message"`
	ReturnObj  *CfCreateTriggerV2024ReturnObjResponse `json:"returnObj"` /*  返回实体  */
}

type CfCreateTriggerV2024ReturnObjResponse struct {
	CreatedAt     *string `json:"createdAt"`     /*  创建时间  */
	UpdatedAt     *string `json:"updatedAt"`     /*  更新时间  */
	Creator       *int32  `json:"creator"`       /*  创建者ID  */
	Editor        *int32  `json:"editor"`        /*  编辑者ID  */
	TriggerName   *string `json:"triggerName"`   /*  触发器名称  */
	TriggerConfig *string `json:"triggerConfig"` /*  触发器配置，JSON字符串，不同触发器类型配置不同，参考本页HttpTriggerConfig等对象。  */
	TriggerType   *string `json:"triggerType"`   /*  触发器类型  */
	Status        *int32  `json:"status"`        /*  1:启用；2:禁用；3:系统禁用  */
	Version       *string `json:"version"`       /*  别名或版本，版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
	Region        *string `json:"region"`        /*  区域ID  */
	FunctionId    *int32  `json:"functionId"`    /*  函数ID  */
	UrlInternet   *string `json:"urlInternet"`   /*  外网URL ，非http触发器无值  */
	UrlIntranet   *string `json:"urlIntranet"`   /*  内网URL，非http触发器无值  */
	IsVersion     *bool   `json:"isVersion"`     /*  是否是版本，版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
}
