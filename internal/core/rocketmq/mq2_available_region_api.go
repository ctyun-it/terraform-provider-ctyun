package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2AvailableRegionApi
/* 查询可创建缓存的资源池列表 */
type Mq2AvailableRegionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2AvailableRegionApi(client *core.CtyunClient) *Mq2AvailableRegionApi {
	return &Mq2AvailableRegionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/availableRegion",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2AvailableRegionApi) Do(ctx context.Context, credential core.Credential, req *Mq2AvailableRegionRequest) (*Mq2AvailableRegionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2AvailableRegionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2AvailableRegionRequest struct{}

type Mq2AvailableRegionResponse struct{}
