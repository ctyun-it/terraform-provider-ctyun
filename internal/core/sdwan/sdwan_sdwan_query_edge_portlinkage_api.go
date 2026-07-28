package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanQueryEdgePortlinkageApi
/* 查询edge portlinkage */
type SdwanSdwanQueryEdgePortlinkageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanQueryEdgePortlinkageApi(client *core.CtyunClient) *SdwanSdwanQueryEdgePortlinkageApi {
	return &SdwanSdwanQueryEdgePortlinkageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/query-edge-portlinkage",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanQueryEdgePortlinkageApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanQueryEdgePortlinkageRequest) (*SdwanSdwanQueryEdgePortlinkageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanQueryEdgePortlinkageResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanQueryEdgePortlinkageRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanSdwanQueryEdgePortlinkageResponse struct {
	StatusCode  int32                                            `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanQueryEdgePortlinkageReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                          `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanQueryEdgePortlinkageReturnObjResponse struct {
	Result       *SdwanSdwanQueryEdgePortlinkageReturnObjResultResponse `json:"result"`       /*  portlinkage信息  */
	TotalCount   int32                                                  `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                                  `json:"currentCount"` /*  当前页数量  */
}

type SdwanSdwanQueryEdgePortlinkageReturnObjResultResponse struct {
	SerialNumber            *string `json:"serialNumber"`            /*  主SN号  */
	SlaveSerialNumber       *string `json:"slaveSerialNumber"`       /*  备SN号  */
	HeartBeatInterface      *string `json:"heartBeatInterface"`      /*  主盒子接口名称  */
	SlaveHeartBeatInterface *string `json:"slaveHeartBeatInterface"` /*  备盒子接口名称  */
}
