package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmAttachVolumeApi
/* 物理机挂载卷
 */type EbmAttachVolumeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmAttachVolumeApi(client *core.CtyunClient) *EbmAttachVolumeApi {
	return &EbmAttachVolumeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebm/attach-volume",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmAttachVolumeApi) Do(ctx context.Context, credential core.Credential, req *EbmAttachVolumeRequest) (*EbmAttachVolumeResponse, error) {
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
	var resp EbmAttachVolumeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmAttachVolumeRequest struct {
	RegionID string `json:"regionID"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string `json:"azName"` /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string `json:"instanceUUID"` /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */VolumeID string `json:"volumeID"` /*  云硬盘ID；<br />获取：<br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7338&data=48">查询云硬盘</a>
	 */
}

type EbmAttachVolumeResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmAttachVolumeReturnObjResponse `json:"returnObj"` /*  返回参数，值为空
	 */
}

type EbmAttachVolumeReturnObjResponse struct{}
