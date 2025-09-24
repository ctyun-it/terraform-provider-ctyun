package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbModifyPgelbSpecApi
/* 保障型负载均衡变配
 */type CtelbModifyPgelbSpecApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbModifyPgelbSpecApi(client *core.CtyunClient) *CtelbModifyPgelbSpecApi {
	return &CtelbModifyPgelbSpecApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/modify-pgelb-spec",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbModifyPgelbSpecApi) Do(ctx context.Context, credential core.Credential, req *CtelbModifyPgelbSpecRequest) (*CtelbModifyPgelbSpecResponse, error) {
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
	var resp CtelbModifyPgelbSpecResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbModifyPgelbSpecRequest struct {
	ClientToken     string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID        string `json:"regionID,omitempty"`        /*  区域ID  */
	ElbID           string `json:"elbID,omitempty"`           /*  负载均衡 ID  */
	SlaName         string `json:"slaName,omitempty"`         /*  lb的规格名称, 支持:elb.s2.small，elb.s3.small，elb.s4.small，elb.s5.small，elb.s2.large，elb.s3.large，elb.s4.large，elb.s5.large  */
	PayVoucherPrice string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位  */
}

type CtelbModifyPgelbSpecResponse struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbModifyPgelbSpecReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbModifyPgelbSpecReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  订单id。  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单编号, 可以为 null。  */
	RegionID      string `json:"regionID,omitempty"`      /*  可用区id。  */
}
