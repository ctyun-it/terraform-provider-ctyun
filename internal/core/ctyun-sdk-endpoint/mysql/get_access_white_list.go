package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetAccessWhiteList struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetAccessWhiteList(client *ctyunsdk.CtyunClient) *TeledbGetAccessWhiteList {
	return &TeledbGetAccessWhiteList{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/rds-manage/rds-access-white-list",
		},
	}
}

type TeledbGetAccessWhiteListRequest struct {
	OuterProdInstID string `json:"outerProdInstId"`
}

type TeledbGetAccessWhiteListRequestHeader struct {
	ProjectID *string `json:"project-id"`
	InstID    *string `json:"inst-id"`
	RegionID  string  `json:"regionId"`
}

type TeledbGetAccessWhiteListResponseReturnObj struct {
	GroupName           string   `json:"groupName"`
	GroupWhiteListCount int32    `json:"groupWhiteListCount"`
	OuterProdInstID     string   `json:"outerProdInstId"`
	CreateTime          int64    `json:"createTime"`
	UpdateTime          int64    `json:"updateTime"`
	WhiteList           []string `json:"whiteList"`
	AccessMachineType   string   `json:"accessMachineType"`
	ID                  int64    `json:"id"`
}

type TeledbGetAccessWhiteListResponse struct {
	StatusCode int32               `json:"statusCode"` // 接口状态码
	Error      string              `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string              `json:"message"`    // 描述信息
	ReturnObj  CompatibleReturnObj `json:"returnObj"`
}

// 自定义兼容类型：支持 [] 或 {}
type CompatibleReturnObj []*TeledbGetAccessWhiteListResponseReturnObj

// 重写 UnmarshalJSON 实现自动兼容
func (c *CompatibleReturnObj) UnmarshalJSON(data []byte) error {
	// 1. 先尝试解析成数组
	var arr []*TeledbGetAccessWhiteListResponseReturnObj
	if err := json.Unmarshal(data, &arr); err == nil {
		*c = arr
		return nil
	}

	// 2. 如果失败，尝试解析成空对象（忽略内容）
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		*c = CompatibleReturnObj{} // 赋值为空切片
		return nil
	}

	// 3. 都不是才返回错误
	return fmt.Errorf("不支持的returnObj格式: %s", string(data))
}

func (this *TeledbGetAccessWhiteList) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetAccessWhiteListRequest, header *TeledbGetAccessWhiteListRequestHeader) (bindResponse *TeledbGetAccessWhiteListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.InstID != nil {
		builder.AddHeader("inst-id", *header.InstID)
	}
	if header.RegionID == "" {
		err = errors.New("查询Mysql白名单，region_id必填")
		return
	}
	if req.OuterProdInstID == "" {
		err = errors.New("查询Mysql白名单，prod_inst_id必填")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	builder.AddParam("outerProdInstId", req.OuterProdInstID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	bindResponse = &TeledbGetAccessWhiteListResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}
