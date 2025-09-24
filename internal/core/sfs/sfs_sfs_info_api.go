package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsInfoApi
/* 根据资源池 ID 和 文件系统resourceID，查询文件系统详情
 */type SfsSfsInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsInfoApi(client *core.CtyunClient) *SfsSfsInfoApi {
	return &SfsSfsInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/info",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsInfoApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsInfoRequest) (*SfsSfsInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("resourceID", req.ResourceID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsInfoRequest struct {
	ResourceID string `json:"resourceID,omitempty"` /*  资源 ID  */
	RegionID   string `json:"regionID,omitempty"`   /*  资源池 ID  */
}

type SfsSfsInfoResponse struct {
	StatusCode  int32                        `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                       `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                       `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsInfoReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                       `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码。  */
	Error       string                       `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式.  */
}

type SfsSfsInfoReturnObjResponse struct {
	Name               string `json:"name"`               /*  弹性文件命名  */
	UID                string `json:"UID"`                /*  弹性文件功能系统唯一 ID  */
	ResourceID         string `json:"resourceID"`         /*  资源 ID（计费资源 ID）  */
	SfsSize            int32  `json:"sfsSize"`            /*  大小（GB）  */
	SfsType            string `json:"sfsType"`            /*  类型，capacity/performance/hdd_e/hpfs_perf/massive  */
	SfsProtocol        string `json:"sfsProtocol"`        /*  挂载协议。2 种，nfs/cifs  */
	SfsStatus          string `json:"sfsStatus"`          /*  弹性文件状态。creating/available/unusable/delete_error/deleting  */
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
