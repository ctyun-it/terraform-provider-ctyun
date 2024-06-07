package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// BandwidthDescribeApi 查询共享带宽详情
// https://www.ctyun.cn/document/10026761/10040766
type BandwidthDescribeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewBandwidthDescribeApi(client *ctyunsdk.CtyunClient) *BandwidthDescribeApi {
	return &BandwidthDescribeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/bandwidth/describe",
		},
	}
}

func (this *BandwidthDescribeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *BandwidthDescribeRequest) (*BandwidthDescribeResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionID", req.RegionId)
	builder.AddParam("projectID", req.ProjectId)
	builder.AddParam("bandwidthID", req.BandwidthId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &bandwidthDescribeRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	var eips []BandwidthDescribeEipsResponse
	for _, each := range result.Eips {
		eips = append(eips, BandwidthDescribeEipsResponse{
			Ip:    each.Ip,
			EipId: each.EipID,
		})
	}
	return &BandwidthDescribeResponse{
		Id:        result.Id,
		Status:    result.Status,
		Bandwidth: result.Bandwidth,
		Name:      result.Name,
		Eips:      eips,
	}, nil
}

type bandwidthDescribeRealResponse struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Bandwidth int    `json:"bandwidth"`
	Name      string `json:"name"`
	Eips      []struct {
		Ip    string `json:"ip"`
		EipID string `json:"eipID"`
	} `json:"eips"`
}

type BandwidthDescribeRequest struct {
	RegionId    string
	ProjectId   string
	BandwidthId string
}

type BandwidthDescribeEipsResponse struct {
	Ip    string
	EipId string
}

type BandwidthDescribeResponse struct {
	Id        string
	Status    string
	Bandwidth int
	Name      string
	Eips      []BandwidthDescribeEipsResponse
}
