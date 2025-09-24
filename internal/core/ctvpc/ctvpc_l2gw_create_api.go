package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwCreateApi
/* 创建l2gw
 */type CtvpcL2gwCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwCreateApi(client *core.CtyunClient) *CtvpcL2gwCreateApi {
	return &CtvpcL2gwCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwCreateApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwCreateRequest) (*CtvpcL2gwCreateResponse, error) {
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
	var resp CtvpcL2gwCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwCreateRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池 ID  */
	Name        string  `json:"name,omitempty"`        /*  二层网关名称，由数字、字母、中文、-、_组成，不能以数字、_和-开头，长度限制2-32个字符  */
	Description *string `json:"description,omitempty"` /*  二层网关描述  */
	VpcID       string  `json:"vpcID,omitempty"`       /*  vpc id  */
	LinkGwType  string  `json:"linkGwType,omitempty"`  /*  隧道连接方式 linegw：云专线  vpn：VPN  */
	LinkGwID    string  `json:"linkGwID,omitempty"`    /*  关联网关  */
	SubnetID    string  `json:"subnetID,omitempty"`    /*  本端隧道子网  */
	Ip          *string `json:"ip,omitempty"`          /*  本端隧道IP  */
	CycleType   string  `json:"cycleType,omitempty"`   /*  订购类型：month（包月） / year（包年） / on_demand（按需）  */
	CycleCount  int32   `json:"cycleCount"`            /*  订购时长，包年和包月时必填  */
	Spec        string  `json:"spec,omitempty"`        /*  规格 STANDARD：标准版  ENHANCED：增强版  BASIC: 基础版  */
	AutoRenew   *bool   `json:"autoRenew"`             /*  是否自动续订  */
}

type CtvpcL2gwCreateResponse struct {
	StatusCode  int32                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcL2gwCreateReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcL2gwCreateReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
	RegionID             *string `json:"regionID,omitempty"`             /*  可用区id。  */
	L2gwID               *string `json:"l2gwID,omitempty"`               /*  ID，当 masterResourceStatus 不为 started, 该值可为 null  */
}
