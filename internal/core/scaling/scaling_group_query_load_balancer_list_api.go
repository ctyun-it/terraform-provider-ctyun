package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupQueryLoadBalancerListApi
/* 查询伸缩组的负载均衡列表<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingGroupQueryLoadBalancerListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupQueryLoadBalancerListApi(client *core.CtyunClient) *ScalingGroupQueryLoadBalancerListApi {
	return &ScalingGroupQueryLoadBalancerListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/query-load-balancer-list",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupQueryLoadBalancerListApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupQueryLoadBalancerListRequest) (*ScalingGroupQueryLoadBalancerListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupQueryLoadBalancerListRequest
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
	var resp ScalingGroupQueryLoadBalancerListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupQueryLoadBalancerListRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID  int64  `json:"groupID,omitempty"`  /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码  */
	PageSize int32  `json:"pageSize,omitempty"` /*  分页查询时设置的每页行数，取值范围:[1~100]，默认值为10  */
}

type ScalingGroupQueryLoadBalancerListResponse struct {
	StatusCode  int32                                               `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupQueryLoadBalancerListReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingGroupQueryLoadBalancerListReturnObjResponse struct {
	TotalCount    int32                                                              `json:"totalCount"`    /*  负载均衡的数量  */
	LoadBalancers []*ScalingGroupQueryLoadBalancerListReturnObjLoadBalancersResponse `json:"loadBalancers"` /*  负载均衡信息列表  */
}

type ScalingGroupQueryLoadBalancerListReturnObjLoadBalancersResponse struct {
	HostGroupName string `json:"hostGroupName"` /*  后端服务组名称  */
	Weight        int32  `json:"weight"`        /*  权重  */
	Port          int32  `json:"port"`          /*  端口号  */
	LbID          string `json:"lbID"`          /*  负载均衡ID  */
	LbName        string `json:"lbName"`        /*  负载均衡名称  */
	HostGroupID   string `json:"hostGroupID"`   /*  后端主机组ID  */
	GroupID       int32  `json:"groupID"`       /*  伸缩组ID  */
}
