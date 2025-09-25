package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwUpdateApi
/* 修改l2gw的属性：名称、描述。
 */type CtvpcL2gwUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwUpdateApi(client *core.CtyunClient) *CtvpcL2gwUpdateApi {
	return &CtvpcL2gwUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwUpdateApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwUpdateRequest) (*CtvpcL2gwUpdateResponse, error) {
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
	var resp CtvpcL2gwUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwUpdateRequest struct {
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池 ID  */
	L2gwID      string  `json:"l2gwID,omitempty"`      /*  l2gw 的 ID  */
	Name        string  `json:"name,omitempty"`        /*  支持拉丁字母、中文、数字，下划线，连字符，必须以中文 / 英文字母开头，不能以数字、_和-、 http: / https: 开头，长度 2 - 32  */
	Description *string `json:"description,omitempty"` /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:"{},./;'[\]·~！@#￥%……&*（） —— -+={}\《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtvpcL2gwUpdateResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
