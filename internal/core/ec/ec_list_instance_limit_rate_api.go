package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcListInstanceLimitRateApi
/* 查询实例限速带宽 */
type EcListInstanceLimitRateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcListInstanceLimitRateApi(client *core.CtyunClient) *EcListInstanceLimitRateApi {
	return &EcListInstanceLimitRateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/instance-limit-rate/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcListInstanceLimitRateApi) Do(ctx context.Context, credential core.Credential, req *EcListInstanceLimitRateRequest) (*EcListInstanceLimitRateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.IbpID != nil && *req.IbpID != "" {
		ctReq.AddParam("ibpID", *req.IbpID)
	}
	if req.LimitID != nil && *req.LimitID != "" {
		ctReq.AddParam("limitID", *req.LimitID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcListInstanceLimitRateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcListInstanceLimitRateRequest struct {
	IbpID   *string `json:"ibpID,omitempty"`   /*  实例带宽包ID(当实例限速带宽ID[limitID]为空时，必填)  */
	LimitID *string `json:"limitID,omitempty"` /*  实例限速带宽ID(当实例带宽包ID[ibpID]为空时，必填)  */
}

type EcListInstanceLimitRateResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                   `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcListInstanceLimitRateReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                   `json:"error"`       /*  错误码，为product.module.code三段式码，详见错误码说明  */
}

type EcListInstanceLimitRateReturnObjResponse struct {
	CurrentCount *int32                                             `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                             `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                             `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcListInstanceLimitRateReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcListInstanceLimitRateReturnObjResultsResponse struct {
	LimitID    *string `json:"limitID"`    /*  实例限速带宽ID  */
	IbpID      *string `json:"ibpID"`      /*  实例带宽包ID  */
	IbpName    *string `json:"ibpName"`    /*  实例带宽包名字  */
	Bandwidth  *int32  `json:"bandwidth"`  /*  带宽，单位MB  */
	SrcIpCidr  *string `json:"srcIpCidr"`  /*  VPC租户网段或主机网段，ipv4 or ipv6  */
	DstIpCidr  *string `json:"dstIpCidr"`  /*  IDC租户网段或主机网段，ipv4 or ipv6   */
	SrcPort    *int32  `json:"srcPort"`    /*  VPC租户端口或主机端口  */
	DstPort    *int32  `json:"dstPort"`    /*  IDC租户端口或主机端口  */
	IpProtocol *int32  `json:"ipProtocol"` /*  标识数据携带的数据是何种协议，标识传输层地址或协议号，如1代表ICMP，6代表TCP，17代表UDP  */
}
