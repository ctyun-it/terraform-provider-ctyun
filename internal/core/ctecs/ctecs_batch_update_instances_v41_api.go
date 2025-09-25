package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsBatchUpdateInstancesV41Api
/* 该接口提供用户更新多台云主机的部分信息的功能<br />目前支持更新云主机的信息为：云主机显示名称（displayName）、云主机名称（instanceName）、云主机描述信息（instanceDescription）<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;如果使用私有镜像创建的云主机执行该操作时，请先检查云主机内部是否安装了QGA（qemu-guest-agent）。不同操作系统请参考：<a href="https://www.ctyun.cn/document/10027726/10747194">Windows系统盘镜像文件安装QGA</a>和<a href="https://www.ctyun.cn/document/10027726/10747147">Linux系统盘镜像文件安装QGA</a><br />
 */type CtecsBatchUpdateInstancesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsBatchUpdateInstancesV41Api(client *core.CtyunClient) *CtecsBatchUpdateInstancesV41Api {
	return &CtecsBatchUpdateInstancesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/batch-update-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsBatchUpdateInstancesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsBatchUpdateInstancesV41Request) (*CtecsBatchUpdateInstancesV41Response, error) {
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
	var resp CtecsBatchUpdateInstancesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsBatchUpdateInstancesV41Request struct {
	RegionID   string                                           `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	UpdateInfo []*CtecsBatchUpdateInstancesV41UpdateInfoRequest `json:"updateInfo"`         /*  批量更新信息列表  */
}

type CtecsBatchUpdateInstancesV41UpdateInfoRequest struct {
	InstanceID          string `json:"instanceID,omitempty"`          /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	DisplayName         string `json:"displayName,omitempty"`         /*  云主机显示名称，长度为2~63个字符<br />注：displayName、instanceName、instanceDescription不可全为空  */
	InstanceName        string `json:"instanceName,omitempty"`        /*  云主机名称。不同操作系统下，云主机名称规则有差异<br />Windows：长度为2-15个字符，允许使用大小写字母、数字或连字符（-）。不能以连字符（-）开头或结尾，不能连续使用连字符（-），也不能仅使用数字；<br />其他操作系统：长度为2-64字符，允许使用点（.）分隔字符成多段，每段允许使用大小写字母、数字或连字符（-），但不能连续使用点号（.）或连字符（-），不能以点号（.）或连字符（-）开头或结尾，也不能仅使用数字<br />注：displayName、instanceName、instanceDescription不可全为空  */
	InstanceDescription string `json:"instanceDescription,omitempty"` /*  云主机描述信息，限制长度为0~255个字符<br />注：displayName、instanceName、instanceDescription不可全为空  */
}

type CtecsBatchUpdateInstancesV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码部分    */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码，详见错误码部分     */
	Message     string                                         `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                         `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsBatchUpdateInstancesV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsBatchUpdateInstancesV41ReturnObjResponse struct {
	UpdateInfo []*CtecsBatchUpdateInstancesV41ReturnObjUpdateInfoResponse `json:"updateInfo"` /*  被更新云主机信息  */
}

type CtecsBatchUpdateInstancesV41ReturnObjUpdateInfoResponse struct {
	InstanceID          string `json:"instanceID,omitempty"`          /*  被更新名称的云主机ID  */
	DisplayName         string `json:"displayName,omitempty"`         /*  更新后的云主机名称  */
	InstanceName        string `json:"instanceName,omitempty"`        /*  被更新云主机名称    */
	InstanceDescription string `json:"instanceDescription,omitempty"` /*  被更新云主机描述信息         */
}
