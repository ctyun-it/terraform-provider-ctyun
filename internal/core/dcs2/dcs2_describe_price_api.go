package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribePriceApi
/* 查询创建分布式缓存Redis实例、续费或变更实例规格等操作产生的费用。
 */type Dcs2DescribePriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribePriceApi(client *core.CtyunClient) *Dcs2DescribePriceApi {
	return &Dcs2DescribePriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/lifeCycleServant/describePrice",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2DescribePriceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribePriceRequest) (*Dcs2DescribePriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribePriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribePriceRequest struct {
	RegionId      string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId    string `json:"prodInstId,omitempty"`    /*  实例ID<br>说明：orderType为BUY无需填写，其他值必填。  */
	OrderType     string `json:"orderType,omitempty"`     /*  订单类型，可选值：<br><br><b>订购</b><li><strong>BUY</strong>：订购<br><br><b>续订</b><li><strong>RENEW</strong>：续订<br><br><b>架构变更</b><li><strong>UPGRADE</strong>：架构变更（说明：单机→主备，集群单机→集群主备时，内存容量不可变更）</li><br><b>纵向扩展</b><li><strong>EXPANSION</strong>：实例/分片规格扩容（适用：Single/Dual/StandardSingle/StandardDual/DirectClusterSingle/DirectCluster/ClusterOriginalProxy/OriginalMultipleReadLvs）<li><strong>CONTRACTION</strong>：实例/分片规格缩容（适用：同上）<br><br><b>横向扩展</b><li><strong>INCREASE_SHARDS</strong>：增加集群分片（适用：DirectClusterSingle/DirectCluster/ClusterOriginalProxy）<li><strong>DECREASE_SHARDS</strong>：减少集群分片（适用：同上）<br><br><b>副本调整</b><li><strong>INCREASE_REPLICAS</strong>：增加数据副本（适用：StandardDual/DirectCluster/ClusterOriginalProxy/OriginalMultipleReadLvs）<li><strong>DECREASE_REPLICAS</strong>：减少数据副本（适用：同上）  */
	Version       string `json:"version,omitempty"`       /*  版本类型，默认BASIC。<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的version值<br>可选值：<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版  */
	Edition       string `json:"edition,omitempty"`       /*  实例类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的seriesCode值  */
	ShardMemSize  string `json:"shardMemSize,omitempty"`  /*  单分片内存(GB)<br/><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中shardMemSizeItems(单分片内存可选值)，若shardMemSizeItems为空则无需填写  */
	ShardCount    int32  `json:"shardCount,omitempty"`    /*  分片数量<li>DirectClusterSingle: 3-256</li><li>DirectCluster: 3-256</li><li>ClusterOriginalProxy: 3-64</li>其他实例类型无需填写此参数  */
	Capacity      string `json:"capacity,omitempty"`      /*  内存容量(GB)<br/><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 计算方式：单分片内存shardMemSize × 分片数量shardCount 或 使用表SeriesInfo中memSizeItems(内存可选值)  */
	CopiesCount   int32  `json:"copiesCount,omitempty"`   /*  副本数量(2-6)<br>OriginalMultipleReadLvs必填，StandardDual/DirectCluster/ClusterOriginalProxy选填  */
	HostType      string `json:"hostType,omitempty"`      /*  主机类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表resItems中resType==ecs的items(主机类型可选值)  */
	DataDiskType  string `json:"dataDiskType,omitempty"`  /*  磁盘类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表resItems中resType==ebs的items(磁盘类型可选值)  */
	DataDiskSize  int32  `json:"dataDiskSize,omitempty"`  /*  存储空间(GB，仅容量型支持)，需为内存5-10倍且为10的倍数  */
	Period        string `json:"period,omitempty"`        /*  订购时长(月)，仅当chargeType=PrePaid时必填，取值范围：1-6,12,24,36  */
	ChargeType    string `json:"chargeType,omitempty"`    /*  计费模式<li>PrePaid：包年包月(需配合period使用)<li>PostPaid：按需计费(默认值)  */
	EngineVersion string `json:"engineVersion,omitempty"` /*  Redis引擎版本，当orderType为BUY时必填。<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的engineTypeItems(引擎版本可选值)  */
	Size          int32  `json:"size,omitempty"`          /*  数量,允许批量订购，允许范围1-100，默认1  */
}

type Dcs2DescribePriceResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribePriceReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribePriceReturnObjResponse struct {
	TotalPrice       float64                                             `json:"totalPrice"`                 /*  订单总价  */
	ServiceTag       string                                              `json:"serviceTag,omitempty"`       /*  serviceTag  */
	FinalPrice       float64                                             `json:"finalPrice"`                 /*  订单最终价格  */
	SubOrderPrices   []*Dcs2DescribePriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"`             /*  子订单价格  */
	VerifyStatusCode string                                              `json:"verifyStatusCode,omitempty"` /*  状态校验  */
}

type Dcs2DescribePriceReturnObjSubOrderPricesResponse struct {
	TotalPrice      float64                                                            `json:"totalPrice"`           /*  子订单总价  */
	ServiceTag      string                                                             `json:"serviceTag,omitempty"` /*  serviceTag  */
	FinalPrice      float64                                                            `json:"finalPrice"`           /*  子订单最终价格  */
	OrderItemPrices []*Dcs2DescribePriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  子订单价格明细  */
}

type Dcs2DescribePriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       string  `json:"itemId,omitempty"`       /*  子订单资源ID  */
	TotalPrice   float64 `json:"totalPrice"`             /*  子订单资源总价  */
	FinalPrice   float64 `json:"finalPrice"`             /*  子订单资源最终价格  */
	ResourceType string  `json:"resourceType,omitempty"` /*  子订单资源类型  */
}
