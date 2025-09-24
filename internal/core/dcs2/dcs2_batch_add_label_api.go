package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2BatchAddLabelApi
/* 批量创建标签
 */type Dcs2BatchAddLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2BatchAddLabelApi(client *core.CtyunClient) *Dcs2BatchAddLabelApi {
	return &Dcs2BatchAddLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/label/batchAdd",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2BatchAddLabelApi) Do(ctx context.Context, credential core.Credential, req *Dcs2BatchAddLabelRequest) (*Dcs2BatchAddLabelResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2BatchAddLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2BatchAddLabelRequest struct {
	RegionId    string                                 /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	LabelKVList []*Dcs2BatchAddLabelLabelKVListRequest `json:"labelKVList"` /*  标签列表对象  */
}

type Dcs2BatchAddLabelLabelKVListRequest struct {
	Key   string `json:"key,omitempty"`   /*  标签键<br>说明：长度限制1-32个字符  */
	Value string `json:"value,omitempty"` /*  标签值<br>说明：长度限制1-32个字符  */
}

type Dcs2BatchAddLabelResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2BatchAddLabelReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2BatchAddLabelReturnObjResponse struct{}
