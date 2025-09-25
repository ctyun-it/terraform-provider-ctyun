package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2PageListLabelApi
/* 查询租户所有标签列表
 */type Dcs2PageListLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2PageListLabelApi(client *core.CtyunClient) *Dcs2PageListLabelApi {
	return &Dcs2PageListLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/label/pageList",
			ContentType:  "",
		},
	}
}

func (a *Dcs2PageListLabelApi) Do(ctx context.Context, credential core.Credential, req *Dcs2PageListLabelRequest) (*Dcs2PageListLabelResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(req.PageIndex), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.LabelKey != "" {
		ctReq.AddParam("labelKey", req.LabelKey)
	}
	if req.LabelVal != "" {
		ctReq.AddParam("labelVal", req.LabelVal)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2PageListLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2PageListLabelRequest struct {
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	PageIndex int32  /*  当前页码（默认值：1）  */
	PageSize  int32  /*  每页行数（默认值：10，范围：1-50）  */
	LabelKey  string /*  标签键（精确匹配查询）  */
	LabelVal  string /*  标签值（精确匹配查询）  */
}

type Dcs2PageListLabelResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2PageListLabelReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2PageListLabelReturnObjResponse struct {
	EndRow   int32                                     `json:"endRow,omitempty"`   /*  标签键结束行号  */
	List     []*Dcs2PageListLabelReturnObjListResponse `json:"list"`               /*  标签键列表  */
	PageNum  int32                                     `json:"pageNum,omitempty"`  /*  请求页码  */
	PageSize int32                                     `json:"pageSize,omitempty"` /*  每页数量  */
	Pages    int32                                     `json:"pages,omitempty"`    /*  总页数  */
	Size     int32                                     `json:"size,omitempty"`     /*  标签键数量  */
	StartRow int32                                     `json:"startRow,omitempty"` /*  标签键开始行号  */
	Total    int32                                     `json:"total,omitempty"`    /*  标签键值对总数  */
}

type Dcs2PageListLabelReturnObjListResponse struct {
	Key  string                                        `json:"key,omitempty"` /*  标签键  */
	Data []*Dcs2PageListLabelReturnObjListDataResponse `json:"data"`          /*  标签键值对象  */
}

type Dcs2PageListLabelReturnObjListDataResponse struct {
	AccountId   string `json:"accountId,omitempty"`   /*  账户ID  */
	Key         string `json:"key,omitempty"`         /*  标签键  */
	LabelId     string `json:"labelId,omitempty"`     /*  标签ID  */
	OperateType string `json:"operateType,omitempty"` /*  操作类型  */
	UserId      string `json:"userId,omitempty"`      /*  用户ID  */
	Value       string `json:"value,omitempty"`       /*  标签值内容  */
}
