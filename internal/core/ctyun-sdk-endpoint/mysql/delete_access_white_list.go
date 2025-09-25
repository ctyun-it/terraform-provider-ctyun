package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbDeleteAccessWhiteList struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbDeleteAccessWhiteList(client *ctyunsdk.CtyunClient) *TeledbDeleteAccessWhiteList {
	return &TeledbDeleteAccessWhiteList{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodDelete,
			UrlPath: "/RDS2/v1/open-api/rds-manage/rds-access-white-list",
		},
	}
}

type TeledbDeleteAccessWhiteListRequest struct {
	OuterProdInstID string `json:"outerProdInstId"`
	GroupName       string `json:"groupName"`
}

type TeledbDeleteAccessWhiteListRequestHeader struct {
	ProjectID *string `json:"project-id"`
	InstID    *string `json:"inst-id"`
	RegionID  string  `json:"regionId"`
}

type TeledbDeleteAccessWhiteListResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}

func (this *TeledbDeleteAccessWhiteList) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbDeleteAccessWhiteListRequest, header *TeledbDeleteAccessWhiteListRequestHeader) (bindResponse *TeledbDeleteAccessWhiteListResponse, err error) {
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
		err = errors.New("删除Mysql白名单，region_id必填")
		return
	}
	if req.OuterProdInstID == "" {
		err = errors.New("删除Mysql白名单，prod_inst_id必填")
		return
	}
	if req.GroupName == "" {
		err = errors.New("删除Mysql白名单，group_name必填")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	builder.AddParam("outerProdInstId", req.OuterProdInstID)
	builder.AddParam("groupName", req.GroupName)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	bindResponse = &TeledbDeleteAccessWhiteListResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}
