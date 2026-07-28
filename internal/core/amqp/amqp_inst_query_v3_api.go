package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpInstQueryV3Api
/* 查询租户实例v3
 */type AmqpInstQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpInstQueryV3Api(client *core.CtyunClient) *AmqpInstQueryV3Api {
	return &AmqpInstQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpInstQueryV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpInstQueryV3Request) (*AmqpInstQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageNum != "" {
		ctReq.AddParam("pageNum", req.PageNum)
	}
	ctReq.AddParam("pageSize", req.PageSize)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpInstQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpInstQueryV3Request struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池id  */
	PageNum  string `json:"pageNum,omitempty"`  /*  分页的页序号  */
	PageSize string `json:"pageSize,omitempty"` /*  分页的大小  */
}

type AmqpInstQueryV3Response struct {
	ReturnObj  *AmqpInstQueryV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                            `json:"message"`    /*  描述状态  */
	StatusCode string                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                            `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpInstQueryV3ReturnObjResponse struct {
	Total int32                                   `json:"total"` /*  集群总数  */
	Data  []*AmqpInstQueryV3ReturnObjDataResponse `json:"data"`  /*  集群信息  */
}

type AmqpInstQueryV3ReturnObjDataResponse struct {
	Prod        string `json:"prod"`        /*  规格  */
	EngineType  string `json:"engineType"`  /*  引擎类型，rabbitmq.cluster表示rabbitmq引擎，pulsar.cluster表示云原生引擎  */
	BillMode    string `json:"billMode"`    /*  计费模式，1表示包年包月，2表示按需计费  */
	ExpireTime  string `json:"expireTime"`  /*  到期时间  */
	CreateTime  string `json:"createTime"`  /*  创建时间  */
	ClusterName string `json:"clusterName"` /*  实例名称  */
	ProdInstId  string `json:"prodInstId"`  /*  实例id  */
	Status      string `json:"status"`      /*  实例状态，1表示运行中，3表示已注销，4表示已退订，5表示变更中，6表示创建中  */
	Cluster     string `json:"cluster"`     /*  实例id  */
}
