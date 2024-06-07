package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// KeypairImportApi 导入RSA密钥对公钥
// https://www.ctyun.cn/document/10026730/10106638
type KeypairImportApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewKeypairImportApi(client *ctyunsdk.CtyunClient) *KeypairImportApi {
	return &KeypairImportApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/keypair/import-keypair",
		},
	}
}

func (this *KeypairImportApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *KeypairImportRequest) (*KeypairImportResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(keypairImportRealRequest{
		RegionID:    req.RegionId,
		KeypairName: req.KeyPairName,
		PublicKey:   req.PublicKey,
		ProjectID:   req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	result := &keypairImportRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	return &KeypairImportResponse{
		PublicKey:   result.PublicKey,
		KeyPairName: result.KeyPairName,
		FingerPrint: result.FingerPrint,
		KeyPairId:   result.KeyPairID,
	}, nil
}

type keypairImportRealRequest struct {
	RegionID    string `json:"regionID"`
	KeypairName string `json:"keyPairName"`
	PublicKey   string `json:"publicKey"`
	ProjectID   string `json:"projectID,omitempty"`
}

type keypairImportRealResponse struct {
	PublicKey   string `json:"publicKey"`
	KeyPairName string `json:"keyPairName"`
	FingerPrint string `json:"fingerPrint"`
	KeyPairID   string `json:"keyPairID"`
}

type KeypairImportRequest struct {
	RegionId    string // 资源池ID
	KeyPairName string // 密钥对名称。只能由数字、字母、-组成,不能以数字和-开头、以-结尾,且长度为2-63字符
	PublicKey   string // 公钥
	ProjectId   string // 企业项目id
}

type KeypairImportResponse struct {
	PublicKey   string // 密钥对的公钥
	KeyPairName string // 密钥对名称
	FingerPrint string // 密钥对的指纹，采用MD5信息摘要算法
	KeyPairId   string // 密钥对的ID
}
