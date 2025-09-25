package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// UserAssociationGroupApi 将用户移入用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=8140&data=114
type UserAssociationGroupApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewUserAssociationGroupApi(client *ctyunsdk.CtyunClient) *UserAssociationGroupApi {
	return &UserAssociationGroupApi{
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/userGroup/userToGroup",
		},
		client: client,
	}
}

func (this *UserAssociationGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserAssociationGroupRequest) (*UserAssociationGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder, err := builder.WriteJson(&userAssociationGroupRealRequest{
		GroupId: req.GroupId,
		UserIds: []userId{{Id: req.UserId}},
	})
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &userAssociationGroupRealResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &UserAssociationGroupResponse{}, nil
}

type userId struct {
	Id string `json:"id"`
}

type userAssociationGroupRealRequest struct {
	GroupId string   `json:"groupId"`
	UserIds []userId `json:"userIds"`
}

type userAssociationGroupRealResponse struct {
}

type UserAssociationGroupRequest struct {
	GroupId string
	UserId  string
}

type UserAssociationGroupResponse struct {
}
