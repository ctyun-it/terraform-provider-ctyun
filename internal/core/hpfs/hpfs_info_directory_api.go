package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsInfoDirectoryApi
/* 查询指定文件系统的指定目录信息
 */type HpfsInfoDirectoryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsInfoDirectoryApi(client *core.CtyunClient) *HpfsInfoDirectoryApi {
	return &HpfsInfoDirectoryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/info-directory",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsInfoDirectoryApi) Do(ctx context.Context, credential core.Credential, req *HpfsInfoDirectoryRequest) (*HpfsInfoDirectoryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	ctReq.AddParam("sfsDirectory", req.SfsDirectory)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsInfoDirectoryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsInfoDirectoryRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池 ID  */
	SfsUID       string `json:"sfsUID,omitempty"`       /*  并行文件唯一ID  */
	SfsDirectory string `json:"sfsDirectory,omitempty"` /*  并行文件目录，目录名仅允许数字、字母、下划线、连接符、中文组成，每级目录最大长度为255字节，最大目录层数为1000，最大全路径长度为4096字节，如果参数为mydir/、mydir、/mydir或/mydir/，则都视为输入/mydir的目录  */
}

type HpfsInfoDirectoryResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                              `json:"message"`     /*  响应描述  */
	Description string                              `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsInfoDirectoryReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsInfoDirectoryReturnObjResponse struct {
	SfsDirectory     string `json:"sfsDirectory"`     /*  并行文件目录  */
	SfsDirectoryMode string `json:"sfsDirectoryMode"` /*  目录权限，每位的范围为0到7。第一位表示目录所有者的权限，第二位表示目录所属用户组的权限，第三位表示其他用户的权限。目录所有者由uid指定，目录所属用户组由gid指定，不是目录所有者且不在目录所属用户组的用户为其他用户。例如：755中第一位7代表该目录所有者对该目录具有读、写、执行权限；第二位5代表该目录所属用户组对该目录具有读、执行权限；第三位5代表其他用户对该目录具有读、执行权限  */
	SfsDirectoryUID  int64  `json:"sfsDirectoryUID"`  /*  目录所有者的用户id，取值范围是0到4,294,967,294（即2^32-2）  */
	SfsDirectoryGID  int64  `json:"sfsDirectoryGID"`  /*  目录所属用户组id，取值范围是0到4,294,967,294（即2^32-2）  */
}
