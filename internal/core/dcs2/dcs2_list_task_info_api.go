package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ListTaskInfoApi
/* 查询数据迁移任务列表
 */type Dcs2ListTaskInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ListTaskInfoApi(client *core.CtyunClient) *Dcs2ListTaskInfoApi {
	return &Dcs2ListTaskInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/transfer/listTaskInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2ListTaskInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ListTaskInfoRequest) (*Dcs2ListTaskInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("pageNum", req.PageNum)
	ctReq.AddParam("pageSize", req.PageSize)
	if req.Status != "" {
		ctReq.AddParam("status", req.Status)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ListTaskInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ListTaskInfoRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	PageNum  string /*  页码（范围：> 0）  */
	PageSize string /*  数量（范围：[1,100]）  */
	Status   string /*  查询指定的任务状态<li>0：所有状态（默认）<li>1：运行中<li>2：成功<li>3：失败  */
}

type Dcs2ListTaskInfoResponse struct {
	StatusCode int32                              `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                             `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ListTaskInfoReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                             `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                             `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                             `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ListTaskInfoReturnObjResponse struct{}
