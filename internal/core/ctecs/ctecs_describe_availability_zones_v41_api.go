package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDescribeAvailabilityZonesV41Api
/* 查询账户指定资源池中可用区的信息
 */type CtecsDescribeAvailabilityZonesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDescribeAvailabilityZonesV41Api(client *core.CtyunClient) *CtecsDescribeAvailabilityZonesV41Api {
	return &CtecsDescribeAvailabilityZonesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/availability-zones/details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDescribeAvailabilityZonesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDescribeAvailabilityZonesV41Request) (*CtecsDescribeAvailabilityZonesV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDescribeAvailabilityZonesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDescribeAvailabilityZonesV41Request struct {
	RegionID string /*  资源池ID  */
}

type CtecsDescribeAvailabilityZonesV41Response struct {
	StatusCode  int32                                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                              `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                              `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                              `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                              `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDescribeAvailabilityZonesV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDescribeAvailabilityZonesV41ReturnObjResponse struct {
	AzList []*CtecsDescribeAvailabilityZonesV41ReturnObjAzListResponse `json:"azList"` /*  可用区列表  */
}

type CtecsDescribeAvailabilityZonesV41ReturnObjAzListResponse struct {
	AzID   string `json:"azID,omitempty"`   /*  可用区ID  */
	AzName string `json:"azName,omitempty"` /*  可用区名称  */
}
