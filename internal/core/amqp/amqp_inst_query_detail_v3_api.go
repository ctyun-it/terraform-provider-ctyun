package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpInstQueryDetailV3Api
/* 查询实例详情V3
 */type AmqpInstQueryDetailV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpInstQueryDetailV3Api(client *core.CtyunClient) *AmqpInstQueryDetailV3Api {
	return &AmqpInstQueryDetailV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/query/detail",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpInstQueryDetailV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpInstQueryDetailV3Request) (*AmqpInstQueryDetailV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpInstQueryDetailV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpInstQueryDetailV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例id  */
}

type AmqpInstQueryDetailV3Response struct {
	ReturnObj  *AmqpInstQueryDetailV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                                  `json:"message"`    /*  描述状态  */
	StatusCode string                                  `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                                  `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpInstQueryDetailV3ReturnObjResponse struct {
	ClusterName   string `json:"clusterName"`   /*  实例名称  */
	Status        string `json:"status"`        /*  运行状态  */
	Prod          string `json:"prod"`          /*  规格  */
	Space         string `json:"space"`         /*  消息存储大小  */
	BillMode      string `json:"billMode"`      /*  付费类型  */
	Network       string `json:"network"`       /*  vpc  */
	Subnet        string `json:"subnet"`        /*  子网  */
	ElasticIp     string `json:"elasticIp"`     /*  公网接入点  */
	SecurityGroup string `json:"securityGroup"` /*  安全组  */
	ExpireTime    string `json:"expireTime"`    /*  到期时间  */
	CreateTime    string `json:"createTime"`    /*  创建时间  */
	DiskType      string `json:"diskType"`      /*  磁盘类型  */
	EngineType    string `json:"engineType"`    /*  引擎类型  */
	Endpoint      string `json:"endpoint"`      /*  安全接入点、负载均衡接入点  */
	SslEndpoint   string `json:"sslEndpoint"`   /*  SSL接入点  */
	ProdInstId    string `json:"prodInstId"`    /*  实例id  */
	NodeCount     string `json:"nodeCount"`     /*  节点数  */
	RegionCode    string `json:"regionCode"`    /*  资源池id  */
	RegionName    string `json:"regionName"`    /*  资源池名称  */
}
