package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EnterpriseProjectGetPolicyApi 查询企业项目用户组策略
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9429&data=114
type EnterpriseProjectGetPolicyApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectGetPolicyApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectGetPolicyApi {
	return &EnterpriseProjectGetPolicyApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/getEpPloy",
		},
	}
}

func (this *EnterpriseProjectGetPolicyApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectGetPolicyRequest) (*EnterpriseProjectGetPolicyResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&enterpriseProjectGetPolicyRealRequest{
		GroupId:   req.GroupId,
		ProjectId: req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp enterpriseProjectGetPolicyRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}

	var list []EnterpriseProjectGetPolicyListResponse
	for _, l := range resp.List {
		var statement []EnterpriseProjectGetPolicyListPloyContentStatementResponse
		for _, s := range l.PloyContent.Statement {
			statement = append(statement, EnterpriseProjectGetPolicyListPloyContentStatementResponse{
				Action: s.Action,
				Effect: s.Effect,
			})
		}

		list = append(list, EnterpriseProjectGetPolicyListResponse{
			Id:          l.Id,
			PloyName:    l.PloyName,
			PloyType:    l.PloyType,
			PloyRange:   l.PloyRange,
			Status:      l.Status,
			ProductName: l.ProductName,
			CreateTime:  l.CreateTime,
			PloyContent: EnterpriseProjectGetPolicyListPloyContentResponse{
				Version:   l.PloyContent.Version,
				Statement: statement,
			},
		})
	}

	return &EnterpriseProjectGetPolicyResponse{
		List: list,
	}, nil
}

type enterpriseProjectGetPolicyRealResponse struct {
	List []struct {
		Id          string `json:"id"`
		PloyName    string `json:"ployName"`
		PloyType    int    `json:"ployType"`
		PloyRange   int    `json:"ployRange"`
		Status      int    `json:"status"`
		PloyContent struct {
			Version   string `json:"Version"`
			Statement []struct {
				Action []string `json:"Action"`
				Effect string   `json:"Effect"`
			} `json:"Statement"`
		} `json:"ployContent"`
		ProductName string `json:"productName"`
		CreateTime  int64  `json:"createTime"`
	} `json:"list"`
}

type EnterpriseProjectGetPolicyListPloyContentStatementResponse struct {
	Action []string
	Effect string
}

type EnterpriseProjectGetPolicyListPloyContentResponse struct {
	Version   string
	Statement []EnterpriseProjectGetPolicyListPloyContentStatementResponse
}

type EnterpriseProjectGetPolicyListResponse struct {
	Id          string
	PloyName    string
	PloyType    int
	PloyRange   int
	Status      int
	ProductName string
	CreateTime  int64
	PloyContent EnterpriseProjectGetPolicyListPloyContentResponse
}

type EnterpriseProjectGetPolicyResponse struct {
	List []EnterpriseProjectGetPolicyListResponse
}

type enterpriseProjectGetPolicyRealRequest struct {
	GroupId   string `json:"groupId"`
	ProjectId string `json:"projectId"`
}

type EnterpriseProjectGetPolicyRequest struct {
	GroupId   string
	ProjectId string
}
