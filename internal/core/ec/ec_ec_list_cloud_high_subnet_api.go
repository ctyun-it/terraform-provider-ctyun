package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcListCloudHighSubnetApi
/* 查询云间高速下所有的子网 */
type EcEcListCloudHighSubnetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListCloudHighSubnetApi(client *core.CtyunClient) *EcEcListCloudHighSubnetApi {
	return &EcEcListCloudHighSubnetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/express-connect/list-subnets",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListCloudHighSubnetApi) Do(ctx context.Context, credential core.Credential, req *EcEcListCloudHighSubnetRequest) (*EcEcListCloudHighSubnetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	ctReq.AddParam("IPVersion", req.IPVersion)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListCloudHighSubnetResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListCloudHighSubnetRequest struct {
	EcID      string `json:"ecID"`      /*  云间高速实例ID  */
	IPVersion string `json:"IPVersion"` /*  ip类型<br/>取值如下:<br/>IPV4：IPV4类型<br/>IPV6:IPV6类型  */
}

type EcEcListCloudHighSubnetResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                   `json:"errorCode"`   /*   业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述   */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述   */
	TraceID     *string                                   `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListCloudHighSubnetReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListCloudHighSubnetReturnObjResponse struct {
	CurrentCount *int32    `json:"currentCount"` /*  当前页记录数  */
	TotalCount   *int32    `json:"totalCount"`   /*  查询的总记录数  */
	TotalPage    *int32    `json:"totalPage"`    /*  总页数  */
	Results      []*string `json:"results"`      /*  返回查询结果，字符串数组，值类型为String   */
}
