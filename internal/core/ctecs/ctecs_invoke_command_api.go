package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsInvokeCommandApi
/* 调用此接口为一台或多台弹性云主机或物理机触发一条云助手命令
 */ /* 说明：仅支持批量为弹性云主机或物理机触发云助手命令，不支持混合触发
 */type CtecsInvokeCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsInvokeCommandApi(client *core.CtyunClient) *CtecsInvokeCommandApi {
	return &CtecsInvokeCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/invoke-command",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsInvokeCommandApi) Do(ctx context.Context, credential core.Credential, req *CtecsInvokeCommandRequest) (*CtecsInvokeCommandResponse, error) {
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
	var resp CtecsInvokeCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsInvokeCommandRequest struct {
	RegionID         string                              `json:"regionID,omitempty"`         /*  资源池ID  */
	InstanceIDs      string                              `json:"instanceIDs,omitempty"`      /*  待执行命令的云主机、物理机ID列表, 使用英文, 分割  */
	CommandID        string                              `json:"commandID,omitempty"`        /*  命令ID  */
	Timeout          int32                               `json:"timeout,omitempty"`          /*  执行命令的超时时间  */
	WorkingDirectory string                              `json:"workingDirectory,omitempty"` /*  命令在云主机中运行目录。说明：若在Windows系统云主机下执行Python命令，需传Python安装全路径  */
	Parameter        *CtecsInvokeCommandParameterRequest `json:"parameter"`                  /*  自定义参数，说明：key仅支持大小写字母(A-a)、数字(0-9)、横线(-)和下划线(_)，key和value均只支持string  */
}

type CtecsInvokeCommandParameterRequest struct{}

type CtecsInvokeCommandResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                               `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                               `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsInvokeCommandReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsInvokeCommandReturnObjResponse struct {
	InvokedID string `json:"invokedID,omitempty"` /*  命令执行ID  */
}
