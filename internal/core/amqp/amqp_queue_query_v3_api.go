package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueueQueryV3Api
/* 查询队列v3
 */type AmqpQueueQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueueQueryV3Api(client *core.CtyunClient) *AmqpQueueQueryV3Api {
	return &AmqpQueueQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/queue/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpQueueQueryV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpQueueQueryV3Request) (*AmqpQueueQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("vhost", req.Vhost)
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	if req.PageNum != "" {
		ctReq.AddParam("pageNum", req.PageNum)
	}
	ctReq.AddParam("pageSize", req.PageSize)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpQueueQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpQueueQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
	Name       string `json:"name,omitempty"`       /*  队列名称  */
	PageNum    string `json:"pageNum,omitempty"`    /*  分页的页序号  */
	PageSize   string `json:"pageSize,omitempty"`   /*  分页的大小  */
}

type AmqpQueueQueryV3Response struct {
	ReturnObj  *AmqpQueueQueryV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                             `json:"message"`    /*  描述状态  */
	StatusCode string                             `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                             `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpQueueQueryV3ReturnObjResponse struct {
	Data *AmqpQueueQueryV3ReturnObjData `json:"data"` /*  返回数据  */
}

type AmqpQueueQueryV3ReturnObjData struct {
	FilteredCount int                                  `json:"filtered_count"`
	ItemCount     int                                  `json:"item_count"`
	Items         []*AmqpQueueQueryV3ReturnObjDataItem `json:"items"`
	Page          int                                  `json:"page"`
	PageCount     int                                  `json:"page_count"`
	PageSize      int                                  `json:"page_size"`
	TotalCount    int                                  `json:"total_count"`
}

type AmqpQueueQueryV3ReturnObjDataItem struct {
	Vhost      string `json:"vhost"`
	Durable    bool   `json:"durable"`
	Name       string `json:"name"`
	AutoDelete bool   `json:"auto_delete"`
	Consumers  int    `json:"consumers"`
}
