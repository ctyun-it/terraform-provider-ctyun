package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsDeleteEbsSnapApi
/* 您可以删除不再需要的快照，以达到节省资源和成本的目的。
 */type EbsDeleteEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsDeleteEbsSnapApi(client *core.CtyunClient) *EbsDeleteEbsSnapApi {
	return &EbsDeleteEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/delete-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsDeleteEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsDeleteEbsSnapRequest) (*EbsDeleteEbsSnapResponse, error) {
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
	var resp EbsDeleteEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsDeleteEbsSnapRequest struct {
	RegionID    string   `json:"regionID,omitempty"` /*  资源池ID，请根据查询资源池列表接口返回值进行传参，获取“regionId”参数。  */
	SnapshotIDs []string `json:"snapshotIDs"`        /*  云硬盘快照ID列表，请根据“查询快照列表接口”，获取snapshotID参数的返回值并进行传参。当refundOrder为true时不校验该字段，将删除所有的快照，该字段传[ ]即可。  */
	RefundOrder bool     `json:"refundOrder"`        /*  是否退订该硬盘下的所有的快照，取值范围：
	●true：将删除所有的快照并删除订单
	●false：只删除快照不删除订单
	默认值为false。  */
	DiskID string `json:"diskID,omitempty"` /*  云硬盘ID。  */
}

type EbsDeleteEbsSnapResponse struct {
	StatusCode  int32                              `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败）。  */
	Message     string                             `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                             `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsDeleteEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                            `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                            `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                            `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsDeleteEbsSnapReturnObjResponse struct {
	SnapshotJobID  *string   `json:"snapshotJobID,omitempty"` /*  删除快照任务ID，仅删除不退订时返回该字段。  */
	NotAllowedList []*string `json:"notAllowedList"`          /*  状态不允许删除的快照ID列表，仅删除不退订时返回该字段。  */
	NotFoundList   []*string `json:"notFoundList"`            /*  不存在的快照ID列表，仅删除不退订时返回该字段。  */
	MasterOrderID  *string   `json:"masterOrderID,omitempty"` /*  主订单ID，退订时返回该参数，普通删除无该字段。  */
	MasterOrderNO  *string   `json:"masterOrderNO,omitempty"` /*  主订单号，退订时返回该参数，普通删除无该字段。  */
	RegionID       *string   `json:"regionID,omitempty"`      /*  资源池ID，退订时返回该参数，普通删除无该字段。  */
}
