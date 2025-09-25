package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVpcVpcpeerBindLabelApi
/* 对等链接绑定标签
 */type CtvpcVpcVpcpeerBindLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpcVpcpeerBindLabelApi(client *core.CtyunClient) *CtvpcVpcVpcpeerBindLabelApi {
	return &CtvpcVpcVpcpeerBindLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/vpcpeer/bind-label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpcVpcpeerBindLabelApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpcVpcpeerBindLabelRequest) (*CtvpcVpcVpcpeerBindLabelResponse, error) {
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
	var resp CtvpcVpcVpcpeerBindLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpcVpcpeerBindLabelRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  区域ID  */
	VpcPeerID  string `json:"vpcPeerID,omitempty"`  /*  对等链接 ID  */
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签 key  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签 取值  */
}

type CtvpcVpcVpcpeerBindLabelResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVpcVpcpeerBindLabelReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVpcVpcpeerBindLabelReturnObjResponse struct {
	VpcPeerID *string `json:"vpcPeerID,omitempty"` /*  对等链接 ID  */
}
