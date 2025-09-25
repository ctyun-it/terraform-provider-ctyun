package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsInstanceAttachSfsV41Api
/* 此接口提供用户实现云主机挂载一个或多个文件系统的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsInstanceAttachSfsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsInstanceAttachSfsV41Api(client *core.CtyunClient) *CtecsInstanceAttachSfsV41Api {
	return &CtecsInstanceAttachSfsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/sfs/attach",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsInstanceAttachSfsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsInstanceAttachSfsV41Request) (*CtecsInstanceAttachSfsV41Response, error) {
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
	var resp CtecsInstanceAttachSfsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsInstanceAttachSfsV41Request struct {
	RegionID    string                                         `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID  string                                         `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	SysInfoList []*CtecsInstanceAttachSfsV41SysInfoListRequest `json:"sysInfoList"`          /*  所绑定的文件系统详细信息  */
}

type CtecsInstanceAttachSfsV41SysInfoListRequest struct {
	FileSysID    string `json:"fileSysID,omitempty"`    /*  文件系统id  */
	FileSysRoute string `json:"fileSysRoute,omitempty"` /*  文件系统地址（固定值，每一个文件都有相对应的文件系统地址）  */
	MountPoint   string `json:"mountPoint,omitempty"`   /*  挂载点，即：本地挂载路径（云主机上用于挂载文件系统的本地路径)<br />linux镜像云主机使用限制：<br />1、单目录的长度不超过255个字符，总长度不能超过4095个字符； <br />2、必须以/开头，由数字，字母，点，下划线，减号组成，通过/（斜杠）分割<br />3、不能是系统路径：/, /bin, /usr, /boot, /dev, /etc, /lib, /lib64, /proc, /run, /sys, /var，/tmp，/sbin<br />推荐在/mnt下新建本地路径作为挂载路径，如：/mnt/docs. <br />windows镜像云主机实用限制：<br />只能输入E~Z（大写）内的单个字母作为盘符  */
	Option       string `json:"option,omitempty"`       /*  挂载参数：<br />linux云主机对应两个挂载参数（vers表示文件系统版本，可选3和4，建议取值3）vers=3,async,nolock,noatime,nodiratime,wsize=1048576,rsize=1048576,timeo=600 vers=4,async,nolock,noatime,nodiratime,wsize=1048576,rsize=1048576,timeo=600<br />windows云主机对应一个挂载参数：net use<br />注：<br />当linux云主机的协议类型选择NFSv3时，对应的option为vers3，当其协议类型选择NFSv4时，option应选vers4（即 protocol和option为对应关系）<br />windows云主机的协议类型CIFS，对应的option为：net use  */
	AutoMount    *bool  `json:"autoMount"`              /*  是否开机自动挂载（true/false），<br />当云主机重启时会自动挂载文件系统，默认为false  */
	Protocol     string `json:"protocol,omitempty"`     /*  文件协议类型，linux协议类型为：NFSv3、NFSv4，<br />windows协议类型为：CIFS  */
}

type CtecsInstanceAttachSfsV41Response struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                      `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                      `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                      `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsInstanceAttachSfsV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsInstanceAttachSfsV41ReturnObjResponse struct {
	JobID string `json:"jobID,omitempty"` /*  任务ID  */
}
