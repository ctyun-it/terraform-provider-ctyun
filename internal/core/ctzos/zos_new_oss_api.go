package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosNewOssApi
/* 对象存储（简称ZOS）需开通服务后才可以使用，通过调用该接口可对指定资源池进行开通。
 */type ZosNewOssApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosNewOssApi(client *core.CtyunClient) *ZosNewOssApi {
	return &ZosNewOssApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/new",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosNewOssApi) Do(ctx context.Context, credential core.Credential, req *ZosNewOssRequest) (*ZosNewOssResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosNewOssResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosNewOssRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  区域 ID，开通集群互联资源池可传递”public“，否则传递资源池ID  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，要求单个云平台账户内唯一。请求后 24 小时内，使用相同的 clientToken 再次请求将忽略 regionID，查询 clientToken 对应开通单的状态  */
}

type ZosNewOssResponse struct {
	StatusCode  int64                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                        `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosNewOssReturnObjResponse   `json:"returnObj"`             /*  响应对象  */
	Description string                        `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                        `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                        `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
	ErrorDetail *ZosNewOssErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的对象存储订单业务相关的错误做明确的错误映射和提升，有唯一对应的 errorCode 。 其他订单侧(bss)的错误，以 oss.order.procFailed 的 errorCode 统一映射返回，并在 errorDetail 中返回订单侧的详细错误信息  */
}

type ZosNewOssReturnObjResponse struct {
	MasterOrderID        string                               `json:"masterOrderID,omitempty"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO        string                               `json:"masterOrderNO,omitempty"`        /*  订单号，可为 null  */
	RegionID             string                               `json:"regionID,omitempty"`             /*  资源所属资源池ID，若为集群互联资源池则返回”public“  */
	MasterResourceID     string                               `json:"masterResourceID,omitempty"`     /*  主资源ID。对象存储场景下，无需关心  */
	MasterResourceStatus string                               `json:"masterResourceStatus,omitempty"` /*  主资源状态。只有主订单资源会返回  */
	Resources            *ZosNewOssReturnObjResourcesResponse `json:"resources"`                      /*  资源明细列表  */
}

type ZosNewOssErrorDetailResponse struct {
	BssErrCode       string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中  */
	BssErrMsg        string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中  */
	BssOrigErr       string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息  */
	BssErrPrefixHint string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息  */
}

type ZosNewOssReturnObjResourcesResponse struct {
	ZosID            string `json:"zosID,omitempty"`            /*  单项资源的变配、续订、退订等需要该资源项的ID  */
	OrderID          string `json:"orderID,omitempty"`          /*  无需关心  */
	StartTime        int64  `json:"startTime,omitempty"`        /*  启动时刻，epoch时戳，毫秒精度  */
	CreateTime       int64  `json:"createTime,omitempty"`       /*  创建时刻，epoch时戳，毫秒精度  */
	UpdateTime       int64  `json:"updateTime,omitempty"`       /*  更新时刻，epoch时戳，毫秒精度  */
	Status           int64  `json:"status,omitempty"`           /*  资源状态，无需关心  */
	IsMaster         *bool  `json:"isMaster"`                   /*  是否是主资源项  */
	ItemValue        int64  `json:"itemValue,omitempty"`        /*  资源规格，对象存储场景下，无需关心  */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型，对象存储服务开通固定为：ZOS_REG  */
	MasterResourceID string `json:"masterResourceID,omitempty"` /*  主资源ID。对象存储场景下，无需关心  */
	MasterOrderID    string `json:"masterOrderID,omitempty"`    /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用 materOrderID 进一步确认订单状态及资源状态  */
}
