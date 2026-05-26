package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryDlqbyIdApi
/* 通过传入MessageID查询指定的死信消息 */
type Mq2QueryDlqbyIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryDlqbyIdApi(client *core.CtyunClient) *Mq2QueryDlqbyIdApi {
	return &Mq2QueryDlqbyIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/dlq/getById",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2QueryDlqbyIdApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2QueryDlqbyIdRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2QueryDlqbyIdApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2QueryDlqbyIdRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	MsgId      string `json:"msgId"`      /*  消息ID  */
}

type Mq2QueryDlqbyIdResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
