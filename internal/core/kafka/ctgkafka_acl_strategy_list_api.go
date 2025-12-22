package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaAclStrategyListApi
/* 查询ACL策略列表。
 */type CtgkafkaAclStrategyListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaAclStrategyListApi(client *core.CtyunClient) *CtgkafkaAclStrategyListApi {
	return &CtgkafkaAclStrategyListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/kafkaAclStrategy/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaAclStrategyListApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaAclStrategyListRequest) (*CtgkafkaAclStrategyListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	if req.PageNum != "" {
		ctReq.AddParam("pageNum", req.PageNum)
	}
	if req.PageSize != "" {
		ctReq.AddParam("pageSize", req.PageSize)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaAclStrategyListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaAclStrategyListRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	Name       string `json:"name,omitempty"`       /*  策略名称，模糊查询。  */
	PageNum    string `json:"pageNum,omitempty"`    /*  分页中的页数，默认1，范围1-40000。  */
	PageSize   string `json:"pageSize,omitempty"`   /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaAclStrategyListResponse struct {
	StatusCode string                                    `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                    `json:"message"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaAclStrategyListReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                                    `json:"error"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaAclStrategyListReturnObjResponse struct {
	Data  []*CtgkafkaAclStrategyListReturnObjDataResponse `json:"data"`  /*  用户列表记录。  */
	Total int32                                           `json:"total"` /*  总记录数。  */
}

type CtgkafkaAclStrategyListReturnObjDataResponse struct {
	Id          int32  `json:"id"`          /*  策略ID。  */
	StrategyId  string `json:"strategyId"`  /*  策略唯一id。  */
	ClusterId   string `json:"clusterId"`   /*  实例ID。  */
	Name        string `json:"name"`        /*  策略名称。  */
	UseNewTopic int32  `json:"useNewTopic"` /*  是否应用到新增主题 <br><li>1：是<br><li>2：否  */
	CreateTime  string `json:"createTime"`  /*  创建时间。  */
}
