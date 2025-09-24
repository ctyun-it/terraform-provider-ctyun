package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsNewSfsApi
/* 创建文件系统
 */type HpfsNewSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsNewSfsApi(client *core.CtyunClient) *HpfsNewSfsApi {
	return &HpfsNewSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/new-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsNewSfsApi) Do(ctx context.Context, credential core.Credential, req *HpfsNewSfsRequest) (*HpfsNewSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsNewSfsRequest
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
	var resp HpfsNewSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsNewSfsRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID，例：100054c0416811e9a6690242ac110002  */
	ProjectID   string `json:"projectID,omitempty"`   /*  资源所属企业项目 ID，默认为"0"  */
	SfsType     string `json:"sfsType,omitempty"`     /*  并行文件类型，hpfs_perf(HPC性能型)  */
	SfsProtocol string `json:"sfsProtocol,omitempty"` /*  协议类型，nfs/hpfs  */
	OnDemand    *bool  `json:"onDemand"`              /*  是否按需下单。true/false，默认为 true  */
	CycleType   string `json:"cycleType,omitempty"`   /*  包周期（subscription）类型，year/month。onDemand 为 false 时，必须指定  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  包周期数。onDemand 为 false 时必须指定。周期最大长度不能超过 3 年  */
	SfsName     string `json:"sfsName,omitempty"`     /*  并行文件名。单账户单资源池下，命名需唯一  */
	SfsSize     int32  `json:"sfsSize,omitempty"`     /*  大小，单位 GB， 起始容量512，步长为512  */
	AzName      string `json:"azName,omitempty"`      /*  多可用区资源池下，必须指定可用区  */
	ClusterName string `json:"clusterName,omitempty"` /*  集群名称，仅资源池支持指定集群时可传入该参数  */
	Baseline    string `json:"baseline,omitempty"`    /*  性能基线（MB/s/TB），仅资源池支持性能基线时可传入该参数  */
	Vpc         string `json:"vpc,omitempty"`         /*  虚拟网 ID，标准协议必填，hpfs协议不作校验  */
	Subnet      string `json:"subnet,omitempty"`      /*  子网 ID，标准协议必填，hpfs协议不作校验  */
}

type HpfsNewSfsResponse struct {
	StatusCode  int32                        `json:"statusCode"`  /*  返回状态码(800为成功，900为处理中/失败，详见errorCode)  */
	Message     string                       `json:"message"`     /*  响应描述  */
	Description string                       `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsNewSfsReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                       `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                       `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsNewSfsReturnObjResponse struct {
	MasterOrderID        string                                  `json:"masterOrderID"`        /*  订单 ID。调用方在拿到 masterOrderID 之后，在若干错误情况下，可以使用 materOrderID 进一步确认订单状态及资源状态  */
	MasterOrderNO        string                                  `json:"masterOrderNO"`        /*  订单号  */
	MasterResourceID     string                                  `json:"masterResourceID"`     /*  主资源 ID  */
	MasterResourceStatus string                                  `json:"masterResourceStatus"` /*  主资源状态。参考主资源状态  */
	RegionID             string                                  `json:"regionID"`             /*  资源所属资源池 ID  */
	Resources            []*HpfsNewSfsReturnObjResourcesResponse `json:"resources"`            /*  资源明细  */
}

type HpfsNewSfsReturnObjResourcesResponse struct {
	ResourceID       string `json:"resourceID"`       /*  单项资源的变配、续订、退订等需要该资源项的 ID。比如某个云主机资源作为主资源，对其挂载  */
	ResourceType     string `json:"resourceType"`     /*  资源类型  */
	OrderID          string `json:"orderID"`          /*  无需关心  */
	StartTime        int64  `json:"startTime"`        /*  启动时刻，epoch 时戳，毫秒精度。例：1589869069561  */
	CreateTime       int64  `json:"createTime"`       /*  创建时刻，epoch 时戳，毫秒精度  */
	UpdateTime       int64  `json:"updateTime"`       /*  更新时刻，epoch 时戳，毫秒精度  */
	Status           int32  `json:"status"`           /*  资源状态  */
	IsMaster         *bool  `json:"isMaster"`         /*  是否是主资源项  */
	ItemValue        int32  `json:"itemValue"`        /*  资源大小  */
	SfsUID           string `json:"sfsUID"`           /*  并行文件内部唯一 ID  */
	SfsStatus        string `json:"sfsStatus"`        /*  并行文件状态  */
	MasterOrderID    string `json:"masterOrderID"`    /*  订单 ID  */
	SfsName          string `json:"sfsName"`          /*  并行文件名字  */
	MasterResourceID string `json:"masterResourceID"` /*  主资源 ID  */
}
