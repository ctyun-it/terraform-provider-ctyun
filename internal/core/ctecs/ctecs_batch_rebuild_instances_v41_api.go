package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsBatchRebuildInstancesV41Api
/* &emsp;&emsp;该接口提供用户重装多台云主机功能，通过填写相应云主机ID、镜像ID对云主机进行重装<br/><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;多台操作：当前接口进行批量操作，重装一台云主机建议使用接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8302&data=87">重装一台云主机</a>进行操作<br/>&emsp;&emsp;异步接口：该接口为异步接口，请求过后会拿到任务ID（jobID），后续可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9271&data=87">查询多个异步任务的结果</a>来查询这些操作是否成功<br />&emsp;&emsp;监控安装：在云服务器创建成功后，3-5分钟内将完成详细监控Agent安装，即开启云服务器CPU，内存，网络，磁盘，进程等指标详细监控，若不开启，则无任何监控数据
 */type CtecsBatchRebuildInstancesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsBatchRebuildInstancesV41Api(client *core.CtyunClient) *CtecsBatchRebuildInstancesV41Api {
	return &CtecsBatchRebuildInstancesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/batch-rebuild-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsBatchRebuildInstancesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsBatchRebuildInstancesV41Request) (*CtecsBatchRebuildInstancesV41Response, error) {
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
	var resp CtecsBatchRebuildInstancesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsBatchRebuildInstancesV41Request struct {
	RegionID    string                                             `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	RebuildInfo []*CtecsBatchRebuildInstancesV41RebuildInfoRequest `json:"rebuildInfo"`        /*  重装信息列表  */
}

type CtecsBatchRebuildInstancesV41RebuildInfoRequest struct {
	InstanceID     string `json:"instanceID,omitempty"`   /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	Password       string `json:"password,omitempty"`     /*  用户密码，满足以下规则：<br />长度在8-30个字符；<br />必须包含大写字母、小写字母、数字以及特殊符号中的三项；<br />特殊符号可选：()`~!@#$%^&*_-+=｜{}[]:;'<>,.?/\且不能以斜线号 / 开头；<br />不能包含3个及以上连续字符；<br />Linux镜像不能包含镜像用户名（root）、用户名的倒序（toor）、用户名大小写变化（如RoOt、rOot等）；<br />Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator等）  */
	ImageID        string `json:"imageID,omitempty"`      /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10030151">镜像概述</a>来了解云主机镜像<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89">查询可以使用的镜像资源</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4765&data=89">创建私有镜像（云主机系统盘）</a><br />注：不填默认以原镜像进行重装  */
	UserData       string `json:"userData,omitempty"`     /*  用户自定义数据，需要以Base64方式编码，Base64编码后的长度限制为1-16384字符  */
	InstanceName   string `json:"instanceName,omitempty"` /*  云主机名称。不同操作系统下，云主机名称规则有差异<br />Windows：长度为2-15个字符，允许使用大小写字母、数字或连字符（-）。不能以连字符（-）开头或结尾，不能连续使用连字符（-），也不能仅使用数字；<br />其他操作系统：长度为2-64字符，允许使用点（.）分隔字符成多段，每段允许使用大小写字母、数字或连字符（-），但不能连续使用点号（.）或连字符（-），不能以点号（.）或连字符（-）开头或结尾，也不能仅使用数字<br />注：如果不填，默认值为原来云主机名称  */
	MonitorService *bool  `json:"monitorService"`         /*  监控参数，支持通过该参数指定云主机在创建后是否开启详细监控，取值范围： <br />false：不开启，<br />true：开启<br />若指定该参数为true或不指定该参数，云主机内默认开启最新详细监控服务<br />若指定该参数为false，默认公共镜像不开启最新监控服务；私有镜像使用镜像中保留的监控服务<br />说明：仅部分资源池支持monitorService参数，详细请参考<a href="https://www.ctyun.cn/document/10026730/10325957">监控Agent概览</a>  */
}

type CtecsBatchRebuildInstancesV41Response struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                          `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsBatchRebuildInstancesV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsBatchRebuildInstancesV41ReturnObjResponse struct {
	JobIDList []*CtecsBatchRebuildInstancesV41ReturnObjJobIDListResponse `json:"jobIDList"` /*  重装任务列表  */
}

type CtecsBatchRebuildInstancesV41ReturnObjJobIDListResponse struct {
	JobID      string `json:"jobID,omitempty"`      /*  重装任务ID  */
	InstanceID string `json:"instanceID,omitempty"` /*  对应任务云主机ID  */
}
