package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpMetadataApi
/* 查询实例交换机、队列以及虚拟机数量
 */type AmqpMetadataApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpMetadataApi(client *core.CtyunClient) *AmqpMetadataApi {
	return &AmqpMetadataApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instances/metadata",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *AmqpMetadataApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *AmqpMetadataRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.AmqpMetadataApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type AmqpMetadataRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例Id，如果填入，则返回指定实例信息  */
}

type AmqpMetadataResponse struct {
	StatusCode string                         `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string                         `json:"message"`    /*  描述状态  */
	ReturnObj  *AmqpMetadataReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}

type AmqpMetadataReturnObjResponse struct {
	Data *AmqpMetadataReturnObjDataResponse `json:"data"` /*  数据  */
}

type AmqpMetadataReturnObjDataResponse struct {
	CurrentExchanges    string `json:"CurrentExchanges"`    /*  交换器数量  */
	CurrentQueues       string `json:"CurrentQueues"`       /*  队列数量  */
	CurrentVirtualHosts string `json:"CurrentVirtualHosts"` /*  虚拟机数量  */
}
