package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeInstanceConfigV2Api
/* 查询分布式缓存Redis实例配置参数
 */type Dcs2DescribeInstanceConfigV2Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeInstanceConfigV2Api(client *core.CtyunClient) *Dcs2DescribeInstanceConfigV2Api {
	return &Dcs2DescribeInstanceConfigV2Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceParam/describeInstanceConfig",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeInstanceConfigV2Api) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeInstanceConfigV2Request) (*Dcs2DescribeInstanceConfigV2Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeInstanceConfigV2Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeInstanceConfigV2Request struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeInstanceConfigV2Response struct {
	StatusCode int32                                          `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                         `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeInstanceConfigV2ReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeInstanceConfigV2ReturnObjResponse struct {
	Maxmemory            string `json:"maxmemory,omitempty"`            /*  指定redis最大内存限制，1-1099511627776整数  */
	Databases            string `json:"databases,omitempty"`            /*  设置数据库的数量  */
	IsOpenSemiSync       string `json:"isOpenSemiSync,omitempty"`       /*  是否打开半同步开关  */
	Appendfsync          string `json:"appendfsync,omitempty"`          /*  指定日志更新条件<li>no<li>everysec<li>always  */
	Dir                  string `json:"dir,omitempty"`                  /*  指定本地数据库存放目录  */
	Logfile              string `json:"logfile,omitempty"`              /*  指定日志文件  */
	AofNewsaveSize       string `json:"aofNewsaveSize,omitempty"`       /*  默认的aof文件重写大小  */
	Maxclients           string `json:"maxclients,omitempty"`           /*  同一时间最大客户端连接数  */
	IsOpenVersion        string `json:"isOpenVersion,omitempty"`        /*  是否开启乐观锁模板  */
	AofNewsavePercentage string `json:"aofNewsavePercentage,omitempty"` /*  默认的aof文件重写内存占比的大小  */
}
