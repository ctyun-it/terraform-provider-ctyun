package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsModifyCommandApi
/* 调用此接口可以修改用户自己创建的云助手命令内容、命令参数等信息
 */type CtecsModifyCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsModifyCommandApi(client *core.CtyunClient) *CtecsModifyCommandApi {
	return &CtecsModifyCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/modify-command",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsModifyCommandApi) Do(ctx context.Context, credential core.Credential, req *CtecsModifyCommandRequest) (*CtecsModifyCommandResponse, error) {
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
	var resp CtecsModifyCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsModifyCommandRequest struct {
	RegionID         string                                       `json:"regionID,omitempty"`         /*  资源池ID  */
	CommandID        string                                       `json:"commandID,omitempty"`        /*  命令ID  */
	CommandName      string                                       `json:"commandName,omitempty"`      /*  命令名称，长度不超过128个字符  */
	Description      string                                       `json:"description,omitempty"`      /*  命令描述，长度不超过512个字符  */
	CommandType      string                                       `json:"commandType,omitempty"`      /*  命令类型，取值范围：<br />Shell：适用于Linux云主机、物理机的Shell命令；<br />Bat：适用于Windows云主机的Bat命令；<br />PowerShell：适用于Windows云主机的PowerShell命令；<br />Python：适用于Python命令  */
	CommandContent   string                                       `json:"commandContent,omitempty"`   /*  加密后的命令内容，base64编码长度不可超过24KB  */
	WorkingDirectory string                                       `json:"workingDirectory,omitempty"` /*  命令在云主机中运行目录  */
	Timeout          int32                                        `json:"timeout,omitempty"`          /*  命令超时时间  */
	EnabledParameter *bool                                        `json:"enabledParameter"`           /*  是否启用自定义参数，若传true，则必须传defaultParameter，若enabledParameter为false，则defaultParameter不能传  */
	DefaultParameter []*CtecsModifyCommandDefaultParameterRequest `json:"defaultParameter"`           /*  自定义参数使能时，修改自定义参数默认值  */
}

type CtecsModifyCommandDefaultParameterRequest struct {
	Key         string `json:"key,omitempty"`         /*  参数名  */
	Description string `json:"description,omitempty"` /*  参数描述  */
	Value       string `json:"value,omitempty"`       /*  参数值  */
}

type CtecsModifyCommandResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                               `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                               `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsModifyCommandReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsModifyCommandReturnObjResponse struct {
	CommandID string `json:"commandID,omitempty"` /*  命令ID  */
}
