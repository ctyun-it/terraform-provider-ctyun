package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupUpdateApi
/* 修改一个弹性伸缩组<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingGroupUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupUpdateApi(client *core.CtyunClient) *ScalingGroupUpdateApi {
	return &ScalingGroupUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/update",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupUpdateApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupUpdateRequest) (*ScalingGroupUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupUpdateRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ScalingGroupUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupUpdateRequest struct {
	RegionID            string                             `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID             int64                              `json:"groupID,omitempty"`         /*  伸缩组ID  */
	Name                string                             `json:"name,omitempty"`            /*  伸缩组名称  */
	MinCount            int32                              `json:"minCount,omitempty"`        /*  最小云主机数，取值范围：[0,50]  */
	MaxCount            int32                              `json:"maxCount,omitempty"`        /*  最大云主机数，取值范围：[minCount,2147483647]  */
	ExpectedCount       *int32                             `json:"expectedCount,omitempty"`   /*  期望云主机数，取值范围：[minCount,maxCount]，非多可用区资源池不支持该参数  */
	UseLb               *int32                             `json:"useLb,omitempty"`           /*  是否使用负载均衡，1：是 2：否  */
	HealthPeriod        int32                              `json:"healthPeriod,omitempty"`    /*  健康检查时间间隔（周期），单位：秒，取值范围：[300,10080]  */
	SecurityGroupIDList []string                           `json:"securityGroupIDList"`       /*  安全组ID列表，非多可用区资源池不使用该参数，其安全组参数在弹性伸缩配置中填写。您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	LbList              []*ScalingGroupUpdateLbListRequest `json:"lbList,omitempty"`          /*  负载均衡列表  */
	SubnetIDList        []string                           `json:"subnetIDList,omitempty"`    /*  子网ID列表。支持一主多辅，最大支持输入5个网卡信息，顺序第一个网卡信息默认主网卡。mazInfo和subnetIDList参数互斥。修改subnetIDList会将mazInfo信息清空。您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94">查询子网列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4812&data=94">创建子网</a>  */
	MoveOutStrategy     int32                              `json:"moveOutStrategy,omitempty"` /*  实例移出策略。<br>取值范围：<br>1：较早创建的配置较早创建的云主机<br>2：较晚创建的配置较晚创建的云主机<br>3：较早创建的云主机<br>4：较晚创建的云主机  */
	RecoveryMode        int32                              `json:"recoveryMode,omitempty"`    /*  实例回收模式。<br>取值范围：<br>1：释放模式。<br>2：停机回收模式。  */
	HealthMode          int32                              `json:"healthMode,omitempty"`      /*  健康检查方式。<br>取值范围：<br>1：云服务器健康检查<br>2：弹性负载均衡健康检查  */
	ConfigID            int32                              `json:"configID,omitempty"`        /*  【Deprecated】伸缩配置ID  */
	ConfigList          []int32                            `json:"configList"`                /*  伸缩配置ID列表，最大支持传入10个伸缩配置。按输入伸缩配置的顺序，决定伸缩配置优先级。<br/>注意：该参数与configID不可同时传入，请尽量选择本参数。 <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5068&data=93">查询弹性伸缩配置</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4995&data=93">创建一个弹性伸缩配置</a>  */
	AzStrategy          int32                              `json:"azStrategy,omitempty"`      /*  扩容策略类型，仅多可用区资源池支持。<br>取值范围：<br>1：均衡分布。<br> 2：优先级分布。  */
}

type ScalingGroupUpdateLbListRequest struct {
	Port        int32  `json:"port,omitempty"`        /*  端口号  */
	Id          string `json:"id,omitempty"`          /*  负载均衡ID  */
	Weight      int32  `json:"weight,omitempty"`      /*  权重  */
	HostGroupID string `json:"hostGroupID,omitempty"` /*  后端主机组ID  */
}

type ScalingGroupUpdateResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                               `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                               `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupUpdateReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingGroupUpdateReturnObjResponse struct {
	GroupID int32 `json:"groupID"` /*  伸缩组ID  */
}
