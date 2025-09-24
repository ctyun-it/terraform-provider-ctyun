package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsGetCommandApi
/* 调用此接口可以查询用户手动创建的云助手命令或者云助手公共命令详情信息
 */type CtecsGetCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsGetCommandApi(client *core.CtyunClient) *CtecsGetCommandApi {
	return &CtecsGetCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/get-command",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsGetCommandApi) Do(ctx context.Context, credential core.Credential, req *CtecsGetCommandRequest) (*CtecsGetCommandResponse, error) {
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
	var resp CtecsGetCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsGetCommandRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID  */
	CommandID string `json:"commandID,omitempty"` /*  命令ID  */
}

type CtecsGetCommandResponse struct {
	StatusCode  int32                             `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                            `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                            `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsGetCommandReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsGetCommandReturnObjResponse struct {
	CommandID        string `json:"commandID,omitempty"`        /*  命令ID  */
	CommandName      string `json:"commandName,omitempty"`      /*  命令名称  */
	Description      string `json:"description,omitempty"`      /*  命令描述  */
	CommandType      string `json:"commandType,omitempty"`      /*  命令类型  */
	CommandContent   string `json:"commandContent,omitempty"`   /*  命令内容明文  */
	WorkingDirectory string `json:"workingDirectory,omitempty"` /*  命令在实例中的运行目录  */
	Timeout          int32  `json:"timeout,omitempty"`          /*  命令超时时间  */
	IsPublic         *bool  `json:"isPublic"`                   /*  是否是公共市场命令  */
	Version          string `json:"version,omitempty"`          /*  公共市场命令的版本，仅公共市场命令有该字段  */
	Owner            string `json:"owner,omitempty"`            /*  公共市场命令的提供者，仅公共市场命令有该字段  */
	EnabledParameter *bool  `json:"enabledParameter"`           /*  是否使能自定义参数  */
	DefaultParameter string `json:"defaultParameter,omitempty"` /*  自定义参数默认值  */
	CreateTime       string `json:"createTime,omitempty"`       /*  创建时间  */
	UpdateTime       string `json:"updateTime,omitempty"`       /*  更新时间  */
}
