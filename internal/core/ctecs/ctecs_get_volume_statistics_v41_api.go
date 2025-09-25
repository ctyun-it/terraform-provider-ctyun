package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsGetVolumeStatisticsV41Api
/* 查询用户云硬盘统计信息<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsGetVolumeStatisticsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsGetVolumeStatisticsV41Api(client *core.CtyunClient) *CtecsGetVolumeStatisticsV41Api {
	return &CtecsGetVolumeStatisticsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/volume/statistics",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsGetVolumeStatisticsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsGetVolumeStatisticsV41Request) (*CtecsGetVolumeStatisticsV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsGetVolumeStatisticsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsGetVolumeStatisticsV41Request struct {
	RegionID  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ProjectID string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
}

type CtecsGetVolumeStatisticsV41Response struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsGetVolumeStatisticsV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsGetVolumeStatisticsV41ReturnObjResponse struct {
	VolumeStatistics *CtecsGetVolumeStatisticsV41ReturnObjVolumeStatisticsResponse `json:"volumeStatistics"` /*  统计明细  */
}

type CtecsGetVolumeStatisticsV41ReturnObjVolumeStatisticsResponse struct {
	TotalCount    int32 `json:"totalCount,omitempty"`    /*  云硬盘总数  */
	RootDiskCount int32 `json:"rootDiskCount,omitempty"` /*  系统盘数量  */
	DataDiskCount int32 `json:"dataDiskCount,omitempty"` /*  数据盘数量  */
	TotalSize     int32 `json:"totalSize,omitempty"`     /*  云硬盘总大小，单位为GB  */
	RootDiskSize  int32 `json:"rootDiskSize,omitempty"`  /*  系统盘大小，单位为GB  */
	DataDiskSize  int32 `json:"dataDiskSize,omitempty"`  /*  数据盘大小，单位为GB  */
}
