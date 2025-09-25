package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwQueryApi
/* 查询l2gw资源。
 */type CtvpcL2gwQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwQueryApi(client *core.CtyunClient) *CtvpcL2gwQueryApi {
	return &CtvpcL2gwQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/l2gw/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwQueryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwQueryRequest) (*CtvpcL2gwQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcL2gwQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwQueryRequest struct {
	RegionID     string  `json:"regionID,omitempty"`     /*  资源池 ID  */
	L2gwID       *string `json:"l2gwID,omitempty"`       /*  l2gw ID  */
	QueryContent *string `json:"queryContent,omitempty"` /*  名称模糊查询  */
	PageNo       int32   `json:"pageNo"`                 /*  列表的页码，默认值为 1。  */
	PageSize     int32   `json:"pageSize"`               /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcL2gwQueryResponse struct {
	StatusCode  int32                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcL2gwQueryReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcL2gwQueryReturnObjResponse struct {
	Results      []*CtvpcL2gwQueryReturnObjResultsResponse `json:"results"`      /*  l2gw组  */
	TotalCount   int32                                     `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                     `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                     `json:"totalPage"`    /*  总页数  */
}

type CtvpcL2gwQueryReturnObjResultsResponse struct {
	L2gwID      *string   `json:"l2gwID,omitempty"`      /*  l2gw 示例 ID  */
	Ip          *string   `json:"ip,omitempty"`          /*  ip  */
	VpcID       *string   `json:"vpcID,omitempty"`       /*  vpcID  */
	Name        *string   `json:"name,omitempty"`        /*  名字  */
	Description *string   `json:"description,omitempty"` /*  描述  */
	ProjectID   *string   `json:"projectID,omitempty"`   /*  企业项目  */
	Spec        *string   `json:"spec,omitempty"`        /*  规格 STANDARD：标准版  ENHANCED：增强版  BASIC: 基础版  */
	LinkCnt     *string   `json:"linkCnt,omitempty"`     /*  二层连接数量  */
	LinkGwID    *string   `json:"linkGwID,omitempty"`    /*  网关ID  */
	LinkGwName  *string   `json:"linkGwName,omitempty"`  /*  网关名字  */
	LinkGwType  *string   `json:"linkGwType,omitempty"`  /*  隧道连接方式 linegw：云专线  vpn：VPN  */
	VpcCidr     *string   `json:"vpcCidr,omitempty"`     /*  vpcCidr  */
	VpcExt      []*string `json:"vpcExt"`                /*  vpc扩展网段  */
	VpcName     *string   `json:"vpcName,omitempty"`     /*  vpc-name  */
	SubnetID    *string   `json:"subnetID,omitempty"`    /*  子网ID  */
	SubnetCidr  *string   `json:"subnetCidr,omitempty"`  /*  子网cidr  */
	SubnetName  *string   `json:"subnetName,omitempty"`  /*  子网name  */
	CreatedAt   *string   `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt   *string   `json:"updatedAt,omitempty"`   /*  更新  */
	ExpiredTime *string   `json:"expiredTime,omitempty"` /*  过期时间  */
	BillMode    *string   `json:"billMode,omitempty"`    /*  账单模式，PACCYCLE：包周期，ONDEMAND：按需  */
	ResourceID  *string   `json:"resourceID,omitempty"`  /*  IT资源ID  */
}
