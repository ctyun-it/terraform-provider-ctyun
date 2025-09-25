package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlListApi(client *ctyunsdk.CtyunClient) *PgsqlListApi {
	return &PgsqlListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method: http.MethodGet,
			//Method:  http.MethodPost,
			UrlPath: "/PG/v1/product/page-paas-product",
		},
	}
}

type PgsqlListRequest struct {
	PageNum      int32   `json:"pageNum"`      // 当前页
	PageSize     int32   `json:"pageSize"`     // 页大小，范围1-500
	ProdInstName *string `json:"prodInstName"` // 实例名称，支持模糊匹配
	LabelName    *string `json:"labelName"`    // 标签名称（一级标签）
	LabelValue   *string `json:"labelValue"`   // 标签值（二级标签）
	ProdInstId   *string `json:"prodInstId"`   // 实例id
	InstanceType *string `json:"instanceType"` // 实例类型，primary:主实例，readonly:只读实例，不传则查询全部实例
}

type PgsqlListRequestHeader struct {
	ProjectID *string `json:"project_id"`
	RegionID  string  `json:"region_id"`
}

type PgsqlListResponsePageInfo struct {
	CreateTime          string   `json:"createTime"`          // 创建时间
	ProdDbEngine        string   `json:"prodDbEngine"`        // 数据库实例引擎
	ProdInstId          string   `json:"prodInstId"`          // 实例id
	ProdInstName        string   `json:"prodInstName"`        // 实例名称
	ProdRunningStatus   int32    `json:"prodRunningStatus"`   // 实例当前的运行状态0:运行中  1:重启中  2:备份中 3:恢复中 1001:已停止 1006:恢复失败  1007:VIP不可用 1008:GATEWAY不可用 1009:主库不可用 1010:备库不可用 1021:实例维护中 2000:开通中 2002:已退订 2005:扩容中 2011:冻结
	Alive               int32    `json:"alive"`               // 实例是否存活, 0:存活，-1：异常
	ProdOrderStatus     int32    `json:"prodOrderStatus"`     // 订单状态，0：正常，1：冻结，2：删除，3：操作中，4：失败, 2005:扩容中
	ProdType            int32    `json:"prodType"`            // 实例部署方式，0：单机部署，1：主备部署
	ReadPort            int32    `json:"readPort"`            // 读端口
	Vip                 string   `json:"vip"`                 // 虚拟ip地址
	WritePort           int32    `json:"writePort"`           // 写端口
	ReadonlyInstnaceIds []string `json:"readonlyInstnaceIds"` // 只读实例id列表
	InstanceType        string   `json:"instanceType"`        // 实例类型，primary:主实例，readonly:只读实例
	ToolType            int32    `json:"toolType"`            // 备份工具类型，1：pg_baseback，2：pgbackrest，3：s3
}

type PgsqlListResponseReturnObj struct {
	PageNum   int32                       `json:"pageNum"`   // 当前页
	PageSize  int32                       `json:"pageSize"`  // 页大小
	Total     int64                       `json:"total"`     // 总记录数
	PageTotal int32                       `json:"pageTotal"` // 总页数
	List      []PgsqlListResponsePageInfo `json:"list"`      // 实例列表

}

type PgsqlListResponse struct {
	StatusCode int32                       `json:"statusCode"` // 接口状态码，参考下方状态码
	Error      string                      `json:"error"`      // 错误码
	Message    *string                     `json:"message"`    // 描述信息
	ReturnObj  *PgsqlListResponseReturnObj `json:"returnObj"`  // 返回对象，包含具体的返回数据
}

func (this *PgsqlListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlListRequest, header *PgsqlListRequestHeader) (listResponse *PgsqlListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.RegionID == "" {
		err = errors.New("missing required field: regionId")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	if req.PageNum == 0 {
		err = errors.New("missing required field: PageNum")
		return
	}
	if req.PageSize == 0 {
		err = errors.New("missing required field: PageSize")
		return
	}
	builder.AddParam("pageNum", fmt.Sprintf("%d", req.PageNum))
	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))
	if req.ProdInstName != nil {
		builder.AddParam("prodInstName", *req.ProdInstName)
	}
	if req.LabelName != nil {
		builder.AddParam("labelName", *req.LabelName)
	}
	if req.LabelValue != nil {
		builder.AddParam("labelValue", *req.LabelValue)
	}
	if req.ProdInstId != nil {
		builder.AddParam("prodInstId", *req.ProdInstId)
	}
	if req.InstanceType != nil {
		builder.AddParam("instanceType", *req.InstanceType)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	listResponse = &PgsqlListResponse{}
	err = resp.Parse(listResponse)
	if err != nil {
		return
	}
	return listResponse, nil
}
