package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbGwlbDeleteApi
/* 删除网关负载均衡
 */type CtelbGwlbDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbDeleteApi(client *core.CtyunClient) *CtelbGwlbDeleteApi {
	return &CtelbGwlbDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbDeleteApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbDeleteRequest) (*CtelbGwlbDeleteResponse, error) {
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
	var resp CtelbGwlbDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbDeleteRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	ProjectID   string `json:"projectID,omitempty"`   /*  企业项目ID，默认"0"  */
	GwLbID      string `json:"gwLbID,omitempty"`      /*  网关负载均衡ID  */
}

type CtelbGwlbDeleteResponse struct {
	StatusCode  int32                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbDeleteReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbDeleteReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  订单id。  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单编号, 可以为 null。  */
	RegionID      string `json:"regionID,omitempty"`      /*  可用区id。  */
}
