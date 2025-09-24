package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmReinstallInstanceApi
/* 物理机重装系统
 */type EbmReinstallInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmReinstallInstanceApi(client *core.CtyunClient) *EbmReinstallInstanceApi {
	return &EbmReinstallInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebm/rebuild",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmReinstallInstanceApi) Do(ctx context.Context, credential core.Credential, req *EbmReinstallInstanceRequest) (*EbmReinstallInstanceResponse, error) {
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
	var resp EbmReinstallInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmReinstallInstanceRequest struct {
	RegionID string `json:"regionID"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string `json:"azName"` /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string `json:"instanceUUID"` /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */Hostname string `json:"hostname"` /*  hostname，linux系统2到63位长度；windows系统2-15位长度；<br/>允许使用大小写字母、数字、连字符'-'，必须以字母开头（大小写均可），不能连续使用'-'，'-'不能用于结尾，不能仅使用数字
	 */Password string `json:"password"` /*  密码 (必须包含大小写字母和（一个数字或者特殊字符）长度8到30位)
	 */ImageUUID string `json:"imageUUID"` /*  物理机镜像UUID<br /><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4577&data=97">查询物理机镜像</a>
	 */SystemVolumeRaidUUID *string `json:"systemVolumeRaidUUID"` /*  本地系统盘raid类型，如果有本地盘则必填<br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=5881&data=97">查询物理机实例本地盘raid信息</a><br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4576&data=97">查询物理机raid</a>
	 */DataVolumeRaidUUID *string `json:"dataVolumeRaidUUID"` /*  本地数据盘raid类型，如果有本地盘则必填<br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=5881&data=97">查询物理机实例本地盘raid信息</a><br/><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4576&data=97">查询物理机raid</a>
	 */RedoRaid bool `json:"redoRaid"` /*  是否重新做raid，如果没有本地盘的必须为false
	 */UserData *string `json:"userData"` /*  用户自定义数据,需要以Base64方式编码,Base64编码后的长度限制为1-16384字符
	 */KeyName *string `json:"keyName"` /*  密钥对名称。满足以下规则：只能由数字、字母、-组成，不能以数字和-开头、以-结尾，且长度为2-63字符
	 */
}

type EbmReinstallInstanceResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmReinstallInstanceReturnObjResponse `json:"returnObj"` /*  返回参数，值为空
	 */
}

type EbmReinstallInstanceReturnObjResponse struct{}
