package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// ServiceListApi 根据条件查询云服务产品
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13942&data=114
type ServiceListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewServiceListApi(client *ctyunsdk.CtyunClient) *ServiceListApi {
	return &ServiceListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/service/queryCtapiServiceByCondition",
		},
	}
}

func (this *ServiceListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ServiceListRequest) (*ServiceListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&serviceListRealRequest{
		ServiceName: req.ServiceName,
		ServiceType: req.ServiceType,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp serviceListRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	var serviceList []ServiceListServiceListResponse
	for _, service := range resp.ServiceList {
		serviceList = append(serviceList, ServiceListServiceListResponse{
			ServiceCode:     service.ServiceCode,
			ServiceType:     service.ServiceType,
			MainServiceName: service.MainServiceName,
			ServiceDesc:     service.ServiceDesc,
			Id:              service.Id,
		})
	}
	return &ServiceListResponse{
		ServiceList: serviceList,
	}, nil
}

type serviceListRealRequest struct {
	ServiceName string `json:"serviceName,omitempty"`
	ServiceType int    `json:"serviceType,omitempty"`
}

type serviceListRealResponse struct {
	ServiceList []struct {
		ServiceCode     string `json:"serviceCode"`
		ServiceType     int    `json:"serviceType"`
		MainServiceName string `json:"mainServiceName"`
		ServiceDesc     string `json:"serviceDesc"`
		Id              int    `json:"id"`
	} `json:"serviceList"`
}

type ServiceListRequest struct {
	ServiceName string
	ServiceType int
}

type ServiceListServiceListResponse struct {
	ServiceCode     string
	ServiceType     int
	MainServiceName string
	ServiceDesc     string
	Id              int
}

type ServiceListResponse struct {
	ServiceList []ServiceListServiceListResponse
}
