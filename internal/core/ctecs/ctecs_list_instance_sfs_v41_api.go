package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListInstanceSfsV41Api
/* 可以根据用户查询对应的文件系统关联的云主机列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsListInstanceSfsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceSfsV41Api(client *core.CtyunClient) *CtecsListInstanceSfsV41Api {
	return &CtecsListInstanceSfsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/sfs/vms-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceSfsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceSfsV41Request) (*CtecsListInstanceSfsV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("fileSysID", req.FileSysID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListInstanceSfsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceSfsV41Request struct {
	RegionID  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	FileSysID string /*  文件系统ID  */
}

type CtecsListInstanceSfsV41Response struct {
	StatusCode  int32                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                    `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                    `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                    `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListInstanceSfsV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsListInstanceSfsV41ReturnObjResponse struct {
	CurrentCount int32                                              `json:"currentCount,omitempty"` /*  当前页码  */
	TotalCount   int32                                              `json:"totalCount,omitempty"`   /*  总记录数  */
	Results      []*CtecsListInstanceSfsV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListInstanceSfsV41ReturnObjResultsResponse struct {
	InstanceID   string                                                 `json:"instanceID,omitempty"`   /*  云主机ID  */
	Option       string                                                 `json:"option,omitempty"`       /*  挂载参数  */
	FileSysRoute string                                                 `json:"fileSysRoute,omitempty"` /*  文件系统绝对路径  */
	MountPoint   string                                                 `json:"mountPoint,omitempty"`   /*  挂载点  */
	AutoMount    *bool                                                  `json:"autoMount"`              /*  是否开机自动挂载  */
	Protocol     string                                                 `json:"protocol,omitempty"`     /*  挂载协议  */
	FileSysID    string                                                 `json:"fileSysID,omitempty"`    /*  文件系统ID  */
	VmInfo       *CtecsListInstanceSfsV41ReturnObjResultsVmInfoResponse `json:"vmInfo"`                 /*  虚机信息  */
}

type CtecsListInstanceSfsV41ReturnObjResultsVmInfoResponse struct {
	InstanceStatus string `json:"instanceStatus,omitempty"` /*  虚机状态  */
	DisplayName    string `json:"displayName,omitempty"`    /*  虚机展示名称  */
	InstanceName   string `json:"instanceName,omitempty"`   /*  名称  */
	ImageName      string `json:"imageName,omitempty"`      /*  镜像名称  */
	InstanceID     string `json:"instanceID,omitempty"`     /*  虚机ID  */
	ProjectID      string `json:"projectID,omitempty"`      /*  项目ID  */
}
