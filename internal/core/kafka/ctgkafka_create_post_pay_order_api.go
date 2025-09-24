package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaCreatePostPayOrderApi
/* 创建按需计费实例。
 */type CtgkafkaCreatePostPayOrderApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaCreatePostPayOrderApi(client *core.CtyunClient) *CtgkafkaCreatePostPayOrderApi {
	return &CtgkafkaCreatePostPayOrderApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/createPostPayOrder",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaCreatePostPayOrderApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaCreatePostPayOrderRequest) (*CtgkafkaCreatePostPayOrderResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.ProjectId != "" {
		ctReq.AddHeader("projectId", req.ProjectId)
	}
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaCreatePostPayOrderRequest
		RegionId  interface{} `json:"regionId,omitempty"`
		ProjectId interface{} `json:"projectId,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaCreatePostPayOrderResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaCreatePostPayOrderRequest struct {
	RegionId        string   `json:"regionId,omitempty"`        /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProjectId       string   `json:"projectId,omitempty"`       /*  企业项目ID(默认值：0)。您可以通过 <a href="https://www.ctyun.cn/document/10017248/10017965">查看企业项目资源</a> 获取企业项目ID。  */
	ClusterName     string   `json:"clusterName,omitempty"`     /*  实例名称。<br>规则：长度4~40个字符，大小写字母开头，只能包含大小写字母、数字及分隔符(-)，大小写字母或数字结尾，实例名称不可重复。  */
	EngineVersion   string   `json:"engineVersion,omitempty"`   /*  实例的引擎版本，默认为3.6。<li>2.8：2.8.x的引擎版本<li>3.6：3.6.x的引擎版本  */
	SpecName        string   `json:"specName,omitempty"`        /*  实例的规格类型，资源池所具备的规格可通过查询产品规格接口获取，默认可选如下：<br>计算增强型的规格可选为：<li>kafka.2u4g.cluster<li>kafka.4u8g.cluster<li>kafka.8u16g.cluster<li>kafka.12u24g.cluster<li>kafka.16u32g.cluster<li>kafka.24u48g.cluster<li>kafka.32u64g.cluster<li>kafka.48u96g.cluster<li>kafka.64u128g.cluster <br>海光-计算增强型的规格可选为：<li>kafka.hg.2u4g.cluster<li>kafka.hg.4u8g.cluster<li>kafka.hg.8u16g.cluster<li>kafka.hg.16u32g.cluster<li>kafka.hg.32u64g.cluster <br>鲲鹏-计算增强型的规格可选为：<li>kafka.kp.2u4g.cluster<li>kafka.kp.4u8g.cluster<li>kafka.kp.8u16g.cluster<li>kafka.kp.16u32g.cluster<li>kafka.kp.32u64g.cluster  */
	NodeNum         int32    `json:"nodeNum,omitempty"`         /*  节点数。单机版为1个，集群版3~50个。  */
	ZoneList        []string `json:"zoneList"`                  /*  实例所在可用区信息。只能填一个（单可用区）或三个（多可用区），可用区信息可调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&isNormal=1&vid=81">资源池可用区查询</a>API接口查询。  */
	DiskType        string   `json:"diskType,omitempty"`        /*  磁盘类型，资源池所具备的磁盘类型可通过查询产品规格接口获取，默认取值：<li>SAS：高IO<li>SSD：超高IO<li>FAST-SSD：极速型SSD  */
	DiskSize        int32    `json:"diskSize,omitempty"`        /*  单个节点的磁盘存储空间，单位为GB，存储空间取值范围100GB ~ 10000，并且为100的倍数。实例总存储空间为diskSize * nodeNum。  */
	VpcId           string   `json:"vpcId,omitempty"`           /*  VPC网络ID。获取方法如下：<li>方法一：登录网络控制台界面，在虚拟私有云的详情页面查找VPC ID。<li>方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94&vid=88">查询VPC列表</a> vpcID字段获取。  */
	SubnetId        string   `json:"subnetId,omitempty"`        /*  VPC子网ID。获取方法如下：<li>方法一：登录网络控制台界面，单击VPC下的子网，进入子网详情页面，查找子网ID。<li>方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94&vid=88">查询子网列表</a> subnetID字段获取。  */
	SecurityGroupId string   `json:"securityGroupId,omitempty"` /*  安全组ID。获取方法如下：<li>方法一：登录网络控制台界面，在安全组的详情页面查找安全组ID。<li>方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/searchCtapi/ctApiDebug?product=18&api=4817&vid=88">查询用户安全组列表</a> id字段获取。  */
	InstanceNum     int32    `json:"instanceNum,omitempty"`     /*  购买数量(1-100，默认1)  */
	EnableIpv6      *bool    `json:"enableIpv6"`                /*  是否启用IPv6，默认为false。<li>true：启用IPv6。<li>false：不启用IPv6，默认值。  */
	PlainPort       int32    `json:"plainPort,omitempty"`       /*  公共接入点(PLAINTEXT)端口，范围在8000到9100之间，默认为8090。  */
	SaslPort        int32    `json:"saslPort,omitempty"`        /*  安全接入点(SASL_PLAINTEXT)端口，范围在8000到9100之间，默认为8092。  */
	SslPort         int32    `json:"sslPort,omitempty"`         /*  SSL接入点(SASL_SSL)端口，范围在8000到9100之间，默认为8098。  */
	HttpPort        int32    `json:"httpPort,omitempty"`        /*  HTTP接入点端口，范围在8000到9100之间，默认为8082。  */
	RetentionHours  int32    `json:"retentionHours,omitempty"`  /*  实例消息保留时长，默认为72小时，可选1~10000小时。  */
}

type CtgkafkaCreatePostPayOrderResponse struct {
	StatusCode string                                       `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                       `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaCreatePostPayOrderReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaCreatePostPayOrderReturnObjResponse struct {
	Data CtgkafkaCreatePostPayOrderReturnObjResponseData `json:"data,omitempty"` /*  返回数据。  */
}

type CtgkafkaCreatePostPayOrderReturnObjResponseData struct {
	Submitted  bool   `json:"submitted"`
	NewOrderId string `json:"newOrderId"`
	NewOrderNo string `json:"newOrderNo"`
	TotalPrice string `json:"totalPrice"`
}
