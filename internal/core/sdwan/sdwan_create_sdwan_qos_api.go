package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanQosApi
/* 增加qos */
type SdwanCreateSdwanQosApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanQosApi(client *core.CtyunClient) *SdwanCreateSdwanQosApi {
	return &SdwanCreateSdwanQosApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/qos/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanQosApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanQosRequest) (*SdwanCreateSdwanQosResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanQosRequest
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
	var resp SdwanCreateSdwanQosResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanQosRequest struct {
	QosName       string `json:"qosName"`       /*  qos策略名称  */
	ProjectID     string `json:"projectID"`     /*  企业项目名称  */
	Description   string `json:"description"`   /*  描述  */
	Bandwidth     string `json:"bandwidth"`     /*  带宽峰值  */
	BandwidthType string `json:"bandwidthType"` /*  本参数表示带宽类型<br/><br/>取值范围:<br/>internet:互联网带宽<br/>sdwan:SD-WAN带宽  */
}

type SdwanCreateSdwanQosResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
