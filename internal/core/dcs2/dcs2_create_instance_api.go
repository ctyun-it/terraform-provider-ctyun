package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2CreateInstanceApi
/* 创建一个或者多个分布式缓存服务Redis基础版、增强版、经典版、容量型实例。
 */type Dcs2CreateInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2CreateInstanceApi(client *core.CtyunClient) *Dcs2CreateInstanceApi {
	return &Dcs2CreateInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/lifeCycleServant/createInstance",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2CreateInstanceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2CreateInstanceRequest) (*Dcs2CreateInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.ProjectID != "" {
		ctReq.AddHeader("projectID", req.ProjectID)
	}
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2CreateInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CreateInstanceRequest struct {
	RegionId           string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProjectID          string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br>说明：默认值为"0"  */
	AutoPay            *bool  `json:"autoPay"`                      /*  是否自动支付(仅对包周期实例有效)：<li>true：自动付费</li><li>false：手动付费(默认值)</li> <br>说明：选择为手动付费时，您需要在控制台的顶部菜单栏进入控制中心，单击费用中心 ，然后单击左侧导航栏的订单管理 > 我的订单，找到目标订单进行支付。  */
	Period             int32  `json:"period,omitempty"`             /*  订购时长(月)，仅当chargeType=PrePaid时必填，取值范围：1-6,12,24,36  */
	ChargeType         string `json:"chargeType,omitempty"`         /*  计费模式：<li>PrePaid：包年包月(需配合period使用)</li><li>PostPaid：按需计费(默认值)</li>  */
	ZoneName           string `json:"zoneName,omitempty"`           /*  主可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=17764&isNormal=1&vid=270">查询可用区信息</a> name字段  */
	SecondaryZoneName  string `json:"secondaryZoneName,omitempty"`  /*  备可用区名称(双/多副本建议填写)<br>默认与主可用区相同  */
	EngineVersion      string `json:"engineVersion,omitempty"`      /*  Redis引擎版本<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的engineTypeItems(引擎版本可选值)  */
	Version            string `json:"version,omitempty"`            /*  版本类型。<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中的version值<br>可选值：<li>BASIC：基础版<li>PLUS：增强版<li>Classic：经典版  */
	Edition            string `json:"edition,omitempty"`            /*  实例类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a>  使用表SeriesInfo中的seriesCode值  */
	HostType           string `json:"hostType,omitempty"`           /*  主机类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表resItems中resType==ecs的items(主机类型可选值)  */
	DataDiskType       string `json:"dataDiskType,omitempty"`       /*  磁盘类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表resItems中resType==ebs的items(磁盘类型可选值)  */
	DataDiskSize       int32  `json:"dataDiskSize,omitempty"`       /*  存储空间(GB，仅容量型支持)，需为内存5-10倍且为10的倍数  */
	MirrorCategoryName string `json:"mirrorCategoryName,omitempty"` /*  操作系统镜像类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表Mirror中attrName值  */
	ShardMemSize       string `json:"shardMemSize,omitempty"`       /*  单分片内存(GB)<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中shardMemSizeItems(单分片内存可选值)，若shardMemSizeItems为空则无需填写  */
	ShardCount         int32  `json:"shardCount,omitempty"`         /*  分片数量<li>DirectClusterSingle: 3-256</li><li>DirectCluster: 3-256</li><li>ClusterOriginalProxy: 3-64</li>其他实例类型无需填写此参数  */
	Capacity           string `json:"capacity,omitempty"`           /*  内存容量(GB)<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 计算方式：单分片内存shardMemSize × 分片数量shardCount 或 使用表SeriesInfo中memSizeItems(内存可选值)  */
	CopiesCount        int32  `json:"copiesCount,omitempty"`        /*  副本数量，取值范围2-6。<li>OriginalMultipleReadLvs：必填</li><li>StandardDual/DirectCluster/ClusterOriginalProxy：选填</li><li>其他实例类型：无需填写</li>  */
	InstanceName       string `json:"instanceName,omitempty"`       /*  实例名称<li>字母开头</li><li>可包含字母/数字/中划线</li><li>长度1-39<li>实例名称不可重复</li>  */
	VpcId              string `json:"vpcId,omitempty"`              /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94&vid=88">查询VPC列表</a> vpcID字段。<br><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94&vid=88">创建VPC</a>  */
	SubnetId           string `json:"subnetId,omitempty"`           /*  子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94&vid=88">查询子网列表</a> subnetID字段。  */
	Secgroups          string `json:"secgroups,omitempty"`          /*  安全组ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/searchCtapi/ctApiDebug?product=18&api=4817&vid=88">查询用户安全组列表</a> id字段。  */
	Password           string `json:"password,omitempty"`           /*  实例密码<li>长度8-26字符</li><li>必须同时包含大写字母、小写字母、数字、英文格式特殊符号(@%^*_+!$-=.) 中的三种类型</li><li>不能有空格</li>  */
	AutoRenew          *bool  `json:"autoRenew"`                    /*  自动续费开关<li>true：开启</li><li>false：关闭(默认)</li>  */
	AutoRenewPeriod    string `json:"autoRenewPeriod,omitempty"`    /*  自动续费周期(月)<br>autoRenew=true时必填，可选：1-6,12,24,36  */
	Size               int32  `json:"size,omitempty"`               /*  购买数量(1-100，默认1)  */
}

type Dcs2CreateInstanceResponse struct {
	Message    string                               `json:"message,omitempty"`    /*  响应信息  */
	StatusCode int32                                `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	ReturnObj  *Dcs2CreateInstanceReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	RequestId  string                               `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                               `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2CreateInstanceReturnObjResponse struct {
	ErrorMessage string  `json:"errorMessage,omitempty"` /*  错误信息  */
	Submitted    *bool   `json:"submitted"`              /*  是否成功提交  */
	NewOrderId   string  `json:"newOrderId,omitempty"`   /*  订单ID  */
	NewOrderNo   string  `json:"newOrderNo,omitempty"`   /*  订单号  */
	TotalPrice   float64 `json:"totalPrice"`             /*  总价  */
}
