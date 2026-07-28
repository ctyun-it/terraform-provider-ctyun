package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanEdgeActiveApi
/* edge激活 */
type SdwanEdgeActiveApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanEdgeActiveApi(client *core.CtyunClient) *SdwanEdgeActiveApi {
	return &SdwanEdgeActiveApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-active/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanEdgeActiveApi) Do(ctx context.Context, credential core.Credential, req *SdwanEdgeActiveRequest) (*SdwanEdgeActiveResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanEdgeActiveRequest
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
	var resp SdwanEdgeActiveResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanEdgeActiveRequest struct {
	EdgeID             string   `json:"edgeID"`                       /*  edge的ID  */
	EdgeCIDR           []string `json:"edgeCIDR"`                     /*  edge的子网  ，值类型为string  */
	DeployMode         string   `json:"deployMode"`                   /*  本参数表示接入方式<br/>取值范围：<br/>inline-mode:串接<br/>dual-arm-mode:双臂旁挂<br/>single-arm-mode:单臂旁挂  */
	DhcpDisable        *string  `json:"dhcpDisable,omitempty"`        /*  串联模式下必填，是否关闭dhcp,默认值为false  */
	LanControlIP       *string  `json:"lanControlIP,omitempty"`       /*  lan侧管理IP  */
	LanBusinessIP      *string  `json:"lanBusinessIP,omitempty"`      /*  lan侧业务IP  */
	LanBusinessIPv6    *string  `json:"lanBusinessIPv6,omitempty"`    /*  lan侧业务IPv6  */
	LanControlIPv6     *string  `json:"lanControlIPv6,omitempty"`     /*  lan侧管理IPv6  */
	LanVrrpIP          *string  `json:"lanVrrpIP,omitempty"`          /*  lan侧业务虚IP  */
	LanVrrpIPv6        *string  `json:"lanVrrpIPv6,omitempty"`        /*  lan侧业务虚IPV6  */
	SlaveLanControlIP  *string  `json:"slaveLanControlIP,omitempty"`  /*  备用设备lan侧管理IP  */
	SlaveLanBusinessIP *string  `json:"slaveLanBusinessIP,omitempty"` /*  备用设备lan侧业务IP  */
	DcID               *string  `json:"dcID,omitempty"`               /*  资源池ID  */
}

type SdwanEdgeActiveResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
