package mysql

import (
	"context"
	"errors"
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
	StatusCode int32                                       `json:"statusCode"` // 接口状态码
	Error      string                                      `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                                      `json:"message"`    // 描述信息
	ReturnObj  []TeledbGetAccessWhiteListResponseReturnObj `json:"returnObj"`
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
