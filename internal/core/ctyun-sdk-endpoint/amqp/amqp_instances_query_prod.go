package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesQueryProdApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesQueryProdApi(client *ctyunsdk.CtyunClient) *AmqpInstancesQueryProdApi {
	return &AmqpInstancesQueryProdApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/queryProd",
		},
	}
}

func (this *AmqpInstancesQueryProdApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesQueryProdRequest) (res *AmqpInstancesQueryProdResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddHeader("regionId", req.RegionId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstancesQueryProdResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesQueryProdRequest struct {
	RegionId string `json:"regionId"`
}

type AmqpInstancesQueryProdResponse struct {
	ReturnObj  *AmqpInstancesQueryProdResponseReturnObj `json:"returnObj"`
	Message    string                                   `json:"message"`
	StatusCode string                                   `json:"statusCode"`
}

type AmqpInstancesQueryProdResponseReturnObj struct {
	Data []AmqpInstancesQueryProdResponseReturnObjData `json:"data"`
}

type AmqpInstancesQueryProdResponseReturnObjData struct {
	FlavorID      string      `json:"flavorID"`
	SpecName      string      `json:"specName"`
	FlavorType    string      `json:"flavorType"`
	FlavorName    string      `json:"flavorName"`
	CpuNum        int32       `json:"cpuNum"`
	MemSize       int32       `json:"memSize"`
	MultiQueue    int32       `json:"multiQueue"`
	Pps           int32       `json:"pps"`
	BandwidthBase float64     `json:"bandwidthBase"`
	BandwidthMax  int32       `json:"bandwidthMax"`
	CpuArch       interface{} `json:"cpuArch"`
	Series        string      `json:"series"`
	AzList        []string    `json:"azList"`
}
