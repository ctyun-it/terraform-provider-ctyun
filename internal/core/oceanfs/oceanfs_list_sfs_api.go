package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// OceanfsListSfsApi
/* 资源池 ID 下，所有的文件系统详情查询
 */type OceanfsListSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsListSfsApi(client *core.CtyunClient) *OceanfsListSfsApi {
	return &OceanfsListSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oceanfs/list-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsListSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsListSfsRequest) (*OceanfsListSfsResponse, error) {
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
	var resp OceanfsListSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsListSfsRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  资源池 ID  */
	ProjectID string `json:"projectID,omitempty"` /*  资源所属企业项目 ID，默认为全部企业项目  */
	PageSize  int32  `json:"pageSize,omitempty"`  /*  每页包含的元素个数  */
	PageNo    int32  `json:"pageNo,omitempty"`    /*  列表的分页页码  */
}

type OceanfsListSfsReturnObjResponse struct {
	List         []*OceanfsListSfsReturnObjItemResponse `json:"list"`
	TotalCount   int32                                  `json:"totalCount"`
	CurrentCount int32                                  `json:"currentCount"`
	Total        int32                                  `json:"total"`
	PageSize     int32                                  `json:"pageSize"`
	PageNo       int32                                  `json:"pageNo"`
}
type OceanfsListSfsReturnObjItemPhySharePathResponse struct {
	VpcID              string `json:"vpcId"`
	VpcName            string `json:"vpcName"`
	SharePath          string `json:"sharePath"`
	SharePathV6        string `json:"sharePathV6"`
	WindowsSharePath   string `json:"windowsSharePath"`
	WindowsSharePathV6 string `json:"windowsSharePathV6"`
}

type OceanfsListSfsReturnObjItemResponse struct {
	SfsName            string                                            `json:"sfsName"`
	SfsUID             string                                            `json:"sfsUid"`
	SfsSize            int32                                             `json:"sfsSize"`
	SfsType            string                                            `json:"sfsType"`
	SfsProtocol        string                                            `json:"sfsProtocol"`
	SfsStatus          string                                            `json:"sfsStatus"`
	UsedSize           int32                                             `json:"usedSize"`
	CreateTime         int64                                             `json:"createTime"`
	UpdateTime         int64                                             `json:"updateTime"`
	ExpireTime         int64                                             `json:"expireTime"`
	ProjectID          string                                            `json:"projectId"`
	OnDemand           bool                                              `json:"onDemand"`
	RegionID           string                                            `json:"regionId"`
	AzName             string                                            `json:"azName"`
	SharePath          string                                            `json:"sharePath"`
	SharePathV6        string                                            `json:"sharePathV6"`
	WindowsSharePath   string                                            `json:"windowsSharePath"`
	WindowsSharePathV6 string                                            `json:"windowsSharePathV6"`
	MountCount         int32                                             `json:"mountCount"`
	CephID             string                                            `json:"cephId"`
	UsedSizeCharge     bool                                              `json:"usedSizeCharge"`
	VpceSharePath      []OceanfsListSfsReturnObjItemPhySharePathResponse `json:"vpceSharePath"`
}

type OceanfsListSfsResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                           `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                           `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                           `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
	ReturnObj   *OceanfsListSfsReturnObjResponse `json:"returnObj"`
}
