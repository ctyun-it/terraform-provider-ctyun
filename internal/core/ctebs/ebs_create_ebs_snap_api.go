package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsCreateEbsSnapApi
/* 云硬盘快照是一种数据备份方式，可以备份或者恢复整个云硬盘的数据，常用于数据备份、制作镜像、应用容灾等场景。
 */ /* 在回滚云硬盘、更换操作系统、数据迁移等重要操作之前，您可以提前创建快照，从而保存指定时刻的云硬盘数据，提高操作的容错率。
 */ /* 注意：
 */ /* 单个租户下的单个云硬盘在首次创建快照时，会通过订单流程创建，之后的快照不再通过订单流程创建。如果租户主动退订快照订单后，再次创建时仍会通过订单流程创建。
 */type EbsCreateEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsCreateEbsSnapApi(client *core.CtyunClient) *EbsCreateEbsSnapApi {
	return &EbsCreateEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/create-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsCreateEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsCreateEbsSnapRequest) (*EbsCreateEbsSnapResponse, error) {
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
	var resp EbsCreateEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsCreateEbsSnapRequest struct {
	ClientToken     *string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。只有通过订单流程创建快照时，clientToken才会保证幂等性。订单流程创建快照场景说明参考接口功能介绍。  */
	RegionID        string  `json:"regionID,omitempty"`        /*  资源池ID。  */
	SnapshotName    string  `json:"snapshotName,omitempty"`    /*  快照名称。仅允许英文字母、数字及_或者-，只能以英文字母开头，且长度为2-63字符。  */
	DiskID          string  `json:"diskID,omitempty"`          /*  云硬盘ID。  */
	RetentionPolicy string  `json:"retentionPolicy,omitempty"` /*  快照保留策略，取值范围：
	●custom：自定义保留天数。
	●forever：永久保留。  */
	RetentionTime int64 `json:"retentionTime"` /*  自定义快照保留天数。取值范围：1-65535。当快照保留策略为custom时该参数为必填，当快照保留策略为forever时，自动设置为65535。  */
}

type EbsCreateEbsSnapResponse struct {
	StatusCode  int64                                `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)。  */
	Message     string                               `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                               `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   EbsCreateEbsSnapReturnObjResponse    `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                              `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                              `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
	ErrorDetail *EbsCreateEbsSnapErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。  其他订单侧(bss)的错误，以ebs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
}

type EbsCreateEbsSnapReturnObjResponse struct {
	MasterOrderID string                                        `json:"masterOrderID,omitempty"` /*  主订单ID。  */
	MasterOrderNO string                                        `json:"sasterOrderNO,omitempty"`
	ClientToken   string                                        `json:"clientToken,omitempty"`
	Resources     []EbsNewEbsSnapshotReturnObjResourcesResponse `json:"resources,omitempty"`
	SnapshotJobID string                                        `json:"snapshotJobID,omitempty"`
}

type EbsCreateEbsSnapErrorDetailResponse struct {
	BssErrCode       *string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中。  */
	BssErrMsg        *string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中。  */
	BssOrigErr       *string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息。  */
	BssErrPrefixHint *string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息。  */
}

type EbsNewEbsSnapshotReturnObjResourcesResponse struct {
	OrderID      string `json:"orderID"` /*  订单ID。  */
	SnapshotID   string `json:"snapshotID"`
	SnapshotName string `json:"snapshotName"`
}
