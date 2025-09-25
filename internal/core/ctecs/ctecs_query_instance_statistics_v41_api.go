package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryInstanceStatisticsV41Api
/* 查询用户云主机统计信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsQueryInstanceStatisticsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryInstanceStatisticsV41Api(client *core.CtyunClient) *CtecsQueryInstanceStatisticsV41Api {
	return &CtecsQueryInstanceStatisticsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/statistics-instance",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryInstanceStatisticsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryInstanceStatisticsV41Request) (*CtecsQueryInstanceStatisticsV41Response, error) {
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
	var resp CtecsQueryInstanceStatisticsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryInstanceStatisticsV41Request struct {
	RegionID  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ProjectID string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
}

type CtecsQueryInstanceStatisticsV41Response struct {
	StatusCode  int32                                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                            `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                            `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                            `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                            `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQueryInstanceStatisticsV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsQueryInstanceStatisticsV41ReturnObjResponse struct {
	InstanceStatistics *CtecsQueryInstanceStatisticsV41ReturnObjInstanceStatisticsResponse `json:"instanceStatistics"` /*  分页明细  */
}

type CtecsQueryInstanceStatisticsV41ReturnObjInstanceStatisticsResponse struct {
	TotalCount          int32 `json:"totalCount,omitempty"`          /*  云主机总数  */
	RunningCount        int32 `json:"RunningCount,omitempty"`        /*  运行中的云主机数量  */
	ShutdownCount       int32 `json:"shutdownCount,omitempty"`       /*  关机数量  */
	ExpireCount         int32 `json:"expireCount,omitempty"`         /*  过期数量  */
	ExpireRunningCount  int32 `json:"expireRunningCount,omitempty"`  /*  过期运行中数量  */
	ExpireShutdownCount int32 `json:"expireShutdownCount,omitempty"` /*  过期已关机数量  */
	CpuCount            int32 `json:"cpuCount,omitempty"`            /*  cpu数量  */
	MemoryCount         int32 `json:"memoryCount,omitempty"`         /*  内存总量，单位为GB  */
}
