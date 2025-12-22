package cda

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CdaCdaListAccountAuthApi
/* 查询账户下已添加的跨账号授权网络实例。 */
type CdaCdaListAccountAuthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCdaCdaListAccountAuthApi(client *core.CtyunClient) *CdaCdaListAccountAuthApi {
	return &CdaCdaListAccountAuthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/cda/accountauth/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CdaCdaListAccountAuthApi) Do(ctx context.Context, credential core.Credential, req *CdaCdaListAccountAuthRequest) (*CdaCdaListAccountAuthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CdaCdaListAccountAuthRequest
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
	var resp CdaCdaListAccountAuthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CdaCdaListAccountAuthRequest struct {
	RegionID      string  `json:"regionID"`                /*  资源池ID  */
	PageNo        *int32  `json:"pageNo,omitempty"`        /*  页数  */
	PageSize      *int32  `json:"pageSize,omitempty"`      /*  每页数据量  */
	VpcID         *string `json:"vpcID,omitempty"`         /*  查询被授权：指定VPC ID,会带跨账号VPC 子网信息  */
	AuthAccountID *string `json:"authAccountID,omitempty"` /*  查询已授权： 账号ID不传<br>查询被授权：账号ID为自己账号ID  */
}

type CdaCdaListAccountAuthResponse struct {
	StatusCode         *int32                                    `json:"statusCode"`         /*  返回状态码(800为成功，900为失败)  */
	Message            *string                                   `json:"message"`            /*  失败时的错误描述，一般为英文描述  */
	Description        *string                                   `json:"description"`        /*  失败时的错误描述，一般为中文描述  */
	ErrorCode          *string                                   `json:"errorCode"`          /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail        *CdaCdaListAccountAuthErrorDetailResponse `json:"errorDetail"`        /*  错误明细  */
	TotalCount         *int32                                    `json:"totalCount"`         /*  总共数量  */
	CurrentCount       *int32                                    `json:"currentCount"`       /*  当前数量  */
	Fuid               *string                                   `json:"fuid"`               /*  跨账号授权实例ID  */
	AccountId          *string                                   `json:"accountId"`          /*  当前账号ID  */
	Account            *string                                   `json:"account"`            /*  当前账号邮箱  */
	VpcId              *string                                   `json:"vpcId"`              /*  授权的VPC ID  */
	VpcName            *string                                   `json:"vpcName"`            /*  授权的VPC Name  */
	VrfName            *string                                   `json:"vrfName"`            /*  授权的VPC给专线网关  */
	AuthAccountId      *string                                   `json:"authAccountId"`      /*  对方账号ID，授权自己VPC给对方账号  */
	AuthAccount        *string                                   `json:"authAccount"`        /*  对方账号邮箱  */
	RegionId           *string                                   `json:"regionId"`           /*  资源池ID  */
	IsSwConfig         *bool                                     `json:"isSwConfig"`         /*  是否下发配置到交换机  */
	LineList           []*string                                 `json:"lineList"`           /*  绑定的物理专线列表  */
	ResourcePool       *string                                   `json:"resourcePool"`       /*  资源池ID  */
	ResourcePoolName   *string                                   `json:"resourcePoolName"`   /*  资源池名字  */
	IsAutomation       *bool                                     `json:"isAutomation"`       /*  是否自动化  */
	Project_id_ecs     *string                                   `json:"project_id_ecs"`     /*  租户ID  */
	LgcreateTime       *string                                   `json:"lgcreateTime"`       /*  创建时间  */
	Fuser_last_updated *string                                   `json:"fuser_last_updated"` /*  上次更新时间  */
	Delete_time        *string                                   `json:"delete_time"`        /*  删除时间  */
}

type CdaCdaListAccountAuthErrorDetailResponse struct{}
