package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateInstanceBackupPolicyApi
/* 创建云主机备份策略<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsCreateInstanceBackupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateInstanceBackupPolicyApi(client *core.CtyunClient) *CtecsCreateInstanceBackupPolicyApi {
	return &CtecsCreateInstanceBackupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-policy/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateInstanceBackupPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtecsCreateInstanceBackupPolicyRequest) (*CtecsCreateInstanceBackupPolicyResponse, error) {
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
	var resp CtecsCreateInstanceBackupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateInstanceBackupPolicyRequest struct {
	RegionID           string        `json:"regionID,omitempty"`           /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyName         string        `json:"policyName,omitempty"`         /*  云主机备份策略名称。满足以下规则：只能由数字、英文字母、中划线-、下划线_、点.组成，长度为2-64字符<br />注：在所有资源池不可重复  */
	CycleType          string        `json:"cycleType,omitempty"`          /*  云主机备份周期类型，取值范围：day（按天备份）week（按星期备份）  */
	CycleDay           int32         `json:"cycleDay,omitempty"`           /*  备份周期（天），取值范围：[1, 30]，默认值为1  <br />注：cycleType为day时需设置  */
	CycleWeek          string        `json:"cycleWeek,omitempty"`          /*  备份周期（星期），星期取值范围：0~6（代表周几，其中0为周日），默认值是0<br />注：只有cycleType为week时需设置；<br />如果一周有多天备份，以逗号隔开（如周日周三进行快照，则填写"0,3"）  */
	Time               string        `json:"time,omitempty"`               /*  备份整点时间，时间取值范围：0~23<br />注：如果一天内多个时间节点备份，以逗号隔开（如11点15点进行快照，则填写"11,15"），默认值12  */
	Status             int32         `json:"status,omitempty"`             /*  备份策略状态，是否启用，取值范围：<br />0（不启用），<br />1（启用）<br />注：默认值0（不启用）  */
	RetentionType      string        `json:"retentionType,omitempty"`      /*  云主机备份保留类型，取值范围：<br />date（按时间保存），<br />num（按数量保存），<br />all（永久保存）  */
	RetentionDay       int32         `json:"retentionDay,omitempty"`       /*  云主机备份保留天数，单位为天，取值范围：[1, 99999] ，默认值1<br />注：retentionType为date时必填  */
	RetentionNum       int32         `json:"retentionNum,omitempty"`       /*  云主机备份保留数量，取值范围：[1, 99999]，默认值1<br />注：retentionType为num时必填  */
	ProjectID          string        `json:"projectID,omitempty"`          /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
	FullBackupInterval int32         `json:"fullBackupInterval,omitempty"` /*  是否启用周期性全量备份。-1代表不开启，默认为-1；取值范围为[-1,100]，即每执行n次增量备份后，执行一次全量备份；若传入为0，代表每一次均为全量备份。  */
	AdvRetentionStatus bool          `json:"advRetentionStatus,omitempty"` /*  是否开启高级保留策略，false（不启用），true(启用)，默认值为false。需校验云主机备份保留类型（retentionType），若保留类型为按数量保存（num），可开启高级保留策略；若保留类型为date（按时间保存）或all（永久保存），不可开启高级保留策略。  */
	AdvRetention       *AdvRetention `json:"advRetention,omitempty"`       /*  高级保留策略内容，只有retentionType为num且advRetentionStatus为true才生效  */
}

type CtecsCreateInstanceBackupPolicyResponse struct {
	StatusCode  int32                                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                            `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                            `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                            `json:"message,omitempty"`     /*  错误信息的英文描述  */
	Description string                                            `json:"description,omitempty"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *CtecsCreateInstanceBackupPolicyReturnObjResponse `json:"returnObj"`             /*  返回对象  */
}

type CtecsCreateInstanceBackupPolicyReturnObjResponse struct {
	Status        int32  `json:"status,omitempty"`        /*  备份策略状态  */
	PolicyName    string `json:"policyName,omitempty"`    /*  备份策略名称  */
	RetentionType string `json:"retentionType,omitempty"` /*  备份保留类型  */
	RetentionDay  int32  `json:"retentionDay,omitempty"`  /*  保留时间，当retentionType为date时返回  */
	RetentionNum  int32  `json:"retentionNum,omitempty"`  /*  保留数量，当retentionType为num时返回  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
	CycleType     string `json:"cycleType,omitempty"`     /*  备份周期类型  */
	CycleDay      int32  `json:"cycleDay,omitempty"`      /*  cycleType为day时返回备份周期值  */
	CycleWeek     string `json:"cycleWeek,omitempty"`     /*  cycleType为week时返回备份周期值  */
	PolicyID      string `json:"policyID,omitempty"`      /*  备份策略ID  */
	Time          string `json:"time,omitempty"`          /*  备份整点时间  */
	ProjectID     string `json:"projectID,omitempty"`     /*  企业项目ID  */
}
