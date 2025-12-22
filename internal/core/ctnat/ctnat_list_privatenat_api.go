package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtnatListPrivatenatApi
/* 查询私网NAT
 */type CtnatListPrivatenatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatListPrivatenatApi(client *core.CtyunClient) *CtnatListPrivatenatApi {
	return &CtnatListPrivatenatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/privatenat/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatListPrivatenatApi) Do(ctx context.Context, credential core.Credential, req *CtnatListPrivatenatRequest) (*CtnatListPrivatenatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.NatGatewayID != "" {
		ctReq.AddParam("natGatewayID", req.NatGatewayID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
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
	var resp CtnatListPrivatenatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatListPrivatenatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要查询的私网NAT的ID。  */
	PageNumber   int32  `json:"pageNumber,omitempty"`   /*  列表的页码，默认值为1。  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtnatListPrivatenatResponse struct {
	StatusCode   int32                                   `json:"statusCode"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                  `json:"message"`      /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  string                                  `json:"description"`  /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    string                                  `json:"errorCode"`    /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtnatListPrivatenatReturnObjResponse `json:"returnObj"`    /*  返回结果  */
	TotalCount   int32                                   `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                   `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                   `json:"totalPage"`    /*  总页数  */
}

type CtnatListPrivatenatReturnObjResponse struct {
	NatGatewayID string `json:"natGatewayID"` /*  nat 网关 id  */
	Name         string `json:"name"`         /*  nat 网关名字  */
	Description  string `json:"description"`  /*  nat 网关描述  */
	VpcID        string `json:"vpcID"`        /*  虚拟私有云 id  */
	VpcName      string `json:"vpcName"`      /*  虚拟私有云名字  */
	SubnetID     string `json:"subnetID"`     /*  子网 id  */
	SubnetName   string `json:"subnetName"`   /*  子网名称  */
	State        string `json:"state"`        /*  私网运行状态: running 表示运行中，freeze表示冻结，expired表示已过期  */
	Spec         string `json:"spec"`         /*  规格取值: small, medium, large, xlarge  */
	ProjectID    string `json:"projectID"`    /*  项目ID  */
	ProjectName  string `json:"projectName"`  /*  项目名称  */
	AzID         string `json:"azID"`         /*  可用区ID  */
	CreateDate   string `json:"createDate"`   /*  创建时间  */
}
