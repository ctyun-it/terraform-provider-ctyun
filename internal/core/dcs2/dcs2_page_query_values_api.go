package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2PageQueryValuesApi
/* 应用数据管理value值分页查询。
 */type Dcs2PageQueryValuesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2PageQueryValuesApi(client *core.CtyunClient) *Dcs2PageQueryValuesApi {
	return &Dcs2PageQueryValuesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisDataMgr/pageQueryValues",
			ContentType:  "",
		},
	}
}

func (a *Dcs2PageQueryValuesApi) Do(ctx context.Context, credential core.Credential, req *Dcs2PageQueryValuesRequest) (*Dcs2PageQueryValuesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupId", req.GroupId)
	ctReq.AddParam("key", req.Key)
	if req.RawType != "" {
		ctReq.AddParam("type", req.RawType)
	}
	if req.Start != 0 {
		ctReq.AddParam("start", strconv.FormatInt(int64(req.Start), 10))
	}
	if req.PerPage != 0 {
		ctReq.AddParam("perPage", strconv.FormatInt(int64(req.PerPage), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2PageQueryValuesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2PageQueryValuesRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	GroupId    string /*  DB编号  */
	Key        string /*  key  */
	RawType    string /*  key类型  */
	Start      int32  /*  开始下标  */
	PerPage    int32  /*  每页显示数量  */
}

type Dcs2PageQueryValuesResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2PageQueryValuesReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2PageQueryValuesReturnObjResponse struct {
	Total int32                                       `json:"total,omitempty"` /*  总数  */
	Rows  []*Dcs2PageQueryValuesReturnObjRowsResponse `json:"rows"`            /*  value值列表  */
}

type Dcs2PageQueryValuesReturnObjRowsResponse struct {
	TotalItems   int32                                                 `json:"totalItems,omitempty"` /*  集合key总数量  */
	TotalPage    int32                                                 `json:"totalPage,omitempty"`  /*  总页数  */
	GroupId      string                                                `json:"groupId,omitempty"`    /*  分组  */
	DbOrder      string                                                `json:"dbOrder,omitempty"`    /*  DB编号  */
	Start        int32                                                 `json:"start,omitempty"`      /*  开始下标  */
	Index        int32                                                 `json:"index,omitempty"`      /*  当前下标  */
	RawType      string                                                `json:"type,omitempty"`       /*  类型  */
	PerPage      int32                                                 `json:"perPage,omitempty"`    /*  每页显示数量  */
	RedisIp      string                                                `json:"redisIp,omitempty"`    /*  key所在节点IP  */
	RedisPort    string                                                `json:"redisPort,omitempty"`  /*  key所在节点端口  */
	Key          string                                                `json:"key,omitempty"`        /*  key值  */
	ByteValueMap *Dcs2PageQueryValuesReturnObjRowsByteValueMapResponse `json:"byteValueMap"`         /*  键值对，KEY是字符串，VALUE是字节数组  */
}

type Dcs2PageQueryValuesReturnObjRowsByteValueMapResponse struct{}
