package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2AddValueApi
/* 向集合类型中插入数据，针对List，Hash，Set，Sorted Set等集合类型。
 */type Dcs2AddValueApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2AddValueApi(client *core.CtyunClient) *Dcs2AddValueApi {
	return &Dcs2AddValueApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisDataMgr/addValue",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2AddValueApi) Do(ctx context.Context, credential core.Credential, req *Dcs2AddValueRequest) (*Dcs2AddValueResponse, error) {
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
	var resp Dcs2AddValueResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2AddValueRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	GroupId    string `json:"groupId,omitempty"`    /*  DB编号  */
	Key        string `json:"key,omitempty"`        /*  key  */
	RawType    string `json:"type,omitempty"`       /*  key类型<li>hash<li>list<li>set<li>zset  */
	NewKey     string `json:"newKey,omitempty"`     /*  hash的Key值<br>说明：type=hash时必填，其他不填  */
	Score      string `json:"score,omitempty"`      /*  score值,type=zset时,必填  */
	NewValue   string `json:"newValue,omitempty"`   /*  集合对象中的newValue值  */
}

type Dcs2AddValueResponse struct {
	StatusCode int32                          `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                         `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2AddValueReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                         `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                         `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                         `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2AddValueReturnObjResponse struct {
	Total int32                              `json:"total,omitempty"` /*  总数  */
	Rows  *Dcs2AddValueReturnObjRowsResponse `json:"rows"`            /*  value值列表  */
}

type Dcs2AddValueReturnObjRowsResponse struct {
	TotalItems int32                                      `json:"totalItems,omitempty"` /*  集合key总数量  */
	NewKey     string                                     `json:"newKey,omitempty"`     /*  新key  */
	OldKey     string                                     `json:"oldKey,omitempty"`     /*  旧key  */
	GroupId    string                                     `json:"groupId,omitempty"`    /*  分组  */
	DbOrder    string                                     `json:"dbOrder,omitempty"`    /*  DB编号  */
	RawType    string                                     `json:"type,omitempty"`       /*  类型  */
	RedisIp    string                                     `json:"redisIp,omitempty"`    /*  节点IP  */
	RedisPort  string                                     `json:"redisPort,omitempty"`  /*  节点端口  */
	Key        string                                     `json:"key,omitempty"`        /*  key值  */
	ValueMap   *Dcs2AddValueReturnObjRowsValueMapResponse `json:"valueMap"`             /*  valueMap {key:value}  */
}

type Dcs2AddValueReturnObjRowsValueMapResponse struct{}
