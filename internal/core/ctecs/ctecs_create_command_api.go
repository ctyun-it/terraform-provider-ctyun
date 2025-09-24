package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateCommandApi
/* 调用此接口可以创建一条云助手命令
 */type CtecsCreateCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateCommandApi(client *core.CtyunClient) *CtecsCreateCommandApi {
	return &CtecsCreateCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/create-command",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateCommandApi) Do(ctx context.Context, credential core.Credential, req *CtecsCreateCommandRequest) (*CtecsCreateCommandResponse, error) {
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
	var resp CtecsCreateCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateCommandRequest struct {
	RegionID         string                                       `json:"regionID,omitempty"`         /*  资源池ID  */
	CommandName      string                                       `json:"commandName,omitempty"`      /*  命令名称，长度不超过128个字符  */
	Description      string                                       `json:"description,omitempty"`      /*  命令描述，长度不超过512个字符  */
	CommandType      string                                       `json:"commandType,omitempty"`      /*  命令类型，取值范围：<br />Shell：适用于Linux云主机、物理机的Shell命令；<br />Bat：适用于Windows云主机的Bat命令；<br />PowerShell：适用于Windows云主机的PowerShell命令；<br />Python：适用于Python命令  */
	CommandContent   string                                       `json:"commandContent,omitempty"`   /*  加密后的命令内容，base64编码长度不可超过24KB  */
	WorkingDirectory string                                       `json:"workingDirectory,omitempty"` /*  命令在实例中运行目录。Linux系统默认路径为 /root;Windows系统默认路径为C:\Windows\System32 <br />说明：若在Windows系统云主机下执行Python脚本命令，需传Python安装全路径。  */
	Timeout          int32                                        `json:"timeout,omitempty"`          /*  命令超时时间，默认值60秒  */
	EnabledParameter *bool                                        `json:"enabledParameter"`           /*  是否启用自定义参数，若传true，则必须传defaultParameter，若enabledParameter为false，则defaultParameter不能传  */
	DefaultParameter []*CtecsCreateCommandDefaultParameterRequest `json:"defaultParameter"`           /*  启用自定义参数功能时，自定义参数的默认取值，json 格式string数组  */
}

type CtecsCreateCommandDefaultParameterRequest struct {
	Key         string `json:"key,omitempty"`         /*  参数名  */
	Description string `json:"description,omitempty"` /*  参数描述  */
	Value       string `json:"value,omitempty"`       /*  参数值  */
}

type CtecsCreateCommandResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                               `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                               `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsCreateCommandReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsCreateCommandReturnObjResponse struct {
	CommandID string `json:"commandID,omitempty"` /*  命令id  */
}
