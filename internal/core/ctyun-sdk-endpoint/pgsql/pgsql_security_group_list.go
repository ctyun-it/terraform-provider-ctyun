package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlSecurityGroupListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlSecurityGroupListApi(client *ctyunsdk.CtyunClient) *PgsqlSecurityGroupListApi {
	return &PgsqlSecurityGroupListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/securityGroup",
		},
	}
}

func (this *PgsqlSecurityGroupListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlSecurityGroupListRequest, header *PgsqlSecurityGroupListRequestHeader) (updateResp *PgsqlSecurityGroupListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	if req.InstID == "" {
		err = errors.New("missing required field: InstID")
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	updateResp = &PgsqlSecurityGroupListResponse{}
	err = resp.Parse(updateResp)
	if err != nil {
		return
	}
	return updateResp, nil
}

type PgsqlSecurityGroupListRequest struct {
	RegionID string `json:"regionId"`
	InstID   string `json:"instId"`
}

type PgsqlSecurityGroupListRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
}

type PgsqlSecurityGroupListReturnObj struct {
	TotalCount int32                                 `json:"totalCount"`
	Data       []PgsqlSecurityGroupListReturnObjData `json:"data"`
	CustomerID string                                `json:"customerId"`
	Cidr       string                                `json:"cidr"`
	TalkOrder  bool                                  `json:"talkOrder"`
}
type PgsqlSecurityGroupListReturnObjData struct {
	SecurityGroupName  string              `json:"securityGroupName"`  // 安全组名称
	ID                 string              `json:"id"`                 // 安全组id
	VmNum              string              `json:"vmNum"`              // 相关云主机
	Origin             string              `json:"origin"`             // 企业项目
	VpcName            string              `json:"vpcName"`            // vpc名称
	CreationTime       string              `json:"creationTime"`       // 创建时间
	RegionID           string              `json:"regionId"`           // 资源池id
	Description        string              `json:"description"`        // 安全组描述信息
	VpcID              string              `json:"vpcId"`              // 虚拟私有云网ID
	SecurityGroupRules []SecurityGroupRule `json:"securityGroupRules"` // 安全组规则信息
}

type SecurityGroupRule struct {
	Direction       string `json:"direction"`       // 出方向-egress、入方向-ingress
	Priority        int32  `json:"priority"`        // 优先级:0~100
	EtherType       string `json:"etherType"`       // IP类型:IPv4、IPv6
	Protocol        string `json:"protocol"`        // 协议: ANY、TCP、UDP、ICMP(v4)
	Range           string `json:"range"`           // 协议: ANY、TCP、UDP、ICMP(v4)
	DestCidrIp      string `json:"destCidrIp"`      // 远端地址:0.0.0.0/0
	Description     string `json:"description"`     // 安全组规则描述信息
	ID              string `json:"id"`              // 唯一标识ID
	SecurityGroupID string `json:"securityGroupId"` // 安全组ID
	Origin          string `json:"origin"`          // 企业项目
	CreateTime      string `json:"createTime"`      // 创建时间，UTC时间
}

type PgsqlSecurityGroupListResponse struct {
	StatusCode int32                            `json:"statusCode"`
	Message    *string                          `json:"message"`
	Error      *string                          `json:"error"`
	ReturnObj  *PgsqlSecurityGroupListReturnObj `json:"returnObj"`
}
