package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsCreateOrderEbsSnapApi
/* 调用本接口为某个云硬盘创建快照订单，即开通快照服务。
 */type EbsCreateOrderEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsCreateOrderEbsSnapApi(client *core.CtyunClient) *EbsCreateOrderEbsSnapApi {
	return &EbsCreateOrderEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/create-order-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsCreateOrderEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsCreateOrderEbsSnapRequest) (*EbsCreateOrderEbsSnapResponse, error) {
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
	var resp EbsCreateOrderEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsCreateOrderEbsSnapRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。若不传clientToken 后续将无法准确获取快照服务是否成功开通。  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID。  */
	DiskID      string  `json:"diskID,omitempty"`      /*  云硬盘ID。  */
}

type EbsCreateOrderEbsSnapResponse struct {
	StatusCode  int32                                     `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                   `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                   `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsCreateOrderEbsSnapReturnObjResponse   `json:"returnObj"`             /*  返回数据结构体。  */
	ErrorCode   *string                                   `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。参考错误码。  */
	Error       *string                                   `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。参考错误码。  */
	ErrorDetail *EbsCreateOrderEbsSnapErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，以ebs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
	Details     *string                                   `json:"details,omitempty"`     /*  可忽略。  */
}

type EbsCreateOrderEbsSnapReturnObjResponse struct {
	MasterOrderID     *string `json:"masterOrderID,omitempty"`     /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO     *string `json:"masterOrderNO,omitempty"`     /*  订单号。  */
	RegionID          *string `json:"regionID,omitempty"`          /*  资源所属资源池ID。  */
	MasterOrderStatus *string `json:"masterOrderStatus,omitempty"` /*  快照服务开通状态（成功/失败）。  */
}

type EbsCreateOrderEbsSnapErrorDetailResponse struct{}
