package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosDeletePolicy1Api
/* 删除策略，若有角色绑定了该策略，则会自动与该角色解绑
 */type ZosDeletePolicy1Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosDeletePolicy1Api(client *core.CtyunClient) *ZosDeletePolicy1Api {
	return &ZosDeletePolicy1Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/delete-policy",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosDeletePolicy1Api) Do(ctx context.Context, credential core.Credential, req *ZosDeletePolicy1Request) (*ZosDeletePolicy1Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosDeletePolicy1Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosDeletePolicy1Request struct {
	RegionID   string `json:"regionID,omitempty"`   /*  区域 ID  */
	PolicyName string `json:"policyName,omitempty"` /*  策略名  */
}

type ZosDeletePolicy1Response struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
