package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryDlqbyTimeApi
/* 查询指定时间段内订阅组内存在的死信消息 */
type Mq2QueryDlqbyTimeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryDlqbyTimeApi(client *core.CtyunClient) *Mq2QueryDlqbyTimeApi {
	return &Mq2QueryDlqbyTimeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/dlq/getByTime",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2QueryDlqbyTimeApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2QueryDlqbyTimeRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2QueryDlqbyTimeApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2QueryDlqbyTimeRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  分页中的页数，默认1  */
	BeginTime  int64  `json:"beginTime"`  /*  开始时间的毫秒时间戳  */
	EndTime    int64  `json:"endTime"`    /*  结束时间的毫秒时间戳  */
}

type Mq2QueryDlqbyTimeResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
