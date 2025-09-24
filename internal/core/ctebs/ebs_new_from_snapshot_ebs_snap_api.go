package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsNewFromSnapshotEbsSnapApi
/* 您可以通过快照快速创建出多个具有相同数据的云硬盘，可用于业务的快速部署。
 */type EbsNewFromSnapshotEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsNewFromSnapshotEbsSnapApi(client *core.CtyunClient) *EbsNewFromSnapshotEbsSnapApi {
	return &EbsNewFromSnapshotEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/new-from-snapshot-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsNewFromSnapshotEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsNewFromSnapshotEbsSnapRequest) (*EbsNewFromSnapshotEbsSnapResponse, error) {
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
	var resp EbsNewFromSnapshotEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsNewFromSnapshotEbsSnapRequest struct {
	SnapshotID  string  `json:"snapshotID,omitempty"`  /*  快照ID。  */
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID。  */
	MultiAttach *bool   `json:"multiAttach"`           /*  是否多云主机挂载，取值范围：
	●true：是
	●false：否
	默认值：false  */
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目ID，默认为”0”  */
	DiskMode  string  `json:"diskMode,omitempty"`  /*  磁盘模式。当前仅支持VBD。
	●VBD（Virtual Block Device）：虚拟块存储设备。  */
	DiskName string `json:"diskName,omitempty"` /*  磁盘命名，单账户单资源池下，命名需唯一。
	仅允许英文字母、数字及_或者-，长度为2-63字符，不能以特殊字符开头。  */
	DiskSize int32 `json:"diskSize"` /*  磁盘大小，单位GB。数据盘的取值范围为：
	●超高IO/高IO/极速型SSD/普通IO：10GB~32768GB
	●XSSD-0：10GB-65536GB
	●XSSD-1：20GB-65536GB
	●XSSD-2：512GB-65536GB  */
	OnDemand *bool `json:"onDemand"` /*  是否按需下单，取值范围：
	●true：是
	●false：否
	默认值：true  */
	CycleType *string `json:"cycleType,omitempty"` /*  包周期类型，取值范围：
	●year：包年。
	●month：包月。
	onDemand为false时，必须指定该参数。  */
	CycleCount int32 `json:"cycleCount"` /*  包周期数。onDemand为false时必须指定。
	周期为年（year）时，周期最大长度不能超过5年。
	周期为月（month）时，周期最大长度不能超过60月。  */
}

type EbsNewFromSnapshotEbsSnapResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为处理中/失败，详见errorCode）。  */
	Message     *string                                       `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                       `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsNewFromSnapshotEbsSnapReturnObjResponse   `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                       `json:"details,omitempty"`     /*  详情描述。  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                                       `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
	ErrorDetail *EbsNewFromSnapshotEbsSnapErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。  其他订单侧(bss)的错误，以ebs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
}

type EbsNewFromSnapshotEbsSnapReturnObjResponse struct {
	MasterOrderID        *string                                                `json:"masterOrderID,omitempty"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO        *string                                                `json:"masterOrderNO,omitempty"`        /*  订单号。  */
	MasterResourceID     *string                                                `json:"masterResourceID,omitempty"`     /*  主资源ID。  */
	MasterResourceStatus *string                                                `json:"masterResourceStatus,omitempty"` /*  主资源状态。只有主订单资源会返回。  */
	RegionID             *string                                                `json:"regionID,omitempty"`             /*  资源所属的资源池ID。  */
	Resources            []*EbsNewFromSnapshotEbsSnapReturnObjResourcesResponse `json:"resources"`                      /*  资源明细列表，参考表resources。  */
}

type EbsNewFromSnapshotEbsSnapErrorDetailResponse struct {
	BssErrCode       *string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中。  */
	BssErrMsg        *string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中。  */
	BssOrigErr       *string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息。  */
	BssErrPrefixHint *string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息。  */
}

type EbsNewFromSnapshotEbsSnapReturnObjResourcesResponse struct {
	DiskID           *string `json:"diskID,omitempty"`           /*  云硬盘ID。  */
	OrderID          *string `json:"orderID,omitempty"`          /*  订单ID。  */
	StartTime        int32   `json:"startTime"`                  /*  启动时刻，epoch时戳，毫秒精度。  */
	CreateTime       int32   `json:"createTime"`                 /*  创建时刻，epoch时戳，毫秒精度。  */
	UpdateTime       int32   `json:"updateTime"`                 /*  更新时刻，epoch时戳，毫秒精度。  */
	Status           int32   `json:"status"`                     /*  资源状态。  */
	IsMaster         *bool   `json:"isMaster"`                   /*  是否是主资源。  */
	ItemValue        int32   `json:"itemValue"`                  /*  资源规格，即为磁盘大小，单位为GB。  */
	DiskName         *string `json:"diskName,omitempty"`         /*  云硬盘名称。  */
	ResourceType     *string `json:"resourceType,omitempty"`     /*  云硬盘资源类型EBS（只有一种）。  */
	MasterOrderID    *string `json:"masterOrderID,omitempty"`    /*  订单ID。  */
	MasterResourceID *string `json:"masterResourceID,omitempty"` /*  主资源ID。  */
}
