package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsRenameSFSApi
/* 指定文件系统重命名
 */type SfsRenameSFSApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsRenameSFSApi(client *core.CtyunClient) *SfsRenameSFSApi {
	return &SfsRenameSFSApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/rename-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsRenameSFSApi) Do(ctx context.Context, credential core.Credential, req *SfsRenameSFSRequest) (*SfsRenameSFSResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsRenameSFSRequest
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
	var resp SfsRenameSFSResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsRenameSFSRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一ID  */
	SfsName  string `json:"sfsName,omitempty"`  /*  文件系统新名称  */
}

type SfsRenameSFSResponse struct {
	StatusCode  int32                          `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                         `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                         `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsRenameSFSReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                         `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                         `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsRenameSFSReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  重命名文件系统的操作ID  */
}
