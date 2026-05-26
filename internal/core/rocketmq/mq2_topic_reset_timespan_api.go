package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicResetTimespanApi
/* 查询主题可重置的时间范围1 */
type Mq2TopicResetTimespanApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicResetTimespanApi(client *core.CtyunClient) *Mq2TopicResetTimespanApi {
	return &Mq2TopicResetTimespanApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/consumer/timeSpan",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2TopicResetTimespanApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2TopicResetTimespanRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2TopicResetTimespanApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2TopicResetTimespanRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
}

type Mq2TopicResetTimespanResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
