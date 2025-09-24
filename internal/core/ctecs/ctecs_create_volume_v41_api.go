package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateVolumeV41Api
/* 支持创建一块按量付费或包年包月云硬盘<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />&emsp;&emsp;计费模式：确认开通云硬盘的计费模式，详细查看<a href="https://www.ctyun.cn/document/10027696/10028345">计费模式</a><br />&emsp;&emsp;地域选择：选择创建云主机的资源池，详细查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a><br /><b>准备工作</b>：<br />&emsp;&emsp;成本估算：了解云硬盘的<a href="https://www.ctyun.cn/document/10027696/10043184">云硬盘计费说明</a><br />&emsp;&emsp;用户配额：确认个人在不同资源池下资源配额，可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9714&data=87">用户配额查询</a>接口进行查询<br />&emsp;&emsp;异步接口：该接口为异步接口，下单成功不代表云硬盘创建成功<br />&emsp;&emsp;企业项目：保证资源隔离
 */type CtecsCreateVolumeV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateVolumeV41Api(client *core.CtyunClient) *CtecsCreateVolumeV41Api {
	return &CtecsCreateVolumeV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/volume/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateVolumeV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreateVolumeV41Request) (*CtecsCreateVolumeV41Response, error) {
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
	var resp CtecsCreateVolumeV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateVolumeV41Request struct {
	RegionID          string                               `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	DiskMode          string                               `json:"diskMode,omitempty"`        /*  磁盘模式，取值范围：<br />FCSAN（光纤通道协议的SAN网络），<br />ISCSI（小型计算机系统接口），<br />VBD（虚拟块存储设备）<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162960">磁盘模式及使用方法</a>了解相关内容。<br />注：默认值为VBD；XSSD类型盘不支持ISCSI和FCSAN   */
	DiskType          string                               `json:"diskType,omitempty"`        /*  云硬盘类型，取值范围：<br />SATA（普通IO），<br />SAS（高IO），<br />SSD（超高IO），<br />FAST-SSD（极速型SSD）<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息<br />注：极速SSD类型盘（FAST-SSD）不支持ISCSI磁盘模式;只有高IO类型（SAS）支持FCSAN磁盘模式；XSSD类型盘不支持多挂载、加密、ISCSI和FCSAN磁盘模式  */
	DiskName          string                               `json:"diskName,omitempty"`        /*  云硬盘命名，单账户单资源池下，命名需唯一  */
	DiskSize          int32                                `json:"diskSize,omitempty"`        /*  云硬盘大小，单位GB，取值范围：[10, 32768]  */
	ClientToken       string                               `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时  */
	AzName            string                               `json:"azName,omitempty"`          /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称（即多可用区，本字段必填），若查询结果中zoneList为空（即为单可用区，本字段不可填写）  */
	MultiAttach       *bool                                `json:"multiAttach"`               /*  是否支持多云主机挂载，注：默认值false；XSSD类型盘不支持多挂载  */
	OnDemand          *bool                                `json:"onDemand"`                  /*  是否按需下单，注：默认值true  */
	CycleType         string                               `json:"cycleType,omitempty"`       /*  订购周期类型，取值范围： <br />MONTH（表示按月），<br />YEAR（表示按年）。注：onDemand为false时，必须指定。  */
	CycleCount        int32                                `json:"cycleCount,omitempty"`      /*  包周期数。注：onDemand为false时必须指定；周期最大长度不能超过5年   */
	IsEncrypt         *bool                                `json:"isEncrypt"`                 /*  是否加密盘，取值范围：<br />true表示加密，<br />false表示未加密。注：默认值false；XSSD类型盘不支持加密  */
	KmsUUID           string                               `json:"kmsUUID,omitempty"`         /*  天翼云自研密钥管理（KMS，Key Management Service）的ID，如果是加密盘（参数isEncrypt为true时），需要提供KMS的uuid，您可以查看<a href="https://www.ctyun.cn/document/10027696/10162638">支持云硬盘加密功能</a>了解云硬盘加密功能  */
	ProjectID         string                               `json:"projectID,omitempty"`       /*  企业项目ID，默认值为"0"，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目   */
	ImageID           string                               `json:"imageID,omitempty"`         /*  镜像ID，如果用镜像创建，只支持数据盘的私有镜像，所创建的数据盘的所在地域要与镜像源一致，容量不可小于镜像对应的磁盘容量，不支持批量创建操作，从镜像创建的数据盘不支持加密、ISCSI和FCSAN高级配置  */
	ProvisionedIops   int32                                `json:"provisionedIops,omitempty"` /*  XSSD类型盘的预配置iops值，最小值为1，其他类型的盘不支持设置。具体取值范围如下：<br />●XSSD-0：（基础IOPS（min{1800+12×容量， 10,000}） + 预配置IOPS） ≤ min{500×容量，50,000}<br />●XSSD-1：（基础IOPS（min{1800+50×容量， 50000}） + 预配置IOPS） ≤ min{500×容量，100000}<br />●XSSD-2：（基础IOPS（min{3000+50×容量， 100000}） + 预配置IOPS） ≤ min{500×容量，1000000}  */
	DeleteSnapWithEbs *bool                                `json:"deleteSnapWithEbs"`         /*  设置全部快照是否随盘删除，取值范围：<br />●true：是<br />●false：否<br />默认值：false  */
	Labels            []*CtecsCreateVolumeV41LabelsRequest `json:"labels"`                    /*  设置云硬盘标签，实际绑定标签的结果请查询云硬盘详情的labels返回值是否如预期  */
	BackupID          string                               `json:"backupID,omitempty"`        /*  云硬盘备份ID参数，有以下限制：<br />●从备份创建盘仅支持VBD模式<br />●新盘容量不能小于备份源盘容量<br />●不支持配置加密属性（自动与备份源盘保持一致）<br />●备份状态必须是可用  */
}

type CtecsCreateVolumeV41LabelsRequest struct {
	Key   string `json:"key,omitempty"`   /*  标签的key值，长度不能超过32个字符  */
	Value string `json:"value,omitempty"` /*  标签的value值，长度不能超过32个字符  */
}

type CtecsCreateVolumeV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中或失败)   */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CtecsCreateVolumeV41ErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。其他订单侧(bss)的错误，以ebs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
	ReturnObj   *CtecsCreateVolumeV41ReturnObjResponse   `json:"returnObj"`             /*  返回参数  */
}

type CtecsCreateVolumeV41ErrorDetailResponse struct {
	BssErrCode       string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中  */
	BssErrMs         string `json:"bssErrMs,omitempty"`         /*  bss错误信息，包含于bss格式化JSON错误信息中  */
	BssOrigErr       string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息  */
	BssErrPrefixHint string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息  */
}

type CtecsCreateVolumeV41ReturnObjResponse struct {
	MasterOrderID        string                                            `json:"masterOrderID,omitempty"`        /*  主订单ID。调用方在拿到masterOrderID之后，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO        string                                            `json:"masterOrderNO,omitempty"`        /*  订单号  */
	MasterResourceID     string                                            `json:"masterResourceID,omitempty"`     /*  主资源ID。云硬盘场景下，无需关心   */
	MasterResourceStatus string                                            `json:"masterResourceStatus,omitempty"` /*  主资源状态。只有主订单资源会返回。unknown/failed/ starting   */
	RegionID             string                                            `json:"regionID,omitempty"`             /*  资源池ID  */
	Resources            []*CtecsCreateVolumeV41ReturnObjResourcesResponse `json:"resources"`                      /*  资源明细列表，注：订单未完成情况下不返回；待订单完成后使用创建时同样的clientToken（24小时后失效）进行查询，则返回该部分内容  */
}

type CtecsCreateVolumeV41ReturnObjResourcesResponse struct {
	DiskID           string `json:"diskID,omitempty"`           /*  资源底层ID，即磁盘ID   */
	OrderID          string `json:"orderID,omitempty"`          /*  无需关心  */
	StartTime        int32  `json:"startTime,omitempty"`        /*  启动时刻，epoch时戳，毫秒精度  */
	CreateTime       int32  `json:"createTime,omitempty"`       /*  创建时刻，epoch时戳，毫秒精度  */
	UpdateTime       int32  `json:"updateTime,omitempty"`       /*  更新时刻，epoch时戳，毫秒精度  */
	Status           int32  `json:"status,omitempty"`           /*  资源状态，无需关心。参考masterResourceStatus   */
	IsMaster         *bool  `json:"isMaster"`                   /*  是否是主资源项  */
	ItemValue        int32  `json:"itemValue,omitempty"`        /*  资源规格，这里指磁盘大小，单位GB   */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型，云硬盘为EBS   */
	DiskName         string `json:"diskName,omitempty"`         /*  云硬盘名称  */
	MasterOrderID    string `json:"masterOrderID,omitempty"`    /*  主订单ID。调用方在拿到masterOrderID之后，可以使用masterOrderID进一步确认订单状态及资源状态  */
	MasterResourceID string `json:"masterResourceID,omitempty"` /*  主资源ID。云硬盘场景下，无需关心  */
}
