package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DestroyInstanceApi
/* 销毁已退订的分布式缓存Redis实例。
 */type Dcs2DestroyInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DestroyInstanceApi(client *core.CtyunClient) *Dcs2DestroyInstanceApi {
	return &Dcs2DestroyInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/lifeCycleServant/destroyInstance",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2DestroyInstanceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DestroyInstanceRequest) (*Dcs2DestroyInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Dcs2DestroyInstanceRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DestroyInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DestroyInstanceRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池ID，可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口，使用resPoolCode字段。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type Dcs2DestroyInstanceResponse struct {
	StatusCode int32                                 `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2DestroyInstanceReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                                `json:"requestId"`  /*  请求 ID。  */
	Code       string                                `json:"code"`       /*  响应码描述。  */
	Error      string                                `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2DestroyInstanceReturnObjResponse struct {
	ErrorMessage      string                                                   `json:"errorMessage"`      /*  错误信息。  */
	OrderPlacedEvents []*Dcs2DestroyInstanceReturnObjOrderPlacedEventsResponse `json:"orderPlacedEvents"` /*  收费项。  */
}

type Dcs2DestroyInstanceReturnObjOrderPlacedEventsResponse struct {
	ErrorMessage string `json:"errorMessage"` /*  错误信息。  */
	Submitted    *bool  `json:"submitted"`    /*  是否提交。  */
	NewOrderId   string `json:"newOrderId"`   /*  订单ID。  */
	NewOrderNo   string `json:"newOrderNo"`   /*  订单号。  */
}
