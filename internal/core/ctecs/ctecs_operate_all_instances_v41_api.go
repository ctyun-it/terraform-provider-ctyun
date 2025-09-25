package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsOperateAllInstancesV41Api
/* 支持全部开机、关机、重启资源池云主机，返回执行开机/关机/重启操作的云主机ID及其对应任务ID<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsOperateAllInstancesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsOperateAllInstancesV41Api(client *core.CtyunClient) *CtecsOperateAllInstancesV41Api {
	return &CtecsOperateAllInstancesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/operate-all-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsOperateAllInstancesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsOperateAllInstancesV41Request) (*CtecsOperateAllInstancesV41Response, error) {
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
	var resp CtecsOperateAllInstancesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsOperateAllInstancesV41Request struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	OperateType string `json:"operateType,omitempty"` /*  操作类型，可选值：start（全部开机）、stop（全部关机）、reboot（全部重启）  */
}

type CtecsOperateAllInstancesV41Response struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息   */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsOperateAllInstancesV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsOperateAllInstancesV41ReturnObjResponse struct {
	JobIDList []*CtecsOperateAllInstancesV41ReturnObjJobIDListResponse `json:"jobIDList"` /*  任务列表  */
}

type CtecsOperateAllInstancesV41ReturnObjJobIDListResponse struct {
	JobID      string `json:"jobID,omitempty"`      /*  任务ID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9271&data=87">查询多个异步任务的结果</a>来查询操作是否成功  */
	InstanceID string `json:"instanceID,omitempty"` /*  对应任务云主机ID  */
}
