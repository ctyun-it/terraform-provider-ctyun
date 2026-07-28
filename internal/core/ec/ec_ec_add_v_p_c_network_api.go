package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcAddVPCNetworkApi
/* 添加VPC网络实例 */
type EcEcAddVPCNetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcAddVPCNetworkApi(client *core.CtyunClient) *EcEcAddVPCNetworkApi {
	return &EcEcAddVPCNetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/vpc-instance/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcAddVPCNetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcAddVPCNetworkRequest) (*EcEcAddVPCNetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcAddVPCNetworkRequest
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
	var resp EcEcAddVPCNetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcAddVPCNetworkRequest struct {
	EcID        string                             `json:"ecID"`                  /*  云间高速实例ID  */
	CgwID       string                             `json:"cgwID"`                 /*  云网关实例ID  */
	RtbID       string                             `json:"rtbID"`                 /*  路由表ID  */
	DcID        string                             `json:"dcID"`                  /*  资源池ID  */
	VpcID       string                             `json:"vpcID"`                 /*  vpc ID  */
	VpcName     *string                            `json:"vpcName,omitempty"`     /*  云间高速侧显示的vpc名称（建议保持和vpc创建时的名称一致）  */
	ExclusiveID *string                            `json:"exclusiveID,omitempty"` /*  专属云资源池ID  */
	RouteLearn  *int32                             `json:"routeLearn,omitempty"`  /*  路由学习开关，开启后云网关自动学习网络实例路由<br/>取值范围:<br/>1:学习<br/>0:不学习<br/>默认学习  */
	RouteSync   *int32                             `json:"routeSync,omitempty"`   /*  路由同步开关，开启后云网关路由自动同步到网络实例<br/>取值范围:<br/>1:同步<br/>0:不同步<br/>默认同步  */
	Subnets     []*EcEcAddVPCNetworkSubnetsRequest `json:"subnets,omitempty"`     /*  subnet列表  */
}

type EcEcAddVPCNetworkSubnetsRequest struct {
	SubnetID   string `json:"subnetID"`   /*  子网ID, 填写错误会导致异步加载失败，导致回滚  */
	IPVersion  string `json:"IPVersion"`  /*  子网类型<br/>取值范围:<br/>"IPV4":IPv4类型<br/>"IPV6":IPv6类型  */
	CIDR       string `json:"CIDR"`       /*  子网CIDR  */
	SubnetName string `json:"subnetName"` /*  子网名称  */
}

type EcEcAddVPCNetworkResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                             `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcAddVPCNetworkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcEcAddVPCNetworkReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /*  操作日志id  */
}
