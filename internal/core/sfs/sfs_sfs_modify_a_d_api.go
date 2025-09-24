package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsModifyADApi
/* 文件系统修改AD域设置
 */type SfsSfsModifyADApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsModifyADApi(client *core.CtyunClient) *SfsSfsModifyADApi {
	return &SfsSfsModifyADApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/modify-ad",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsModifyADApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsModifyADRequest) (*SfsSfsModifyADResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsModifyADRequest
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
	var resp SfsSfsModifyADResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsModifyADRequest struct {
	RegionID       string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID         string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
	IsAnonymousAcc *bool  `json:"isAnonymousAcc"`     /*  是否允许匿名访问。true：允许匿名访问。false（默认）：不允许匿名访问。
	注：本参数和keytab信息至少传一项  */
	Keytab string `json:"keytab,omitempty"` /*  keytab 文件内容通过 base64 加密后的字符串。
	注：必须与keytabMd5同时传  */
	KeytabMd5 string `json:"keytabMd5,omitempty"` /*  keytab 文件内容通过 MD5 加密后的字符串。
	注：必须与keytab同时传  */
	KeytabName string `json:"keytabName,omitempty"` /*  keytab 文件名称  */
}

type SfsSfsModifyADResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
