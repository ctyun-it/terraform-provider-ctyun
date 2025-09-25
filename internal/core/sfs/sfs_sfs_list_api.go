package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListApi
/* 资源池 ID 下，所有的文件系统详情查询
 */type SfsSfsListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListApi(client *core.CtyunClient) *SfsSfsListApi {
	return &SfsSfsListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListRequest) (*SfsSfsListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
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
	var resp SfsSfsListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池 ID  */
	ProjectID string `json:"projectID,omitempty"` /*  资源所属企业项目 ID，默认为"0"  */
	PageSize  int32  `json:"pageSize,omitempty"`  /*  每页包含的元素个数  */
	PageNo    int32  `json:"pageNo,omitempty"`    /*  列表的分页页码  */
}

type SfsSfsListResponse struct {
	StateCode   int32                        `json:"stateCode"`   /*  返回状态码(800为成功，900为失败)  */
	Message     string                       `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                       `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SfsSfsListReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                       `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码。  */
}

type SfsSfsListReturnObjResponse struct {
	List         []*SfsSfsListReturnObjListResponse `json:"list"`         /*  返回的文件列表  */
	TotalCount   int32                              `json:"totalCount"`   /*  资源池下用户弹性文件总数  */
	CurrentCount int32                              `json:"currentCount"` /*  当前页码下查询回来的用户弹性文件数  */
}

type SfsSfsListReturnObjListResponse struct {
	Name               string `json:"name"`               /*  弹性文件命名  */
	UID                string `json:"UID"`                /*  弹性文件功能系统唯一 ID  */
	ResourceID         string `json:"resourceID"`         /*  资源 ID（计费资源 ID）  */
	SfsSize            int32  `json:"sfsSize"`            /*  大小（GB）  */
	SfsType            string `json:"sfsType"`            /*  类型，capacity/performance/hdd_e/hpfs_perf/massive  */
	SfsProtocol        string `json:"sfsProtocol"`        /*  挂载协议。2 种，nfs/cifs  */
	SfsStatus          string `json:"sfsStatus"`          /*  弹性文件状态  */
	UsedSize           int32  `json:"usedSize"`           /*  已用大小（MB）  */
	CreateTime         int64  `json:"createTime"`         /*  创建时刻，epoch 时戳，精度毫秒  */
	UpdateTime         int64  `json:"updateTime"`         /*  更新时刻，epoch 时戳，精度毫秒  */
	ExpireTime         int64  `json:"expireTime"`         /*  过期时刻，epoch 时戳，精度毫秒  */
	ProjectID          string `json:"projectID"`          /*  资源所属企业项目 ID  */
	IsEncrypt          *bool  `json:"isEncrypt"`          /*  是否加密盘  */
	KmsUUID            string `json:"kmsUUID"`            /*  加密盘密钥 UUID  */
	OnDemand           *bool  `json:"onDemand"`           /*  是否按需订购  */
	RegionID           string `json:"regionID"`           /*  资源池 ID  */
	AzName             string `json:"azName"`             /*  多可用区下的可用区名字  */
	SharePath          string `json:"sharePath"`          /*  linux 主机共享路径  */
	SharePathV6        string `json:"sharePathV6"`        /*  linux 主机 IPv6 共享路径  */
	WindowsSharePath   string `json:"windowsSharePath"`   /*  win 主机共享路径  */
	WindowsSharePathV6 string `json:"windowsSharePathV6"` /*  win 主机 IPv6 共享路径  */
	MountCount         int32  `json:"mountCount"`         /*  挂载点数量  */
	CephID             string `json:"cephID"`             /*  ceph底层的id  */
}
