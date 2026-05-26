package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryByTimeApi
/* 查询指定时间段内Topic的消息 */
type Mq2QueryByTimeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryByTimeApi(client *core.CtyunClient) *Mq2QueryByTimeApi {
	return &Mq2QueryByTimeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/message/queryByTime",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2QueryByTimeApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2QueryByTimeRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2QueryByTimeApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2QueryByTimeRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  主题名字  */
	BeginTime  int64  `json:"beginTime"`  /*  开始时间的毫秒时间戳  */
	EndTime    int64  `json:"endTime"`    /*  结束时间的毫秒时间戳  */
}

type Mq2QueryByTimeResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
