package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsInfoByNameApi
/* i根据文件系统name和资源池ID，查询文件系统详情
 */type SfsSfsInfoByNameApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsInfoByNameApi(client *core.CtyunClient) *SfsSfsInfoByNameApi {
	return &SfsSfsInfoByNameApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/info-by-name",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsInfoByNameApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsInfoByNameRequest) (*SfsSfsInfoByNameResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("name", req.Name)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsInfoByNameResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsInfoByNameRequest struct {
	Name     string `json:"name,omitempty"`     /*  文件系统名称  */
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
}

type SfsSfsInfoByNameResponse struct {
	StatusCode int32                              `json:"statusCode"` /*  返回状态码(800为成功，900为失败)  */
	Message    string                             `json:"message"`    /*  失败时的错误描述，一般为英文描述  */
	MsgDesc    string                             `json:"msgDesc"`    /*  失败时的错误描述，一般为中文描述  */
	MsgCode    string                             `json:"msgCode"`    /*  业务细分码，为product.module.code三段式码  */
	ReturnObj  *SfsSfsInfoByNameReturnObjResponse `json:"returnObj"`  /*  参考详情对象  */
}

type SfsSfsInfoByNameReturnObjResponse struct {
	Name               string `json:"name"`               /*  弹性文件命名  */
	UID                string `json:"UID"`                /*  弹性文件功能系统唯一ID  */
	ResourceID         string `json:"resourceID"`         /*  资源ID（计费资源ID）  */
	SfsSize            int32  `json:"sfsSize"`            /*  大小（MB）  */
	SfsType            string `json:"sfsType"`            /*  类型，capacity/performance  */
	SfsProtocol        string `json:"sfsProtocol"`        /*  挂载协议。nfs/cifs/smb  */
	SfsStatus          string `json:"sfsStatus"`          /*  1  */
	UsedSize           int32  `json:"usedSize"`           /*  已用大小,MB  */
	CreateTime         int64  `json:"createTime"`         /*  创建时刻，epoch时戳，精度毫秒  */
	UpdateTime         int64  `json:"updateTime"`         /*  更新时刻，epoch时戳，精度毫秒  */
	ExpireTime         int64  `json:"expireTime"`         /*  过期时刻，epoch时戳，精度毫秒  */
	ProjectID          string `json:"projectID"`          /*  资源所属企业项目ID  */
	IsEncrypt          *bool  `json:"isEncrypt"`          /*  是否加密盘  */
	KmsUUID            string `json:"kmsUUID"`            /*  加密盘密钥UUID  */
	OnDemand           *bool  `json:"onDemand"`           /*  是否按需订购  */
	RegionID           string `json:"regionID"`           /*  资源池ID  */
	AzName             string `json:"azName"`             /*  多可用区下的可用区名字  */
	SharePath          string `json:"sharePath"`          /*  linux主机共享路径  */
	SharePathV6        string `json:"sharePathV6"`        /*  linux主机IPv6共享路径  */
	WindowsSharePath   string `json:"windowsSharePath"`   /*  win主机共享路径  */
	WindowsSharePathV6 string `json:"windowsSharePathV6"` /*  win主机IPv6共享路径  */
	MountCount         int32  `json:"mountCount"`         /*  挂载点数量  */
}
