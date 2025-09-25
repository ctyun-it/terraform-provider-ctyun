package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpDateSubnetApi
/* 修改子网的属性：名称、描述。
 */type CtvpcUpDateSubnetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpDateSubnetApi(client *core.CtyunClient) *CtvpcUpDateSubnetApi {
	return &CtvpcUpDateSubnetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/update-subnet",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpDateSubnetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpDateSubnetRequest) (*CtvpcUpDateSubnetResponse, error) {
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
	var resp CtvpcUpDateSubnetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpDateSubnetRequest struct {
	RegionID string    `json:"regionID,omitempty"` /*  资源池 ID  */
	SubnetID string    `json:"subnetID,omitempty"` /*  子网 的 ID  */
	Name     *string   `json:"name,omitempty"`     /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	DnsList  []*string `json:"dnsList"`            /*  子网 dns 列表, 最多同时支持 4 个 dns 地址  */
}

type CtvpcUpDateSubnetResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
