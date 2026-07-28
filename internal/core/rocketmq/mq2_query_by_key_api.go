package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryByKeyApi
/* 根据Message_Key查询消息 */
type Mq2QueryByKeyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryByKeyApi(client *core.CtyunClient) *Mq2QueryByKeyApi {
	return &Mq2QueryByKeyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/message/queryByKey",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2QueryByKeyApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2QueryByKeyRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2QueryByKeyApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2QueryByKeyRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  主题名字  */
	Key        string `json:"key"`        /*  用于查询的消息KEY值  */
}

type Mq2QueryByKeyResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
