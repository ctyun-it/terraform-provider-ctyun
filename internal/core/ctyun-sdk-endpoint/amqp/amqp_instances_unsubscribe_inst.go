package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesUnsubscribeInstApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesUnsubscribeInstApi(client *ctyunsdk.CtyunClient) *AmqpInstancesUnsubscribeInstApi {
	return &AmqpInstancesUnsubscribeInstApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/unsubscribeInst",
		},
	}
}

func (this *AmqpInstancesUnsubscribeInstApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesUnsubscribeInstRequest) (res *AmqpInstancesUnsubscribeInstResponse, err error) {
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
	res = &AmqpInstancesUnsubscribeInstResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesUnsubscribeInstRequest struct {
	ProdInstId string `json:"prodInstId"`
	RegionId   string `json:"regionId"`
}

type AmqpInstancesUnsubscribeInstResponse struct {
	ReturnObj  *AmqpInstancesUnsubscribeInstResponseReturnObj `json:"returnObj"`
	Message    string                                         `json:"message"`
	StatusCode string                                         `json:"statusCode"`
}

type AmqpInstancesUnsubscribeInstResponseReturnObj struct {
}

type AmqpInstancesUnsubscribeInstResponseReturnObjData struct {
}
