package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsBatchUpdateInstancesApi 包周期付费云主机标记到期转按需
type EcsBatchUpdateInstancesApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsBatchUpdateInstancesApi(client *ctyunsdk.CtyunClient) *EcsBatchUpdateInstancesApi {
	return &EcsBatchUpdateInstancesApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/batch-update-instances",
		},
	}
}

func (this *EcsBatchUpdateInstancesApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsBatchUpdateInstancesRequest) (*EcsBatchUpdateInstancesResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	var updateInfoReq []ecsBatchUpdateInstancesUpdateInfoRealRequest
	for _, request := range req.UpdateInfo {
		updateInfoReq = append(updateInfoReq, ecsBatchUpdateInstancesUpdateInfoRealRequest{
			InstanceID:  request.InstanceId,
			DisplayName: request.DisplayName,
		})
	}
	_, err := builder.WriteJson(&ecsBatchUpdateInstancesRealRequest{
		RegionID:   req.RegionId,
		UpdateInfo: updateInfoReq,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsBatchUpdateInstancesOrderInfoRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var orderInfoResp []EcsBatchUpdateInstancesOrderInfoResponse
	for _, info := range realResponse.UpdateInfo {
		orderInfoResp = append(orderInfoResp, EcsBatchUpdateInstancesOrderInfoResponse{
			Id:          info.ID,
			DisplayName: info.DisplayName,
		})
	}
	return &EcsBatchUpdateInstancesResponse{
		UpdateInfo: orderInfoResp,
	}, nil
}

type ecsBatchUpdateInstancesUpdateInfoRealRequest struct {
	InstanceID  string `json:"instanceID"`
	DisplayName string `json:"displayName"`
}

type ecsBatchUpdateInstancesRealRequest struct {
	RegionID   string                                         `json:"regionID"`
	UpdateInfo []ecsBatchUpdateInstancesUpdateInfoRealRequest `json:"updateInfo"`
}

type ecsBatchUpdateInstancesOrderInfoRealResponse struct {
	UpdateInfo []struct {
		ID          string `json:"ID"`
		DisplayName string `json:"displayName"`
	} `json:"updateInfo"`
}

type EcsBatchUpdateInstancesRequest struct {
	RegionId   string
	UpdateInfo []EcsBatchUpdateInstancesUpdateInfoRequest
}

type EcsBatchUpdateInstancesUpdateInfoRequest struct {
	InstanceId  string
	DisplayName string
}

type EcsBatchUpdateInstancesOrderInfoResponse struct {
	Id          string
	DisplayName string
}

type EcsBatchUpdateInstancesResponse struct {
	UpdateInfo []EcsBatchUpdateInstancesOrderInfoResponse
}
