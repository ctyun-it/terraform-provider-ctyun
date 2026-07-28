package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2AccumulateInfoApi
/* 查询订阅组的消息消费堆积 */
type Mq2AccumulateInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2AccumulateInfoApi(client *core.CtyunClient) *Mq2AccumulateInfoApi {
	return &Mq2AccumulateInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/consumer/accumulate",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2AccumulateInfoApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2AccumulateInfoRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2AccumulateInfoApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2AccumulateInfoRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
}

type Mq2AccumulateInfoResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
