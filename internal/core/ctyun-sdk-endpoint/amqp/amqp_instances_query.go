package amqp

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesQueryApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesQueryApi(client *ctyunsdk.CtyunClient) *AmqpInstancesQueryApi {
	return &AmqpInstancesQueryApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/query",
		},
	}
}

func (this *AmqpInstancesQueryApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesQueryRequest) (res *AmqpInstancesQueryResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddParam("pageNum", fmt.Sprintf("%d", req.PageNum))
	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))
	builder.AddHeader("regionId", req.RegionId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstancesQueryResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesQueryRequest struct {
	RegionId string `json:"regionId"`
	PageNum  int32  `json:"pageNum"`
	PageSize int32  `json:"pageSize"`
}

type AmqpInstancesQueryResponse struct {
	ReturnObj  *AmqpInstancesQueryResponseReturnObj `json:"returnObj"`
	Message    string                               `json:"message"`
	StatusCode string                               `json:"statusCode"`
}

type AmqpInstancesQueryResponseReturnObj struct {
	Total int32                                      `json:"total"`
	Data  []*AmqpInstancesQueryResponseReturnObjData `json:"data"`
}

type AmqpInstancesQueryResponseReturnObjData struct {
	Cluster       string      `json:"cluster"`       // 实例id
	Subnet        string      `json:"subnet"`        // 子网名称？
	Prod          string      `json:"prod"`          // 规格
	EngineType    string      `json:"engineType"`    // 引擎类型
	BillMode      string      `json:"billMode"`      // 账单
	SecurityGroup string      `json:"securityGroup"` // 安全组名称
	ProdType      interface{} `json:"prodType"`
	Network       string      `json:"network"`     // vpc名称?
	ExpireTime    string      `json:"expireTime"`  // 过期时间
	CreateTime    string      `json:"createTime"`  // 创建时间
	ClusterName   string      `json:"clusterName"` // 实例名称
	ProdInstId    string      `json:"prodInstId"`  // 实例id
	Status        int32       `json:"status"`      // 状态
}
