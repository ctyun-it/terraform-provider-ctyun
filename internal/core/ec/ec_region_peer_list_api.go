package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcRegionPeerListApi
/* 查询已创建的云间高速 */
type EcRegionPeerListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcRegionPeerListApi(client *core.CtyunClient) *EcRegionPeerListApi {
	return &EcRegionPeerListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/region-peer/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcRegionPeerListApi) Do(ctx context.Context, credential core.Credential, req *EcRegionPeerListRequest) (*EcRegionPeerListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.EcID == "" {
		err := fmt.Errorf("missing required field 'EcID'")
		return nil, err
	}
	if req.PacketID != nil && *req.PacketID != "" {
		ctReq.AddParam("packetID", *req.PacketID)
	}
	if req.CgwID != nil && *req.CgwID != "" {
		ctReq.AddParam("cgwID", *req.CgwID)
	}
	ctReq.AddParam("ecID", req.EcID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcRegionPeerListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcRegionPeerListRequest struct {
	EcID     string  `json:"ecID"`
	PacketID *string `json:"packetID,omitempty"`
	CgwID    *string `json:"cgwID,omitempty"`
}

type EcRegionPeerListResponse struct {
	StatusCode  *int32                             `json:"statusCode"`  /*  返回状态码,<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                            `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcRegionPeerListReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcRegionPeerListReturnObjResponse struct {
	CurrentCount *int32                                      `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                      `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                      `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcRegionPeerListReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcRegionPeerListReturnObjResultsResponse struct {
	PeerID     *string `json:"peerID"`
	PeerName   *string `json:"peerName"`
	SrcCgwID   *string `json:"srcCgwID"`
	DstCgwID   *string `json:"dstCgwID"`
	SrcCgwName *string `json:"srcCgwName"`
	DstCgwName *string `json:"dstCgwName"`
	SrcDcID    *string `json:"srcDcID"`
	DstDcID    *string `json:"dstDcID"`
	SrcDcName  *string `json:"srcDcName"`
	DstDcName  *string `json:"dstDcName"`
	PeerType   *int32  `json:"peerType"`
	EcID       *string `json:"ecID"`
	PacketID   *string `json:"packetID"`
	PacketName *string `json:"packetName"`
	Rate       *int32  `json:"rate"`
	Status     *string `json:"status"`
	UpdateDate *string `json:"updateDate"`
}
