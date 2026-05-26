package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2NodeExtendApi
/* 实例节点扩容 */
type Mq2NodeExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2NodeExtendApi(client *core.CtyunClient) *Mq2NodeExtendApi {
	return &Mq2NodeExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/nodeExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2NodeExtendApi) Do(ctx context.Context, credential core.Credential, req *Mq2NodeExtendRequest) (*Mq2NodeExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2NodeExtendRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2NodeExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2NodeExtendRequest struct {
	RegionId      string `json:"regionId"`      /*  资源池编码  */
	ProdInstId    string `json:"prodInstId"`    /*  实例id  */
	ExtendNodeNum int32  `json:"extendNodeNum"` /*  扩容后实例的节点数量，对应取值为等于代理数*2，范围为[4,32]  */
}

type Mq2NodeExtendResponse struct{}
