package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateEdgeSNApi
/* 录入智能网关SN信息 */
type SdwanUpdateEdgeSNApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateEdgeSNApi(client *core.CtyunClient) *SdwanUpdateEdgeSNApi {
	return &SdwanUpdateEdgeSNApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-sn/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateEdgeSNApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateEdgeSNRequest) (*SdwanUpdateEdgeSNResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateEdgeSNRequest
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
	var resp SdwanUpdateEdgeSNResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateEdgeSNRequest struct {
	RegionID         string                                    `json:"regionID"`                   /*  资源池ID  */
	EdgeID           string                                    `json:"edgeID"`                     /*  智能网关ID  */
	SerialNumberDict *SdwanUpdateEdgeSNSerialNumberDictRequest `json:"serialNumberDict,omitempty"` /*  SN结构  */
}

type SdwanUpdateEdgeSNSerialNumberDictRequest struct {
	MasterSN string  `json:"masterSN"`          /*  主edge的sn  */
	SlaveSN  *string `json:"slaveSN,omitempty"` /*  备edge的sn  */
}

type SdwanUpdateEdgeSNResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                               `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                               `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                               `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanUpdateEdgeSNReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanUpdateEdgeSNReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
