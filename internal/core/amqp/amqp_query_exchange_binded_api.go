package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryExchangeBindedApi
/* 查交换机被绑定
 */type AmqpQueryExchangeBindedApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryExchangeBindedApi(client *core.CtyunClient) *AmqpQueryExchangeBindedApi {
	return &AmqpQueryExchangeBindedApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/binding/exchange/destination",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *AmqpQueryExchangeBindedApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *AmqpQueryExchangeBindedRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.AmqpQueryExchangeBindedApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type AmqpQueryExchangeBindedRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
	Exchange   string `json:"exchange,omitempty"`   /*  交换器名称  */
}

type AmqpQueryExchangeBindedResponse struct {
	StatusCode string                                      `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string                                      `json:"message"`    /*  描述状态  */
	ReturnObj  []*AmqpQueryExchangeBindedReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}

type AmqpQueryExchangeBindedReturnObjResponse struct{}
