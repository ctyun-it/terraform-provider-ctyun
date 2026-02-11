package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsEcsUpdateAutoRenewConfigApi
/* 对原包周期付费类型（包年包月）的云主机自动续订配置进行修改，包括开启自动续订或关闭自动续订。<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;若实例到期前10天内系统已下发自动续订任务，此时取消自动续订操作仅对后续周期生效。具体表现为：当前已触发的续订任务仍将正常执行完成，而自动续订功能将在实例下一次到期时正式停用。 */
type CtecsEcsUpdateAutoRenewConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsEcsUpdateAutoRenewConfigApi(client *core.CtyunClient) *CtecsEcsUpdateAutoRenewConfigApi {
	return &CtecsEcsUpdateAutoRenewConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/update-auto-renew-config",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsEcsUpdateAutoRenewConfigApi) Do(ctx context.Context, credential core.Credential, req *CtecsEcsUpdateAutoRenewConfigRequest) (*CtecsEcsUpdateAutoRenewConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtecsEcsUpdateAutoRenewConfigRequest
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
	var resp CtecsEcsUpdateAutoRenewConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtecsEcsUpdateAutoRenewConfigRequest struct {
	RegionID            string  `json:"regionID"`                      /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceIDList      string  `json:"instanceIDList"`                /*  云主机ID列表，多台使用英文逗号分割，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	AutoRenewStatus     int32   `json:"autoRenewStatus"`               /*  是否自动续订，取值范围：<br />0（不续费），<br />1（自动续费），<br />注：如填写值为1，而不填写autoRenewCycleType和autoRenewCycleCount，则默认开启自动续订且自动续订周期为1个月；如填写值为0，则autoRenewCycleType和autoRenewCycleCount填写无效  */
	AutoRenewCycleType  *string `json:"autoRenewCycleType,omitempty"`  /*  自动续订周期类型，取值范围：<br />MONTH：按月，<br />YEAR：按年。<br />注：不填默认值（缺省值）为MONTH  */
	AutoRenewCycleCount *int32  `json:"autoRenewCycleCount,omitempty"` /*  订购时长，该参数需要与cycleType一同使用<br />注：当autoRenewCycleType为MONTH时，取值范围为[1, 11]；当autoRenewCycleType为YEAR时，取值范围为[1, 5]。即按月最大可自动续订11个月，按年最大自动续订5年  */
}

type CtecsEcsUpdateAutoRenewConfigResponse struct {
	StatusCode  *int32                                          `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   *string                                         `json:"errorCode"`   /*  错误码，为product.module.code三段式码  */
	Error       *string                                         `json:"error"`       /*  错误码，为product.module.code三段式码  */
	Message     *string                                         `json:"message"`     /*  英文描述信息  */
	Description *string                                         `json:"description"` /*  中文描述信息  */
	ReturnObj   *CtecsEcsUpdateAutoRenewConfigReturnObjResponse `json:"returnObj"`   /*  成功时返回，通常为空  */
}

type CtecsEcsUpdateAutoRenewConfigReturnObjResponse struct{}
