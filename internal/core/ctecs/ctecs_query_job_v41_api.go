package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryJobV41Api
/* 该接口通过一个或多个异步任务的jobID查询任务执行的结果<br/><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br/>&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br/><b>注意事项：</b><br/>&emsp;&emsp;异步任务结果查询：请先通过异步接口得到对应的任务ID（jobID），注：异步任务查询不支持查询订单接口（即涉及付费的接口，如创建云主机）<br/>&emsp;&emsp;多个任务查询：当前接口可以查询同一资源池内多个任务结果，查询单个任务结果，推荐使用接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5543&data=87">查询一个异步任务的结果</a>来查询
 */type CtecsQueryJobV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryJobV41Api(client *core.CtyunClient) *CtecsQueryJobV41Api {
	return &CtecsQueryJobV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/job/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryJobV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryJobV41Request) (*CtecsQueryJobV41Response, error) {
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
	var resp CtecsQueryJobV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryJobV41Request struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	JobIDs   string `json:"jobIDs,omitempty"`   /*  异步任务ID列表，以英文逗号分隔每个ID  */
}

type CtecsQueryJobV41Response struct {
	StatusCode  int32                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败），默认值：800  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                             `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                             `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                             `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQueryJobV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsQueryJobV41ReturnObjResponse struct {
	JobList []*CtecsQueryJobV41ReturnObjJobListResponse `json:"jobList"` /*  异步任务列表  */
}

type CtecsQueryJobV41ReturnObjJobListResponse struct {
	JobID     string `json:"jobID,omitempty"`     /*  异步任务ID  */
	JobStatus int32  `json:"jobStatus,omitempty"` /*  任务执行状态，取值范围：<br />0：执行中,<br />1：执行成功,<br />2：执行失败  */
}
