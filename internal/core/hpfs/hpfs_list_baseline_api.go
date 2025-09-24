package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListBaselineApi
/* 查询对应资源池 ID 下，指定存储类型的性能基线列表，若资源池不支持性能基线，则该接口会报错
 */type HpfsListBaselineApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListBaselineApi(client *core.CtyunClient) *HpfsListBaselineApi {
	return &HpfsListBaselineApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-baseline",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListBaselineApi) Do(ctx context.Context, credential core.Credential, req *HpfsListBaselineRequest) (*HpfsListBaselineResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsType", req.SfsType)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	if req.ClusterName != "" {
		ctReq.AddParam("clusterName", req.ClusterName)
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
	var resp HpfsListBaselineResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListBaselineRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	SfsType     string `json:"sfsType,omitempty"`     /*  类型，hpfs_perf(HPC性能型)  */
	AzName      string `json:"azName,omitempty"`      /*  多可用区下的可用区名字，4.0资源池必填  */
	ClusterName string `json:"clusterName,omitempty"` /*  集群名称  */
	PageNo      int32  `json:"pageNo,omitempty"`      /*  列表的分页页码 ，默认值为1  */
	PageSize    int32  `json:"pageSize,omitempty"`    /*  每页包含的元素个数范围(1-50)，默认值为10  */
}

type HpfsListBaselineResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                             `json:"message"`     /*  响应描述  */
	Description string                             `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListBaselineReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListBaselineReturnObjResponse struct {
	BaselineList []*HpfsListBaselineReturnObjBaselineListResponse `json:"baselineList"` /*  返回的性能基线列表  */
	TotalCount   int32                                            `json:"totalCount"`   /*  指定条件下性能基线总数  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页码下查询回来的基线数  */
	PageSize     int32                                            `json:"pageSize"`     /*  每页包含的元素个数范围(1-50)  */
	PageNo       int32                                            `json:"pageNo"`       /*  列表的分页页码  */
}

type HpfsListBaselineReturnObjBaselineListResponse struct {
	Baseline     string   `json:"baseline"`     /*  性能基线（MB/s/TB）  */
	StorageType  string   `json:"storageType"`  /*  支持类型，hpfs_perf(HPC性能型)  */
	AzName       string   `json:"azName"`       /*  多可用区下可用区名称  */
	ClusterNames []string `json:"clusterNames"` /*  集群列表  */
}
