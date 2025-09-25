package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsCreateSubdirSfsApi
/* 根据资源池ID 和 文件系统sfsUID ，在该文件系统下创建子目录
 */type SfsCreateSubdirSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsCreateSubdirSfsApi(client *core.CtyunClient) *SfsCreateSubdirSfsApi {
	return &SfsCreateSubdirSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/create-subdir-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsCreateSubdirSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsCreateSubdirSfsRequest) (*SfsCreateSubdirSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsCreateSubdirSfsRequest
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
	var resp SfsCreateSubdirSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsCreateSubdirSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
	SubDIR   string `json:"subDIR,omitempty"`   /*  文件系统子目录路径, 目录名只能以斜杠“/”开始，只允许包含字母、数字、下划线、连字符“-”和中文，字母区分大小写，每层目录的目录名长度要求1-63字符之间  */
}

type SfsCreateSubdirSfsResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                               `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                               `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsCreateSubdirSfsReturnObjResponse `json:"returnObj"`   /*  返回队形  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsCreateSubdirSfsReturnObjResponse struct {
	Subdir_name string `json:"subdir_name"` /*  已创建子目录路径列表  */
	SfsUID      string `json:"sfsUID"`      /*  弹性文件功能系统唯一 ID  */
}
