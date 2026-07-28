package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryVhostApi
/* 查询虚拟机
 */type AmqpQueryVhostApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryVhostApi(client *core.CtyunClient) *AmqpQueryVhostApi {
	return &AmqpQueryVhostApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/vhost/query",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *AmqpQueryVhostApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *AmqpQueryVhostRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.AmqpQueryVhostApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type AmqpQueryVhostRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
}

type AmqpQueryVhostResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
