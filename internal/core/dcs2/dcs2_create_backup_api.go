package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2CreateBackupApi
/* 手动备份实例数据。
 */type Dcs2CreateBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2CreateBackupApi(client *core.CtyunClient) *Dcs2CreateBackupApi {
	return &Dcs2CreateBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisMgr/createBackup",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2CreateBackupApi) Do(ctx context.Context, credential core.Credential, req *Dcs2CreateBackupRequest) (*Dcs2CreateBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Dcs2CreateBackupRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2CreateBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CreateBackupRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	Remark     string `json:"remark,omitempty"`     /*  备注信息，不超过128个字符。  */
}

type Dcs2CreateBackupResponse struct {
	StatusCode int32                              `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                             `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2CreateBackupReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                             `json:"requestId"`  /*  请求 ID。  */
	Code       string                             `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                             `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2CreateBackupReturnObjResponse struct {
	RestoreName string `json:"restoreName"` /*  备份名，格式为YYYYMMDDHHMMSS  */
}
