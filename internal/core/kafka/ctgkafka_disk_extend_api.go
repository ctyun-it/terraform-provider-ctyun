package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaDiskExtendApi
/* 磁盘扩容。
 */type CtgkafkaDiskExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaDiskExtendApi(client *core.CtyunClient) *CtgkafkaDiskExtendApi {
	return &CtgkafkaDiskExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/diskExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaDiskExtendApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaDiskExtendRequest) (*CtgkafkaDiskExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaDiskExtendRequest
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
	var resp CtgkafkaDiskExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaDiskExtendRequest struct {
	RegionId       string `json:"regionId,omitempty"`       /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId     string `json:"prodInstId,omitempty"`     /*  实例ID。  */
	DiskExtendSize string `json:"diskExtendSize,omitempty"` /*  每个节点扩容后的存储空间，单位GB，范围为当前每个节点存储空间 ~ 10000，并且为100的倍数。  */
	AutoPay        *bool  `json:"autoPay"`                  /*  是否自动支付，当实例为按需计费模式不生效。<br><li>true：自动付费(默认值)<br><li>false：手动付费 <br>说明：选择为手动付费时，您需要在控制台的顶部菜单栏进入控制中心，单击费用中心 ，然后单击左侧导航栏的订单管理 > 我的订单，找到目标订单进行支付。  */
}

type CtgkafkaDiskExtendResponse struct {
	StatusCode string                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                               `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaDiskExtendReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaDiskExtendReturnObjResponse struct {
}
