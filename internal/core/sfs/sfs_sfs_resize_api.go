package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsResizeApi
/* 修改文件系统大小
 */type SfsSfsResizeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsResizeApi(client *core.CtyunClient) *SfsSfsResizeApi {
	return &SfsSfsResizeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/resize",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsResizeApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsResizeRequest) (*SfsSfsResizeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsResizeRequest
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
	var resp SfsSfsResizeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsResizeRequest struct {
	SfsSize     int32  `json:"sfsSize,omitempty"`     /*  变配后的大小，单位 GB  */
	ResourceID  string `json:"resourceID,omitempty"`  /*  文件系统资源 ID  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
}

type SfsSfsResizeResponse struct {
	StatusCode  int32                          `json:"statusCode"`  /*  返回状态码(800为成功，900为订单处理中/失败)  */
	Message     string                         `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SfsSfsResizeReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                         `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
}

type SfsSfsResizeReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID"` /*  订单 ID。调用方在拿到 masterOrderID 之后，在若干错误情况下，可以使用 masterOrderID 进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO"` /*  订单号  */
	RegionID      string `json:"regionID"`      /*  资源所属资源池 ID  */
}
