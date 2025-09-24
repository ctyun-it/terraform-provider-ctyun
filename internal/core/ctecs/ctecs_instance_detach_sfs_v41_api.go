package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsInstanceDetachSfsV41Api
/* 此接口提供用户实现云主机卸载一个或多个文件系统的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsInstanceDetachSfsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsInstanceDetachSfsV41Api(client *core.CtyunClient) *CtecsInstanceDetachSfsV41Api {
	return &CtecsInstanceDetachSfsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/sfs/detach",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsInstanceDetachSfsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsInstanceDetachSfsV41Request) (*CtecsInstanceDetachSfsV41Response, error) {
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
	var resp CtecsInstanceDetachSfsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsInstanceDetachSfsV41Request struct {
	RegionID    string                                         `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID  string                                         `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	SysInfoList []*CtecsInstanceDetachSfsV41SysInfoListRequest `json:"sysInfoList"`          /*  所解绑的文件系统详细信息  */
	ForceDel    *bool                                          `json:"forceDel"`             /*  是否强制解绑 (true/false)，默认非强制  */
}

type CtecsInstanceDetachSfsV41SysInfoListRequest struct {
	FileSysRoute string `json:"fileSysRoute,omitempty"` /*  文件系统地址（固定值，每一个文件都有相对应的文件系统地址）  */
	MountPoint   string `json:"mountPoint,omitempty"`   /*  挂载点  */
}

type CtecsInstanceDetachSfsV41Response struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                      `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                      `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                      `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsInstanceDetachSfsV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsInstanceDetachSfsV41ReturnObjResponse struct {
	JobID string `json:"jobID,omitempty"` /*  任务ID  */
}
