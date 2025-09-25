package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2UnlockDBInstanceWriteApi
/* 解锁实例写操作，开启后，实例可写，该操作为幂等操作。
 */type Dcs2UnlockDBInstanceWriteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2UnlockDBInstanceWriteApi(client *core.CtyunClient) *Dcs2UnlockDBInstanceWriteApi {
	return &Dcs2UnlockDBInstanceWriteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceManageMgrServant/unlockDBInstanceWrite",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2UnlockDBInstanceWriteApi) Do(ctx context.Context, credential core.Credential, req *Dcs2UnlockDBInstanceWriteRequest) (*Dcs2UnlockDBInstanceWriteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2UnlockDBInstanceWriteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2UnlockDBInstanceWriteRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
}

type Dcs2UnlockDBInstanceWriteResponse struct {
	StatusCode int32                                       `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                      `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2UnlockDBInstanceWriteReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                      `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                      `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2UnlockDBInstanceWriteReturnObjResponse struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	LockStatus *bool  `json:"lockStatus"`           /*  锁定状态<li>true：锁定<li>false：解锁  */
}
