package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupInstanceMonitorApi
/* 获取弹性伸缩云主机数量监控数据<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingGroupInstanceMonitorApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupInstanceMonitorApi(client *core.CtyunClient) *ScalingGroupInstanceMonitorApi {
	return &ScalingGroupInstanceMonitorApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/instance-monitor",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupInstanceMonitorApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupInstanceMonitorRequest) (*ScalingGroupInstanceMonitorResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupInstanceMonitorRequest
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
	var resp ScalingGroupInstanceMonitorResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupInstanceMonitorRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID   int32  `json:"groupID,omitempty"`   /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	TimeRange int32  `json:"timeRange,omitempty"` /*  时间范围，默认值为60分钟  */
	TimeFrom  string `json:"timeFrom,omitempty"`  /*  开始时间  */
	TimeTill  string `json:"timeTill,omitempty"`  /*  结束时间  */
	Period    int32  `json:"period,omitempty"`    /*  监控周期，单位：分钟  */
}

type ScalingGroupInstanceMonitorResponse struct {
	StatusCode  int32                                         `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                        `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                        `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                        `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupInstanceMonitorReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                        `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingGroupInstanceMonitorReturnObjResponse struct {
	VList []*ScalingGroupInstanceMonitorReturnObjVListResponse `json:"vList"` /*  弹性伸缩云主机数量监控列表, 类型元素是dict，包括时间戳和云主机数量  */
}

type ScalingGroupInstanceMonitorReturnObjVListResponse struct {
	Timestamp int32 `json:"timestamp"` /*  时间戳  */
	Value     int32 `json:"value"`     /*  监控值  */
}
