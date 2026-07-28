package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqDiskExtendApi
/* 磁盘扩容。
 */type RocketmqDiskExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqDiskExtendApi(client *core.CtyunClient) *RocketmqDiskExtendApi {
	return &RocketmqDiskExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/diskExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *RocketmqDiskExtendApi) Do(ctx context.Context, credential core.Credential, req *RocketmqDiskExtendRequest) (*RocketmqDiskExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*RocketmqDiskExtendRequest
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
	var resp RocketmqDiskExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqDiskExtendRequest struct {
	RegionId       string `json:"regionId,omitempty"`       /*  实例的资源池 ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId     string `json:"prodInstId,omitempty"`     /*  实例 ID。  */
	DiskExtendSize int32  `json:"diskExtendSize,omitempty"` /*  每个节点扩容后的存储空间，单位 G  */
	AutoPay        bool   `json:"autoPay"`                  /*  是否自动支付。<br>- true：自动付费 (默认值)<br>- false：手动付费 <br>说明：选择为手动付费时，您需要在控制台的顶部菜单栏进入控制中心，单击费用中心，然后单击左侧导航栏的订单管理 > 我的订单，找到目标订单进行支付。  */
}

type RocketmqDiskExtendResponse struct {
	StatusCode int32                        `json:"statusCode"` // 接口系统层面状态码。成功："800"，失败："900"
	Message    string                       `json:"message"`    // 接口调用状态描述。成功时为"success"，失败时为具体失败信息
	ReturnObj  *RocketmqDiskExtendReturnObj `json:"returnObj"`  // 核心返回对象。成功时包含订单数据，失败时为空对象
	Error      string                       `json:"error"`      // 错误码。仅失败时返回，描述具体错误信息
}

type RocketmqDiskExtendReturnObj struct {
	Data *OrderData `json:"data"` // 订单核心数据。仅接口调用成功时返回
}

//type OrderData struct {
//	Submitted  bool   `json:"submitted"`  // 订单是否提交成功标识
//	NewOrderId string `json:"newOrderId"` // 系统生成的订单唯一标识 ID
//	NewOrderNo string `json:"newOrderNo"` // 系统生成的业务订单编号
//	TotalPrice string `json:"totalPrice"` // 订单总价格（单位：元）
//}
