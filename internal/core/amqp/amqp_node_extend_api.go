package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpNodeExtendApi
/* 节点扩容。
 */type AmqpNodeExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpNodeExtendApi(client *core.CtyunClient) *AmqpNodeExtendApi {
	return &AmqpNodeExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/nodeExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpNodeExtendApi) Do(ctx context.Context, credential core.Credential, req *AmqpNodeExtendRequest) (*AmqpNodeExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpNodeExtendRequest
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
	var resp AmqpNodeExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpNodeExtendRequest struct {
	RegionId      string `json:"regionId,omitempty"`      /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId    string `json:"prodInstId,omitempty"`    /*  实例ID。  实例需是集群版。  */
	ExtendNodeNum int32  `json:"extendNodeNum,omitempty"` /*  扩容后的节点数，取值范围5、7、9  */
	AutoPay       *bool  `json:"autoPay"`                 /*  是否自动支付，当实例为按需计费模式不生效。<br>- true：自动付费(默认值)<br>- false：手动付费 <br>说明：选择为手动付费时，您需要在控制台的顶部菜单栏进入控制中心，单击费用中心 ，然后单击左侧导航栏的订单管理 > 我的订单，找到目标订单进行支付。  */
}

type AmqpNodeExtendResponse struct {
	StatusCode string                           `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                           `json:"message"`    /*  描述状态。  */
	ReturnObj  *AmqpNodeExtendReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                           `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type AmqpNodeExtendReturnObjResponse struct {
	Data *AmqpNodeExtendReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type AmqpNodeExtendReturnObjDataResponse struct {
	Submitted  *bool   `json:"submitted"`  /*  是否成功提交。  */
	NewOrderId string  `json:"newOrderId"` /*  订单ID。  */
	NewOrderNo string  `json:"newOrderNo"` /*  订单号。  */
	TotalPrice float64 `json:"totalPrice"` /*  总价。  */
}
