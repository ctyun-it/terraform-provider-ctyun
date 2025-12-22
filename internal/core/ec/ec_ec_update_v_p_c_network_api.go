package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcUpdateVPCNetworkApi
/* 修改VPC网络实例 */
type EcEcUpdateVPCNetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcUpdateVPCNetworkApi(client *core.CtyunClient) *EcEcUpdateVPCNetworkApi {
	return &EcEcUpdateVPCNetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/vpc-instance/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcUpdateVPCNetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcUpdateVPCNetworkRequest) (*EcEcUpdateVPCNetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcUpdateVPCNetworkRequest
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
	var resp EcEcUpdateVPCNetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcUpdateVPCNetworkRequest struct {
	EcID    string                                `json:"ecID"`    /*  云间高速实例ID  */
	VpcID   string                                `json:"vpcID"`   /*  vpc ID  */
	Subnets []*EcEcUpdateVPCNetworkSubnetsRequest `json:"subnets"` /*  subnet列表  */
}

type EcEcUpdateVPCNetworkSubnetsRequest struct {
	SubnetID   string `json:"subnetID"`   /*  子网ID, 填写错误会导致异步加载失败，导致回滚  */
	IPVersion  string `json:"IPVersion"`  /*  子网类型<br/>取值范围:<br/>"IPV4":IPv4类型<br/>"IPV6":IPv6类型  */
	CIDR       string `json:"CIDR"`       /*  子网CIDR  */
	SubnetName string `json:"subnetName"` /*  子网名称  */
}

type EcEcUpdateVPCNetworkResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcUpdateVPCNetworkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcUpdateVPCNetworkReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /*  操作日志id  */
}
