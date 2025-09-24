package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2StatisticApi
/* 查询当前租户下处于“运行中”状态的分布式缓存Redis实例的统计信息。
 */type Dcs2StatisticApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2StatisticApi(client *core.CtyunClient) *Dcs2StatisticApi {
	return &Dcs2StatisticApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/statistic",
			ContentType:  "",
		},
	}
}

func (a *Dcs2StatisticApi) Do(ctx context.Context, credential core.Credential, req *Dcs2StatisticRequest) (*Dcs2StatisticResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.IncludeFailure != "" {
		ctReq.AddParam("includeFailure", req.IncludeFailure)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2StatisticResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2StatisticRequest struct {
	RegionId       string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	IncludeFailure string /*  否  */
}

type Dcs2StatisticResponse struct {
	StatusCode int32                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2StatisticReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2StatisticReturnObjResponse struct {
	Statistics *Dcs2StatisticReturnObjStatisticsResponse `json:"statistics"` /*  返回数据对象，数据见statistics  */
}

type Dcs2StatisticReturnObjStatisticsResponse struct {
	InstanceId  string `json:"instanceId,omitempty"`  /*  实例id  */
	InputKbps   string `json:"inputKbps,omitempty"`   /*  实例网络入流量，单位：kbit/s。  */
	OutputKbps  string `json:"outputKbps,omitempty"`  /*  实例网络出流量，单位：kbit/s。  */
	Keys        int32  `json:"keys,omitempty"`        /*  存储的KEY数量  */
	UsedMemory  int32  `json:"usedMemory,omitempty"`  /*  使用中的内存，单位：MB。  */
	MaxMemory   int32  `json:"maxMemory,omitempty"`   /*  最大的内存，单位：MB。  */
	CmdGetCount int32  `json:"cmdGetCount,omitempty"` /*  GET命令调用次数  */
	CmdSetCount int32  `json:"cmdSetCount,omitempty"` /*  SET命令调用次数  */
	UsedCpu     int32  `json:"usedCpu,omitempty"`     /*  实例进程消耗cpu时间的累计值，单位：秒。  */
}
