package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsDeleteMountpointSfsApi
/* 根据文件系统sfsUID及挂载点ID ，删除文件系统指定挂载点
 */type SfsDeleteMountpointSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsDeleteMountpointSfsApi(client *core.CtyunClient) *SfsDeleteMountpointSfsApi {
	return &SfsDeleteMountpointSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/delete-mountpoint-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsDeleteMountpointSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsDeleteMountpointSfsRequest) (*SfsDeleteMountpointSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsDeleteMountpointSfsRequest
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
	var resp SfsDeleteMountpointSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsDeleteMountpointSfsRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池 ID  */
	SfsUID       string `json:"sfsUID,omitempty"`       /*  弹性文件功能系统唯一 ID  */
	MountPointID string `json:"mountPointID,omitempty"` /*  文件系统挂载点ID  */
}

type SfsDeleteMountpointSfsResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                   `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                   `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
	ReturnObj   *SfsDeleteMountpointSfsReturnObjResponse `json:"returnObj"`   /*  参考[响应示例]  */
}

type SfsDeleteMountpointSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  删除文件系统挂载点的操作号  */
}
