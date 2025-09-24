package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyInstanceSpecApi
/* 变更分布式缓存Redis实例规格：支持实例类型变更；规格/分片规格扩容；规格/分片规格缩容；增加分片；减少分片；增加副本；减少副本。
 */type Dcs2ModifyInstanceSpecApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyInstanceSpecApi(client *core.CtyunClient) *Dcs2ModifyInstanceSpecApi {
	return &Dcs2ModifyInstanceSpecApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/lifeCycleServant/modifyInstanceSpec",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyInstanceSpecApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyInstanceSpecRequest) (*Dcs2ModifyInstanceSpecResponse, error) {
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
	var resp Dcs2ModifyInstanceSpecResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyInstanceSpecRequest struct {
	RegionId          string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId        string `json:"prodInstId,omitempty"`        /*  实例ID  */
	OrderType         string `json:"orderType,omitempty"`         /*  实例变更类型，可选值：<br><br><b>架构变更</b><li><strong>UPGRADE</strong>：架构变更（说明：单机→主备，集群单机→集群主备时，内存容量不可变更）</li><br><b>纵向扩展</b><li><strong>EXPANSION</strong>：实例/分片规格扩容（适用：Single/Dual/StandardSingle/StandardDual/DirectClusterSingle/DirectCluster/ClusterOriginalProxy/OriginalMultipleReadLvs）<li><strong>CONTRACTION</strong>：实例/分片规格缩容（适用：同上）<br><br><b>横向扩展</b><li><strong>INCREASE_SHARDS</strong>：增加集群分片（适用：DirectClusterSingle/DirectCluster/ClusterOriginalProxy）<li><strong>DECREASE_SHARDS</strong>：减少集群分片（适用：同上）<br><br><b>副本调整</b><li><strong>INCREASE_REPLICAS</strong>：增加数据副本（适用：StandardDual/DirectCluster/ClusterOriginalProxy/OriginalMultipleReadLvs）<li><strong>DECREASE_REPLICAS</strong>：减少数据副本（适用：同上）  */
	AutoPay           *bool  `json:"autoPay"`                     /*  是否自动支付(仅对包周期实例有效)：<li>true：自动付费<li>false：手动付费(默认值)<br>选择为手动付费时，您需要在控制台的顶部菜单栏进入控制中心，单击费用中心 ，然后单击左侧导航栏的订单管理 > 我的订单，找到目标订单进行支付。  */
	SecondaryZoneName string `json:"secondaryZoneName,omitempty"` /*  备可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=17764&isNormal=1&vid=270">查询可用区信息</a> name字段<br>新增副本时必填<br>实例仅允许最多两个可用区。  */
	Capacity          string `json:"capacity,omitempty"`          /*  目标实例的总容量（单位：G）<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 计算方式：单分片内存shardMemSize × 分片数量shardCount 或 使用表SeriesInfo中memSizeItems(内存可选值)  */
	Version           string `json:"version,omitempty"`           /*  版本类型，默认BASIC<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的version值。<br>说明：当orderType为UPGRADE类型升级时必填。  */
	Edition           string `json:"edition,omitempty"`           /*  实例类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的seriesCode值  */
	ShardMemSize      string `json:"shardMemSize,omitempty"`      /*  目标实例的单分片内存(GB)<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中shardMemSizeItems(单分片内存可选值)，若shardMemSizeItems为空则无需填写  */
	ShardCount        int32  `json:"shardCount,omitempty"`        /*  目标实例的分片数量<li>DirectClusterSingle: 3-256<li>DirectCluster: 3-256<li>ClusterOriginalProxy: 3-64<br><br>其他实例类型无需填写此参数  */
	CopiesCount       int32  `json:"copiesCount,omitempty"`       /*  副本数量(2-6)<br>OriginalMultipleReadLvs必填，StandardDual/DirectCluster/ClusterOriginalProxy选填  */
}

type Dcs2ModifyInstanceSpecResponse struct {
	Message    string                                   `json:"message,omitempty"`    /*  响应信息  */
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	ReturnObj  *Dcs2ModifyInstanceSpecReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	RequestId  string                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifyInstanceSpecReturnObjResponse struct {
	ErrorMessage string  `json:"errorMessage,omitempty"` /*  错误信息  */
	Submitted    *bool   `json:"submitted"`              /*  是否提交成功  */
	NewOrderId   string  `json:"newOrderId,omitempty"`   /*  新订单id  */
	NewOrderNo   string  `json:"newOrderNo,omitempty"`   /*  新订单号  */
	TotalPrice   float64 `json:"totalPrice"`             /*  总价  */
}
