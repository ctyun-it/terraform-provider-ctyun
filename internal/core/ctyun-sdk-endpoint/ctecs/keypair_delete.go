package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// KeypairDeleteApi 删除SSH密钥对
// https://www.ctyun.cn/document/10026730/10040173
type KeypairDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewKeypairDeleteApi(client *ctyunsdk.CtyunClient) *KeypairDeleteApi {
	return &KeypairDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/keypair/delete",
		},
	}
}

func (this *KeypairDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *KeypairDeleteRequest) (*KeypairDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(keypairDeleteRealRequest{
		RegionID:    req.RegionId,
		KeyPairName: req.KeyPairName,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	result := &KeypairDeleteRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	return &KeypairDeleteResponse{
		KeyPairName: result.KeyPairName,
	}, nil
}

type keypairDeleteRealRequest struct {
	RegionID    string `json:"regionID"`
	KeyPairName string `json:"keyPairName"`
}

type KeypairDeleteRealResponse struct {
	KeyPairName string `json:"keyPairName"`
}

type KeypairDeleteRequest struct {
	RegionId    string // 资源池ID
	KeyPairName string // 密钥对名称。
}

type KeypairDeleteResponse struct {
	KeyPairName string // 密钥对名称
}
