package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupQueryInstanceListApi
/* 查询伸缩组内云主机的列表，并列出云主机的信息<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingGroupQueryInstanceListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupQueryInstanceListApi(client *core.CtyunClient) *ScalingGroupQueryInstanceListApi {
	return &ScalingGroupQueryInstanceListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/query-instance-list",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupQueryInstanceListApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupQueryInstanceListRequest) (*ScalingGroupQueryInstanceListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupQueryInstanceListRequest
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
	var resp ScalingGroupQueryInstanceListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupQueryInstanceListRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID  int64  `json:"groupID,omitempty"`  /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码  */
	PageSize int32  `json:"pageSize,omitempty"` /*  分页查询时设置的每页行数，取值范围:[1~100]，默认值为10  */
}

type ScalingGroupQueryInstanceListResponse struct {
	StatusCode  int32                                           `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupQueryInstanceListReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                          `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingGroupQueryInstanceListReturnObjResponse struct {
	InstanceList []*ScalingGroupQueryInstanceListReturnObjInstanceListResponse `json:"instanceList"` /*  云主机信息列表  */
	TotalCount   int32                                                         `json:"totalCount"`   /*  总计  */
}

type ScalingGroupQueryInstanceListReturnObjInstanceListResponse struct {
	InstanceID    string `json:"instanceID"`    /*  云主机ID  */
	RegionID      string `json:"regionID"`      /*  资源池ID  */
	GroupID       int64  `json:"groupID"`       /*  伸缩组ID  */
	ZabbixName    string `json:"zabbixName"`    /*  【Deprecated】监控设备ID  */
	ProjectIDEcs  string `json:"projectIDEcs"`  /*  企业项目ID  */
	CreateDate    string `json:"createDate"`    /*  创建时间  */
	Id            int32  `json:"id"`            /*  实例ID  */
	Status        int32  `json:"status"`        /*  伸缩活动状态。<br>取值范围：<br>1：已启用。<br>2：正在移入。<br>3：正在移出。  */
	InstanceName  string `json:"instanceName"`  /*  云主机名称  */
	ExecutionMode int32  `json:"executionMode"` /*  执行方式。<br>取值范围：<br>1：自动执行策略。<br>2：手动执行策略。<br>3：手动移入实例。<br>4：手动移出实例。<br>5：新建伸缩组满足最小数。<br>6：修改伸缩组满足最大最小限制。<br>7：健康检查移入。<br>8：健康检查移出。  */
	HealthStatus  int32  `json:"healthStatus"`  /*  健康检查状态。<br>取值范围：<br>1：正常。<br>2：异常。<br>3：初始化。  */
	ConfigName    string `json:"configName"`    /*  伸缩配置名称  */
	ConfigID      int64  `json:"configID"`      /*  伸缩配置ID  */
	ActiveID      int32  `json:"activeID"`      /*  伸缩活动ID  */
	ProtectStatus int32  `json:"protectStatus"` /*  保护状态。<br>取值范围：<br>1：已保护。<br>2：未保护。  */
	JoinDate      string `json:"joinDate"`      /*  加入时间  */
	AzName        string `json:"azName"`        // 可用区
}
