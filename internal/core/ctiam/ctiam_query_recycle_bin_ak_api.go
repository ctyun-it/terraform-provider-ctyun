package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryRecycleBinAkApi
/* 查询回收站ak */
type CtiamQueryRecycleBinAkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryRecycleBinAkApi(client *core.CtyunClient) *CtiamQueryRecycleBinAkApi {
	return &CtiamQueryRecycleBinAkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/credential/queryRecycleBinAk",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryRecycleBinAkApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryRecycleBinAkRequest) (*CtiamQueryRecycleBinAkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryRecycleBinAkRequest
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
	var resp CtiamQueryRecycleBinAkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryRecycleBinAkRequest struct {
	UserIdList []string `json:"userIdList"` /*  用户列表  */
}

type CtiamQueryRecycleBinAkResponse struct {
	StatusCode *string                                  `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                  `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamQueryRecycleBinAkReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                                  `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryRecycleBinAkReturnObjResponse struct {
	RecycleBinAkList []*CtiamQueryRecycleBinAkReturnObjRecycleBinAkListResponse `json:"recycleBinAkList"` /*  回收站ak  */
	UserId           *string                                                    `json:"userId"`           /*  用户id  */
}

type CtiamQueryRecycleBinAkReturnObjRecycleBinAkListResponse struct {
	Ak            *string `json:"ak"`            /*  ak  */
	AkCreatedTime *string `json:"akCreatedTime"` /*  ak创建时间  */
	RecoveryTime  *string `json:"recoveryTime"`  /*  回收时间  */
	ClearTime     *string `json:"clearTime"`     /*  清理时间  */
}
