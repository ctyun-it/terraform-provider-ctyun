package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicUpdateApi
/* 配置Topic的读写模式 */
type Mq2TopicUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicUpdateApi(client *core.CtyunClient) *Mq2TopicUpdateApi {
	return &Mq2TopicUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/update",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2TopicUpdateApi) Do(ctx context.Context, credential core.Credential, req *Mq2TopicUpdateRequest) (*Mq2TopicUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2TopicUpdateRequest
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
	var resp Mq2TopicUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicUpdateRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	Topic      string `json:"topic"`      /*  主题名称  */
	Perm       int32  `json:"perm"`       /*  设置该Topic的读写模式。取值说明如下：6：同时支持读写;4：禁写; 2：禁读  */
	Remark     string `json:"remark"`     /*  设置该Topic的读写模式。取值说明如下：6：同时支持读写;4：禁写; 2：禁读  */
}

type Mq2TopicUpdateResponse struct {
	StatusCode int32                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                          `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2TopicUpdateReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                          `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2TopicUpdateReturnObjResponse struct{}
