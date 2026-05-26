package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2GroupSubDetailV3Api
/* 查看订阅组订阅信息 */
type Mq2GroupSubDetailV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GroupSubDetailV3Api(client *core.CtyunClient) *Mq2GroupSubDetailV3Api {
	return &Mq2GroupSubDetailV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/group/subDetail",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2GroupSubDetailV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GroupSubDetailV3Request) (*Mq2GroupSubDetailV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2GroupSubDetailV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GroupSubDetailV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  消费组名  */
}

type Mq2GroupSubDetailV3Response struct {
	StatusCode *string                               `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败："900"  */
	Message    *string                               `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2GroupSubDetailV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                               `json:"error"`      /*  错误码，描述错误信息  */
}

type Mq2GroupSubDetailV3ReturnObjResponse struct {
	Data *Mq2GroupSubDetailV3ReturnObjDataResponse `json:"data"` /*  订阅相关数据对象  */
}

type Mq2GroupSubDetailV3ReturnObjDataResponse struct {
	ProdInstId          *string                                                      `json:"prodInstId"`          /*  实例id  */
	GroupName           *string                                                      `json:"groupName"`           /*  消费组名称  */
	SubscriptionDataMap *Mq2GroupSubDetailV3ReturnObjDataSubscriptionDataMapResponse `json:"subscriptionDataMap"` /*  订阅数据映射表（没有消费组实例在线时为空）  */
}

type Mq2GroupSubDetailV3ReturnObjDataSubscriptionDataMapResponse struct{}
