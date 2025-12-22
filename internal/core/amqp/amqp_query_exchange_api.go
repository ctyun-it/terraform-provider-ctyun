package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryExchangeApi
/* 查询交换器
 */type AmqpQueryExchangeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryExchangeApi(client *core.CtyunClient) *AmqpQueryExchangeApi {
	return &AmqpQueryExchangeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/exchange/query",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *AmqpQueryExchangeApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *AmqpQueryExchangeRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.AmqpQueryExchangeApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type AmqpQueryExchangeRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
	Name       string `json:"name,omitempty"`       /*  交换器名称（模糊匹配）  */
	PageNum    int32  `json:"pageNum,omitempty"`    /*  当前页面（默认1）  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  分页大小（默认100）  */
}

type AmqpQueryExchangeResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
