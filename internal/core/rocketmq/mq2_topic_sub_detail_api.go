package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicSubDetailApi
/* 查看Topic的订阅信息 */
type Mq2TopicSubDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicSubDetailApi(client *core.CtyunClient) *Mq2TopicSubDetailApi {
	return &Mq2TopicSubDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/topic/subDetail",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2TopicSubDetailApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2TopicSubDetailRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2TopicSubDetailApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2TopicSubDetailRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  int32  `json:"topicName"`  /*  topic名字  */
}

type Mq2TopicSubDetailResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
