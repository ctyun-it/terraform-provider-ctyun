package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCloudAssistantDescribeInvocationResultsApi
/* 查询一条或多条云助手命令在弹性云主机、物理机中执行结果
 */type CtecsCloudAssistantDescribeInvocationResultsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCloudAssistantDescribeInvocationResultsApi(client *core.CtyunClient) *CtecsCloudAssistantDescribeInvocationResultsApi {
	return &CtecsCloudAssistantDescribeInvocationResultsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/describe-invocation-results",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCloudAssistantDescribeInvocationResultsApi) Do(ctx context.Context, credential core.Credential, req *CtecsCloudAssistantDescribeInvocationResultsRequest) (*CtecsCloudAssistantDescribeInvocationResultsResponse, error) {
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
	var resp CtecsCloudAssistantDescribeInvocationResultsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCloudAssistantDescribeInvocationResultsRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID  */
	CommandID string `json:"commandID,omitempty"` /*  命令ID  */
	InvokedID string `json:"invokedID,omitempty"` /*  命令执行ID  */
	PageNo    int32  `json:"pageNo,omitempty"`    /*  当前页码，默认值为1  */
	PageSize  int32  `json:"pageSize,omitempty"`  /*  分页查询时设置的每页行数，最大值为100，默认为10  */
}

type CtecsCloudAssistantDescribeInvocationResultsResponse struct {
	StatusCode  int32                                                          `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                                                         `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsCloudAssistantDescribeInvocationResultsReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsCloudAssistantDescribeInvocationResultsReturnObjResponse struct {
	Results    *CtecsCloudAssistantDescribeInvocationResultsReturnObjResultsResponse `json:"results"`              /*  命令执行结果集合  */
	PageNo     int32                                                                 `json:"pageNo,omitempty"`     /*  当前页码  */
	TotalCount int32                                                                 `json:"totalCount,omitempty"` /*  命令总个数  */
	PageSize   int32                                                                 `json:"pageSize,omitempty"`   /*  每页行数  */
}

type CtecsCloudAssistantDescribeInvocationResultsReturnObjResultsResponse struct {
	InvocationStatus   string `json:"invocationStatus,omitempty"`   /*  单台云主机的命令运行状态，可能值：<br />Pending：系统正在校验或发送命令；<br />Running：命令正在云主机上运行；<br />Success：命令执行完成，且退出码为0；<br />Failed：命令执行完成，且退出码非0；  */
	CreateTime         string `json:"createTime,omitempty"`         /*  命令执行创建时间  */
	UpdateTime         string `json:"updateTime,omitempty"`         /*  命令执行完成时间  */
	InvokedID          string `json:"invokedID,omitempty"`          /*  命令执行ID  */
	CommandID          string `json:"commandID,omitempty"`          /*  命令ID  */
	InstanceID         string `json:"instanceID,omitempty"`         /*  云主机ID  */
	Output             string `json:"output,omitempty"`             /*  命令执行后的输出信息  */
	ExitCode           int32  `json:"exitCode,omitempty"`           /*  命令退出码  */
	ErrorInfo          string `json:"errorInfo,omitempty"`          /*  命令执行失败原因详情  */
	InvokeRecordStatus string `json:"invokeRecordStatus,omitempty"` /*  单个命令执行任务的总状态，取值范围：<br />Pending：未执行，当有云主机中命令状态为Pending，则总的执行状态为未执行；<br />Running：运行中，有云主机中命令进程为运行中，则总的执行状态为运行中；<br />Finished：已完成。所有云主机命令进程全部完成执行；<br />Failed：执行失败，有云主机中命令进程为执行失败，则总的状态为Failed。  */
}
