package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcListAuthVPCBindCloudHighApi
/* 查看VPC绑定云间高速列表 */
type EcEcListAuthVPCBindCloudHighApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListAuthVPCBindCloudHighApi(client *core.CtyunClient) *EcEcListAuthVPCBindCloudHighApi {
	return &EcEcListAuthVPCBindCloudHighApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/vpc/list-express-connect",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListAuthVPCBindCloudHighApi) Do(ctx context.Context, credential core.Credential, req *EcEcListAuthVPCBindCloudHighRequest) (*EcEcListAuthVPCBindCloudHighResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcListAuthVPCBindCloudHighRequest
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
	var resp EcEcListAuthVPCBindCloudHighResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListAuthVPCBindCloudHighRequest struct {
	VpcID string `json:"vpcID"` /*  vpc ID  */
}

type EcEcListAuthVPCBindCloudHighResponse struct {
	StatusCode  *int32                                         `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                        `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                        `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                        `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcListAuthVPCBindCloudHighReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListAuthVPCBindCloudHighReturnObjResponse struct {
	CurrentCount *int32                                                  `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                                  `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                                  `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListAuthVPCBindCloudHighReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcListAuthVPCBindCloudHighReturnObjResultsResponse struct {
	EcID   *string `json:"ecID"`   /*  云间高速ID  */
	EcName *string `json:"ecName"` /*  云间高速名称  */
}
