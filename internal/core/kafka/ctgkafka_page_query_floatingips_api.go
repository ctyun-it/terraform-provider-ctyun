package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaPageQueryFloatingipsApi
/* 查询可绑定的弹性IP。
 */type CtgkafkaPageQueryFloatingipsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaPageQueryFloatingipsApi(client *core.CtyunClient) *CtgkafkaPageQueryFloatingipsApi {
	return &CtgkafkaPageQueryFloatingipsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/pageQueryFloatingips",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaPageQueryFloatingipsApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaPageQueryFloatingipsRequest) (*CtgkafkaPageQueryFloatingipsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(req.PageNum), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaPageQueryFloatingipsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaPageQueryFloatingipsRequest struct {
	RegionId string `json:"regionId,omitempty"` /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	PageNum  int32  `json:"pageNum,omitempty"`  /*  分页中的页数，默认1，范围1-40000。  */
	PageSize int32  `json:"pageSize,omitempty"` /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaPageQueryFloatingipsResponse struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释。  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}
