package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQuerySnapshotListV41Api
/* 查询云主机快照列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsQuerySnapshotListV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQuerySnapshotListV41Api(client *core.CtyunClient) *CtecsQuerySnapshotListV41Api {
	return &CtecsQuerySnapshotListV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQuerySnapshotListV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQuerySnapshotListV41Request) (*CtecsQuerySnapshotListV41Response, error) {
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
	var resp CtecsQuerySnapshotListV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQuerySnapshotListV41Request struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ProjectID      string `json:"projectID,omitempty"`      /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	PageNo         int32  `json:"pageNo,omitempty"`         /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	InstanceID     string `json:"instanceID,omitempty"`     /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	SnapshotStatus string `json:"snapshotStatus,omitempty"` /*  云主机快照状态，取值范围：<br />pending：创建中，<br /> available：可用，<br /> restoring：恢复中，<br /> error：错误<br>注：该参数大小写敏感  */
	SnapshotID     string `json:"snapshotID,omitempty"`     /*  云主机快照ID，<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8349&data=87">查询云主机快照列表</a><br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8352&data=87">创建云主机快照</a>  */
	QueryContent   string `json:"queryContent,omitempty"`   /*  模糊查询查询内容，（匹配字段：instanceID、snapshotID、snapshotName）  */
	SnapshotName   string `json:"snapshotName,omitempty"`   /*  云主机快照名称。满足以下规则：不能使用中文，且长度为2-63字符  */
}

type CtecsQuerySnapshotListV41Response struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                      `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                      `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                      `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQuerySnapshotListV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsQuerySnapshotListV41ReturnObjResponse struct {
	CurrentCount int32                                                `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsQuerySnapshotListV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsQuerySnapshotListV41ReturnObjResultsResponse struct {
	SnapshotID          string                                                      `json:"snapshotID,omitempty"`          /*  云主机快照ID  */
	InstanceID          string                                                      `json:"instanceID,omitempty"`          /*  云主机ID  */
	InstanceName        string                                                      `json:"instanceName,omitempty"`        /*  云主机名称  */
	AzName              string                                                      `json:"azName,omitempty"`              /*  可用区名称  */
	SnapshotName        string                                                      `json:"snapshotName,omitempty"`        /*  云主机快照名称  */
	InstanceStatus      string                                                      `json:"instanceStatus,omitempty"`      /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	SnapshotStatus      string                                                      `json:"snapshotStatus,omitempty"`      /*  云主机快照状态，<br />pending：创建中, <br />available：可用， <br />restoring：恢复中，<br />error：错误  */
	SnapshotDescription string                                                      `json:"snapshotDescription,omitempty"` /*  云主机快照描述  */
	ProjectID           string                                                      `json:"projectID,omitempty"`           /*  企业项目ID  */
	CreatedTime         string                                                      `json:"createdTime,omitempty"`         /*  创建时间  */
	UpdatedTime         string                                                      `json:"updatedTime,omitempty"`         /*  更新时间  */
	ImageID             string                                                      `json:"imageID,omitempty"`             /*  云主机镜像ID  */
	Memory              int32                                                       `json:"memory,omitempty"`              /*  云主机内存大小，单位 MB  */
	Cpu                 int32                                                       `json:"cpu,omitempty"`                 /*  云主机cpu核数  */
	FlavorID            string                                                      `json:"flavorID,omitempty"`            /*  云主机规格ID  */
	Members             []*CtecsQuerySnapshotListV41ReturnObjResultsMembersResponse `json:"members"`                       /*  云主机的云硬盘及其快照详细信息  */
}

type CtecsQuerySnapshotListV41ReturnObjResultsMembersResponse struct {
	DiskType           string `json:"diskType,omitempty"`           /*  云硬盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO，<br />FAST-SSD：极速型SSD  */
	DiskID             string `json:"diskID,omitempty"`             /*  云硬盘ID  */
	DiskName           string `json:"diskName,omitempty"`           /*  云硬盘名称  */
	IsBootable         *bool  `json:"isBootable"`                   /*  是否是可启动磁盘,取值范围：<br />false：非启动盘，<br />true：可启动盘  */
	IsEncrypt          *bool  `json:"isEncrypt"`                    /*  是否加密盘，取值范围：<br />false：不加密，<br /> true：加密，<br /> 默认值为false  */
	DiskSize           int32  `json:"diskSize,omitempty"`           /*  云硬盘大小  */
	DiskSnapshotID     string `json:"diskSnapshotID,omitempty"`     /*  云硬盘快照ID  */
	DiskSnapshotStatus string `json:"diskSnapshotStatus,omitempty"` /*  云硬盘快照状态,详见枚举值表格  */
}
