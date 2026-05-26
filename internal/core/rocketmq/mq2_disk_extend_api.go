package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2DiskExtendApi
/* 实例磁盘扩容 */
type Mq2DiskExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2DiskExtendApi(client *core.CtyunClient) *Mq2DiskExtendApi {
	return &Mq2DiskExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/diskExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2DiskExtendApi) Do(ctx context.Context, credential core.Credential, req *Mq2DiskExtendRequest) (*Mq2DiskExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2DiskExtendRequest
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
	var resp Mq2DiskExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2DiskExtendRequest struct {
	RegionId       string `json:"regionId"`       /*  资源池编码  */
	ProdInstId     string `json:"prodInstId"`     /*  实例id  */
	DiskExtendSize string `json:"diskExtendSize"` /*  每个节点扩容后的存储空间，单位G  */
}

type Mq2DiskExtendResponse struct {
	StatusCode *string                         `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                         `json:"message"`    /*  接口调用状态描述。成功时为"success"，失败时为具体失败信息  */
	ReturnObj  *Mq2DiskExtendReturnObjResponse `json:"returnObj"`  /*  核心返回对象。成功时包含订单数据，失败时为空对象  */
	Error      *string                         `json:"error"`      /*  错误码。仅失败时返回，描述具体错误信息  */
}

type Mq2DiskExtendReturnObjResponse struct {
	Data *Mq2DiskExtendReturnObjDataResponse `json:"data"` /*  订单核心数据。仅接口调用成功时返回  */
}

type Mq2DiskExtendReturnObjDataResponse struct {
	Submitted  *bool   `json:"submitted"`  /*  订单是否提交成功标识  */
	NewOrderId *string `json:"newOrderId"` /*  系统生成的订单唯一标识ID  */
	NewOrderNo *string `json:"newOrderNo"` /*  系统生成的业务订单编号  */
	TotalPrice *string `json:"totalPrice"` /*  订单总价格（单位：元）  */
}
