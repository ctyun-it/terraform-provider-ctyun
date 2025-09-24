package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsSendFileApi
/* 调用此接口可以上传文件到弹性云主机、物理机内部
 */ /* 说明：仅支持批量上传文件到弹性云主机或物理机内部，不支持混合上传，仅支持Linux系统
 */type CtecsSendFileApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsSendFileApi(client *core.CtyunClient) *CtecsSendFileApi {
	return &CtecsSendFileApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/send-file",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsSendFileApi) Do(ctx context.Context, credential core.Credential, req *CtecsSendFileRequest) (*CtecsSendFileResponse, error) {
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
	var resp CtecsSendFileResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsSendFileRequest struct {
	RegionID        string `json:"regionID,omitempty"`        /*  资源池ID  */
	InstanceIDs     string `json:"instanceIDs,omitempty"`     /*  待下发文件的弹性云主机、物理机ID列表, 使用 , 分割  */
	FileName        string `json:"fileName,omitempty"`        /*  文件名称，长度不超过128个字符  */
	Description     string `json:"description,omitempty"`     /*  文件描述，长度不超过512个字符  */
	FileContent     string `json:"fileContent,omitempty"`     /*  加密的文件内容，base64编码长度不可超过24KB  */
	TargetDirectory string `json:"targetDirectory,omitempty"` /*  下发文件的目标路径  */
	FileOwner       string `json:"fileOwner,omitempty"`       /*  文件所属用户，只针对linux实例，默认root  */
	FileGroup       string `json:"fileGroup,omitempty"`       /*  文件用户组，只针对linux实例，默认root  */
	FileMode        string `json:"fileMode,omitempty"`        /*  文件权限，只针对linux实例，默认644  */
	Overwrite       *bool  `json:"overwrite"`                 /*  是否覆盖，如果目标路径下同名文件已经存在，true：覆盖，false：不覆盖。默认false  */
}

type CtecsSendFileResponse struct {
	StatusCode  int32                           `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsSendFileReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsSendFileReturnObjResponse struct {
	InvokedID string `json:"invokedID,omitempty"` /*  执行ID  */
}
