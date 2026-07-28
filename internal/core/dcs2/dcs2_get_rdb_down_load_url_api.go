package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2GetRdbDownLoadUrlApi
/* 获取指定实例的备份文件下载链接，下载备份文件。 */
type Dcs2GetRdbDownLoadUrlApi struct {
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
		if resp.StatusCode == 0 {
			resp.StatusCode = 900
			return &resp, err
		}
		return nil, err
	}
	return &resp, nil
}

type Dcs2GetRdbDownLoadUrlRequest struct {
	RegionId    string `json:"regionId"`    /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	ProdInstId  string `json:"prodInstId"`  /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	RestoreName string `json:"restoreName"` /*  备份名，可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7778&isNormal=1&vid=270">查询备份文件信息</a> 接口，使用Restore表restoreName字段。  */
	IpType      string `json:"ipType"`      /*  网络类型，可选值：<li>publicIp：公网IP。<li>privateIp：私网IP。  */
}

type Dcs2GetRdbDownLoadUrlResponse struct {
	StatusCode int32       `json:"statusCode"`          /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    *string     `json:"message,omitempty"`   /*  响应信息。  */
	ReturnObj  interface{} `json:"returnObj"`           /*  返回数据对象，包含备份文件下载链接：<br>key格式：Redis节点名<br>value格式：备份文件下载URL  */
	RequestId  *string     `json:"requestId,omitempty"` /*  请求 ID。  */
	Code       *string     `json:"code,omitempty"`      /*  响应码，仅表示请求是否执行。  */
	Error      *string     `json:"error,omitempty"`     /*  错误码，参见错误码说明。  */
}
