package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeAvailableResourceApi
/* 在创建分布式缓存Redis实例时，需要配置订购的实例类型属性version、edition、hostType、dataDiskType、capacity、shardMemSize、shardCount、copiesCount，可通过该接口查询。
 */type Dcs2DescribeAvailableResourceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeAvailableResourceApi(client *core.CtyunClient) *Dcs2DescribeAvailableResourceApi {
	return &Dcs2DescribeAvailableResourceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/lifeCycleServant/describeAvailableResource",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeAvailableResourceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeAvailableResourceRequest) (*Dcs2DescribeAvailableResourceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeAvailableResourceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeAvailableResourceRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
}

type Dcs2DescribeAvailableResourceResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeAvailableResourceReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeAvailableResourceReturnObjResponse struct {
	SeriesInfoList []*Dcs2DescribeAvailableResourceReturnObjSeriesInfoListResponse `json:"seriesInfoList"` /*  产品系列信息  */
	MirrorArray    []*Dcs2DescribeAvailableResourceReturnObjMirrorArrayResponse    `json:"mirrorArray"`    /*  系统镜像信息  */
}

type Dcs2DescribeAvailableResourceReturnObjSeriesInfoListResponse struct {
	Version           string                                                                  `json:"version,omitempty"`    /*  版本类型<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版<li>Capacity：容量型  */
	SeriesCode        string                                                                  `json:"seriesCode,omitempty"` /*  产品系列编码  */
	SeriesName        string                                                                  `json:"seriesName,omitempty"` /*  产品系列名称  */
	SeriesId          int64                                                                   `json:"seriesId,omitempty"`   /*  产品系列ID  */
	ResItems          []*Dcs2DescribeAvailableResourceReturnObjSeriesInfoListResItemsResponse `json:"resItems"`             /*  资源类型信息  */
	EngineTypeItems   []string                                                                `json:"engineTypeItems"`      /*  引擎版本  */
	MemSizeItems      []string                                                                `json:"memSizeItems"`         /*  内存容量可选值(GB)<br>说明：version为Classic和Capacity有值  */
	ShardMemSizeItems []string                                                                `json:"shardMemSizeItems"`    /*  单分片内存可选值(GB)。<br>说明：version为BASIC和PLUS有值  */
}

type Dcs2DescribeAvailableResourceReturnObjMirrorArrayResponse struct {
	AttrVal  string `json:"attrVal,omitempty"`  /*  操作系统  */
	Sort     int32  `json:"sort,omitempty"`     /*  排序  */
	AttrName string `json:"attrName,omitempty"` /*  操作系统名称  */
	Status   int32  `json:"status,omitempty"`   /*  状态,1：正常，其他表示异常  */
}

type Dcs2DescribeAvailableResourceReturnObjSeriesInfoListResItemsResponse struct {
	ResType string   `json:"resType,omitempty"` /*  资源类型<li>ecs：云服务器<li>ebs：磁盘  */
	ResName string   `json:"resName,omitempty"` /*  资源名称  */
	Items   []string `json:"items"`             /*  资源类型可选值<br>说明：以实际返回为准<br><br>云服务器<li>S7：通用型<li>C7：计算型<li>M7：内存型<li>HS1：海光通用型<li>HC1：海光计算增强型<li>KS1：鲲鹏通用型<li>KC1：鲲鹏计算增强型  <br><br>磁盘<li>SATA：普通IO<li>SAS：高IO<li>SSD：超高IO<li>FAST-SSD：极速型SSD  */
}
