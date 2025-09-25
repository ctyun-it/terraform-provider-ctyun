package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupCheckApi
/* 检查弹性伸缩组是否可以修改
 */type ScalingGroupCheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupCheckApi(client *core.CtyunClient) *ScalingGroupCheckApi {
	return &ScalingGroupCheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/check",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupCheckApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupCheckRequest) (*ScalingGroupCheckResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupCheckRequest
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
	var resp ScalingGroupCheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupCheckRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	GroupID  int64  `json:"groupID,omitempty"`  /*  伸缩组ID  */
}

type ScalingGroupCheckResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	Message     string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupCheckReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
}

type ScalingGroupCheckReturnObjResponse struct {
	Result int32 `json:"result"` /*  伸缩组是否可以修改<br>取值范围：<br>1：可以修改<br>2：不可以修改  */
}
