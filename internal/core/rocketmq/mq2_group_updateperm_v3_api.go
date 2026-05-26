package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2GroupUpdatepermV3Api
/* 配置订阅组读取权限 */
type Mq2GroupUpdatepermV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GroupUpdatepermV3Api(client *core.CtyunClient) *Mq2GroupUpdatepermV3Api {
	return &Mq2GroupUpdatepermV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/group/updatePerm",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2GroupUpdatepermV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GroupUpdatepermV3Request) (*Mq2GroupUpdatepermV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2GroupUpdatepermV3Request
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
	var resp Mq2GroupUpdatepermV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GroupUpdatepermV3Request struct {
	RegionId       string   `json:"regionId"`       /*  资源池编码  */
	ProdInstId     string   `json:"prodInstId"`     /*  实例ID  */
	GroupName      string   `json:"groupName"`      /*  订阅组名字  */
	BrokerNameList []string `json:"brokerNameList"` /*  brokerNameList  */
	Enable         bool     `json:"enable"`         /*  是否开启  */
	Remark         string   `json:"remark"`         /*  备注  */
}

type Mq2GroupUpdatepermV3Response struct {
	StatusCode int32                                  `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2GroupUpdatepermV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                                `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
	Remark     *string                                `json:"remark"`     /*  备注  */
}

type Mq2GroupUpdatepermV3ReturnObjResponse struct{}
