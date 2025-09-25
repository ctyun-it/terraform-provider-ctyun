package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreatePgelbApi
/* 保障型负载均衡创建
 */type CtelbCreatePgelbApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreatePgelbApi(client *core.CtyunClient) *CtelbCreatePgelbApi {
	return &CtelbCreatePgelbApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-pgelb",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreatePgelbApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreatePgelbRequest) (*CtelbCreatePgelbResponse, error) {
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
	var resp CtelbCreatePgelbResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreatePgelbRequest struct {
	ClientToken      string `json:"clientToken,omitempty"`      /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID         string `json:"regionID,omitempty"`         /*  区域ID  */
	ProjectID        string `json:"projectID,omitempty"`        /*  企业项目 ID，默认为'0'  */
	VpcID            string `json:"vpcID,omitempty"`            /*  vpc的ID  */
	SubnetID         string `json:"subnetID,omitempty"`         /*  子网的ID  */
	Name             string `json:"name,omitempty"`             /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	EipID            string `json:"eipID,omitempty"`            /*  弹性公网IP的ID。当resourceType=external为必填  */
	SlaName          string `json:"slaName,omitempty"`          /*  lb的规格名称, 支持:elb.s2.small，elb.s3.small，elb.s4.small，elb.s5.small，elb.s2.large，elb.s3.large，elb.s4.large，elb.s5.large  */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型。internal：内网负载均衡，external：公网负载均衡  */
	PrivateIpAddress string `json:"privateIpAddress,omitempty"` /*  负载均衡的私有IP地址，不指定则自动分配  */
	CycleType        string `json:"cycleType,omitempty"`        /*  订购类型：month（包月） / year（包年） / on_demand （按需)  */
	CycleCount       int32  `json:"cycleCount,omitempty"`       /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年，当 cycleType = on_demand 可以不传  */
	PayVoucherPrice  string `json:"payVoucherPrice,omitempty"`  /*  代金券金额，支持到小数点后两位  */
}

type CtelbCreatePgelbResponse struct {
	StatusCode  int32                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbCreatePgelbReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       string                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbCreatePgelbReturnObjResponse struct {
	MasterOrderID        string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     string `json:"masterResourceID,omitempty"`     /*  资源 ID 可以为 null。  */
	RegionID             string `json:"regionID,omitempty"`             /*  可用区id。  */
	ElbID                string `json:"elbID,omitempty"`                /*  负载均衡 ID  */
}
