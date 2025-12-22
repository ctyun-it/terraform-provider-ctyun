package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetEpPloyApi
/* 查询企业项目用户组策略 */
type CtiamGetEpPloyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetEpPloyApi(client *core.CtyunClient) *CtiamGetEpPloyApi {
	return &CtiamGetEpPloyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/getEpPloy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetEpPloyApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetEpPloyRequest) (*CtiamGetEpPloyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetEpPloyRequest
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
	var resp CtiamGetEpPloyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetEpPloyRequest struct {
	ProjectId string `json:"projectId"` /*  企业项目id  */
	GroupId   string `json:"groupId"`   /*  用户组id  */
}

type CtiamGetEpPloyResponse struct {
	StatusCode *string                          `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                          `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamGetEpPloyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                          `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetEpPloyReturnObjResponse struct {
	List []*CtiamGetEpPloyReturnObjListResponse `json:"list"` /*  策略列表  */
}

type CtiamGetEpPloyReturnObjListResponse struct {
	Id          *string                                         `json:"id"`          /*  id  */
	PloyName    *string                                         `json:"ployName"`    /*  策略名称  */
	PloyType    *int32                                          `json:"ployType"`    /*  策略类型（1：系统策略，2：自定义策略）  */
	PloyRange   *int32                                          `json:"ployRange"`   /*  策略范围（1：项目级，2：全局）  */
	PloyContent *CtiamGetEpPloyReturnObjListPloyContentResponse `json:"ployContent"` /*  策略内容  */
	ProductName *string                                         `json:"productName"` /*  产品名称  */
}

type CtiamGetEpPloyReturnObjListPloyContentResponse struct {
	Version   *string                                                    `json:"Version"`   /*  版本  */
	Statement []*CtiamGetEpPloyReturnObjListPloyContentStatementResponse `json:"Statement"` /*  策略  */
}

type CtiamGetEpPloyReturnObjListPloyContentStatementResponse struct {
	Action []*string `json:"Action"` /*  策略三元组集合  */
	Effect *string   `json:"Effect"` /*  策略操作  */
}
