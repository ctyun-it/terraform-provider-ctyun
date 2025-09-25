package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListSfsByClusterApi
/* 查询指定集群的并行文件列表
 */type HpfsListSfsByClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListSfsByClusterApi(client *core.CtyunClient) *HpfsListSfsByClusterApi {
	return &HpfsListSfsByClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-sfs-by-cluster",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListSfsByClusterApi) Do(ctx context.Context, credential core.Credential, req *HpfsListSfsByClusterRequest) (*HpfsListSfsByClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	ctReq.AddParam("clusterName", req.ClusterName)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsListSfsByClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListSfsByClusterRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	AzName      string `json:"azName,omitempty"`      /*  可用区名称，多可用区下必填  */
	ClusterName string `json:"clusterName,omitempty"` /*  集群名称  */
	ProjectID   string `json:"projectID,omitempty"`   /*  资源所属企业项目 ID，默认为"0"  */
	PageSize    int32  `json:"pageSize,omitempty"`    /*  每页包含的元素个数范围(1-50)，默认值为10  */
	PageNo      int32  `json:"pageNo,omitempty"`      /*  列表的分页页码，默认值为1  */
}

type HpfsListSfsByClusterResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述  */
	Description string                                 `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListSfsByClusterReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListSfsByClusterReturnObjResponse struct {
	List         []*HpfsListSfsByClusterReturnObjListResponse `json:"list"`         /*  文件系统详情列表  */
	Total        int32                                        `json:"total"`        /*  指定条件下用户并行文件总数  */
	TotalCount   int32                                        `json:"totalCount"`   /*  指定条件下用户并行文件总数  */
	CurrentCount int32                                        `json:"currentCount"` /*  当前页码的元素个数  */
	PageSize     int32                                        `json:"pageSize"`     /*  每页个数  */
	PageNo       int32                                        `json:"pageNo"`       /*  当前页数  */
}

type HpfsListSfsByClusterReturnObjListResponse struct {
	SfsName       string   `json:"sfsName"`       /*  并行文件命名  */
	SfsUID        string   `json:"sfsUID"`        /*  并行文件唯一 ID  */
	SfsSize       int32    `json:"sfsSize"`       /*  大小（GB）  */
	SfsType       string   `json:"sfsType"`       /*  类型，hpfs_perf(HPC性能型)  */
	SfsProtocol   string   `json:"sfsProtocol"`   /*  挂载协议，nfs/hpfs  */
	SfsStatus     string   `json:"sfsStatus"`     /*  并行文件状态  */
	UsedSize      int32    `json:"usedSize"`      /*  已用大小（MB）  */
	CreateTime    int64    `json:"createTime"`    /*  创建时刻，epoch 时戳，精度毫秒  */
	UpdateTime    int64    `json:"updateTime"`    /*  更新时刻，epoch 时戳，精度毫秒  */
	ProjectID     string   `json:"projectID"`     /*  资源所属企业项目 ID  */
	OnDemand      *bool    `json:"onDemand"`      /*  是否按需订购  */
	RegionID      string   `json:"regionID"`      /*  资源池 ID  */
	AzName        string   `json:"azName"`        /*  多可用区下的可用区名字  */
	ClusterName   string   `json:"clusterName"`   /*  集群名称  */
	Baseline      string   `json:"baseline"`      /*  性能基线（MB/s/TB）  */
	HpfsSharePath string   `json:"hpfsSharePath"` /*  HPFS文件系统共享路径(Linux)  */
	SecretKey     string   `json:"secretKey"`     /*  HPC型挂载需要的密钥  */
	DataflowList  []string `json:"dataflowList"`  /*  HPFS文件系统下的数据流动策略ID列表  */
	DataflowCount int32    `json:"dataflowCount"` /*  HPFS文件系统下的数据流动策略数量  */
}
