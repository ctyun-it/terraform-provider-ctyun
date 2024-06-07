package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// KeypairCreateApi 创建一对SSH密钥对
// https://www.ctyun.cn/document/10026730/10106637
type KeypairCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewKeypairCreateApi(client *ctyunsdk.CtyunClient) *KeypairCreateApi {
	return &KeypairCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/keypair/create-keypair",
		},
	}
}

func (this *KeypairCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *KeypairCreateRequest) (*KeypairCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(keypairCreateRealRequest{
		RegionID:    req.RegionId,
		KeypairName: req.KeyPairName,
		ProjectID:   req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	result := &KeypairCreateRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	return &KeypairCreateResponse{
		PublicKey:   result.PublicKey,
		PrivateKey:  result.PrivateKey,
		KeyPairName: result.KeyPairName,
		FingerPrint: result.FingerPrint,
		KeyPairId:   result.KeyPairID,
	}, nil
}

type keypairCreateRealRequest struct {
	RegionID    string `json:"regionID"`
	KeypairName string `json:"keypairName"`
	ProjectID   string `json:"projectID,omitempty"`
}

type KeypairCreateRealResponse struct {
	PublicKey   string `json:"publicKey"`
	PrivateKey  string `json:"privateKey"`
	KeyPairName string `json:"keyPairName"`
	FingerPrint string `json:"fingerPrint"`
	KeyPairID   string `json:"keyPairID"`
}

type KeypairCreateRequest struct {
	RegionId    string // 资源池ID
	KeyPairName string // 密钥对名称。只能由数字、字母、-组成,不能以数字和-开头、以-结尾,且长度为2-63字符
	ProjectId   string // 企业项目id
}

type KeypairCreateResponse struct {
	PublicKey   string // 密钥对的公钥
	PrivateKey  string // 密钥对的私钥
	KeyPairName string // 密钥对名称
	FingerPrint string // 密钥对的指纹，采用MD5信息摘要算法
	KeyPairId   string // 密钥对的ID
}
