package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosHeadObjectApi
/* 查询对象是否存在。
 */type ZosHeadObjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosHeadObjectApi(client *core.CtyunClient) *ZosHeadObjectApi {
	return &ZosHeadObjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/head-object",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosHeadObjectApi) Do(ctx context.Context, credential core.Credential, req *ZosHeadObjectRequest) (*ZosHeadObjectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("key", req.Key)
	if req.VersionID != "" {
		ctReq.AddParam("versionID", req.VersionID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosHeadObjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosHeadObjectRequest struct {
	Bucket    string /*  桶名  */
	RegionID  string /*  区域 ID  */
	Key       string /*  对象名  */
	VersionID string /*  版本ID，在开启多版本时可使用  */
}

type ZosHeadObjectResponse struct {
	ReturnObj   *ZosHeadObjectReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	StatusCode  int64                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                          `json:"message,omitempty"`     /*  状态描述  */
	Description string                          `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                          `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosHeadObjectReturnObjResponse struct {
	DeleteMarker              *bool                                   `json:"deleteMarker"`                        /*  指定检索到的对象是（true）还是不是（false）删除标记。如果为 false，则此响应标头不会出现在响应中  */
	AcceptRanges              string                                  `json:"acceptRanges,omitempty"`              /*  表示指定了一个字节范围  */
	Expiration                string                                  `json:"expiration,omitempty"`                /*  如果配置了对象过期（请参阅 PUT Bucket 生命周期），则响应包含此标头。它包括提供对象过期信息的过期日期和规则 ID 键值对。 rule-id 的值是 URL 编码的  */
	Restore                   string                                  `json:"restore,omitempty"`                   /*  如果对象是存档对象（其存储类为 GLACIER 的对象），则如果存档恢复正在进行（请参阅 RestoreObject 或存档副本已恢复），则响应将包含此标头。如果存档副本已恢复，则标头值指示对象存储计划何时删除对象副本。例如：x-amz-restore: concurrent-request="false", expiry-date="Fri, 21 Dec 2012 00:00:00 GMT"。如果对象恢复正在进行中，标头返回值 concurrent-request="true"  */
	ArchiveStatus             string                                  `json:"archiveStatus,omitempty"`             /*  head 对象的归档状态。可能的值为 ARCHIVE_ACCESS，DEEP_ARCHIVE_ACCESS  */
	LastModified              string                                  `json:"lastModified,omitempty"`              /*  对象创建日期， ISO8601 格式字符串  */
	ContentLength             int64                                   `json:"contentLength,omitempty"`             /*  内容长度  */
	ETag                      string                                  `json:"ETag,omitempty"`                      /*  ETag  */
	MissingMeta               int64                                   `json:"missingMeta,omitempty"`               /*  这设置为 x-amz-meta 标头中未返回的元数据条目数。如果您使用像 SOAP 这样支持比 REST API 更灵活的元数据的 API 创建元数据，就会发生这种情况。例如，使用 SOAP，您可以创建其值不是合法 HTTP 标头的元数据  */
	VersionID                 string                                  `json:"versionID,omitempty"`                 /*  版本ID，在开启多版本时可使用  */
	CacheControl              string                                  `json:"cacheControl,omitempty"`              /*  指定沿请求/响应链的缓存行为  */
	ContentDisposition        string                                  `json:"contentDisposition,omitempty"`        /*  指定对象的表示信息  */
	ContentEncoding           string                                  `json:"contentEncoding,omitempty"`           /*  内容编码  */
	ContentLanguage           string                                  `json:"contentLanguage,omitempty"`           /*  内容语言  */
	ContentType               string                                  `json:"contentType,omitempty"`               /*  内容类型，枚举值可参考HTTP content-type类型  */
	Expires                   string                                  `json:"expires,omitempty"`                   /*  对象不再可缓存的日期和时间。 ISO8601 格式字符串。  */
	WebsiteRedirectLocation   string                                  `json:"websiteRedirectLocation,omitempty"`   /*  网站重定向位置  */
	ServerSideEncryption      string                                  `json:"serverSideEncryption,omitempty"`      /*  服务端加密算法，AES256，aws:kms  */
	Metadata                  *ZosHeadObjectReturnObjMetadataResponse `json:"metadata"`                            /*  与 S3 中的对象一起存储的元数据映射  */
	SSECustomerAlgorithm      string                                  `json:"SSECustomerAlgorithm,omitempty"`      /*  如果请求使用客户提供的加密密钥进行服务器端加密，则响应将包含此参数，以确认所使用的加密算法。  */
	SSECustomerKeyMD5         string                                  `json:"SSECustomerKeyMD5,omitempty"`         /*  如果请求使用客户提供的加密密钥进行服务器端加密，则响应将包含此标头以提供往返消息的完整性验证  */
	SSEKMSKeyID               string                                  `json:"SSEKMSKeyID,omitempty"`               /*  SSEKMSKeyID  */
	BucketKeyEnabled          *bool                                   `json:"bucketKeyEnabled"`                    /*  指示对象是否通过服务端加密  */
	StorageClass              string                                  `json:"storageClass,omitempty"`              /*  存储类，可能的值有：STANDARD（标准存储）、STANDARD_IA（低频存储）、GLACIER（归档存储）  */
	RequestCharged            string                                  `json:"requestCharged,omitempty"`            /*  如果存在，则表明请求者已成功为请求收费  */
	ReplicationStatus         string                                  `json:"replicationStatus,omitempty"`         /*  复制状态，如COMPLETE，PENDING，FAILED，REPLICA  */
	PartsCount                int64                                   `json:"partsCount,omitempty"`                /*  此对象拥有的分段数  */
	ObjectLockMode            string                                  `json:"objectLockMode,omitempty"`            /*  对象锁定模式，GOVERNANCE，COMPLIANCE  */
	ObjectLockRetainUntilDate string                                  `json:"objectLockRetainUntilDate,omitempty"` /*  对象锁定保留期到期的日期和时间。ISO8601 格式字符串  */
	ObjectLockLegalHoldStatus string                                  `json:"objectLockLegalHoldStatus,omitempty"` /*  指定此对象的合法保留是否有效，可能的值是 ON，OFF  */
}

type ZosHeadObjectReturnObjMetadataResponse struct{}
