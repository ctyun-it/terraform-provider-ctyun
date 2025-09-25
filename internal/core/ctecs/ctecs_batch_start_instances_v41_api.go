package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsBatchStartInstancesV41Api
/* 该接口提供用户开启多台云主机功能<br/><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;多台操作：当前接口进行批量操作，开启一台云主机建议使用接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8296&data=87">开启一台云主机</a>进行操作<br/>&emsp;&emsp;异步接口：该接口为异步接口，请求过后会拿到任务ID（jobID），后续可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9271&data=87">查询多个异步任务的结果</a>来查询这些操作是否成功
 */type CtecsBatchStartInstancesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsBatchStartInstancesV41Api(client *core.CtyunClient) *CtecsBatchStartInstancesV41Api {
	return &CtecsBatchStartInstancesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/batch-start-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsBatchStartInstancesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsBatchStartInstancesV41Request) (*CtecsBatchStartInstancesV41Response, error) {
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
	var resp CtecsBatchStartInstancesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsBatchStartInstancesV41Request struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceIDList string `json:"instanceIDList,omitempty"` /*  云主机ID列表，多台使用英文逗号分割，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a><br />注：批量开启云主机的最大数量为50台  */
}

type CtecsBatchStartInstancesV41Response struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsBatchStartInstancesV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsBatchStartInstancesV41ReturnObjResponse struct {
	JobIDList []*CtecsBatchStartInstancesV41ReturnObjJobIDListResponse `json:"jobIDList"` /*  开机任务列表  */
}

type CtecsBatchStartInstancesV41ReturnObjJobIDListResponse struct {
	JobID      string `json:"jobID,omitempty"`      /*  开机任务ID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9271&data=87">查询多个异步任务的结果</a>来查询操作是否成功  */
	InstanceID string `json:"instanceID,omitempty"` /*  对应任务云主机ID  */
}
