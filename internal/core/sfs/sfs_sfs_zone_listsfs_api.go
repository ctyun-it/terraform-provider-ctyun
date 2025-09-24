package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsZoneListsfsApi
/* 文件系统的区域（zone）查询
 */type SfsSfsZoneListsfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsZoneListsfsApi(client *core.CtyunClient) *SfsSfsZoneListsfsApi {
	return &SfsSfsZoneListsfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/zonelist",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsZoneListsfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsZoneListsfsRequest) (*SfsSfsZoneListsfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsZoneListsfsRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsZoneListsfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsZoneListsfsRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  可用区所在的地域 ID  */
	FileSystemType string `json:"fileSystemType,omitempty"` /*  文件系统类型。取值范围：performance-性能型、capacity-标准型、hdd_e-标准型专属  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  分页查询时每页包含的地域数  */
	PageNo         int32  `json:"pageNo,omitempty"`         /*  列表的分页页码  */
}

type SfsSfsZoneListsfsResponse struct {
	StatusCode  string `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string `json:"description"` /*  失败时的错误描述，一般为英文描述  */
	ReturnObj   string `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码。  */
}
