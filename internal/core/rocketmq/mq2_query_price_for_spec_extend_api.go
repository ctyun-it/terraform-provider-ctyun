package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryPriceForSpecExtendApi
/* 规格变更查价 */
type Mq2QueryPriceForSpecExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryPriceForSpecExtendApi(client *core.CtyunClient) *Mq2QueryPriceForSpecExtendApi {
	return &Mq2QueryPriceForSpecExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/queryPriceForSpecExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2QueryPriceForSpecExtendApi) Do(ctx context.Context, credential core.Credential, req *Mq2QueryPriceForSpecExtendRequest) (*Mq2QueryPriceForSpecExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2QueryPriceForSpecExtendRequest
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
	var resp Mq2QueryPriceForSpecExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryPriceForSpecExtendRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	CpuNum     int32  `json:"cpuNum"`     /*  扩容后的cpu核数  */
	MemSize    int32  `json:"memSize"`    /*  扩容后的内存大小，单位G  */
}

type Mq2QueryPriceForSpecExtendResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}
