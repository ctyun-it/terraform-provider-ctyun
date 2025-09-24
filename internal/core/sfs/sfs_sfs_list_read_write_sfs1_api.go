package sfs

import (
	"context"
	"encoding/json"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
	"strings"
)

// SfsSfsListReadWriteSfs1Api
/* 查询指定文件系统或全部文件系统的读写权限信息
 */type SfsSfsListReadWriteSfs1Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListReadWriteSfs1Api(client *core.CtyunClient) *SfsSfsListReadWriteSfs1Api {
	return &SfsSfsListReadWriteSfs1Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-rw-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListReadWriteSfs1Api) Do(ctx context.Context, credential core.Credential, req *SfsSfsListReadWriteSfs1Request) (*SfsSfsListReadWriteSfs1Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.SfsUID != "" {
		ctReq.AddParam("sfsUID", req.SfsUID)
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsListReadWriteSfs1Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 实现自定义反序列化逻辑
func (b *FlexibleBool) UnmarshalJSON(data []byte) error {
	// 尝试解析为布尔值
	if err := json.Unmarshal(data, &b.Value); err == nil {
		return nil
	}

	// 尝试解析为字符串
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// 将字符串转换为布尔值
	value := strings.ToLower(s) == "true"
	b.Value = &value
	return nil
}

type SfsSfsListReadWriteSfs1Request struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的元素个数，默认10  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  列表的分页页码，默认1  */
}

type SfsSfsListReadWriteSfs1Response struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                    `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                    `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListReadWriteSfs1ReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsListReadWriteSfs1ReturnObjResponse struct {
	List         []*SfsSfsListReadWriteSfs1ReturnObjListResponse `json:"list"`         /*  查询的详情列表  */
	TotalCount   int32                                           `json:"totalCount"`   /*  文件系统读写权限信息总数  */
	CurrentCount int32                                           `json:"currentCount"` /*  当前页码的元素个数  */
	PageSize     int32                                           `json:"pageSize"`     /*  每页个数  */
	PageNo       int32                                           `json:"pageNo"`       /*  当前页数  */
}

type SfsSfsListReadWriteSfs1ReturnObjListResponse struct {
	SfsUID   string       `json:"sfsUID"`   /*  弹性文件功能系统唯一 ID  */
	ReadOnly FlexibleBool `json:"readOnly"` /*  是否是只读  */
}

// 新增自定义类型处理 readOnly 字段
type FlexibleBool struct {
	Value *bool
}
