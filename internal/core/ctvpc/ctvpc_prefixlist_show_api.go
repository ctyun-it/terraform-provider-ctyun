package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcPrefixlistShowApi
/* 查询前缀列表详情
 */type CtvpcPrefixlistShowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrefixlistShowApi(client *core.CtyunClient) *CtvpcPrefixlistShowApi {
	return &CtvpcPrefixlistShowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/prefixlist/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrefixlistShowApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrefixlistShowRequest) (*CtvpcPrefixlistShowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("prefixListID", req.PrefixListID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcPrefixlistShowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrefixlistShowRequest struct {
	RegionID     string /*  资源池ID  */
	PrefixListID string /*  prefixlistID  */
}

type CtvpcPrefixlistShowResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcPrefixlistShowReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcPrefixlistShowReturnObjResponse struct {
	PrefixlistID *string `json:"prefixlistID,omitempty"` /*  prefixlist id  */
	Name         *string `json:"name,omitempty"`         /*  prefixlist 名称  */
	Limit        int32   `json:"limit"`                  /*  前缀列表支持的最大条目容量  */
	AddressType  *string `json:"addressType,omitempty"`  /*  地址类型，4：ipv4，6：ipv6  */
	Description  *string `json:"description,omitempty"`  /*  描述  */
	CreatedAt    *string `json:"createdAt,omitempty"`    /*  创建时间  */
	UpdatedAt    *string `json:"updatedAt,omitempty"`    /*  更新时间  */
}
