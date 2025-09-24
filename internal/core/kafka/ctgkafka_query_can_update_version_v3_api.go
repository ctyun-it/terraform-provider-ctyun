package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaQueryCanUpdateVersionV3Api
/* 查询可升级实例版本。
 */type CtgkafkaQueryCanUpdateVersionV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQueryCanUpdateVersionV3Api(client *core.CtyunClient) *CtgkafkaQueryCanUpdateVersionV3Api {
	return &CtgkafkaQueryCanUpdateVersionV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/queryCanUpdateVersion",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaQueryCanUpdateVersionV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQueryCanUpdateVersionV3Request) (*CtgkafkaQueryCanUpdateVersionV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaQueryCanUpdateVersionV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQueryCanUpdateVersionV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type CtgkafkaQueryCanUpdateVersionV3Response struct {
	StatusCode string                                            `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                            `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaQueryCanUpdateVersionV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                            `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaQueryCanUpdateVersionV3ReturnObjResponse struct {
	Data []string `json:"data"` /*  返回可升级的版本号列表。  */
}
