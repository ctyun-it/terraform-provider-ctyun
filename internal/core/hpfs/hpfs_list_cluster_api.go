package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListClusterApi
/* 查询对应资源池 ID 下集群列表
 */type HpfsListClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListClusterApi(client *core.CtyunClient) *HpfsListClusterApi {
	return &HpfsListClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-cluster",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListClusterApi) Do(ctx context.Context, credential core.Credential, req *HpfsListClusterRequest) (*HpfsListClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.SfsType != "" {
		ctReq.AddParam("sfsType", req.SfsType)
	}
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	if req.EbmDeviceType != "" {
		ctReq.AddParam("ebmDeviceType", req.EbmDeviceType)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsListClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListClusterRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池 ID  */
	SfsType       string `json:"sfsType,omitempty"`       /*  类型，hpfs_perf(HPC性能型)  */
	AzName        string `json:"azName,omitempty"`        /*  多可用区下的可用区名字  */
	EbmDeviceType string `json:"ebmDeviceType,omitempty"` /*  裸金属设备规格  */
	PageNo        int32  `json:"pageNo,omitempty"`        /*  列表的分页页码，默认值为1  */
	PageSize      int32  `json:"pageSize,omitempty"`      /*  每页包含的元素个数范围(1-50)，默认值为10  */
}

type HpfsListClusterResponse struct {
	StatusCode  int32                             `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                            `json:"message"`     /*  响应描述  */
	Description string                            `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListClusterReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                            `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                            `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListClusterReturnObjResponse struct {
	ClusterList  []*HpfsListClusterReturnObjClusterListResponse `json:"clusterList"`  /*  返回的集群列表  */
	TotalCount   int32                                          `json:"totalCount"`   /*  某资源池指定条件下集群总数  */
	CurrentCount int32                                          `json:"currentCount"` /*  当前页码下查询回来的集群数  */
	PageSize     int32                                          `json:"pageSize"`     /*  每页包含的元素个数范围(1-50)  */
	PageNo       int32                                          `json:"pageNo"`       /*  列表的分页页码  */
}

type HpfsListClusterReturnObjClusterListResponse struct {
	ClusterName     string   `json:"clusterName"`     /*  集群名称  */
	RemainingStatus *bool    `json:"remainingStatus"` /*  该集群是否可以售卖  */
	StorageType     string   `json:"storageType"`     /*  集群的存储类型  */
	AzName          string   `json:"azName"`          /*  多可用区下的可用区名字  */
	ProtocolType    []string `json:"protocolType"`    /*  集群支持的协议列表  */
	Baselines       []string `json:"baselines"`       /*  集群支持的性能基线列表（仅当资源池支持性能基线时返回）  */
	NetworkType     string   `json:"networkType"`     /*  集群的网络类型（tcp/o2ib）  */
	EbmDeviceTypes  []string `json:"ebmDeviceTypes"`  /*  集群支持的裸金属设备规格列表  */
}
