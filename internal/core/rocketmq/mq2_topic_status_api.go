package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicStatusApi
/* 查询Topic状态 */
type Mq2TopicStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicStatusApi(client *core.CtyunClient) *Mq2TopicStatusApi {
	return &Mq2TopicStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/topic/status",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2TopicStatusApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2TopicStatusRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2TopicStatusApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2TopicStatusRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  topic名字  */
}

type Mq2TopicStatusResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
