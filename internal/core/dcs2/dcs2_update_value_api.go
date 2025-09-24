package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2UpdateValueApi
/* 更新集合类型对象中的oldKey对应的value，oldKey不存在则新增oldKey。针对List，Hash，Set，Sorted Set等集合类型。
 */type Dcs2UpdateValueApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2UpdateValueApi(client *core.CtyunClient) *Dcs2UpdateValueApi {
	return &Dcs2UpdateValueApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisDataMgr/updateValue",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2UpdateValueApi) Do(ctx context.Context, credential core.Credential, req *Dcs2UpdateValueRequest) (*Dcs2UpdateValueResponse, error) {
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
	var resp Dcs2UpdateValueResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2UpdateValueRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	GroupId    string `json:"groupId,omitempty"`    /*  DB编号  */
	Key        string `json:"key,omitempty"`        /*  key  */
	RawType    string `json:"type,omitempty"`       /*  key类型<li>hash<li>list<li>set<li>zset  */
	OldKey     string `json:"oldKey,omitempty"`     /*  hash的Key值,type=hash时必填。  */
	NewValue   string `json:"newValue,omitempty"`   /*  newValue值  */
	OldValue   string `json:"oldValue,omitempty"`   /*  oldValue值。type=set或zset时必填  */
	Index      string `json:"index,omitempty"`      /*  索引<br>说明：type=list时填写,用于确定修改的位置，其他不填  */
}

type Dcs2UpdateValueResponse struct {
	StatusCode int32                             `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                            `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2UpdateValueReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                            `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                            `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                            `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2UpdateValueReturnObjResponse struct {
	Total int32                                   `json:"total,omitempty"` /*  总数  */
	Rows  []*Dcs2UpdateValueReturnObjRowsResponse `json:"rows"`            /*  value值列表  */
}

type Dcs2UpdateValueReturnObjRowsResponse struct {
	TotalItems int32                                         `json:"totalItems,omitempty"` /*  集合key总数量  */
	NewKey     string                                        `json:"newKey,omitempty"`     /*  新key  */
	OldKey     string                                        `json:"oldKey,omitempty"`     /*  旧key  */
	GroupId    string                                        `json:"groupId,omitempty"`    /*  分组  */
	DbOrder    string                                        `json:"dbOrder,omitempty"`    /*  DB编号  */
	RawType    string                                        `json:"type,omitempty"`       /*  类型  */
	RedisIp    string                                        `json:"redisIp,omitempty"`    /*  节点IP  */
	RedisPort  string                                        `json:"redisPort,omitempty"`  /*  节点端口  */
	Key        string                                        `json:"key,omitempty"`        /*  key值  */
	ValueMap   *Dcs2UpdateValueReturnObjRowsValueMapResponse `json:"valueMap"`             /*  valueMap {key:value}  */
}

type Dcs2UpdateValueReturnObjRowsValueMapResponse struct{}
