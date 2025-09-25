package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreatePrivateZoneRecordApi
/* 创建内网 DNS 记录
 */type CtvpcCreatePrivateZoneRecordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreatePrivateZoneRecordApi(client *core.CtyunClient) *CtvpcCreatePrivateZoneRecordApi {
	return &CtvpcCreatePrivateZoneRecordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone-record/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreatePrivateZoneRecordApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreatePrivateZoneRecordRequest) (*CtvpcCreatePrivateZoneRecordResponse, error) {
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
	var resp CtvpcCreatePrivateZoneRecordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreatePrivateZoneRecordRequest struct {
	ClientToken string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string   `json:"regionID,omitempty"`    /*  资源池ID  */
	ZoneID      string   `json:"zoneID,omitempty"`      /*  zoneID  */
	RawType     string   `json:"type,omitempty"`        /*  支持: A / CNAME / MX / AAAA / TXT, 大小写不敏感  */
	ValueList   []string `json:"valueList"`             /*  最多同时支持 8 个，最多同时支持 8 个，当 type = A 时，valueList 中必须是 IPv4 地址；当 type = CNAME 时，valueList 只能存在一个元素, 数组中元素格式要求如下：域名以点号分隔成多个字符串, 单个字符串由各国文字的特定字符集、字母、数字、连字符（-）组成，字母不区分大小写，连字符（-）不得出现在字符串的头部或者尾部, 单个字符串长度不超过63个字符, 总长度不超过 254；当 type = AAAA 时，valueList 中必须是 IPv6 地址；当 type = TXT 时，valueList 数组中元素格式要求如下：支持数字，字符，符号：~!@#$%^&*()_+-={}[]:;',,./<>?，空格，且元素中的值和 zone record name 拼接起来长度不能超过 256；当 type = MX 时，valueList 数组中元素格式要求如下：priority dnsname，priority 的取值在 0 - 65535，dnsname 域名以点号分隔成多个字符串, 单个字符串由字母、数字、连字符（-）组成，字母不区分大小写，连字符（-）不得出现在字符串的头部或者尾部, 单个字符串长度不超过63个字符, 总长度不超过 254：一个例子：0 ctyun.cn  */
	TTL         int32    `json:"TTL"`                   /*  zone ttl, 单位秒。default is 300，大于等于300，小于等于2147483647  */
	Name        *string  `json:"name,omitempty"`        /*  dns 记录集的 name 长度 + dns 记录的 name 长度，总和不超过 256, 支持数字、英文字母、连字符、*, 不能以连字符结尾或开头  */
}

type CtvpcCreatePrivateZoneRecordResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreatePrivateZoneRecordReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                        `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreatePrivateZoneRecordReturnObjResponse struct {
	ZoneRecordID *string `json:"zoneRecordID,omitempty"` /*  名称  */
	Name         *string `json:"name,omitempty"`         /*  zone record名称  */
}
