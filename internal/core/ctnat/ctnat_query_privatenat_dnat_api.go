package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtnatQueryPrivatenatDnatApi
/* 查询DNAT
 */type CtnatQueryPrivatenatDnatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatQueryPrivatenatDnatApi(client *core.CtyunClient) *CtnatQueryPrivatenatDnatApi {
	return &CtnatQueryPrivatenatDnatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/privatenat/list-dnats",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatQueryPrivatenatDnatApi) Do(ctx context.Context, credential core.Credential, req *CtnatQueryPrivatenatDnatRequest) (*CtnatQueryPrivatenatDnatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("natGatewayID", req.NatGatewayID)
	ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtnatQueryPrivatenatDnatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatQueryPrivatenatDnatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要查询的私网NAT的ID。  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  分页查询时每页的行数，最大值为50，默认值为10。  */
}

type CtnatQueryPrivatenatDnatResponse struct {
	StatusCode   int32                                        `json:"statusCode"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                       `json:"message"`      /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  string                                       `json:"description"`  /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    string                                       `json:"errorCode"`    /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtnatQueryPrivatenatDnatReturnObjResponse `json:"returnObj"`    /*  见下表  */
	TotalCount   int32                                        `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                        `json:"currentCount"` /*  分页查询时当前页的行数。  */
	TotalPage    int32                                        `json:"totalPage"`    /*  总页数  */
}

type CtnatQueryPrivatenatDnatReturnObjResponse struct {
	DnatID       string `json:"dnatID"`       /*  DNAT规则的ID  */
	ExternalIP   string `json:"externalIP"`   /*  中转IP  */
	ExternalPort int32  `json:"externalPort"` /*  外部端口  */
	InternalIP   string `json:"internalIP"`   /*  内部IP  */
	InternalPort int32  `json:"internalPort"` /*  内部端口  */
	PortID       string `json:"portID"`       /*  对应的网卡ID  */
	PortName     string `json:"portName"`     /*  网卡名称  */
	DeviceID     string `json:"deviceID"`     /*  网卡对应的设备ID  */
	Protocol     string `json:"protocol"`     /*  协议: tcp/udp  */
	State        string `json:"state"`        /*  DNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期  */
	CreatedAt    string `json:"createdAt"`    /*  创建时间  */
	Description  string `json:"description"`  /*  描述  */
}
