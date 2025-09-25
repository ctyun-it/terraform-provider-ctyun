package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsExecuteSnapshotPolicyV41Api
/* 该接口提供用户立即执行云主机快照策略的功能，通过该接口可以立即对当前快照策略绑定的云主机执行快照<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsExecuteSnapshotPolicyV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsExecuteSnapshotPolicyV41Api(client *core.CtyunClient) *CtecsExecuteSnapshotPolicyV41Api {
	return &CtecsExecuteSnapshotPolicyV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot-policy/execute",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsExecuteSnapshotPolicyV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsExecuteSnapshotPolicyV41Request) (*CtecsExecuteSnapshotPolicyV41Response, error) {
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
	var resp CtecsExecuteSnapshotPolicyV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsExecuteSnapshotPolicyV41Request struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty"` /*  云主机快照策略ID，32字节<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9600&data=87">查询云主机快照策略列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9588&data=87">创建云主机快照策略</a>  */
}

type CtecsExecuteSnapshotPolicyV41Response struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                          `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsExecuteSnapshotPolicyV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsExecuteSnapshotPolicyV41ReturnObjResponse struct {
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty"` /*  执行的云主机快照策略ID  */
}
