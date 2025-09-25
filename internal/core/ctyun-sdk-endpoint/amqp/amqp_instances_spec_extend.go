package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesSpecExtendApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesSpecExtendApi(client *ctyunsdk.CtyunClient) *AmqpInstancesSpecExtendApi {
	return &AmqpInstancesSpecExtendApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/specExtend",
		},
	}
}

func (this *AmqpInstancesSpecExtendApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesSpecExtendRequest) (res *AmqpInstancesSpecExtendResponse, err error) {
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
	res = &AmqpInstancesSpecExtendResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesSpecExtendRequest struct {
	ProdInstId string `json:"prodInstId"`
	SpecName   string `json:"specName"`
	AutoPay    bool   `json:"autoPay"`
	RegionId   string `json:"regionId"`
}

type AmqpInstancesSpecExtendResponse struct {
	ReturnObj  *AmqpInstancesSpecExtendResponseReturnObj `json:"returnObj"`
	Message    string                                    `json:"message"`
	StatusCode string                                    `json:"statusCode"`
}

type AmqpInstancesSpecExtendResponseReturnObj struct {
}

type AmqpInstancesSpecExtendResponseReturnObjData struct {
}
