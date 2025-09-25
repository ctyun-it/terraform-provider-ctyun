package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeNodeMonitorValuesApi
/* 查询某一个类簇的调用次数，包含string、hash、keys、list、set 等。
 */type Dcs2DescribeNodeMonitorValuesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeNodeMonitorValuesApi(client *core.CtyunClient) *Dcs2DescribeNodeMonitorValuesApi {
	return &Dcs2DescribeNodeMonitorValuesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/monitorManageMgrServant/describeNodeMonitorValues",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeNodeMonitorValuesApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeNodeMonitorValuesRequest) (*Dcs2DescribeNodeMonitorValuesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("type", req.RawType)
	ctReq.AddParam("nodeName", req.NodeName)
	ctReq.AddParam("startTime", req.StartTime)
	ctReq.AddParam("endTime", req.EndTime)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeNodeMonitorValuesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeNodeMonitorValuesRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	RawType    string /*  监控类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7780&isNormal=1&vid=270">查询命令调用类族</a> returnObj表data字段  */
	NodeName   string /*  节点名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7752&isNormal=1&vid=270">获取redis节点名列表</a> node表nodeName字段  */
	StartTime  string /*  开始时间  */
	EndTime    string /*  结束时间  */
}

type Dcs2DescribeNodeMonitorValuesResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeNodeMonitorValuesReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeNodeMonitorValuesReturnObjResponse struct {
	Total int32                                                 `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeNodeMonitorValuesReturnObjRowsResponse `json:"rows"`            /*  白名单列表  */
}

type Dcs2DescribeNodeMonitorValuesReturnObjRowsResponse struct {
	Metric *Dcs2DescribeNodeMonitorValuesReturnObjRowsMetricResponse   `json:"metric"` /*  监控指标  */
	Values []*Dcs2DescribeNodeMonitorValuesReturnObjRowsValuesResponse `json:"values"` /*  键值对数组  */
}

type Dcs2DescribeNodeMonitorValuesReturnObjRowsMetricResponse struct {
	Tenant_id      string `json:"tenant_id,omitempty"`      /*  租户ID  */
	Carms_obj_id   string `json:"carms_obj_id,omitempty"`   /*  实例ID  */
	Instance_name  string `json:"instance_name,omitempty"`  /*  实例名称  */
	Redis_url      string `json:"redis_url,omitempty"`      /*  节点IP  */
	Carms_obj_name string `json:"carms_obj_name,omitempty"` /*  实例名称  */
	Carms_obj_type string `json:"carms_obj_type,omitempty"` /*  组件类型，固定值  */
	User_name      string `json:"user_name,omitempty"`      /*  实例名  */
	Cache_exporter string `json:"cache_exporter,omitempty"` /*  Exporter类型  */
	Proxy_name     string `json:"proxy_name,omitempty"`     /*  代理节点名  */
	Monitor_point  string `json:"monitor_point,omitempty"`  /*  监控节点  */
	Instance_id    string `json:"instance_id,omitempty"`    /*  实例ID  */
	Region_name    string `json:"region_name,omitempty"`    /*  资源池名称  */
	Cmd            string `json:"cmd,omitempty"`            /*  命令  */
	Region_code    string `json:"region_code,omitempty"`    /*  资源池编码  */
	Tenant_code    string `json:"tenant_code,omitempty"`    /*  租户编码  */
}

type Dcs2DescribeNodeMonitorValuesReturnObjRowsValuesResponse struct {
	Value []string `json:"value"` /*  时间、值字符串数组  */
}
