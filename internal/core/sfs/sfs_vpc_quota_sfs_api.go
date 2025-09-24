package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsVpcQuotaSfsApi
/* 查询指定文件系统vpc配额的总量，已使用容量，剩余容量。
 */type SfsVpcQuotaSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsVpcQuotaSfsApi(client *core.CtyunClient) *SfsVpcQuotaSfsApi {
	return &SfsVpcQuotaSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/vpc-quota-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsVpcQuotaSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsVpcQuotaSfsRequest) (*SfsVpcQuotaSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsVpcQuotaSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsVpcQuotaSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一ID  */
}

type SfsVpcQuotaSfsResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                           `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                           `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsVpcQuotaSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                           `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsVpcQuotaSfsReturnObjResponse struct {
	TotalQuota  int32 `json:"totalQuota"`  /*  文件系统总体可用vpc配额  */
	UsedQuota   int32 `json:"usedQuota"`   /*  文件系统已用vpc配额个数  */
	RemainQuota int32 `json:"remainQuota"` /*  文件系统剩余可用vpc配额个数  */
}
