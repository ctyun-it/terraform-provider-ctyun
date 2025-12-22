package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpExchangeQueryV3Api
/* 查询交换器v3
 */type AmqpExchangeQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpExchangeQueryV3Api(client *core.CtyunClient) *AmqpExchangeQueryV3Api {
	return &AmqpExchangeQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/exchange/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpExchangeQueryV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpExchangeQueryV3Request) (*AmqpExchangeQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.Vhost != "" {
		ctReq.AddParam("vhost", req.Vhost)
	}
	if req.PageNum != "" {
		ctReq.AddParam("pageNum", req.PageNum)
	}
	if req.PageSize != "" {
		ctReq.AddParam("pageSize", req.PageSize)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpExchangeQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpExchangeQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  虚拟机名称  */
	PageNum    string `json:"pageNum,omitempty"`    /*  当前页面（默认1）  */
	PageSize   string `json:"pageSize,omitempty"`   /*  分页大小（默认100）  */
	Name       string `json:"name,omitempty"`       /*  交换器名称（模糊匹配）  */
}

type AmqpExchangeQueryV3Response struct {
	ReturnObj  *AmqpExchangeQueryV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                                `json:"message"`    /*  描述状态  */
	StatusCode string                                `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                                `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpExchangeQueryV3ReturnObjResponse struct {
	Data *AmqpExchangeQueryV3ReturnObjDataResponse `json:"data"` /*  返回数据  */
}

type AmqpExchangeQueryV3ReturnObjDataResponse struct {
	Filtered_count int32                                            `json:"filtered_count"` /*  过滤后的总交换器  */
	Item_count     int32                                            `json:"item_count"`     /*  当前页面上的交换器数量  */
	Items          []*AmqpExchangeQueryV3ReturnObjDataItemsResponse `json:"items"`          /*  交换器详细信息  */
	Page           int32                                            `json:"page"`           /*  当前页数  */
	Page_count     int32                                            `json:"page_count"`     /*  总页数  */
	Page_size      int32                                            `json:"page_size"`      /*  分页设置的每个页面的最多队列数  */
	Total_count    int32                                            `json:"total_count"`    /*  总交换器数  */
}

type AmqpExchangeQueryV3ReturnObjDataItemsResponse struct {
	Auto_delete bool   `json:"auto_delete"` /*  是否自动删除  */
	Durable     bool   `json:"durable"`     /*  是否持久化，默认都是持久化  */
	Internal    bool   `json:"internal"`
	Name        string `json:"name"`  /*  交换器名称  */
	RawType     string `json:"type"`  /*  交换器类型  */
	Vhost       string `json:"vhost"` /*  虚拟机名称  */
	Argument    struct {
		XDelayedType string `json:"x-delayed-type"`
	} `json:"argument"` /*  交换器参数  */
}
