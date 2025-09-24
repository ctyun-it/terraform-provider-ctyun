package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2GetClientIPInfoApi
/* 获取节点客户端列表。
 */type Dcs2GetClientIPInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2GetClientIPInfoApi(client *core.CtyunClient) *Dcs2GetClientIPInfoApi {
	return &Dcs2GetClientIPInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/resourceMonitor/getClientIPInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2GetClientIPInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2GetClientIPInfoRequest) (*Dcs2GetClientIPInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("nodeName", req.NodeName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2GetClientIPInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2GetClientIPInfoRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	NodeName   string /*  节点名称<li>实例类型为Proxy集群、经典集群版时传入proxyName<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7741&isNormal=1&vid=270">查询实例的逻辑拓扑</a> AccessNode表proxyName字段<li>其他实例类型传入Redis节点名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7752&isNormal=1&vid=270">获取redis节点名列表</a> node表nodeName字段  */
}

type Dcs2GetClientIPInfoResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2GetClientIPInfoReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2GetClientIPInfoReturnObjResponse struct {
	InstanceId     string                                                `json:"instanceId,omitempty"`   /*  实例名称(实例Id)  */
	Skuname        string                                                `json:"skuname,omitempty"`      /*  规格名称  */
	Nodeaddr       string                                                `json:"nodeaddr,omitempty"`     /*  节点连接地址  */
	Total          int32                                                 `json:"total,omitempty"`        /*  客户端数量  */
	TotalActives   int32                                                 `json:"totalActives,omitempty"` /*  活跃客户端数量  */
	ClientInfoList []*Dcs2GetClientIPInfoReturnObjClientInfoListResponse `json:"clientInfoList"`         /*  客户端信息  */
	Clientmap      *Dcs2GetClientIPInfoReturnObjClientmapResponse        `json:"clientmap"`              /*  客户端map  */
	Cmdmap         *Dcs2GetClientIPInfoReturnObjCmdmapResponse           `json:"cmdmap"`                 /*  命令Map  */
	Dbmap          *Dcs2GetClientIPInfoReturnObjDbmapResponse            `json:"dbmap"`                  /*  数据库Map  */
}

type Dcs2GetClientIPInfoReturnObjClientInfoListResponse struct {
	ConnectionId string `json:"connectionId,omitempty"` /*  连接ID  */
	Addr         string `json:"addr,omitempty"`         /*  客户端连接地址  */
	Cmd          string `json:"cmd,omitempty"`          /*  命令  */
	Db           string `json:"db,omitempty"`           /*  数据库序号  */
	Idle         string `json:"idle,omitempty"`         /*  过期时间  */
	Age          string `json:"age,omitempty"`          /*  连接时长  */
}

type Dcs2GetClientIPInfoReturnObjClientmapResponse struct{}

type Dcs2GetClientIPInfoReturnObjCmdmapResponse struct{}

type Dcs2GetClientIPInfoReturnObjDbmapResponse struct{}
