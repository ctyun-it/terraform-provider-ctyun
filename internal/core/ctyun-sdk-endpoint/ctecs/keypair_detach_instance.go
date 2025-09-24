package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// KeypairDetachApi 为云主机解绑SSH密钥对
// https://www.ctyun.cn/document/10026730/10106639
type KeypairDetachApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewKeypairDetachApi(client *ctyunsdk.CtyunClient) *KeypairDetachApi {
	return &KeypairDetachApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/keypair/detach-instance",
		},
	}
}

func (this *KeypairDetachApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *KeypairDetachRequest) (*KeypairDetachResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&keypairDetachRealRequest{
		RegionID:    req.RegionId,
		KeyPairName: req.KeyPairName,
		InstanceID:  req.InstanceId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &KeypairDetachResponse{}, nil
}

type keypairDetachRealRequest struct {
	RegionID    string `json:"regionID"`
	KeyPairName string `json:"keyPairName"`
	InstanceID  string `json:"instanceID"`
}

type KeypairDetachRequest struct {
	RegionId    string // 区域id
	KeyPairName string // 密钥对名称
	InstanceId  string // 云主机ID
}

type KeypairDetachResponse struct {
}
