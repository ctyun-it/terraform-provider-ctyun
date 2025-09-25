package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsSetProtectSwitchApi
/* 为文件系统设置覆盖保护开关，若某文件系统要作为跨域复制的目标文件系统，则需要先禁用复制覆盖保护的开关
 */type SfsSfsSetProtectSwitchApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsSetProtectSwitchApi(client *core.CtyunClient) *SfsSfsSetProtectSwitchApi {
	return &SfsSfsSetProtectSwitchApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/set-protect-switch",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsSetProtectSwitchApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsSetProtectSwitchRequest) (*SfsSfsSetProtectSwitchResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsSetProtectSwitchRequest
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
	var resp SfsSfsSetProtectSwitchResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsSetProtectSwitchRequest struct {
	SfsUID        string `json:"sfsUID,omitempty"`   /*  文件系统唯一ID  */
	RegionID      string `json:"regionID,omitempty"` /*  文件系统所在资源池ID  */
	ProtectSwitch bool   `json:"protectSwitch"`      /*  是否关闭了目的文件系统覆盖保护。只有设置为False，才表示关闭覆盖保护，可以下发创建跨域复制  */
}

type SfsSfsSetProtectSwitchResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     string                                   `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                   `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
	ReturnObj   *SfsSfsSetProtectSwitchReturnObjResponse `json:"returnObj"`   /*  参考[结果码]  */
}

type SfsSfsSetProtectSwitchReturnObjResponse struct {
	Status string `json:"status"` /*  返回状态，OK为成功，其他为失败  */
}
