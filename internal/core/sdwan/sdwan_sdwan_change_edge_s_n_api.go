package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanChangeEdgeSNApi
/* 迁移sn,sn切换后，旧设备的入云、组网和静态路由规划可平滑迁移到新设备，其他高级配置，如ACL，QOS等，再切换完成后可能会失效，需要冲洗为新设备下发相应规则。请做好业务规划 */
type SdwanSdwanChangeEdgeSNApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanChangeEdgeSNApi(client *core.CtyunClient) *SdwanSdwanChangeEdgeSNApi {
	return &SdwanSdwanChangeEdgeSNApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/change-edge-sn",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanChangeEdgeSNApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanChangeEdgeSNRequest) (*SdwanSdwanChangeEdgeSNResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanChangeEdgeSNRequest
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
	var resp SdwanSdwanChangeEdgeSNResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanChangeEdgeSNRequest struct {
	EdgeID   string  `json:"edgeID"`            /*  智能网关ID  */
	MasterSN string  `json:"masterSN"`          /*  迁移的主edge的sn  */
	SlaveSN  *string `json:"slaveSN,omitempty"` /*  迁移的备edge的sn  */
}

type SdwanSdwanChangeEdgeSNResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanChangeEdgeSNReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanChangeEdgeSNReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
