package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetRdsParameterTemplateDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetRdsParameterTemplateDetailApi(client *ctyunsdk.CtyunClient) *TeledbGetRdsParameterTemplateDetailApi {
	return &TeledbGetRdsParameterTemplateDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/parameter/describe-parameters",
		},
	}
}

func (this *TeledbGetRdsParameterTemplateDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetRdsParameterTemplateDetailRequest, header *TeledbGetRdsParameterTemplateDetailRequestHeader) (GetParameterTemplateDetailResp *TeledbGetRdsParameterTemplateDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)
	if req.PageNow == 0 {
		err = errors.New("page_no 为空")
		return
	}
	if req.PageSize == 0 {
		err = errors.New("page_size 为空")
		return
	}

	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))
	builder.AddParam("pageNow", fmt.Sprintf("%d", req.PageNow))

	if req.OuterProdInstId == "" {
		err = errors.New("id不能为空")
		return
	}
	builder.AddParam("outerProdInstId", fmt.Sprintf("%s", req.OuterProdInstId))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetParameterTemplateDetailResp = &TeledbGetRdsParameterTemplateDetailResponse{}
	err = resp.Parse(GetParameterTemplateDetailResp)
	if err != nil {
		return
	}
	return GetParameterTemplateDetailResp, nil
}

type TeledbGetRdsParameterTemplateDetailRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` //
	PageNow         int32  `json:"page_now"`
	PageSize        int32  `json:"page_size"`
}

type TeledbGetRdsParameterTemplateDetailRequestHeader struct {
	ProjectID *string `json:"Project-Id"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
	InstID    string  `json:"inst-id"`
}

type TeledbGetRdsParameterTemplateDetailResponse struct {
	StatusCode int32 `json:"statusCode"` // 接口状态码
	//Error      *string                                               `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message   string                                                `json:"message"` // 描述信息
	ReturnObj *TeledbGetRdsParameterTemplateDetailResponseReturnObj `json:"returnObj"`
}
type TeledbGetRdsParameterTemplateDetailResponseReturnObj struct {
	NavigatePageNums  []int                                                        `json:"navigatepageNums"`
	StartRow          int32                                                        `json:"startRow"`
	HasNextPage       bool                                                         `json:"hasNextPage"`
	PrePage           int32                                                        `json:"prePage"`
	NextPage          int32                                                        `json:"nextPage"`
	EndRow            int32                                                        `json:"endRow"`
	PageSize          int32                                                        `json:"pageSize"`
	List              []TeledbGetRdsParameterTemplateDetailResponseReturnObjDetail `json:"list"`
	PageNum           int32                                                        `json:"pageNum"`
	NavigatePages     int32                                                        `json:"navigatePages"`
	NavigateFirstPage int32                                                        `json:"navigateFirstPage"`
	Total             int32                                                        `json:"total"`
	Pages             int32                                                        `json:"pages"`
	Size              int32                                                        `json:"size"`
	IsLastPage        bool                                                         `json:"isLastPage"`
	HasPreviousPage   bool                                                         `json:"hasPreviousPage"`
	NavigateLastPage  int32                                                        `json:"navigateLastPage"`
	IsFirstPage       bool                                                         `json:"isFirstPage"`
}

type TeledbGetRdsParameterTemplateDetailResponseReturnObjDetail struct {
	ParameterGroupName string `json:"parameterGroupName"`
	ValueType          string `json:"valuetype"`
	DescriptionEn      string `json:"descriptionEn"`
	Restart            string `json:"restart"`
	Description        string `json:"description"`
	ProdInstId         int64  `json:"prodInstId"`
	ID                 int64  `json:"id"`
	ParameterName      string `json:"parameterName"`
	ParameterValue     string `json:"parameterValue"`
	PermitValue        string `json:"permitValue"`
	UserID             int64  `json:"userId"`
}
