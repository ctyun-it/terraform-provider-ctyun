package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2InstanceQueryPriceForDiskExtendApi
/* 磁盘扩容查价 */
type Mq2InstanceQueryPriceForDiskExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2InstanceQueryPriceForDiskExtendApi(client *core.CtyunClient) *Mq2InstanceQueryPriceForDiskExtendApi {
	return &Mq2InstanceQueryPriceForDiskExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/queryPriceForDiskExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2InstanceQueryPriceForDiskExtendApi) Do(ctx context.Context, credential core.Credential, req *Mq2InstanceQueryPriceForDiskExtendRequest) (*Mq2InstanceQueryPriceForDiskExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2InstanceQueryPriceForDiskExtendRequest
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
	var resp Mq2InstanceQueryPriceForDiskExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2InstanceQueryPriceForDiskExtendRequest struct {
	RegionId       string `json:"regionId"`       /*  资源池编码  */
	ProdInstId     string `json:"prodInstId"`     /*  实例ID  */
	DiskExtendSize int32  `json:"diskExtendSize"` /*  每个节点扩容后的存储空间，单位G  */
}

type Mq2InstanceQueryPriceForDiskExtendResponse struct {
	StatusCode *string                                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                              `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2InstanceQueryPriceForDiskExtendReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                                              `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2InstanceQueryPriceForDiskExtendReturnObjResponse struct {
	TotalPrice      *string                                                               `json:"totalPrice"`      /*  子订单总价格  */
	ServiceTag      *string                                                               `json:"serviceTag"`      /*  服务标签  */
	FinalPrice      *string                                                               `json:"finalPrice"`      /*  子订单最终价格  */
	OrderItemPrices []*Mq2InstanceQueryPriceForDiskExtendReturnObjOrderItemPricesResponse `json:"orderItemPrices"` /*  订单子项价格列表  */
}

type Mq2InstanceQueryPriceForDiskExtendReturnObjOrderItemPricesResponse struct {
	ItemId       *string `json:"itemId"`       /*  项目id  */
	TotalPrice   *string `json:"totalPrice"`   /*  项目总价格  */
	FinalPrice   *string `json:"finalPrice"`   /*  项目最终价格  */
	ResourceType *string `json:"resourceType"` /*  资源类型，PAAS_DMQ为rocketmq实例，PAAS_DMQ_EBS为云硬盘  */
}
