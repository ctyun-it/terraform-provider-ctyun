package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryInstanceApi
/* 查询租户实例
 */type AmqpQueryInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryInstanceApi(client *core.CtyunClient) *AmqpQueryInstanceApi {
	return &AmqpQueryInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instances/query",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *AmqpQueryInstanceApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *AmqpQueryInstanceRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.AmqpQueryInstanceApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type AmqpQueryInstanceRequest struct {
	PageNum    string `json:"pageNum,omitempty"`    /*  分页中的页数，默认1  */
	PageSize   string `json:"pageSize,omitempty"`   /*  分页中的每页大小，默认10  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例Id，如果填入，则返回指定实例信息  */
}

type AmqpQueryInstanceResponse struct {
	StatusCode string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string                              `json:"message"`    /*  描述状态  */
	ReturnObj  *AmqpQueryInstanceReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}

type AmqpQueryInstanceReturnObjResponse struct {
	Total string                                    `json:"total"` /*  总数  */
	Data  []*AmqpQueryInstanceReturnObjDataResponse `json:"data"`  /*  数据  */
}

type AmqpQueryInstanceReturnObjDataResponse struct {
	Cluster       string `json:"cluster"`       /*  实例ID  */
	ClusterName   string `json:"clusterName"`   /*  实例名称  */
	Status        int32  `json:"status"`        /*  状态 1正常，2暂停，3注销  */
	ProdType      string `json:"prodType"`      /*  产品规格类型：10009002基础版，10009001高级版  */
	CreateTime    int64  `json:"createTime"`    /*  创建时间，时间戳  */
	ExpireTime    int64  `json:"expireTime"`    /*  过期时间，时间戳  */
	SecurityGroup string `json:"securityGroup"` /*  安全组  */
	Network       string `json:"network"`       /*  网络  */
	Subnet        string `json:"subnet"`        /*  子网  */
}
