package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2QueryGroupKeyApi
/* 查询分组KEY
 */type Dcs2QueryGroupKeyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryGroupKeyApi(client *core.CtyunClient) *Dcs2QueryGroupKeyApi {
	return &Dcs2QueryGroupKeyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisDataMgr/queryGroupKey",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryGroupKeyApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryGroupKeyRequest) (*Dcs2QueryGroupKeyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupId", req.GroupId)
	if req.Count != 0 {
		ctReq.AddParam("count", strconv.FormatInt(int64(req.Count), 10))
	}
	if req.Cursor != "" {
		ctReq.AddParam("cursor", req.Cursor)
	}
	if req.Key != "" {
		ctReq.AddParam("key", req.Key)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryGroupKeyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryGroupKeyRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	GroupId    string /*  DB编号  */
	Count      int32  /*  查询数量，默认10  */
	Cursor     string /*  游标  */
	Key        string /*  匹配的key值  */
}

type Dcs2QueryGroupKeyResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2QueryGroupKeyReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2QueryGroupKeyReturnObjResponse struct {
	Cursor string   `json:"cursor,omitempty"` /*  游标  */
	Total  int32    `json:"total,omitempty"`  /*  总数  */
	Count  int32    `json:"count,omitempty"`  /*  查询数量  */
	Rows   []string `json:"rows"`             /*  key 数组  */
}
