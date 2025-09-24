package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsNewDirectoryApi
/* 指定并行文件创建目录并设置权限
 */ /* 此请求是异步处理，返回800代表请求下发成功，具体结果请使用【查询并行文件目录信息】确定是否创建成功
 */type HpfsNewDirectoryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsNewDirectoryApi(client *core.CtyunClient) *HpfsNewDirectoryApi {
	return &HpfsNewDirectoryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/new-directory",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsNewDirectoryApi) Do(ctx context.Context, credential core.Credential, req *HpfsNewDirectoryRequest) (*HpfsNewDirectoryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsNewDirectoryRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsNewDirectoryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsNewDirectoryRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池 ID  */
	SfsUID           string `json:"sfsUID,omitempty"`           /*  并行文件唯一ID  */
	SfsDirectory     string `json:"sfsDirectory,omitempty"`     /*  并行文件目录，目录名仅允许数字、字母、下划线、连接符、中文组成，每级目录最大长度为255字节，最大全路径长度为4096字节，最大目录层数为1000，如果参数为mydir/、mydir、/mydir或/mydir/，则都视为输入/mydir的目录  */
	SfsDirectoryMode string `json:"sfsDirectoryMode,omitempty"` /*  目录权限，默认值是755，若传入则必须为三位，每位的范围为0到7。第一位表示目录所有者的权限，第二位表示目录所属用户组的权限，第三位表示其他用户的权限。目录所有者由uid指定，目录所属用户组由gid指定，不是目录所有者且不在目录所属用户组的用户为其他用户。例如：755中第一位7代表该目录所有者对该目录具有读、写、执行权限；第二位5代表该目录所属用户组对该目录具有读、执行权限；第三位5代表其他用户对该目录具有读、执行权限  */
	SfsDirectoryUID  int64  `json:"sfsDirectoryUID,omitempty"`  /*  目录所有者的用户id，默认值是0，取值范围是0到4,294,967,294（即2^32-2）  */
	SfsDirectoryGID  int64  `json:"sfsDirectoryGID,omitempty"`  /*  目录所属用户组id，默认值是0，取值范围是0到4,294,967,294（即2^32-2）  */
}

type HpfsNewDirectoryResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述  */
	Description string `json:"description"` /*  响应描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}
