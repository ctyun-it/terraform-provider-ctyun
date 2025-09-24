package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdatedhcpoptionsetsApi
/* 更新dhcpoptionsets
 */type CtvpcUpdatedhcpoptionsetsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdatedhcpoptionsetsApi(client *core.CtyunClient) *CtvpcUpdatedhcpoptionsetsApi {
	return &CtvpcUpdatedhcpoptionsetsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/dhcpoptionsets/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdatedhcpoptionsetsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdatedhcpoptionsetsRequest) (*CtvpcUpdatedhcpoptionsetsResponse, error) {
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
	var resp CtvpcUpdatedhcpoptionsetsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdatedhcpoptionsetsRequest struct {
	RegionID         string    `json:"regionID,omitempty"`         /*  资源池 ID  */
	DhcpOptionSetsID string    `json:"dhcpOptionSetsID,omitempty"` /*  集合ID  */
	Name             *string   `json:"name,omitempty"`             /*  集合名，支持拉丁字母、中文、数字，下划线，连字符，必须以中文 / 英文字母开头，不能以数字、_和-、 http: / https: 开头，长度 2 - 32  */
	Description      *string   `json:"description,omitempty"`      /*  描述  */
	DomainName       *string   `json:"domainName,omitempty"`       /*  	整个域名的总长度不能超过 255 个字符，每个子域名（包括顶级域名）的长度不能超过 63 个字符，域名中的字符集包括大写字母、小写字母、数字和连字符（减号），连字符不能位于域名的开头  */
	DnsList          []*string `json:"dnsList"`                    /*  服务ip地址列表，最多只能4个IP地址  */
}

type CtvpcUpdatedhcpoptionsetsResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcUpdatedhcpoptionsetsReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcUpdatedhcpoptionsetsReturnObjResponse struct {
	DhcpOptionSetsID *string `json:"dhcpOptionSetsID,omitempty"` /*  ID  */
}
