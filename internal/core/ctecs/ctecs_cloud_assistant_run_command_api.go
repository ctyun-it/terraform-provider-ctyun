package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCloudAssistantRunCommandApi
/* 调用此接口在一台或多台弹性云主机或物理机中执行一段Shell、PowerShell、Bat或Python类型的脚本命令
 */type CtecsCloudAssistantRunCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCloudAssistantRunCommandApi(client *core.CtyunClient) *CtecsCloudAssistantRunCommandApi {
	return &CtecsCloudAssistantRunCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/run-command",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCloudAssistantRunCommandApi) Do(ctx context.Context, credential core.Credential, req *CtecsCloudAssistantRunCommandRequest) (*CtecsCloudAssistantRunCommandResponse, error) {
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
	var resp CtecsCloudAssistantRunCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCloudAssistantRunCommandRequest struct {
	RegionID         string                                                  `json:"regionID,omitempty"`         /*  资源池ID  */
	InstanceIDs      string                                                  `json:"instanceIDs,omitempty"`      /*  待执行命令的弹性云主机、物理机ID列表, 使用英文 , 分割（当前仅支持同时下发弹性云主机ID或同时下发物理机ID，不支持混合下发）  */
	CommandName      string                                                  `json:"commandName,omitempty"`      /*  命令名称，长度不超过128个字符  */
	Description      string                                                  `json:"description,omitempty"`      /*  命令描述，长度不超过512个字符  */
	CommandType      string                                                  `json:"commandType,omitempty"`      /*  命令类型，取值范围：<br />Shell：适用于Linux云主机、物理机的Shell命令；<br />Bat：适用于Windows云主机的Bat命令；<br />PowerShell：适用于Windows云主机的PowerShell命令；<br />Python：适用于Python命令。<br />说明：当前物理机云助手还不支持windows系统。  */
	CommandContent   string                                                  `json:"commandContent,omitempty"`   /*  加密后的命令内容，长度不可超过24KB  */
	WorkingDirectory string                                                  `json:"workingDirectory,omitempty"` /*  命令在实例中运行目录：<br />Linux系统默认路径为 /tmp;<br />Windows系统默认路径为 C:\Windows\System32。<br />说明：若在Windows系统云主机下执行Python命令，需要传Python安装全路径。  */
	Timeout          int32                                                   `json:"timeout,omitempty"`          /*  命令超时时间。默认值60秒  */
	SaveCommand      *bool                                                   `json:"saveCommand"`                /*  是否保存命令，默认false  */
	EnabledParameter *bool                                                   `json:"enabledParameter"`           /*  是否启用自定义参数，默认值为false，若传true，则必须传defaultParameter，若enabledParameter为false，则defaultParameter和parameter都不能传  */
	DefaultParameter []*CtecsCloudAssistantRunCommandDefaultParameterRequest `json:"defaultParameter"`           /*  启用自定义参数功能时，自定义参数的默认取值，json 格式map数组，说明：key仅支持大小写字母(A-a)、数字(0-9)、横线(-)和下划线(_)  */
	Parameter        *CtecsCloudAssistantRunCommandParameterRequest          `json:"parameter"`                  /*  自定义参数，说明：key仅支持大小写字母(A-a)、数字(0-9)、横线(-)和下划线(_)，key和value均只支持string  */
}

type CtecsCloudAssistantRunCommandDefaultParameterRequest struct {
	Key         string `json:"key,omitempty"`         /*  参数名  */
	Description string `json:"description,omitempty"` /*  参数描述  */
	Value       string `json:"value,omitempty"`       /*  参数值  */
}

type CtecsCloudAssistantRunCommandParameterRequest struct{}

type CtecsCloudAssistantRunCommandResponse struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsCloudAssistantRunCommandReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsCloudAssistantRunCommandReturnObjResponse struct {
	CommandID string `json:"commandID,omitempty"` /*  命令ID
	说明：saveCommand参数传false时不会生成commandID  */
	InvokedID string `json:"invokedID,omitempty"` /*  命令执行ID  */
}
