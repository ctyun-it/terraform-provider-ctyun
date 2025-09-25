package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaListTagV3Api
/* 查询标签列表v3
 */type CtgkafkaListTagV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaListTagV3Api(client *core.CtyunClient) *CtgkafkaListTagV3Api {
	return &CtgkafkaListTagV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/resourceTag/listTag",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaListTagV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaListTagV3Request) (*CtgkafkaListTagV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.TagName != "" {
		ctReq.AddParam("tagName", req.TagName)
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
	var resp CtgkafkaListTagV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaListTagV3Request struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池编码  */
	TagName  string `json:"tagName,omitempty"`  /*  标签名称  */
	PageNum  string `json:"pageNum,omitempty"`  /*  分页中的页数，默认1  */
	PageSize string `json:"pageSize,omitempty"` /*  分页中的每页大小，默认10  */
}

type CtgkafkaListTagV3Response struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message,omitempty"`    /*  状态信息  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息  */
}
