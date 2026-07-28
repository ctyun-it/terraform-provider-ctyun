package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtnatListPrivatenatCidrsApi
/* 查询中转网段
 */type CtnatListPrivatenatCidrsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatListPrivatenatCidrsApi(client *core.CtyunClient) *CtnatListPrivatenatCidrsApi {
	return &CtnatListPrivatenatCidrsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/privatenat/list-cidrs",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatListPrivatenatCidrsApi) Do(ctx context.Context, credential core.Credential, req *CtnatListPrivatenatCidrsRequest) (*CtnatListPrivatenatCidrsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("natGatewayID", req.NatGatewayID)
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtnatListPrivatenatCidrsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatListPrivatenatCidrsRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要查询的私网NAT的ID。  */
	PageNumber   int32  `json:"pageNumber,omitempty"`   /*  列表的页码，默认值为1。  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtnatListPrivatenatCidrsResponse struct {
	StatusCode   int32                                        `json:"statusCode"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                       `json:"message"`      /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  string                                       `json:"description"`  /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    string                                       `json:"errorCode"`    /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtnatListPrivatenatCidrsReturnObjResponse `json:"returnObj"`    /*  返回结果  */
	TotalCount   int32                                        `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                        `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                        `json:"totalPage"`    /*  总页数  */
}

type CtnatListPrivatenatCidrsReturnObjResponse struct {
	Name string `json:"name"` /*  中转网段名称  */
	Cidr string `json:"cidr"` /*  对用网段  */
}
