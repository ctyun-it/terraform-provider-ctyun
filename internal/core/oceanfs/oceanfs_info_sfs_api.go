package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsInfoSfsApi
/* 根据资源池 ID 和海量文件的sfsUID，查询文件系统详情
 */type OceanfsInfoSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsInfoSfsApi(client *core.CtyunClient) *OceanfsInfoSfsApi {
	return &OceanfsInfoSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oceanfs/info-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsInfoSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsInfoSfsRequest) (*OceanfsInfoSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sfsUID", req.SfsUID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp OceanfsInfoSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsInfoSfsRequest struct {
	SfsUID   string `json:"sfsUID,omitempty"`   /*  海量文件功能系统唯一 ID  */
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
}

type OceanfsInfoSfsReturnObjVpceSharePathResponse struct {
	VpcID              string `json:"vpcID"`
	WindowsSharePathV6 string `json:"windowsSharePathV6"`
	WindowsSharePath   string `json:"windowsSharePath"`
	VpcName            string `json:"vpcName"`
	SharePath          string `json:"sharePath"`
	SharePathV6        string `json:"sharePathV6"`
}

type OceanfsInfoSfsReturnObjResponse struct {
	UpdateTime         int64                                          `json:"updateTime"`
	ProjectID          string                                         `json:"projectID"`
	UsedSize           int32                                          `json:"usedSize"`
	SfsSize            int32                                          `json:"sfsSize"`
	UsedSizeCharge     bool                                           `json:"used_size_charge"`
	OnDemand           bool                                           `json:"onDemand"`
	RegionID           string                                         `json:"regionID"`
	VpceSharePath      []OceanfsInfoSfsReturnObjVpceSharePathResponse `json:"vpceSharePath"`
	SharePath          string                                         `json:"sharePath"`
	MountCount         int32                                          `json:"mountCount"`
	SfsType            string                                         `json:"sfsType"`
	WindowsSharePath   string                                         `json:"windowsSharePath"`
	CreateTime         int64                                          `json:"createTime"`
	ExpireTime         int64                                          `json:"expireTime"`
	SfsProtocol        string                                         `json:"sfsProtocol"`
	AzName             string                                         `json:"azName"`
	CephID             string                                         `json:"cephID"`
	SfsUID             string                                         `json:"sfsUID"`
	SharePathV6        string                                         `json:"sharePathV6"`
	SfsStatus          string                                         `json:"sfsStatus"`
	WindowsSharePathV6 string                                         `json:"windowsSharePathV6"`
	ProtectSwitch      string                                         `json:"protectSwitch"`
	SfsName            string                                         `json:"sfsName"`
}

type OceanfsInfoSfsResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                           `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                           `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                           `json:"error"`       /*  业务细分码，为 product.module.code 三段式码大驼峰形式  */
	ReturnObj   *OceanfsInfoSfsReturnObjResponse `json:"returnObj"`   /*  returnObj  */
}
