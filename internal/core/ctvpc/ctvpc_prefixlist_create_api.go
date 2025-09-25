package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcPrefixlistCreateApi
/* 创建 prefixlist
 */type CtvpcPrefixlistCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrefixlistCreateApi(client *core.CtyunClient) *CtvpcPrefixlistCreateApi {
	return &CtvpcPrefixlistCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/prefixlist/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrefixlistCreateApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrefixlistCreateRequest) (*CtvpcPrefixlistCreateResponse, error) {
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
	var resp CtvpcPrefixlistCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrefixlistCreateRequest struct {
	RegionID        string                                         `json:"regionID,omitempty"` /*  资源池ID  */
	Name            string                                         `json:"name,omitempty"`     /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Limit           int32                                          `json:"limit"`              /*  前缀列表支持的最大条目容量，创建后将无法修改,限制1-200条，具体以账户配额为准,不能小于前缀列表规则个数  */
	AddressType     int32                                          `json:"addressType"`        /*  地址类型，4：ipv4，6：ipv6  */
	PrefixListRules []*CtvpcPrefixlistCreatePrefixListRulesRequest `json:"prefixListRules"`    /*  接口业务数据  */
}

type CtvpcPrefixlistCreatePrefixListRulesRequest struct {
	Cidr string `json:"cidr,omitempty"` /*  前缀列表条目,cidr  */
}

type CtvpcPrefixlistCreateResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcPrefixlistCreateReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcPrefixlistCreateReturnObjResponse struct {
	PrefixListID *string `json:"prefixListID,omitempty"` /*  prefixlist id  */
}
