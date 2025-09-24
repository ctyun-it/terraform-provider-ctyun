package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsBatchUpdateInstancesPasswordV41Api
/* 该接口提供用户更新多台云主机的密码的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;如果使用私有镜像创建的云主机执行该操作时，请先检查云主机内部是否安装了QGA（qemu-guest-agent）。不同操作系统请参考：<a href="https://www.ctyun.cn/document/10027726/10747194">Windows系统盘镜像文件安装QGA</a>和<a href="https://www.ctyun.cn/document/10027726/10747147">Linux系统盘镜像文件安装QGA</a><br />
 */type CtecsBatchUpdateInstancesPasswordV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsBatchUpdateInstancesPasswordV41Api(client *core.CtyunClient) *CtecsBatchUpdateInstancesPasswordV41Api {
	return &CtecsBatchUpdateInstancesPasswordV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/batch-reset-password",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsBatchUpdateInstancesPasswordV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsBatchUpdateInstancesPasswordV41Request) (*CtecsBatchUpdateInstancesPasswordV41Response, error) {
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
	var resp CtecsBatchUpdateInstancesPasswordV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsBatchUpdateInstancesPasswordV41Request struct {
	RegionID      string                                                      `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	UpdatePwdInfo []*CtecsBatchUpdateInstancesPasswordV41UpdatePwdInfoRequest `json:"updatePwdInfo"`      /*  批量更新密码信息列表  */
}

type CtecsBatchUpdateInstancesPasswordV41UpdatePwdInfoRequest struct {
	InstanceID  string `json:"instanceID,omitempty"`  /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	NewPassword string `json:"newPassword,omitempty"` /*  新的用户密码，满足以下规则：<br />长度在8～30个字符；<br />必须包含大写字母、小写字母、数字以及特殊符号中的三项；<br />特殊符号可选：()`~!@#$%^&*_-+=｜{}[]:;'<>,.?/且不能以斜线号 / 开头；<br />不能包含3个及以上连续字符；<br />Linux镜像不能包含镜像用户名（root）、用户名的倒序（toor）、用户名大小写变化（如RoOt、rOot等）；<br />Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator等）  */
}

type CtecsBatchUpdateInstancesPasswordV41Response struct {
	StatusCode  int32                                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                                 `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码部分  */
	Error       string                                                 `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码，详见错误码部分  */
	Message     string                                                 `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                                 `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsBatchUpdateInstancesPasswordV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsBatchUpdateInstancesPasswordV41ReturnObjResponse struct {
	InstanceIDList string `json:"instanceIDList,omitempty"` /*  被更新密码的云主机ID列表，使用英文逗号分割，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}
