package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2QueryManagerApi
/* 查询分组KEY对应的value列表
 */type Dcs2QueryManagerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryManagerApi(client *core.CtyunClient) *Dcs2QueryManagerApi {
	return &Dcs2QueryManagerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisDataMgr/queryManager",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryManagerApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryManagerRequest) (*Dcs2QueryManagerResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupId", req.GroupId)
	ctReq.AddParam("key", req.Key)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryManagerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryManagerRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	GroupId    string /*  DB编号  */
	Key        string /*  key值  */
}

type Dcs2QueryManagerResponse struct {
	StatusCode int32                              `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                             `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2QueryManagerReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                             `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                             `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                             `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2QueryManagerReturnObjResponse struct {
	Total int32                                    `json:"total,omitempty"` /*  总数  */
	Rows  []*Dcs2QueryManagerReturnObjRowsResponse `json:"rows"`            /*  value值列表  */
}

type Dcs2QueryManagerReturnObjRowsResponse struct {
	GroupId   string  `json:"groupId,omitempty"`   /*  分组  */
	DbOrder   string  `json:"dbOrder,omitempty"`   /*  DB编号  */
	RedisIp   string  `json:"redisIp,omitempty"`   /*  节点IP  */
	RedisPort string  `json:"redisPort,omitempty"` /*  节点端口  */
	Key       string  `json:"key,omitempty"`       /*  key  */
	ByteValue []int32 `json:"byteValue"`           /*  字节数组，内存  */
}
