package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ConsumerGroupOutputtpsApi
/* 查询订阅组一段时间内消费tps信息 */
type Mq2ConsumerGroupOutputtpsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ConsumerGroupOutputtpsApi(client *core.CtyunClient) *Mq2ConsumerGroupOutputtpsApi {
	return &Mq2ConsumerGroupOutputtpsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/group/outputTps",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2ConsumerGroupOutputtpsApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2ConsumerGroupOutputtpsRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2ConsumerGroupOutputtpsApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2ConsumerGroupOutputtpsRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  topic名字  */
	BrokerName string `json:"brokerName"` /*  Broker名字  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
}

type Mq2ConsumerGroupOutputtpsResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
