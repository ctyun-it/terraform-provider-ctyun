package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2SpecExtendApi
/* 变更实例节点机器规格 */
type Mq2SpecExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2SpecExtendApi(client *core.CtyunClient) *Mq2SpecExtendApi {
	return &Mq2SpecExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/specExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2SpecExtendApi) Do(ctx context.Context, credential core.Credential, req *Mq2SpecExtendRequest) (*Mq2SpecExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2SpecExtendRequest
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
	var resp Mq2SpecExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2SpecExtendRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例id  */
	CpuNum     int32  `json:"cpuNum"`     /*  扩容后的cpu核数  */
	MemSize    int32  `json:"memSize"`    /*  扩容后的内存大小，单位G  */
}

type Mq2SpecExtendResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象  */
}
