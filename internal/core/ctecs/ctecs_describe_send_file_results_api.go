package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDescribeSendFileResultsApi
/* 调用此接口可以查询上传到弹性云主机或物理机的文件的结果
 */type CtecsDescribeSendFileResultsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDescribeSendFileResultsApi(client *core.CtyunClient) *CtecsDescribeSendFileResultsApi {
	return &CtecsDescribeSendFileResultsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/describe-send-file-results",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDescribeSendFileResultsApi) Do(ctx context.Context, credential core.Credential, req *CtecsDescribeSendFileResultsRequest) (*CtecsDescribeSendFileResultsResponse, error) {
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
	var resp CtecsDescribeSendFileResultsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDescribeSendFileResultsRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID  */
	FileName  string `json:"fileName,omitempty"`  /*  文件名称  */
	InvokedID string `json:"invokedID,omitempty"` /*  执行ID  */
	PageNo    int32  `json:"pageNo,omitempty"`    /*  当前页码，默认值为1  */
	PageSize  int32  `json:"pageSize,omitempty"`  /*  分页查询时设置的每页行数，最大值为100，默认为10  */
}

type CtecsDescribeSendFileResultsResponse struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明。  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述。  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述。  */
	ReturnObj   *CtecsDescribeSendFileResultsReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据。  */
}

type CtecsDescribeSendFileResultsReturnObjResponse struct {
	Results    *CtecsDescribeSendFileResultsReturnObjResultsResponse `json:"results"`              /*  命令执行结果集合  */
	PageNo     int32                                                 `json:"pageNo,omitempty"`     /*  当前页码  */
	TotalCount int32                                                 `json:"totalCount,omitempty"` /*  命令总个数  */
	PageSize   int32                                                 `json:"pageSize,omitempty"`   /*  每页行数  */
}

type CtecsDescribeSendFileResultsReturnObjResultsResponse struct {
	CreateTime             string                                                               `json:"createTime,omitempty"`             /*  文件下发任务创建时间  */
	InvocationRecordStatus string                                                               `json:"invocationRecordStatus,omitempty"` /*  文件下发总的执行状态，取值范围：<br />Pending：下发中，有一台云主机未下发，则总的状态为下发中；<br />Running：运行中，有云主机中文件下发进程为运行中，则总的执行状态为运行中；<br />Finished：已完成，所有云主机文件下发全部完成执行；<br />Failed：执行失败，所有云主机文件下发全部执行失败  */
	InvokedID              string                                                               `json:"invokedID,omitempty"`              /*  命令执行ID  */
	FileName               string                                                               `json:"fileName,omitempty"`               /*  文件名称，长度不超过128个字符  */
	Description            string                                                               `json:"description,omitempty"`            /*  文件描述，长度不超过512个字符  */
	FileContent            string                                                               `json:"fileContent,omitempty"`            /*  文件内容  */
	TargetDirectory        string                                                               `json:"targetDirectory,omitempty"`        /*  下发文件的目标路径  */
	Timeout                int32                                                                `json:"timeout,omitempty"`                /*  下发文件的超时时间。默认值60秒  */
	FileOwner              string                                                               `json:"fileOwner,omitempty"`              /*  文件所属用户，只针对linux实例，默认root  */
	FileGroup              string                                                               `json:"fileGroup,omitempty"`              /*  文件用户组，只针对linux实例，默认root  */
	FileMode               string                                                               `json:"fileMode,omitempty"`               /*  文件权限，只针对linux实例，默认0644  */
	Overwrite              *bool                                                                `json:"overwrite"`                        /*  是否覆盖，如果目标路径下同名文件已经存在，true：覆盖；false：不覆盖。默认false  */
	VmCount                int32                                                                `json:"vmCount,omitempty"`                /*  文件下发的实例数量  */
	InvokeInstances        *CtecsDescribeSendFileResultsReturnObjResultsInvokeInstancesResponse `json:"invokeInstances"`                  /*  实例列表  */
}

type CtecsDescribeSendFileResultsReturnObjResultsInvokeInstancesResponse struct {
	CreateTime       string `json:"createTime,omitempty"`       /*  文件下发任务创建时间  */
	UpdateTime       string `json:"updateTime,omitempty"`       /*  文件下发任务完成时间  */
	InvocationStatus string `json:"invocationStatus,omitempty"` /*  单台云主机的文件下发执行状态，可能值：<br />Pending：系统正在下发文件；<br />Running：文件下发中；<br />Success：文件下发成功；<br />Failed：文件下发失败；  */
	ExitCode         string `json:"exitCode,omitempty"`         /*  命令退出码  */
	ErrorInfo        string `json:"errorInfo,omitempty"`        /*  错误信息  */
	InstanceID       string `json:"instanceID,omitempty"`       /*  实例ID  */
}
