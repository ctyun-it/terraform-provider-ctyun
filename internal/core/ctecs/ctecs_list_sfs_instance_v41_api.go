package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListSfsInstanceV41Api
/* 可以根据虚机查询绑定的文件系统列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsListSfsInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListSfsInstanceV41Api(client *core.CtyunClient) *CtecsListSfsInstanceV41Api {
	return &CtecsListSfsInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/sfs/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListSfsInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListSfsInstanceV41Request) (*CtecsListSfsInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListSfsInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListSfsInstanceV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}

type CtecsListSfsInstanceV41Response struct {
	StatusCode  int32                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                    `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                    `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                    `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListSfsInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsListSfsInstanceV41ReturnObjResponse struct {
	TotalCount int32                                              `json:"totalCount,omitempty"` /*  总记录数  */
	Results    []*CtecsListSfsInstanceV41ReturnObjResultsResponse `json:"results"`              /*  分页明细  */
}

type CtecsListSfsInstanceV41ReturnObjResultsResponse struct {
	FileSysID     string `json:"fileSysID,omitempty"`     /*  文件系统ID  */
	FileSysName   string `json:"fileSysName,omitempty"`   /*  文件系统名称  */
	FileSysStatus string `json:"fileSysStatus,omitempty"` /*  文件系统状态  */
	StorageType   string `json:"storageType,omitempty"`   /*  存储类型  */
	ShareProtocol string `json:"shareProtocol,omitempty"` /*  存储协议  */
	ExpireTime    string `json:"expireTime,omitempty"`    /*  到期时间  */
	SharePath     string `json:"sharePath,omitempty"`     /*  共享路径（文件提供的共享目录）  */
	CephID        string `json:"cephID,omitempty"`        /*  文件系统ID（底层）  */
	ResourceID    string `json:"resourceID,omitempty"`    /*  文件系统资源ID（IT）  */
}
