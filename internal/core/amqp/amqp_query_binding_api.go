package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryBindingApi
/* 绑定查询
 */type AmqpQueryBindingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryBindingApi(client *core.CtyunClient) *AmqpQueryBindingApi {
	return &AmqpQueryBindingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/binding/query",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *AmqpQueryBindingApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *AmqpQueryBindingRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.AmqpQueryBindingApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type AmqpQueryBindingRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
}

type AmqpQueryBindingResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
}
