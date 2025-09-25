package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeOffLineKeyAnalysisTaskInfoApi
/* 查询分布式缓存Redis实例离线全量key分析报告详情。
 */type Dcs2DescribeOffLineKeyAnalysisTaskInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeOffLineKeyAnalysisTaskInfoApi(client *core.CtyunClient) *Dcs2DescribeOffLineKeyAnalysisTaskInfoApi {
	return &Dcs2DescribeOffLineKeyAnalysisTaskInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/keyAnalysisMgrServant/describeOffLineKeyAnalysisTaskInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeOffLineKeyAnalysisTaskInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeOffLineKeyAnalysisTaskInfoRequest) (*Dcs2DescribeOffLineKeyAnalysisTaskInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("taskId", req.TaskId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeOffLineKeyAnalysisTaskInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	TaskId     string /*  任务ID  */
}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoResponse struct {
	StatusCode int32                                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjResponse struct {
	TaskId                         string                                                          `json:"taskId,omitempty"`               /*  任务ID  */
	Rdbtype                        string                                                          `json:"rdbtype,omitempty"`              /*  rdb文件类型<li>0：上一个备份文件<li>1：新建备份文件<li>2：历史备份文件  */
	UsedMemory                     int64                                                           `json:"usedMemory,omitempty"`           /*  已使用内存  */
	Maxmemory                      int64                                                           `json:"maxmemory,omitempty"`            /*  总内存​  */
	UsedMemoryPercentage           float64                                                         `json:"usedMemoryPercentage"`           /*  已使用内存比例%  */
	Totalkeys                      int64                                                           `json:"totalkeys,omitempty"`            /*  总key数量  */
	Total_set_expirekeys           int64                                                           `json:"total_set_expirekeys,omitempty"` /*  设置过期key总数量  */
	Total_set_expirekeysPercentage float64                                                         `json:"total_set_expirekeysPercentage"` /*  设置过期key的比例  */
	ExpiredKeys                    int64                                                           `json:"expiredKeys,omitempty"`          /*  已过期key数量  */
	EvictedKeys                    int64                                                           `json:"evictedKeys,omitempty"`          /*  已逐出key数量  */
	StartTime                      string                                                          `json:"startTime,omitempty"`            /*  开始时间  */
	EndTime                        string                                                          `json:"endTime,omitempty"`              /*  结束时间  */
	Rdbinfo                        *Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoResponse `json:"rdbinfo"`                        /*  该对象返回格式随key类型变化  */
}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoResponse struct {
	LargestKeyPrefixes *Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoLargestKeyPrefixesResponse `json:"LargestKeyPrefixes"`        /*  最大key占用  */
	TypeBytes          *Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoTypeBytesResponse          `json:"TypeBytes"`                 /*  类型直接占用大小  */
	TypeNum            *Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoTypeNumResponse            `json:"TypeNum"`                   /*  key类型数量大小  */
	CurrentInstance    string                                                                            `json:"CurrentInstance,omitempty"` /*  CurrentInstance  */
	TotleBytes         int64                                                                             `json:"TotleBytes,omitempty"`      /*  总字节数  */
	TotleNum           int64                                                                             `json:"TotleNum,omitempty"`        /*  key总数  */
	LenLevelCount      *Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoLenLevelCountResponse      `json:"LenLevelCount"`             /*  LenLevelCount  */
}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoLargestKeyPrefixesResponse struct {
	Bytes int64  `json:"Bytes,omitempty"` /*  字节  */
	Type  string `json:"Type,omitempty"`  /*  key类型  */
	Num   int64  `json:"Num,omitempty"`   /*  数量  */
	Key   string `json:"Key,omitempty"`   /*  前缀  */
}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoTypeBytesResponse struct{}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoTypeNumResponse struct{}

type Dcs2DescribeOffLineKeyAnalysisTaskInfoReturnObjRdbinfoLenLevelCountResponse struct{}
