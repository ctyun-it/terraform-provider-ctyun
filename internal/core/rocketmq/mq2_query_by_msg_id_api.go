package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryByMsgIdApi
/* 通过MsgId查询消息 */
type Mq2QueryByMsgIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryByMsgIdApi(client *core.CtyunClient) *Mq2QueryByMsgIdApi {
	return &Mq2QueryByMsgIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/message/queryByMsgId",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2QueryByMsgIdApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2QueryByMsgIdRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2QueryByMsgIdApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2QueryByMsgIdRequest struct {
	MsgId      string `json:"msgId"`      /*  消息ID  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
}

type Mq2QueryByMsgIdResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
