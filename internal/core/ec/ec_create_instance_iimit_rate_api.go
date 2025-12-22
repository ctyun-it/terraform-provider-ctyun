package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcCreateInstanceIimitRateApi
/* 创建实例限速带宽 */
type EcCreateInstanceIimitRateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcCreateInstanceIimitRateApi(client *core.CtyunClient) *EcCreateInstanceIimitRateApi {
	return &EcCreateInstanceIimitRateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/instance-limit-rate/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcCreateInstanceIimitRateApi) Do(ctx context.Context, credential core.Credential, req *EcCreateInstanceIimitRateRequest) (*EcCreateInstanceIimitRateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcCreateInstanceIimitRateRequest
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
	var resp EcCreateInstanceIimitRateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcCreateInstanceIimitRateRequest struct {
	IbpID      string `json:"ibpID"`      /*  实例带宽包ID  */
	Bandwidth  int32  `json:"bandwidth"`  /*  带宽，单位MB  */
	SrcIpCidr  string `json:"srcIpCidr"`  /*  VPC租户网段或主机网段，ipv4 or ipv6  */
	DstIpCidr  string `json:"dstIpCidr"`  /*  IDC租户网段或主机网段，ipv4 or ipv6  */
	SrcPort    int32  `json:"srcPort"`    /*  VPC租户端口或主机端口  */
	DstPort    int32  `json:"dstPort"`    /*  IDC租户端口或主机端口  */
	IpProtocol int32  `json:"ipProtocol"` /*  标识数据携带的数据是何种协议，标识传输层地址或协议号，如1代表ICMP，6代表TCP，17代表UDP  */
}

type EcCreateInstanceIimitRateResponse struct {
	StatusCode  *int32                                      `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                     `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcCreateInstanceIimitRateReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcCreateInstanceIimitRateReturnObjResponse struct {
	LimitID *string `json:"limitID"` /*  实例限速带宽ID  */
	Message *string `json:"message"` /*  更新结果  */
}
