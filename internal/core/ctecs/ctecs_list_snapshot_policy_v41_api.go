package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListSnapshotPolicyV41Api
/* 该接口提供用户查询云主机快照策略列表的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />&emsp;&emsp;匹配查找：可以通过部分字段进行匹配筛选数据，无符合条件的为空
 */type CtecsListSnapshotPolicyV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListSnapshotPolicyV41Api(client *core.CtyunClient) *CtecsListSnapshotPolicyV41Api {
	return &CtecsListSnapshotPolicyV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot-policy/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListSnapshotPolicyV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListSnapshotPolicyV41Request) (*CtecsListSnapshotPolicyV41Response, error) {
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
	var resp CtecsListSnapshotPolicyV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListSnapshotPolicyV41Request struct {
	RegionID             string `json:"regionID,omitempty"`             /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PageNo               int32  `json:"pageNo,omitempty"`               /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize             int32  `json:"pageSize,omitempty"`             /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	SnapshotPolicyStatus int32  `json:"snapshotPolicyStatus,omitempty"` /*  快照策略状态，是否启用，取值范围：<br />0：不启用，<br />1：启用<br />注：默认值1（启用）  */
	QueryContent         string `json:"queryContent,omitempty"`         /*  模糊匹配查询内容（匹配字段：snapshotPolicyID、snapshotPolicyName）  */
}

type CtecsListSnapshotPolicyV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                       `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                       `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListSnapshotPolicyV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListSnapshotPolicyV41ReturnObjResponse struct {
	CurrentCount       int32                                                            `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount         int32                                                            `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage          int32                                                            `json:"totalPage,omitempty"`    /*  总页数  */
	SnapshotPolicyList []*CtecsListSnapshotPolicyV41ReturnObjSnapshotPolicyListResponse `json:"snapshotPolicyList"`     /*  分页明细  */
}

type CtecsListSnapshotPolicyV41ReturnObjSnapshotPolicyListResponse struct {
	SnapshotPolicyID     string `json:"snapshotPolicyID,omitempty"`     /*  云主机快照策略ID  */
	SnapshotPolicyStatus int32  `json:"snapshotPolicyStatus,omitempty"` /*  快照策略状态，是否启用，取值范围：<br />0：不启用，<br />1：启用<br />注：默认值1（启用）  */
	SnapshotPolicyName   string `json:"snapshotPolicyName,omitempty"`   /*  云主机快照策略名称，满足以下规则：长度为2~63字符，由数字、字母、-组成，只能以字母开头，以数字和字母结尾  */
	SnapshotTime         string `json:"snapshotTime,omitempty"`         /*  快照整点时间，时间取值范围：<br />0~23<br />注：如果一天内多个时间节点备份，以逗号隔开（如11点15点进行快照，则填写"11,15"），默认值0  */
	RetentionType        string `json:"retentionType,omitempty"`        /*  云主机快照保留类型，取值范围：<br />date：按时间保存，<br />num：按数量保存  */
	RetentionDay         string `json:"retentionDay,omitempty"`         /*  云主机快照保留天数，快照保留类型为date时返回，否则为空字符串  */
	RetentionNum         string `json:"retentionNum,omitempty"`         /*  云主机快照保留数量，快照保留类型为num时返回，否则为空字符串  */
	CycleType            string `json:"cycleType,omitempty"`            /*  云主机快照周期类型，取值范围：<br />day：天，<br />week：周  */
	CycleDay             int32  `json:"cycleDay,omitempty"`             /*  周期天数，周期类型为day时返回，表示多少天进行快照  */
	CycleWeek            string `json:"cycleWeek,omitempty"`            /*  周期星期，周期类型为week时返回，表示周几进行快照，由逗号拼接，由0~6组成，0表示周日  */
	ResourceCount        int32  `json:"resourceCount,omitempty"`        /*  绑定云主机数量  */
}
