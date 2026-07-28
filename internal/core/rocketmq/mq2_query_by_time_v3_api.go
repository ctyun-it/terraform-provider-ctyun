package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Mq2QueryByTimeV3Api
/* 查询指定时间段内Topic的消息 */
type Mq2QueryByTimeV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryByTimeV3Api(client *core.CtyunClient) *Mq2QueryByTimeV3Api {
	return &Mq2QueryByTimeV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/message/queryByTime",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2QueryByTimeV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2QueryByTimeV3Request) (*Mq2QueryByTimeV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("beginTime", strconv.FormatInt(int64(req.BeginTime), 10))
	ctReq.AddParam("endTime", strconv.FormatInt(int64(req.EndTime), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2QueryByTimeV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryByTimeV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例id  */
	TopicName  string `json:"topicName"`  /*  主题名字  */
	BeginTime  int64  `json:"beginTime"`  /*  开始时间的毫秒时间戳  */
	EndTime    int64  `json:"endTime"`    /*  结束时间的毫秒时间戳  */
}

type Mq2QueryByTimeV3Response struct {
	StatusCode *string                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                            `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2QueryByTimeV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                            `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2QueryByTimeV3ReturnObjResponse struct{}
