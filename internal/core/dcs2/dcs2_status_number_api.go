package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2StatusNumberApi
/* 查询该租户在当前区域下不同状态的实例数。
 */type Dcs2StatusNumberApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2StatusNumberApi(client *core.CtyunClient) *Dcs2StatusNumberApi {
	return &Dcs2StatusNumberApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/statusNumber",
			ContentType:  "",
		},
	}
}

func (a *Dcs2StatusNumberApi) Do(ctx context.Context, credential core.Credential, req *Dcs2StatusNumberRequest) (*Dcs2StatusNumberResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.IncludeFailure != "" {
		ctReq.AddParam("includeFailure", req.IncludeFailure)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2StatusNumberResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2StatusNumberRequest struct {
	RegionId       string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	IncludeFailure string /*  否  */
}

type Dcs2StatusNumberResponse struct {
	StatusCode int32                              `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                             `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2StatusNumberReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                             `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                             `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                             `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2StatusNumberReturnObjResponse struct {
	RunningCount      int32 `json:"runningCount,omitempty"`      /*  运行中的实例数  */
	CreatingCount     int32 `json:"creatingCount,omitempty"`     /*  创建中的实例数  */
	FrozenCount       int32 `json:"frozenCount,omitempty"`       /*  冻结中的实例数  */
	ChangingCount     int32 `json:"changingCount,omitempty"`     /*  变更中的实例数  */
	CreateFailedCount int32 `json:"createFailedCount,omitempty"` /*  创建失败的实例数  */
	RestartingCount   int32 `json:"restartingCount,omitempty"`   /*  重启中的实例数  */
	ErrorCount        int32 `json:"errorCount,omitempty"`        /*  异常的实例数  */
}
