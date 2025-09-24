package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmCreateInstanceApi
/* 创建物理机
 */type EbmCreateInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmCreateInstanceApi(client *core.CtyunClient) *EbmCreateInstanceApi {
	return &EbmCreateInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebm/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmCreateInstanceApi) Do(ctx context.Context, credential core.Credential, req *EbmCreateInstanceRequest) (*EbmCreateInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmCreateInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmCreateInstanceRequest struct {
	RegionID string `json:"regionID"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string `json:"azName"` /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */DeviceType string `json:"deviceType"` /*  物理机套餐类型<br/><a href="https://www.ctyun.cn/document/10027724/10040107">查询资源池内物理机套餐</a><br /><a href="https://www.ctyun.cn/document/10027724/10040124">查询指定物理机的套餐信息</a>
	 */Name string `json:"name"` /*  物理机名称，长度为2-31位
	 */Hostname string `json:"hostname"` /*  hostname，linux系统2到63位长度；windows系统2-15位长度；<br/>允许使用大小写字母、数字、连字符'-'，必须以字母开头（大小写均可），不能连续使用'-'，'-'不能用于结尾，不能仅使用数字；<br/>支持模式串{R:x}，表示生成数字[x,x+n-1]，其中n表示购买实例的数量，1 ≤ x ≤ 9799且x只能为整数。<br/>例子：填写server{R:3}pm，购买1台时，实例主机名为server0003pm；购买2台时，实例主机名分别为server0003pm，server0004pm )
	 */ImageUUID string `json:"imageUUID"` /*  物理机镜像UUID<br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4577&data=97">查询物理机镜像</a>
	 */Password *string `json:"password"` /*  密码  (必须包含大小写字母和（一个数字或者特殊字符）长度8到30位)，未传入有效的keyName时必须传入password
	 */SystemVolumeRaidUUID *string `json:"systemVolumeRaidUUID"` /*  本地系统盘raid类型，如果有本地盘则必填<br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=5881&data=97">查询物理机实例本地盘raid信息</a><br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4576&data=97">查询物理机raid</a>
	 */DataVolumeRaidUUID *string `json:"dataVolumeRaidUUID"` /*  本地数据盘raid类型，如果有本地盘则必填<br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=5881&data=97">查询物理机实例本地盘raid信息</a><br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4576&data=97">查询物理机raid</a>
	 */VpcID string `json:"vpcID"` /*  主网卡虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br /> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94">查询VPC列表</a><br /><a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94">创建VPC</a><br />注：在多可用区类型资源池下，vpcID通常为“vpc-”开头，非多可用区类型资源池vpcID为uuid格式
	 */ExtIP string `json:"extIP"` /*  是否使用弹性公网IP ，取值范围:[1=自动分配,0=不使用,2=使用已有]
	 */ProjectID *string `json:"projectID"` /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"
	 */IpType *string `json:"ipType"` /*  弹性IP版本 ，取值范围:[ipv4=v4地址,ipv6=v6地址]，默认值:ipv4
	 */BandWidth int32 `json:"bandWidth"` /*  带宽 ，单位为Mbit/s，取值范围:[1~2000]，默认值:100
	 */BandWidthType *string `json:"bandWidthType"` /*  带宽类型
	 */PublicIP *string `json:"publicIP"` /*  弹性公网IP的ID，您可以查看<a href="https://www.ctyun.cn/document/10026753/10026909">产品定义-弹性IP</a>来了解弹性公网IP <br /><a href="https://www.ctyun.cn/document/10026753/10040758">查询指定地域已创建的弹性 IP</a><br /> <a href="https://www.ctyun.cn/document/10026753/10040759">创建弹性 IP</a>
	 */SecurityGroupID *string `json:"securityGroupID"` /*  安全组ID，套餐中smartNicExist属性为true可支持安全组。创建弹性裸金属必须传入安全组ID，标准裸金属不支持传入安全组ID。您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br /><a href="https://www.ctyun.cn/document/10026755/10040907">查询用户安全组列表</a><br /><a href="https://www.ctyun.cn/document/10026755/10040938">创建安全组</a> <br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4574&data=97">查询物理机套餐</a><br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=5117&data=97">物理机查询对应套餐信息</a>
	 */DiskList []*EbmCreateInstanceDiskListRequest `json:"diskList"` /*  云盘信息列表，套餐中supportCloud为true表示支持云盘，若不支持则无需填写。<br />您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4574&data=97">查询资源池内物理机套餐</a>和<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=5117&data=97">查询指定物理机的套餐信息</a>了解当前物理机是否支持云盘<br/ >您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7338&data=48">查询云硬盘列表</a>获取云硬盘信息
	 */NetworkCardList []*EbmCreateInstanceNetworkCardListRequest `json:"networkCardList"` /*  网卡
	 */PayVoucherPrice float32 `json:"payVoucherPrice"` /*  代金券，满足以下规则：两位小数，不足两位自动补0，超过两位小数无效；不可为负数；字段为0时表示不使用代金券
	 */AutoRenewStatus int32 `json:"autoRenewStatus"` /*  是否自动续订，默认非自动续订。取值范围：<br/>0（不续费），<br/>1（自动续费），<br/>注：按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年
	 */InstanceChargeType *string `json:"instanceChargeType"` /*  实例计费类型，默认为ORDER_ON_CYCLE（包年包月）
	 */CycleCount int32 `json:"cycleCount"` /*  订购时长
	 */CycleType string `json:"cycleType"` /*  订购周期类型 ，取值范围:[MONTH=按月,YEAR=按年]
	 */OrderCount int32 `json:"orderCount"` /*  购买数量
	 */ClientToken string `json:"clientToken"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时
	 */
}

type EbmCreateInstanceDiskListRequest struct {
	DiskType string `json:"diskType"` /*  磁盘类型 system或data，套餐中cloudBoot为true表示支持云盘系统盘
	 */DiskMode *string `json:"diskMode"` /*  磁盘属性(VBD)
	 */Title *string `json:"title"` /*  磁盘名称 ，长度2~64,不支持中文
	 */RawType string `json:"type"` /*  磁盘分类，取值范围:[SAS=SAS盘,SATA=SATA盘,SSD=SSD盘]
	 */Size int32 `json:"size"` /*  云硬盘容量，单位为GB；系统盘容量取值范围：[100, 2048]，数据盘容量取值范围：[10, 32768]
	 */
}

type EbmCreateInstanceNetworkCardListRequest struct {
	Title *string `json:"title"` /*  网卡名称
	 */FixedIP *string `json:"fixedIP"` /*  浮动IP，内网IPv4地址，注：不可使用已占用IP
	 */Master bool `json:"master"` /*  是否主节点(True代表主节点)
	 */Ipv6 *string `json:"ipv6"` /*  内网IPv6地址
	 */SubnetID string `json:"subnetID"` /*  子网UUID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义<br /> <a href="https://www.ctyun.cn/document/10026755/10040797">查询子网列表</a><br /><a href="https://www.ctyun.cn/document/10026755/10040804">创建子网</a><br/>注：在多可用区类型资源池下，subnetID通常以“subnet-”开头；非多可用区类型资源池subnetID为uuid格式
	 */
}

type EbmCreateInstanceResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmCreateInstanceReturnObjResponse `json:"returnObj"` /*  返回参数，参考returnObj
	 */
}

type EbmCreateInstanceReturnObjResponse struct {
	RegionID *string `json:"regionID"` /*  资源池ID
	 */MasterOrderID *string `json:"masterOrderID"` /*  订单ID
	 */MasterOrderNO *string `json:"masterOrderNO"` /*  订单号
	 */
}
