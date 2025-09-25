package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2GetRdbDownLoadUrlApi
/* 获取指定实例的备份文件下载链接，下载备份文件。
 */type Dcs2GetRdbDownLoadUrlApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2GetRdbDownLoadUrlApi(client *core.CtyunClient) *Dcs2GetRdbDownLoadUrlApi {
	return &Dcs2GetRdbDownLoadUrlApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisMgr/getRdbDownLoadUrl",
			ContentType:  "",
		},
	}
}

func (a *Dcs2GetRdbDownLoadUrlApi) Do(ctx context.Context, credential core.Credential, req *Dcs2GetRdbDownLoadUrlRequest) (*Dcs2GetRdbDownLoadUrlResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("restoreName", req.RestoreName)
	ctReq.AddParam("ipType", req.IpType)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2GetRdbDownLoadUrlResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2GetRdbDownLoadUrlRequest struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId  string /*  实例ID  */
	RestoreName string /*  备份名<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7778&isNormal=1&vid=270">查询备份文件信息</a> Restore表restoreName字段  */
	IpType      string /*  网络类型<li>publicIp：公网IP<li>privateIp：私网IP  */
}

type Dcs2GetRdbDownLoadUrlResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                  `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2GetRdbDownLoadUrlReturnObjResponse `json:"returnObj"`            /*  返回数据对象，包含备份文件下载链接：<br>key格式：Redis节点名<br>value格式：备份文件下载URL  */
	RequestId  string                                  `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                  `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2GetRdbDownLoadUrlReturnObjResponse struct{}
