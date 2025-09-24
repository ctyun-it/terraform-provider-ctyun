package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbGetListApi 查询mongodb实例列表
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=62&api=8708&data=78&isNormal=1&vid=72
type MongodbGetListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbGetListApi(client *ctyunsdk.CtyunClient) *MongodbGetListApi {
	return &MongodbGetListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS2/v1/openApi/getAllInstances",
		},
	}
}

type MongodbGetListHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` //项目id
	RegionID  string  `json:"regionId"`            //资源区regionId，比如实例在资源区A，则需要填写A资源区的regionId
}

func (this *MongodbGetListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbGetListRequest, header *MongodbGetListHeaders) (*MongodbGetListResponse, error) {
	var err error
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.RegionID == "" {
		err = errors.New("regionId is empty")
		return nil, err
	}
	builder.AddHeader("regionId", header.RegionID)
	if req.ProdInstName != nil {
		builder.AddParam("prodInstName", *req.ProdInstName)
	}
	if req.ResDbEngine != nil {
		builder.AddParam("resDbEngine", *req.ResDbEngine)
	}
	if req.LabelIds != nil {
		builder.AddParam("labelIds", *req.LabelIds)
	}
	if err != nil {
		return nil, err
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return nil, err
	}
	var listResponse MongodbGetListResponse
	err = resp.Parse(&listResponse)
	if err != nil {
		return nil, err
	}
	return &listResponse, nil
}

type MongodbGetListRequest struct {
	PageNow      int32   `json:"pageNow"`                // 当前页
	PageSize     int32   `json:"pageSize"`               // 单页记录条数
	ResDbEngine  *string `json:"resDbEngine,omitempty"`  // 数据库引擎
	ProdInstName *string `json:"prodInstName,omitempty"` // 实例名称
	LabelIds     *string `json:"labelIds,omitempty"`     // 标签id
}

type MongodbGetListResponse struct {
	StatusCode int32                            `json:"statusCode"`        // 返回码
	Message    *string                          `json:"message,omitempty"` // 返回消息
	ReturnObj  *MongodbGetListResponseReturnObj `json:"returnObj"`         // 分页信息
	Error      *string                          `json:"error,omitempty"`   // 错误码（失败时才返回）
}

type MongodbGetListResponseReturnObj struct {
	LastPage        int32                                    `json:"lastPage"`
	StartRow        int32                                    `json:"startRow"`        // 当前页的第一行的索引
	HasNextPage     bool                                     `json:"hasNextPage"`     // 是否有下一页数据可用
	PrePage         int32                                    `json:"prePage"`         // 前一页的页码
	NextPage        int32                                    `json:"nextPage"`        // 下一页的页码
	EndRow          int32                                    `json:"endRow"`          // 当前页的最后一行的索引
	PageSize        int32                                    `json:"pageSize"`        // 每页包含的数据条目数量
	PageNum         int32                                    `json:"pageNum"`         // 当前页的页码
	NavigatePages   int32                                    `json:"navigatePages"`   // 导航页数的数量
	Total           int32                                    `json:"total"`           // 总数据条目数量
	Pages           int32                                    `json:"pages"`           // 总页数
	FirstPage       int32                                    `json:"firstPage"`       // 第第一页的页码
	Size            int32                                    `json:"size"`            // 当前页包含的数据条目数量
	IsLastPage      bool                                     `json:"isLastPage"`      // 是否是最后一页
	HasPreviousPage bool                                     `json:"hasPreviousPage"` // 是否有前一页数据可用
	IsFirstPage     bool                                     `json:"isFirstPage"`     // 是否是第一页
	List            []MongodbGetListResponseReturnDetailList `json:"list"`
}

type MongodbGetListResponseReturnDetailList struct {
	ProdOrderStatus             int32   `json:"prodOrderStatus"`             //0->订单正常, 1->订单冻结, 2->订单注销, 3->施工中, 4->施工失败
	SubNetID                    string  `json:"subNetId"`                    //子网ID
	MaintainTime                string  `json:"maintainTime"`                //可维护时间
	Subnet                      string  `json:"subnet"`                      //子网
	LogStatus                   bool    `json:"logStatus"`                   //实例日志审计状态
	OrderId                     int64   `json:"orderId"`                     //订单ID
	NetName                     string  `json:"netName"`                     //专有网络
	VersionNum                  *string `json:"versionNum"`                  // 版本号
	SecurityGroupId             string  `json:"securityGroupId"`             // 安全组ID
	ParameterConfigsvrGroupUsed *string `json:"parameterConfigsvrGroupUsed"` // 参数配置
	DiskSize                    int32   `json:"diskSize"`                    // 存储空间大小
	TplName                     string  `json:"tplName"`                     // 模板名称
	ProdInstSetName             string  `json:"prodInstSetName"`             // 实例对应的SET名
	Released                    bool    `json:"released"`                    // 实例是否已被释放
	SecurityGroup               string  `json:"securityGroup"`               // 安全组
	ProdType                    int32   `json:"prodType"`                    // 0:单机, 2:副本集(三节点), 4:副本集(五节点), 6:副本集(七节点), 10:分片集群
	ExpireTime                  int64   `json:"expireTime"`                  // 到期时间
	ProdInstId                  string  `json:"prodInstId"`                  // 实例id
	ProjectName                 string  `json:"projectName"`                 // 企业项目名称
	ProjectId                   string  `json:"projectId"`                   // 企业项目id
	DestroyedTime               string  `json:"destroyedTime"`               // 实例销毁时间
	ProdInstFlag                string  `json:"prodInstFlag"`                // 规定为“实例ID 实例名称”
	ProdDbEngine                string  `json:"prodDbEngine"`                // dds数据库产品的版本
	BillMode                    int32   `json:"billMode"`                    // 1:包周期计费, 2:按需计费
	ProdId                      string  `json:"prodId"`                      // 产品表示
	RestoreTime                 string  `json:"restoreTime"`                 // 实例恢复时间
	ProdRunningStatus           int32   `json:"prodRunningStatus"`           // 实例运行状态: 实例运行状态: 0->运行正常, 1->重启中, 2-备份操作中,3->恢复操作中,4->转换ssl,5->异常,6->修改参数组中,7->已冻结,8->已注销,9->施工中,10->施工失败,11->扩容中,12->主备切换中
	DiskUsed                    *string `json:"diskUsed"`                    // 磁盘空间
	ParameterGroupUsed          string  `json:"parameterGroupUsed"`          // 参数组名称，标明参数组的版本
	VpcId                       string  `json:"vpcId"`                       // vpc网络ID
	InnodbThreadConcurrency     int64   `json:"innodbThreadConcurrency"`     // 线程数
	DiskType                    string  `json:"diskType"`                    // 存储类型
	ProdBillType                int32   `json:"prodBillType"`                // 0:按月计费, 1:按天计费, 2:按年计费, 3:按流量计费
	MachineSpec                 string  `json:"machineSpec"`                 // CPU内存规格
	ProdInstName                string  `json:"prodInstName"`                // 实例名称
	InnodbBufferPoolSize        string  `json:"innodbBufferPoolSize"`        // 缓存池大小
	UsedSpace                   *string `json:"usedSpace"`                   // 已使用空间
	UserId                      int64   `json:"userId"`                      // 用户id
	ProdBillTime                int32   `json:"prodBillTime"`                // 购买时长
	Destroyed                   bool    `json:"destroyed"`                   // 实例是否已经被销毁
	CreateTime                  int64   `json:"createTime"`                  // 创建时间
	TenantId                    int64   `json:"tenantId"`                    // 租户id
	OuterId                     string  `json:"outerId"`                     // 产品ID
	TplCode                     string  `json:"tplCode"`                     // 模板编码
	ReadPort                    *int32  `json:"readPort"`                    // 读端口
}
