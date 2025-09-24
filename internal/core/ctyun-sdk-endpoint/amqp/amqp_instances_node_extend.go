package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesNodeExtendApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesNodeExtendApi(client *ctyunsdk.CtyunClient) *AmqpInstancesNodeExtendApi {
	return &AmqpInstancesNodeExtendApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/nodeExtend",
		},
	}
}

func (this *AmqpInstancesNodeExtendApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesNodeExtendRequest) (res *AmqpInstancesNodeExtendResponse, err error) {
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
	res = &AmqpInstancesNodeExtendResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesNodeExtendRequest struct {
	ProdInstId    string `json:"prodInstId"`
	ExtendNodeNum int32  `json:"extendNodeNum"`
	AutoPay       bool   `json:"autoPay"`
	RegionId      string `json:"regionId"`
}

type AmqpInstancesNodeExtendResponse struct {
	ReturnObj  *AmqpInstancesNodeExtendResponseReturnObj `json:"returnObj"`
	Message    string                                    `json:"message"`
	StatusCode string                                    `json:"statusCode"`
}

type AmqpInstancesNodeExtendResponseReturnObj struct {
}

type AmqpInstancesNodeExtendResponseReturnObjData struct {
}
