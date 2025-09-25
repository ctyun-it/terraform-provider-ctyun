package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2EnableRenameCommandApi
/* 启用危险命令重命名功能并配置重命名规则。
 */type Dcs2EnableRenameCommandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2EnableRenameCommandApi(client *core.CtyunClient) *Dcs2EnableRenameCommandApi {
	return &Dcs2EnableRenameCommandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2.1/securityMgrServant/enableRenameCommand",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2EnableRenameCommandApi) Do(ctx context.Context, credential core.Credential, req *Dcs2EnableRenameCommandRequest) (*Dcs2EnableRenameCommandResponse, error) {
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
	var resp Dcs2EnableRenameCommandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2EnableRenameCommandRequest struct {
	RegionId       string                                        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId     string                                        `json:"prodInstId,omitempty"` /*  实例ID  */
	RenameCommands *Dcs2EnableRenameCommandRenameCommandsRequest `json:"renameCommands"`       /*  重命名的命令<br>参数为JSON对象,格式为键值对（原命令名:重命名后命令名）.<br>支持命令：COMMAND,KEYS,FLUSHDB,FLUSHALL,HGETALL,SCAN,HSCAN,SSCAN,ZSCAN,SHUTDOWN,CONFIG"  */
}

type Dcs2EnableRenameCommandRenameCommandsRequest struct{}

type Dcs2EnableRenameCommandResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2EnableRenameCommandReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2EnableRenameCommandReturnObjResponse struct {
	RenameStatus *bool `json:"renameStatus"` /*  危险命令开关状态<li>false：未开启</li><li>true：开启</li>  */
}
