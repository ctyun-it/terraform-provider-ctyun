package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesDiskExtendApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesDiskExtendApi(client *ctyunsdk.CtyunClient) *AmqpInstancesDiskExtendApi {
	return &AmqpInstancesDiskExtendApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/diskExtend",
		},
	}
}

func (this *AmqpInstancesDiskExtendApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesDiskExtendRequest) (res *AmqpInstancesDiskExtendResponse, err error) {
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
	res = &AmqpInstancesDiskExtendResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesDiskExtendRequest struct {
	ProdInstId     string `json:"prodInstId"`
	DiskExtendSize int32  `json:"diskExtendSize"`
	AutoPay        bool   `json:"autoPay"`
	RegionId       string `json:"regionId"`
}

type AmqpInstancesDiskExtendResponse struct {
	ReturnObj  *AmqpInstancesDiskExtendResponseReturnObj `json:"returnObj"`
	Message    string                                    `json:"message"`
	StatusCode string                                    `json:"statusCode"`
}

type AmqpInstancesDiskExtendResponseReturnObj struct {
}

type AmqpInstancesDiskExtendResponseReturnObjData struct {
}
