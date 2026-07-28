package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQuerySensitiveEventApi
type CtiamQuerySensitiveEventApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQuerySensitiveEventApi(client *core.CtyunClient) *CtiamQuerySensitiveEventApi {
	return &CtiamQuerySensitiveEventApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/sensitive/querySensitiveEvent",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQuerySensitiveEventApi) Do(ctx context.Context, credential core.Credential, req *CtiamQuerySensitiveEventRequest) (*CtiamQuerySensitiveEventResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQuerySensitiveEventRequest
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
	var resp CtiamQuerySensitiveEventResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQuerySensitiveEventRequest struct {
	PageNum   int32  `json:"pageNum"`   /*  页码  */
	PageSize  int32  `json:"pageSize"`  /*  每页条数  */
	StartTime string `json:"startTime"` /*  开始时间（时区 GMT+8）  */
	EndTime   string `json:"endTime"`   /*  结束时间（时区 GMT+8）  */
}

type CtiamQuerySensitiveEventResponse struct {
	StatusCode *string                                    `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                    `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                    `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQuerySensitiveEventReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamQuerySensitiveEventReturnObjResponse struct {
	List []*CtiamQuerySensitiveEventReturnObjListResponse `json:"list"` /*  敏感事件  */
}

type CtiamQuerySensitiveEventReturnObjListResponse struct {
	TraceId             *string `json:"traceId"`             /*  traceId  */
	TraceName           *string `json:"traceName"`           /*  权限码  */
	Code                *int32  `json:"code"`                /*  请求码  */
	AccountId           *string `json:"accountId"`           /*  账号id  */
	OperateUser         *string `json:"operateUser"`         /*  操作用户id  */
	Level               *int32  `json:"level"`               /*  等级  */
	SourceIp            *string `json:"sourceIp"`            /*  源ip  */
	RequestMessage      *string `json:"requestMessage"`      /*  请求数据  */
	ResponseMessage     *string `json:"responseMessage"`     /*  响应体  */
	CreateTime          *string `json:"createTime"`          /*  创建时间  */
	Email               *string `json:"email"`               /*  邮箱  */
	Identity            *string `json:"identity"`            /*  身份  */
	ReferencedResources *string `json:"referencedResources"` /*  关联资源  */
	PageNum             *int32  `json:"pageNum"`             /*  页码  */
	PageSize            *int32  `json:"pageSize"`            /*  每页条数  */
	Total               *int32  `json:"total"`               /*  总条数  */
	Pages               *int32  `json:"pages"`               /*  总页数  */
}
