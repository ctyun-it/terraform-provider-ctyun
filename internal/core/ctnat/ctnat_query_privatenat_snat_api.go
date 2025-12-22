package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtnatQueryPrivatenatSnatApi
/* 查询SNAT
 */type CtnatQueryPrivatenatSnatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatQueryPrivatenatSnatApi(client *core.CtyunClient) *CtnatQueryPrivatenatSnatApi {
	return &CtnatQueryPrivatenatSnatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/privatenat/list-snats",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatQueryPrivatenatSnatApi) Do(ctx context.Context, credential core.Credential, req *CtnatQueryPrivatenatSnatRequest) (*CtnatQueryPrivatenatSnatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("natGatewayID", req.NatGatewayID)
	if req.SnatID != "" {
		ctReq.AddParam("snatID", req.SnatID)
	}
	ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtnatQueryPrivatenatSnatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatQueryPrivatenatSnatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要查询的私网NAT的ID。  */
	SnatID       string `json:"snatID,omitempty"`       /*  SNAT规则的ID  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtnatQueryPrivatenatSnatResponse struct {
	StatusCode   int32                                        `json:"statusCode"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                       `json:"message"`      /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  string                                       `json:"description"`  /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    string                                       `json:"errorCode"`    /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtnatQueryPrivatenatSnatReturnObjResponse `json:"returnObj"`    /*  返回结果  */
	TotalCount   int32                                        `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                        `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                        `json:"totalPage"`    /*  总页数  */
}

type CtnatQueryPrivatenatSnatReturnObjResponse struct {
	SnatID        string   `json:"snatID"`        /*  SNAT规则的ID  */
	SrcCIDR       string   `json:"srcCIDR"`       /*  源地址段  */
	SrcVpcName    string   `json:"srcVpcName"`    /*  源vpc名称  */
	SrcSubnetID   string   `json:"srcSubnetID"`   /*  源Subnet的ID  */
	SrcSubnetName string   `json:"srcSubnetName"` /*  源Subnet名称  */
	Addresses     []string `json:"addresses"`     /*  中转IP地址  */
	Description   string   `json:"description"`   /*  描述  */
	State         string   `json:"state"`         /*  SNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期  */
}
