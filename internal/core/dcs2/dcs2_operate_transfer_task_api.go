package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2OperateTransferTaskApi
/* 结束运行中的迁移任务；删除已经完成（成功或者失败）的迁移任务记录
 */type Dcs2OperateTransferTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2OperateTransferTaskApi(client *core.CtyunClient) *Dcs2OperateTransferTaskApi {
	return &Dcs2OperateTransferTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/transfer/operateTransferTask",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2OperateTransferTaskApi) Do(ctx context.Context, credential core.Credential, req *Dcs2OperateTransferTaskRequest) (*Dcs2OperateTransferTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2OperateTransferTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2OperateTransferTaskRequest struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	TaskId      string `json:"taskId,omitempty"`      /*  任务ID<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15103&data=81&isNormal=1&vid=270">查询数据迁移任务列表</a>  */
	OperateType int32  `json:"operateType,omitempty"` /*  操作类型<li>2：结束运行中的任务<li>3：删除成功或者失败的任务记录  */
}

type Dcs2OperateTransferTaskResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2OperateTransferTaskReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2OperateTransferTaskReturnObjResponse struct{}
