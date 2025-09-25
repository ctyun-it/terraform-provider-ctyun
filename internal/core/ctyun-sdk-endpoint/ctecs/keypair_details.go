package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// KeypairDetailApi 查询一个或多个密钥对
// https://www.ctyun.cn/document/10026730/10106636
type KeypairDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewKeypairDetailApi(client *ctyunsdk.CtyunClient) *KeypairDetailApi {
	return &KeypairDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/keypair/details",
		},
	}
}

func (this *KeypairDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *KeypairDetailRequest) (*KeypairDetailResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(keypairDetailRealRequest{
		RegionID:     req.RegionId,
		ProjectID:    req.ProjectId,
		KeypairName:  req.KeyPairName,
		QueryContent: req.QueryContent,
		PageNo:       req.PageNo,
		PageSize:     req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	result := &keypairDetailRealResponse{}
	err = response.ParseByStandardModelWithCheck(&result)

	var keypairDetail []KeypairDetailResults
	for _, detail := range result.Results {
		keypairDetail = append(keypairDetail, KeypairDetailResults{
			PublicKey:   detail.PublicKey,
			KeyPairName: detail.KeyPairName,
			FingerPrint: detail.FingerPrint,
			KeyPairId:   detail.KeyPairID,
			ProjectId:   detail.ProjectID,
		})
	}
	return &KeypairDetailResponse{
		CurrentCount: result.CurrentCount,
		TotalCount:   result.TotalCount,
		Results:      keypairDetail,
	}, err
}

type keypairDetailRealRequest struct {
	RegionID     string `json:"regionID"`
	ProjectID    string `json:"projectID,omitempty"`
	KeypairName  string `json:"keyPairName"`
	QueryContent string `json:"queryContent"`
	PageNo       int    `json:"pageNo"`
	PageSize     int    `json:"pageSize"`
}

type keypairDetailRealResponse struct {
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	Results      []struct {
		PublicKey   string `json:"publicKey"`
		KeyPairName string `json:"keyPairName"`
		FingerPrint string `json:"fingerPrint"`
		KeyPairID   string `json:"keyPairID"`
		ProjectID   string `json:"projectID,omitempty"`
	} `json:"results"`
}

type KeypairDetailResults struct {
	PublicKey   string
	KeyPairName string
	FingerPrint string
	KeyPairId   string
	ProjectId   string
}

type KeypairDetailRequest struct {
	RegionId     string // 资源池ID
	ProjectId    string // 企业项目id
	KeyPairName  string // 密钥对名称。只能由数字、字母、-组成,不能以数字和-开头、以-结尾,且长度为2-63字符
	QueryContent string
	PageNo       int
	PageSize     int
}

type KeypairDetailResponse struct {
	CurrentCount int // 当前页记录数目
	TotalCount   int // 总记录数
	Results      []KeypairDetailResults
}
