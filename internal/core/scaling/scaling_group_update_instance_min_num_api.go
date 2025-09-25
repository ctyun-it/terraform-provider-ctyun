package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupUpdateInstanceMinNumApi
/* 修改伸缩组最小云主机数
 */type ScalingGroupUpdateInstanceMinNumApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupUpdateInstanceMinNumApi(client *core.CtyunClient) *ScalingGroupUpdateInstanceMinNumApi {
	return &ScalingGroupUpdateInstanceMinNumApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/update-instance-min-num",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupUpdateInstanceMinNumApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupUpdateInstanceMinNumRequest) (*ScalingGroupUpdateInstanceMinNumResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupUpdateInstanceMinNumRequest
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
	var resp ScalingGroupUpdateInstanceMinNumResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupUpdateInstanceMinNumRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	GroupID  int32  `json:"groupID,omitempty"`  /*  伸缩组ID  */
	MinCount int32  `json:"minCount,omitempty"` /*  最小云主机数，取值范围：[0,50]  */
}

type ScalingGroupUpdateInstanceMinNumResponse struct {
	StatusCode  int32                                              `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	Message     string                                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupUpdateInstanceMinNumReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
}

type ScalingGroupUpdateInstanceMinNumReturnObjResponse struct {
	GroupID int32 `json:"groupID"` /*  伸缩组ID  */
}
