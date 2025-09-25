package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsDeleteProtocolServiceApi
/* 删除协议服务
 */type HpfsDeleteProtocolServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsDeleteProtocolServiceApi(client *core.CtyunClient) *HpfsDeleteProtocolServiceApi {
	return &HpfsDeleteProtocolServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/delete-protocol-service",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsDeleteProtocolServiceApi) Do(ctx context.Context, credential core.Credential, req *HpfsDeleteProtocolServiceRequest) (*HpfsDeleteProtocolServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsDeleteProtocolServiceRequest
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
	var resp HpfsDeleteProtocolServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsDeleteProtocolServiceRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID  */
	ProtocolServiceID string `json:"protocolServiceID,omitempty"` /*  协议服务唯一ID  */
}

type HpfsDeleteProtocolServiceResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为处理中/失败，详见errorCode)  */
	Message     string                                      `json:"message"`     /*  响应描述  */
	Description string                                      `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsDeleteProtocolServiceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                      `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsDeleteProtocolServiceReturnObjResponse struct {
	RegionID  string                                                 `json:"regionID"`  /*  资源所属资源池 ID  */
	Resources []*HpfsDeleteProtocolServiceReturnObjResourcesResponse `json:"resources"` /*  资源明细  */
}

type HpfsDeleteProtocolServiceReturnObjResourcesResponse struct {
	SfsUID            string `json:"sfsUID"`            /*  并行文件内部唯一 ID  */
	ProtocolServiceID string `json:"protocolServiceID"` /*  协议服务唯一ID  */
}
