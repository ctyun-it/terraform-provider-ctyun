package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// KeypairAttachApi 绑定SSH密钥对到Linux云主机
// https://www.ctyun.cn/document/10026730/10106635
type KeypairAttachApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewKeypairAttachApi(client *ctyunsdk.CtyunClient) *KeypairAttachApi {
	return &KeypairAttachApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/keypair/attach-instance",
		},
	}
}

func (this *KeypairAttachApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *KeypairAttachRequest) (*KeypairAttachResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&keypairAttachRealRequest{
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
	return &KeypairAttachResponse{}, nil
}

type keypairAttachRealRequest struct {
	RegionID    string `json:"regionID"`
	KeyPairName string `json:"keyPairName"`
	InstanceID  string `json:"instanceID"`
}

type KeypairAttachRequest struct {
	RegionId    string // 区域id
	KeyPairName string // 密钥对名称
	InstanceId  string // 云主机ID
}

type KeypairAttachResponse struct {
}
