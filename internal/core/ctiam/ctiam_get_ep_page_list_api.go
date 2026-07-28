package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetEpPageListApi
/* 查询企业项目列表 */
type CtiamGetEpPageListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetEpPageListApi(client *core.CtyunClient) *CtiamGetEpPageListApi {
	return &CtiamGetEpPageListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/getEpPageList",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetEpPageListApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetEpPageListRequest) (*CtiamGetEpPageListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetEpPageListRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetEpPageListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetEpPageListRequest struct {
	AccountId   string `json:"accountId"`   /*  账号id  */
	CurrentPage int32  `json:"currentPage"` /*  当前页  */
	PageSize    int32  `json:"pageSize"`    /*  每页显示条数  */
}

type CtiamGetEpPageListResponse struct {
	StatusCode *string                              `json:"statusCode"` /*  状态码  */
	ReturnObj  *CtiamGetEpPageListReturnObjResponse `json:"returnObj"`  /*  回传参数  */
	Message    *string                              `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string                              `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetEpPageListReturnObjResponse struct {
	CurrentPage    *int32                                           `json:"currentPage"`    /*  当前页  */
	PageSize       *int32                                           `json:"pageSize"`       /*  每页显示条数  */
	RecordCount    *int32                                           `json:"recordCount"`    /*  总记录数  */
	PageCount      *int32                                           `json:"pageCount"`      /*  总页数  */
	BeginPageIndex *int32                                           `json:"beginPageIndex"` /*  页码列表的开始索引（包含）  */
	EndPageIndex   *int32                                           `json:"endPageIndex"`   /*  页码列表的结束索引（包含）  */
	RecordList     []*CtiamGetEpPageListReturnObjRecordListResponse `json:"recordList"`     /*  企业项目信息列表  */
}

type CtiamGetEpPageListReturnObjRecordListResponse struct {
	Id          *string `json:"id"`          /*  企业项目id  */
	ProjectName *string `json:"projectName"` /*  企业名称  */
	Status      *int32  `json:"status"`      /*  企业项目状态  */
	HwProjectId *string `json:"hwProjectId"` /*  华为项目id  */
	Description *string `json:"description"` /*  企业项目描述  */
	CreateTime  *string `json:"createTime"`  /*  创建时间  */
	UpdateTime  *string `json:"updateTime"`  /*  修改时间  */
	AccountId   *string `json:"accountId"`   /*  账号id  */
}
