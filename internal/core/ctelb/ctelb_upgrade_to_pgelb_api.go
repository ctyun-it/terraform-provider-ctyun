package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbUpgradeToPgelbApi
/* 升级为保障型负载均衡
 */type CtelbUpgradeToPgelbApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbUpgradeToPgelbApi(client *core.CtyunClient) *CtelbUpgradeToPgelbApi {
	return &CtelbUpgradeToPgelbApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/upgrade-to-pgelb",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbUpgradeToPgelbApi) Do(ctx context.Context, credential core.Credential, req *CtelbUpgradeToPgelbRequest) (*CtelbUpgradeToPgelbResponse, error) {
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
	var resp CtelbUpgradeToPgelbResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbUpgradeToPgelbRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	ElbID       string `json:"elbID,omitempty"`       /*  负载均衡 ID  */
	SlaName     string `json:"slaName,omitempty"`     /*  lb的规格名称, 支持:elb.s2.small，elb.s3.small，elb.s4.small，elb.s5.small，elb.s2.large，elb.s3.large，elb.s4.large，elb.s5.large  */
	CycleType   string `json:"cycleType,omitempty"`   /*  订购类型：month（包月） / year（包年）  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
}

type CtelbUpgradeToPgelbResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbUpgradeToPgelbReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbUpgradeToPgelbReturnObjResponse struct {
	MasterOrderID        string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     string `json:"masterResourceID,omitempty"`     /*  资源 ID 可以为 null。  */
	RegionID             string `json:"regionID,omitempty"`             /*  可用区id。  */
}
