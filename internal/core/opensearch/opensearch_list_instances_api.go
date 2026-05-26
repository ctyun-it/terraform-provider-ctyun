package opensearch

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OpensearchListInstancesApi
/* 查询 OpenSearch 实例列表<br />
<b>准备工作：</b><br />
&emsp;&emsp;构造请求：在调用前需要了解如何构造请求<br />
&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用<br />
*/type OpensearchListInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOpensearchListInstancesApi(client *core.CtyunClient) *OpensearchListInstancesApi {
	return &OpensearchListInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/os/openapi/v1/cluster/selectInstancePage",
			ContentType:  "application/json",
		},
	}
}

func (a *OpensearchListInstancesApi) Do(ctx context.Context, credential core.Credential, req *ListInstancesRequest) (*ListInstancesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ListInstancesRequest
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
	var resp ListInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ListInstancesRequest struct {
	RegionID         string  `json:"regionId"`                   /*  资源池 id  */
	PageIndex        int     `json:"pageIndex"`                  /*  当前页  */
	PageSize         int     `json:"pageSize"`                   /*  每页大小  */
	ClusterName      string  `json:"clusterName,omitempty"`      /*  实例名称（可选）*/
	ClusterType      int32   `json:"clusterType"`                /*  实例类型：1: OpenSearch, 2: Elasticsearch  */
	ProjectID        string  `json:"projectId,omitempty"`        /*  企业项目编码（可选）*/
	ClusterStateList []int32 `json:"clusterStateList,omitempty"` /*  需要查询的实例状态列表（可选）*/
	IsQueryNet       *bool   `json:"isQueryNet,omitempty"`       /*  是否查询实例图形化节点公网地址（可选）*/
}

type ListInstancesResponse struct {
	StatusCode int32           `json:"statusCode"` /*  状态码，成功："200"，失败："500"  */
	Error      string          `json:"error"`      /*  错误码，请求成功时，不返回该字段  */
	Message    string          `json:"message"`    /*  用来简述当前接口调用状态以及必要提示信息  */
	ReturnObj  InstancesReturn `json:"returnObj"`  /*  返回结果  */
}

type InstancesReturn struct {
	Total     int             `json:"total"`     /*  总实例数量  */
	Records   []ClusterRecord `json:"records"`   /*  集群列表  */
	PageIndex int             `json:"pageIndex"` /*  当前页  */
	PageSize  int             `json:"pageSize"`  /*  每页大小  */
}

type ClusterRecord struct {
	ClusterID            string   `json:"clusterId"`            /*  实例 id  */
	ClusterName          string   `json:"clusterName"`          /*  实例名称  */
	ClusterState         string   `json:"clusterState"`         /*  实例状态  */
	ClusterStateType     int32    `json:"clusterStateType"`     /*  实例状态类型  */
	ProcessType          *int32   `json:"processType"`          /*  创建中的进度  */
	ReleaseVersion       string   `json:"releaseVersion"`       /*  发布版本  */
	IndexNum             string   `json:"indexNum"`             /*  索引数量  */
	StorageUsage         string   `json:"storageUsage"`         /*  存储总用量  */
	StoragePercent       float64  `json:"storagePercent"`       /*  存储百分比  */
	ClusterType          int32    `json:"clusterType"`          /*  实例类型  */
	ClusterTypeName      string   `json:"clusterTypeName"`      /*  类型名称  */
	ClusterTypeVersion   string   `json:"clusterTypeVersion"`   /*  实例版本  */
	PayType              string   `json:"payType"`              /*  付费类型  */
	PayTypeValue         int32    `json:"payTypeValue"`         /*  付费类型值  */
	ClusterDueTime       int64    `json:"clusterDueTime"`       /*  实例到期时间  */
	CreateTime           int64    `json:"createTime"`           /*  创建时间  */
	ComponentUrl         *string  `json:"componentUrl"`         /*  组件访问链接  */
	Ipv6ComponentUrl     *string  `json:"ipv6ComponentUrl"`     /*  组件 ipv6 访问链接  */
	CerebroIpv4Url       *string  `json:"cerebroIpv4Url"`       /*  Cerebro IPv4 访问链接  */
	CerebroIpv6Url       *string  `json:"cerebroIpv6Url"`       /*  Cerebro IPv6 访问链接  */
	ClusterMessage       *string  `json:"clusterMessage"`       /*  错误原因  */
	EnableIpv6           string   `json:"enableIpv6"`           /*  是否开启 ipv6  */
	ComponentName        string   `json:"componentName"`        /*  实例组件列表  */
	ProjectID            string   `json:"projectId"`            /*  企业项目编码  */
	ProjectName          string   `json:"projectName"`          /*  企业项目名称  */
	LogstashClusterNames []string `json:"logstashClusterNames"` /*  Logstash 集群名称列表  */
}
