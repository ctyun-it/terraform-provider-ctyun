package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ExcuteCommandApi
/* 执行web-cli命令。
 */type Dcs2ExcuteCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ExcuteCommandApi(client *core.CtyunClient) *Dcs2ExcuteCommandApi {
	return &Dcs2ExcuteCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisDataMgr/excuteCommand",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ExcuteCommandApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ExcuteCommandRequest) (*Dcs2ExcuteCommandResponse, error) {
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
	var resp Dcs2ExcuteCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ExcuteCommandRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	GroupName  string `json:"groupName,omitempty"`  /*  DB编号  */
	Command    string `json:"command,omitempty"`    /*  命令  */
}

type Dcs2ExcuteCommandResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ExcuteCommandReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ExcuteCommandReturnObjResponse struct {
	ResultVO *Dcs2ExcuteCommandReturnObjResultVOResponse `json:"resultVO"` /*  执行结果  */
}

type Dcs2ExcuteCommandReturnObjResultVOResponse struct {
	RawType   int32  `json:"type,omitempty"`      /*  结果类型<li>1：字符串<li>2：列表<li>3：scan命令列表  */
	Index     string `json:"index,omitempty"`     /*  对应scan命令的游标  */
	ResultStr string `json:"resultStr,omitempty"` /*  字符串结果  */
}
