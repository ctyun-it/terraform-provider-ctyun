package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type AmqpInstancesCreatePostPayOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstancesCreatePostPayOrderApi(client *ctyunsdk.CtyunClient) *AmqpInstancesCreatePostPayOrderApi {
	return &AmqpInstancesCreatePostPayOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/createPostPayOrder",
		},
	}
}

func (this *AmqpInstancesCreatePostPayOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstancesCreatePostPayOrderRequest) (res *AmqpInstancesCreatePostPayOrderResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddHeader("regionId", req.RegionId)
	builder.AddHeader("projectId", req.ProjectId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpInstancesCreatePostPayOrderResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstancesCreatePostPayOrderRequest struct {
	RegionId        string   `json:"regionId,omitempty"`        /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProjectId       string   `json:"projectId,omitempty"`       /*  企业项目ID(默认值：0)。您可以通过 <a href="https://www.ctyun.cn/document/10017248/10017965">查看企业项目资源</a> 获取企业项目ID。  */
	ClusterName     string   `json:"clusterName,omitempty"`     /*  实例名称。<br>规则：长度4~40个字符，大小写字母开头，只能包含大小写字母、数字及分隔符(-)，大小写字母或数字结尾，实例名称不可重复。  */
	SpecName        string   `json:"specName,omitempty"`        /*  实例的规格类型，资源池所具备的规格可通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=55&api=20202&data=39&isNormal=1&vid=38">查询产品规格</a>接口获取，集群版可选如下：<br>计算增强型的规格可选为：<br>- rabbitmq.2u4g.cluster<br>- rabbitmq.4u8g.cluster<br>- rabbitmq.8u16g.cluster<br>- rabbitmq.12u24g.cluster<br>- rabbitmq.16u32g.cluster<br>- rabbitmq.24u48g.cluster<br>- rabbitmq.32u64g.cluster<br>- rabbitmq.48u96g.cluster<br>- rabbitmq.64u128g.cluster <br>海光-计算增强型的规格可选为：<br>- rabbitmq.hg.2u4g.cluster<br>- rabbitmq.hg.4u8g.cluster<br>- rabbitmq.hg.8u16g.cluster<br>- rabbitmq.hg.16u32g.cluster<br>- rabbitmq.hg.32u64g.cluster <br>鲲鹏-计算增强型的规格可选为：<br>- rabbitmq.kp.2u4g.cluster<br>- rabbitmq.kp.4u8g.cluster<br>- rabbitmq.kp.8u16g.cluster<br>- rabbitmq.kp.16u32g.cluster<br>- rabbitmq.kp.32u64g.cluster  */
	NodeNum         int32    `json:"nodeNum,omitempty"`         /*  节点数。设置1为单机版，设置3、5、7、9为集群版。  */
	ZoneList        []string `json:"zoneList"`                  /*  实例所在可用区信息。只能填一个（单可用区）或三个（多可用区），可用区信息可调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&isNormal=1&vid=81">资源池可用区查询</a>API接口查询。  */
	DiskType        string   `json:"diskType,omitempty"`        /*  磁盘类型，资源池所具备的磁盘类型可通过查询产品规格接口获取，默认取值：<br>- SAS：高IO<br>- SSD：超高IO<br>- FAST-SSD：极速型SSD  */
	DiskSize        int32    `json:"diskSize,omitempty"`        /*  单个节点的磁盘存储空间，单位为GB，存储空间取值范围100GB ~ 10000，并且为100的倍数。实例总存储空间为diskSize * nodeNum。  */
	VpcId           string   `json:"vpcId,omitempty"`           /*  VPC网络ID。获取方法如下：<br>- 方法一：登录网络控制台界面，在虚拟私有云的详情页面查找VPC ID。<br>- 方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94&vid=88">查询VPC列表</a> vpcID字段获取。  */
	SubnetId        string   `json:"subnetId,omitempty"`        /*  VPC子网ID。获取方法如下：<br>- 方法一：登录网络控制台界面，单击VPC下的子网，进入子网详情页面，查找子网ID。<br>- 方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94&vid=88">查询子网列表</a> subnetID字段获取。  */
	SecurityGroupId string   `json:"securityGroupId,omitempty"` /*  安全组ID。获取方法如下：<br>- 方法一：登录网络控制台界面，在安全组的详情页面查找安全组ID。<br>- 方法二：您可以通过 <a href="https://eop.ctyun.cn/ebp/searchCtapi/ctApiDebug?product=18&api=4817&vid=88">查询用户安全组列表</a> id字段获取。  */
	EnableIpv6      *bool    `json:"enableIpv6"`                /*  是否启用IPv6，默认为false。<br>- true：启用IPv6。<br>- false：不启用IPv6，默认值。  */
}

type AmqpInstancesCreatePostPayOrderResponse struct {
	ReturnObj  *AmqpInstancesCreatePostPayOrderResponseReturnObj `json:"returnObj"`
	Message    string                                            `json:"message"`
	StatusCode string                                            `json:"statusCode"`
}

type AmqpInstancesCreatePostPayOrderResponseReturnObj struct {
	Data AmqpInstancesCreatePostPayOrderResponseReturnObjData `json:"data"`
}

type AmqpInstancesCreatePostPayOrderResponseReturnObjData struct {
	Submitted  bool   `json:"submitted"`
	NewOrderId string `json:"newOrderId"`
	NewOrderNo string `json:"newOrderNo"`
}
