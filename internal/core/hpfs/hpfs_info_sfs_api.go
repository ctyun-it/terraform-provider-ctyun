package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsInfoSfsApi
/* 根据资源池 ID 和 sfsUID，查询文件系统详情
 */type HpfsInfoSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsInfoSfsApi(client *core.CtyunClient) *HpfsInfoSfsApi {
	return &HpfsInfoSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/info-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsInfoSfsApi) Do(ctx context.Context, credential core.Credential, req *HpfsInfoSfsRequest) (*HpfsInfoSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sfsUID", req.SfsUID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsInfoSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsInfoSfsRequest struct {
	SfsUID   string `json:"sfsUID,omitempty"`   /*  并行文件唯一 ID  */
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
}

type HpfsInfoSfsResponse struct {
	StatusCode  int32                         `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                        `json:"message"`     /*  响应描述  */
	Description string                        `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsInfoSfsReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                        `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                        `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsInfoSfsReturnObjResponse struct {
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
