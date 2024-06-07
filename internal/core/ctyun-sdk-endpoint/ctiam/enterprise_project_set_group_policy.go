package ctiam

import (
	"context"
	"net/http"
	"strings"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EnterpriseProjectSetGroupPolicyApi 设置企业项目所属用户组及策略
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9431&data=114
type EnterpriseProjectSetGroupPolicyApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectSetGroupPolicyApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectSetGroupPolicyApi {
	return &EnterpriseProjectSetGroupPolicyApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/setEpGroupPloy",
		},
	}
}

func (this *EnterpriseProjectSetGroupPolicyApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectSetGroupPolicyRequest) (*EnterpriseProjectSetGroupPolicyResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&enterpriseProjectSetGroupPolicyRealRequest{
		GroupId:   req.GroupId,
		ProjectId: req.ProjectId,
		PloyIds:   strings.Join(req.PloyIds, ","),
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &EnterpriseProjectSetGroupPolicyResponse{}, nil
}

type EnterpriseProjectSetGroupPolicyResponse struct {
}

type enterpriseProjectSetGroupPolicyRealRequest struct {
	GroupId   string `json:"groupId"`
	ProjectId string `json:"projectId"`
	PloyIds   string `json:"ployIds"`
}

type EnterpriseProjectSetGroupPolicyRequest struct {
	GroupId   string
	ProjectId string
	PloyIds   []string
}
