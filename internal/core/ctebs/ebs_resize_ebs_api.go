package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsResizeEbsApi
/* 当云硬盘空间不足时，您可以扩大云硬盘的容量，也就是云硬盘扩容。此接口也支持随云主机订购的系统盘及数据盘的扩容。
 */ /* 当您扩容成功后，需要将扩容部分的容量划分至原有分区内，或者对扩容部分的云硬盘分配新的分区。详见：
 */ /* <a href="https://www.ctyun.cn/document/10027696/10029076">云硬盘扩容概述</a>
 */type EbsResizeEbsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsResizeEbsApi(client *core.CtyunClient) *EbsResizeEbsApi {
	return &EbsResizeEbsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/resize-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsResizeEbsApi) Do(ctx context.Context, credential core.Credential, req *EbsResizeEbsRequest) (*EbsResizeEbsResponse, error) {
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
	var resp EbsResizeEbsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsResizeEbsRequest struct {
	DiskSize int32 `json:"diskSize"` /*  变配后的云硬盘大小，数据盘的取值范围为：
	●超高IO/高IO/极速型SSD/普通IO：10GB~32768GB
	●XSSD-0：10GB-65536GB
	●XSSD-1：20GB-65536GB
	●XSSD-2：512GB-65536GB
	系统盘的取值范围为：
	●超高IO/高IO/极速型SSD/普通IO：40GB~2048GB
	●XSSD-0：40GB-2048GB
	●XSSD-1：40GB-2048GB
	●XSSD-2：512GB-2048GB  */
	DiskID      string  `json:"diskID,omitempty"`      /*  待扩容的云硬盘ID。  */
	RegionID    *string `json:"regionID,omitempty"`    /*  如本地语境支持保存regionID，那么建议传递。  */
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。<br/>要求单个云平台账户内唯一。  */
}

type EbsResizeEbsResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码(800为成功，900为处理中/失败)。  */
	Message     *string                          `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                          `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsResizeEbsReturnObjResponse   `json:"returnObj"`             /*  返回结构体。  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                          `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
	ErrorDetail *EbsResizeEbsErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。<br> 其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
}

type EbsResizeEbsReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单ID。调用方在拿到masterOrderID之后，<br/>在若干错误情况下，可以使用masterOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单号。  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  主资源ID。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  主资源状态。只有主订单资源会返回。  */
}

type EbsResizeEbsErrorDetailResponse struct {
	BssErrCode       *string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中。  */
	BssErrMsg        *string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中。  */
	BssOrigErr       *string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息。  */
	BssErrPrefixHint *string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息。  */
}
