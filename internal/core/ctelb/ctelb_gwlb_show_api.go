package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbGwlbShowApi
/* 查看网关负载均衡详情
 */type CtelbGwlbShowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbShowApi(client *core.CtyunClient) *CtelbGwlbShowApi {
	return &CtelbGwlbShowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbShowApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbShowRequest) (*CtelbGwlbShowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	ctReq.AddParam("gwLbID", req.GwLbID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbGwlbShowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbShowRequest struct {
	RegionID  string /*  资源池 ID  */
	ProjectID string /*  企业项目ID，默认"0"  */
	GwLbID    string /*  网关负载均衡ID  */
}

type CtelbGwlbShowResponse struct {
	StatusCode  int32                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbShowReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbShowReturnObjResponse struct {
	GwLbID           string `json:"gwLbID,omitempty"`           /*  网关负载均衡 ID  */
	Name             string `json:"name,omitempty"`             /*  名字  */
	Description      string `json:"description,omitempty"`      /*  描述  */
	VpcID            string `json:"vpcID,omitempty"`            /*  虚拟私有云 ID  */
	SubnetID         string `json:"subnetID,omitempty"`         /*  子网 ID  */
	PortID           string `json:"portID,omitempty"`           /*  网卡 ID  */
	Ipv6Enabled      *bool  `json:"ipv6Enabled"`                /*  是否开启 ipv6  */
	PrivateIpAddress string `json:"privateIpAddress,omitempty"` /*  私有 IP 地址  */
	Ipv6Address      string `json:"ipv6Address,omitempty"`      /*  ipv6 地址  */
	SlaName          string `json:"slaName,omitempty"`          /*  规格  */
	DeleteProtection *bool  `json:"deleteProtection"`           /*  是否开启删除保护  */
	CreatedAt        string `json:"createdAt,omitempty"`        /*  创建时间  */
	UpdatedAt        string `json:"updatedAt,omitempty"`        /*  更新时间  */
}
