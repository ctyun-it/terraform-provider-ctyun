package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaQuotaClientListApi
/* 查询用户/客户端流控配置。
 */type CtgkafkaQuotaClientListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQuotaClientListApi(client *core.CtyunClient) *CtgkafkaQuotaClientListApi {
	return &CtgkafkaQuotaClientListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/quota/user-client/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaQuotaClientListApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQuotaClientListRequest) (*CtgkafkaQuotaClientListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaQuotaClientListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQuotaClientListRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type CtgkafkaQuotaClientListResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                    `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaQuotaClientListReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaQuotaClientListReturnObjResponse struct {
	User             string `json:"user,omitempty"`             /*  用户名，defaultUser=false时有值  */
	ClientId         string `json:"clientId,omitempty"`         /*  客户端ID，defaultClient=false时有值。  */
	DefaultUser      *bool  `json:"defaultUser"`                /*  是否是默认用户，true表示对所有用户限流。  */
	DefaultClient    *bool  `json:"defaultClient"`              /*  是否是默认客户端，true表示对所有客户端限流。  */
	ProducerByteRate int64  `json:"producerByteRate,omitempty"` /*  生产者流控上限，单位字节。  */
	ConsumerByteRate int64  `json:"consumerByteRate,omitempty"` /*  消费者流控上限，单位字节。  */
}
