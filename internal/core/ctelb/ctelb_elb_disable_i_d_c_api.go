package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbElbDisableIDCApi
/* 关闭IDC
 */type CtelbElbDisableIDCApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbElbDisableIDCApi(client *core.CtyunClient) *CtelbElbDisableIDCApi {
	return &CtelbElbDisableIDCApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/disable-idc",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbElbDisableIDCApi) Do(ctx context.Context, credential core.Credential, req *CtelbElbDisableIDCRequest) (*CtelbElbDisableIDCResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("vpcID", req.VpcID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbElbDisableIDCResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbElbDisableIDCRequest struct {
	RegionID string /*  区域ID  */
	VpcID    string /*  虚拟私有云 ID  */
}

type CtelbElbDisableIDCResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbElbDisableIDCReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbElbDisableIDCReturnObjResponse struct {
	VpcID string `json:"vpcID,omitempty"` /*  虚拟私有云 id  */
}
