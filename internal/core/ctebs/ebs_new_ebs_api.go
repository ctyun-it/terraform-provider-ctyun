package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsNewEbsApi
/* 支持按需/包年包月创建云硬盘。
 */type EbsNewEbsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsNewEbsApi(client *core.CtyunClient) *EbsNewEbsApi {
	return &EbsNewEbsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/new-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsNewEbsApi) Do(ctx context.Context, credential core.Credential, req *EbsNewEbsRequest) (*EbsNewEbsResponse, error) {
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
	var resp EbsNewEbsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsNewEbsRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID。  */
	MultiAttach *bool   `json:"multiAttach"`           /*  是否多云主机挂载，取值范围：
	●true：是
	●false：否
	默认值：false  */
	IsEncrypt *bool `json:"isEncrypt"` /*  是否加密盘，取值范围：
	●true：是
	●false：否
	默认值：false
	共享盘、ISCSI模式磁盘、极速型SSD类型盘、XSSD系列盘不支持加密。  */
	KmsUUID   *string `json:"kmsUUID,omitempty"`   /*  如果是加密盘，需要提供kms的uuid。  */
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目ID，默认为0。  */
	DiskMode  string  `json:"diskMode,omitempty"`  /*  云硬盘磁盘模式，分为：
	●VBD（Virtual Block Device）：虚拟块存储设备。
	●ISCSI （Internet Small Computer System Interface）：小型计算机系统接口。
	●FCSAN（Fibre Channel SAN）：光纤通道协议的SAN网络。  */
	DiskType string `json:"diskType,omitempty"` /*  云硬盘类型及相关限制：
	●SATA：普通IO
	●SAS：高IO，只有该类型支持FCSAN模式
	●SSD：超高IO
	●FAST-SSD：极速型SSD，不支持ISCSI模式
	●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘，不支持加密，不支持ISCSI模式或FCSAN模式  */
	DiskName string `json:"diskName,omitempty"` /*  磁盘命名，单账户单资源池下，命名需唯一。
	仅允许英文字母、数字及_或者-，长度为2-63字符，不能以特殊字符开头。  */
	DiskSize int32 `json:"diskSize"` /*  云硬盘大小，单位GB。数据盘的取值范围：
	●超高IO/高IO/极速型SSD/普通IO：10GB~32768GB
	●XSSD-0：10GB-65536GB
	●XSSD-1：20GB-65536GB
	●XSSD-2：512GB-65536GB  */
	OnDemand *bool `json:"onDemand"` /*  是否按需下单，取值范围：
	●true：是
	●false：否
	默认值：true  */
	CycleType *string `json:"cycleType,omitempty"` /*  包周期类型，year：年，month：月。
	onDemand为false时，必须指定year和month的值。  */
	CycleCount int32 `json:"cycleCount"` /*  包周期数。onDemand为false时必须指定。
	周期为年（year）时，周期最大长度不能超过5年。
	周期为月（month）时，周期最大长度不能超过60月。  */
	ImageID *string `json:"imageID,omitempty"` /*  镜像ID，如果用镜像创建，只支持数据盘的私有镜像，所创建的数据盘的所在地域要与镜像源一致，容量不可小于镜像对应的磁盘容量。
	不支持批量创建操作，从镜像创建的数据盘不支持加密、ISCSI和FCSAN高级配置。  */
	AzName          *string `json:"azName,omitempty"` /*  多可用区资源池下，必须指定可用区。  */
	ProvisionedIops int32   `json:"provisionedIops"`  /*  XSSD类型盘的预配置iops值，最小值为1，其他类型的盘不支持设置。具体取值范围如下：
	●XSSD-0：（基础IOPS（min{1800+12×容量， 10,000}）   + 预配置IOPS） ≤ min{500×容量，50,000}
	●XSSD-1：（基础IOPS（min{1800+50×容量， 50000}） + 预配置IOPS） ≤ min{500×容量，100000}
	●XSSD-2：（基础IOPS（min{3000+50×容量， 100000}） + 预配置IOPS） ≤ min{500×容量，1000000}  */
	DeleteSnapWithEbs *bool `json:"deleteSnapWithEbs"` /*  设置全部快照是否随盘删除，取值范围：
	●true：是
	●false：否
	默认值：false
	*/
	Labels   []*EbsNewEbsLabelsRequest `json:"labels"`             /*  设置云硬盘标签，实际绑定标签的结果请查询云硬盘详情的labels返回值是否如预期。  */
	BackupID *string                   `json:"backupID,omitempty"` /*  云硬盘备份ID参数，有以下限制：
	●从备份创建盘仅支持VBD模式。
	●新盘容量不能小于备份源盘容量。
	●不支持配置加密属性（自动与备份源盘保持一致）。
	●备份状态必须是可用。  */
}

type EbsNewEbsLabelsRequest struct {
	Key   string `json:"key,omitempty"`   /*  标签的key值，长度不能超过32个字符。  */
	Value string `json:"value,omitempty"` /*  标签的value值，长度不能超过32个字符。  */
}

type EbsNewEbsResponse struct {
	StatusCode  int32                         `json:"statusCode"`            /*  返回状态码(800为成功，900为处理中/失败)。  */
	Message     *string                       `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                       `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsNewEbsReturnObjResponse   `json:"returnObj"`             /*  返回数据结构体。  */
	ErrorCode   *string                       `json:"errorCode,omitempty"`   /*  业务细分码，为Product.Module.Code三段式码.请参考错误码。  */
	Error       *string                       `json:"error,omitempty"`       /*  业务细分码，为Product.Module.Code三段式大驼峰码. 请参考错误码。  */
	ErrorDetail *EbsNewEbsErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。<br> 其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
}

type EbsNewEbsReturnObjResponse struct {
	MasterOrderID        *string                                `json:"masterOrderID,omitempty"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。参考错误码。  */
	MasterOrderNO        *string                                `json:"masterOrderNO,omitempty"`        /*  订单号。  */
	MasterResourceID     *string                                `json:"masterResourceID,omitempty"`     /*  主资源ID。  */
	MasterResourceStatus *string                                `json:"masterResourceStatus,omitempty"` /*  主资源状态。只有主订单资源会返回。  */
	RegionID             *string                                `json:"regionID,omitempty"`             /*  资源所属资源池ID。  */
	Resources            []*EbsNewEbsReturnObjResourcesResponse `json:"resources"`                      /*  资源明细列表，参考表resources。  */
}

type EbsNewEbsErrorDetailResponse struct{}

type EbsNewEbsReturnObjResourcesResponse struct {
	DiskID           *string `json:"diskID,omitempty"`           /*  资源底层ID，即云硬盘ID。  */
	OrderID          *string `json:"orderID,omitempty"`          /*  订单ID。  */
	StartTime        int64   `json:"startTime"`                  /*  启动时刻，epoch时戳，毫秒精度。  */
	CreateTime       int64   `json:"createTime"`                 /*  创建时刻，epoch时戳，毫秒精度。  */
	UpdateTime       int64   `json:"updateTime"`                 /*  更新时刻，epoch时戳，毫秒精度。  */
	Status           int32   `json:"status"`                     /*  资源状态。  */
	IsMaster         *bool   `json:"isMaster"`                   /*  是否是主资源项。  */
	ItemValue        int32   `json:"itemValue"`                  /*  资源规格，即为磁盘大小，单位为GB。  */
	ResourceType     *string `json:"resourceType,omitempty"`     /*  资源类型，云硬盘为EBS。  */
	DiskName         *string `json:"diskName,omitempty"`         /*  云硬盘名称。  */
	MasterOrderID    *string `json:"masterOrderID,omitempty"`    /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。  */
	MasterResourceID *string `json:"masterResourceID,omitempty"` /*  主资源ID。  */
}
