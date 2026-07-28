package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanConfigEdgePortlinkageApi
/* 配置edge portlinkage */
type SdwanSdwanConfigEdgePortlinkageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanConfigEdgePortlinkageApi(client *core.CtyunClient) *SdwanSdwanConfigEdgePortlinkageApi {
	return &SdwanSdwanConfigEdgePortlinkageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/config-edge-portlinkage",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanConfigEdgePortlinkageApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanConfigEdgePortlinkageRequest) (*SdwanSdwanConfigEdgePortlinkageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanConfigEdgePortlinkageRequest
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
	var resp SdwanSdwanConfigEdgePortlinkageResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanConfigEdgePortlinkageRequest struct {
	EdgeID             string `json:"edgeID"`             /*  edge的ID  */
	InterfaceName      string `json:"interfaceName"`      /*  edge接口名称  */
	SlaveInterfaceName string `json:"slaveInterfaceName"` /*  edge备接口名称  */
}

type SdwanSdwanConfigEdgePortlinkageResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	Result      *string `json:"result"`      /*  portlinkage配置结果  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
