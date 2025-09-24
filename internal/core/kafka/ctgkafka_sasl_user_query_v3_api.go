package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaSaslUserQueryV3Api
/* 查询SASL用户列表。
 */type CtgkafkaSaslUserQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaSaslUserQueryV3Api(client *core.CtyunClient) *CtgkafkaSaslUserQueryV3Api {
	return &CtgkafkaSaslUserQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/saslUser/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaSaslUserQueryV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaSaslUserQueryV3Request) (*CtgkafkaSaslUserQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.Username != "" {
		ctReq.AddParam("username", req.Username)
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
	var resp CtgkafkaSaslUserQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaSaslUserQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	Username   string `json:"username,omitempty"`   /*  用户名称，模糊查询。  */
	PageNum    string `json:"pageNum,omitempty"`    /*  分页中的页数，默认1，范围1-40000。  */
	PageSize   string `json:"pageSize,omitempty"`   /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaSaslUserQueryV3Response struct {
	StatusCode string                                    `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                    `json:"message,omitempty"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaSaslUserQueryV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaSaslUserQueryV3ReturnObjResponse struct {
	Data  []*CtgkafkaSaslUserQueryV3ReturnObjDataResponse `json:"data"`            /*  用户列表记录。  */
	Total int32                                           `json:"total,omitempty"` /*  总记录数。  */
}

type CtgkafkaSaslUserQueryV3ReturnObjDataResponse struct {
	Id          int32  `json:"id,omitempty"`          /*  用户ID。  */
	Username    string `json:"username,omitempty"`    /*  用户名。  */
	Description string `json:"description,omitempty"` /*  用户描述。  */
	Ctime       string `json:"ctime,omitempty"`       /*  创建时间。  */
}
