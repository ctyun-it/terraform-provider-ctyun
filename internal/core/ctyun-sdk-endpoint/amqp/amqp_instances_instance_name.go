package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesInstanceNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesInstanceNameApi(client *ctyunsdk.CtyunClient) *AmqpInstancesInstanceNameApi {
	return &AmqpInstancesInstanceNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/instanceName",
		},
	}
}

func (this *AmqpInstancesInstanceNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesInstanceNameRequest) (res *AmqpInstancesInstanceNameResponse, err error) {
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
	res = &AmqpInstancesInstanceNameResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesInstanceNameRequest struct {
	ProdInstId   string `json:"prodInstId"`
	InstanceName string `json:"instanceName"`
	RegionId     string `json:"regionId"`
}

type AmqpInstancesInstanceNameResponse struct {
	ReturnObj  *AmqpInstancesInstanceNameResponseReturnObj `json:"returnObj"`
	Message    string                                      `json:"message"`
	StatusCode string                                      `json:"statusCode"`
}

type AmqpInstancesInstanceNameResponseReturnObj struct {
}

type AmqpInstancesInstanceNameResponseReturnObjData struct {
}
