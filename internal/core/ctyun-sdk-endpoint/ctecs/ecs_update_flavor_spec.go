package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsUpdateFlavorSpecApi 云主机修改规格
// https://www.ctyun.cn/document/10026730/10106612
type EcsUpdateFlavorSpecApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsUpdateFlavorSpecApi(client *ctyunsdk.CtyunClient) *EcsUpdateFlavorSpecApi {
	return &EcsUpdateFlavorSpecApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/update-flavor-spec",
		},
	}
}

func (this *EcsUpdateFlavorSpecApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsUpdateFlavorSpecRequest) (*EcsUpdateFlavorSpecResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsUpdateFlavorSpecRealRequest{
		RegionID:    req.RegionId,
		InstanceID:  req.InstanceId,
		FlavorID:    req.FlavorId,
		ClientToken: req.ClientToken,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsUpdateFlavorSpecRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsUpdateFlavorSpecResponse{
		MasterOrderNo: realResponse.MasterOrderNO,
		RegionId:      realResponse.RegionID,
		MasterOrderId: realResponse.MasterOrderID,
	}, nil
}

type ecsUpdateFlavorSpecRealRequest struct {
	RegionID    string `json:"regionID"`
	InstanceID  string `json:"instanceID"`
	FlavorID    string `json:"flavorID"`
	ClientToken string `json:"clientToken"`
}

type EcsUpdateFlavorSpecRequest struct {
	RegionId    string
	InstanceId  string
	FlavorId    string
	ClientToken string
}

type ecsUpdateFlavorSpecRealResponse struct {
	MasterOrderNO string `json:"masterOrderNO"`
	RegionID      string `json:"regionID"`
	MasterOrderID string `json:"masterOrderID"`
}

type EcsUpdateFlavorSpecResponse struct {
	MasterOrderNo string
	RegionId      string
	MasterOrderId string
}
