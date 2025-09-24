package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupListApi
/* 查询伸缩组列表<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingGroupListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupListApi(client *core.CtyunClient) *ScalingGroupListApi {
	return &ScalingGroupListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/list",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupListApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupListRequest) (*ScalingGroupListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupListRequest
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
	var resp ScalingGroupListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupListRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID   int64  `json:"groupID,omitempty"`   /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	ProjectID string `json:"projectID,omitempty"` /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目  */
	PageNo    int32  `json:"pageNo,omitempty"`    /*  页码  */
	Page      int32  `json:"page,omitempty"`      /*  【Deprecated】页码  */
	PageSize  int32  `json:"pageSize,omitempty"`  /*  分页查询时设置的每页行数，取值范围:[1~100]，默认值为10  */
}

type ScalingGroupListResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupListReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingGroupListReturnObjResponse struct {
	NumberOfAll   int32                                             `json:"numberOfAll"`   /*  伸缩组总数量  */
	TotalCount    int32                                             `json:"totalCount"`    /*  【Deprecated】本次查询返回的伸缩组数量  */
	ScalingGroups []*ScalingGroupListReturnObjScalingGroupsResponse `json:"scalingGroups"` /*  伸缩组信息列表  */
}

type ScalingGroupListReturnObjScalingGroupsResponse struct {
	ConfigID            int32    `json:"configID"`            /*  【Deprecated】伸缩组配置ID，仅伸缩配置数目为1时有值  */
	ConfigList          []int32  `json:"configList"`          /*  伸缩组配置ID列表  */
	RecoveryMode        int32    `json:"recoveryMode"`        /*  实例回收模式。<br>取值范围：<br>1： 释放模式。<br> 2： 停机回收模式。  */
	HealthPeriod        int32    `json:"healthPeriod"`        /*  健康检查时间间隔（周期），单位：秒，取值范围：[300,10080]  */
	MaxCount            int32    `json:"maxCount"`            /*  最大云主机数，取值范围：[minCount,2147483647]  */
	MinCount            int32    `json:"minCount"`            /*  最小云主机数，取值范围：[0,50]  */
	ExpectedCount       int32    `json:"expectedCount"`       /*  期望云主机数，取值范围：[minCount,maxCount]，非多可用区资源池不支持该参数  */
	MoveOutStrategy     int32    `json:"moveOutStrategy"`     /*  实例移出策略。<br>取值范围：<br>1： 较早创建的配置较早创建的云主机。<br> 2： 较晚创建的配置较晚创建的云主机。<br> 3： 较早创建的云主机。<br> 4： 较晚创建的云主机。  */
	CreateDate          string   `json:"createDate"`          /*  创建时间  */
	GroupID             int64    `json:"groupID"`             /*  伸缩组ID  */
	UpdateDate          string   `json:"updateDate"`          /*  更新时间  */
	HealthMode          int32    `json:"healthMode"`          /*  健康检查方式。<br>取值范围：<br>1： 云服务器健康检查。<br> 2： 弹性负载均衡健康检查。  */
	UseLb               int32    `json:"useLb"`               /*  是否使用负载均衡，1：是 2：否  */
	ZabbixName          string   `json:"zabbixName"`          /*  【Deprecated】监控设备ID  */
	SubnetIDList        []string `json:"subnetIDList"`        /*  子网ID列表  */
	VpcCidr             string   `json:"vpcCidr"`             /*  虚拟私有云网段  */
	Status              int32    `json:"status"`              /*  伸缩组状态。<br>取值范围：<br>1： 启用。<br> 2： 停用。  */
	VpcName             string   `json:"vpcName"`             /*  虚拟私有云名称  */
	InstanceCount       int32    `json:"instanceCount"`       /*  伸缩组包含云主机数量  */
	ProjectIDEcs        string   `json:"projectIDEcs"`        /*  企业项目ID  */
	ConfigName          string   `json:"configName"`          /*  【Deprecated】伸缩配置名称，仅伸缩配置数目为1时有值  */
	Name                string   `json:"name"`                /*  伸缩组名称  */
	SecurityGroupIDList []string `json:"securityGroupIDList"` /*  多可用区资源池安全组ID列表。非多可用区资源池该值为空数组，可在查询弹性伸缩配置接口查询。  */
	VpcID               string   `json:"vpcID"`               /*  虚拟私有云ID  */
	AzStrategy          int32    `json:"azStrategy"`          /*  扩容策略类型，仅多可用区资源池支持。非多可用区资源池该值为空。<br>取值范围：<br>1：均衡分布。<br> 2：优先级分布。  */
	DeleteProtection    *bool    `json:"deleteProtection"`    /*  是否开启伸缩组保护。<br>取值范围：<br>true：开启伸缩组保护，此时不能删除该伸缩组。<br>false：关闭伸缩组保护，可删除该伸缩组。  */
}
