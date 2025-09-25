package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesQueryDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesQueryDetailApi(client *ctyunsdk.CtyunClient) *AmqpInstancesQueryDetailApi {
	return &AmqpInstancesQueryDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/query/detail",
		},
	}
}

func (this *AmqpInstancesQueryDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesQueryDetailRequest) (res *AmqpInstancesQueryDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddHeader("regionId", req.RegionId)
	builder.AddParam("prodInstId", req.ProdInstId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstancesQueryDetailResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesQueryDetailRequest struct {
	RegionId   string `json:"regionId"`
	ProdInstId string `json:"prodInstId"`
}

type AmqpInstancesQueryDetailResponse struct {
	ReturnObj  *AmqpInstancesQueryDetailResponseReturnObj `json:"returnObj"`
	Message    string                                     `json:"message"`
	StatusCode string                                     `json:"statusCode"`
}

type AmqpInstancesQueryDetailResponseReturnObj struct {
	Data *AmqpInstancesQueryDetailResponseReturnObjData `json:"data"`
}

type AmqpInstancesQueryDetailResponseReturnObjData struct {
	Id            int         `json:"id"`
	TenantName    string      `json:"tenantName"`
	TenantCode    string      `json:"tenantCode"`
	UserId        string      `json:"userId"`
	Cluster       string      `json:"cluster"`
	ClusterName   string      `json:"clusterName"`
	Status        int32       `json:"status"`
	ProdType      interface{} `json:"prodType"`
	Prod          string      `json:"prod"`
	TopicsNum     interface{} `json:"topicsNum"`
	Space         string      `json:"space"`
	BillMode      string      `json:"billMode"`
	Network       string      `json:"network"`
	Subnet        string      `json:"subnet"`
	ElasticIp     string      `json:"elasticIp"`
	SecurityGroup string      `json:"securityGroup"`
	DiskType      string      `json:"diskType"`
	EngineType    string      `json:"engineType"`
	OperOrderSrc  int         `json:"operOrderSrc"`
	RegionCode    string      `json:"regionCode"`
	RegionName    string      `json:"regionName"`
	Endpoint      string      `json:"endpoint"`
	SslEndpoint   string      `json:"sslEndpoint"`
	ProdInstId    string      `json:"prodInstId"`
	ExpireTime    string      `json:"expireTime"`
	CreateTime    string      `json:"createTime"`
	NodeCount     int32       `json:"nodeCount"`
}
