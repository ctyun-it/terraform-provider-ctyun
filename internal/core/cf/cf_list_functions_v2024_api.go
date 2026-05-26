package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListFunctionsV2024Api
/* 分页查询函数列表 */
type CfListFunctionsV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListFunctionsV2024Api(client *core.CtyunClient) *CfListFunctionsV2024Api {
	return &CfListFunctionsV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListFunctionsV2024Api) Do(ctx context.Context, credential core.Credential, req *CfListFunctionsV2024Request) (*CfListFunctionsV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListFunctionsV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListFunctionsV2024Request struct {
	RegionId  string  `json:"regionId"`            /*  资源池id  */
	PageIndex *int32  `json:"pageIndex,omitempty"` /*  页码，默认为1  */
	PageSize  *int32  `json:"pageSize,omitempty"`  /*  分页大小，默认为10  */
	Search    *string `json:"search,omitempty"`    /*  模糊查询的关键字，默认为空  */
}

type CfListFunctionsV2024Response struct {
	StatusCode *int32                                 `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                                `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListFunctionsV2024ReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListFunctionsV2024ReturnObjResponse struct {
	Data       []*CfListFunctionsV2024ReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListFunctionsV2024ReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListFunctionsV2024ReturnObjDataResponse struct {
	FunctionId            *string                                                         `json:"functionId"`            /*  函数 id  */
	DeployInfo            *CfListFunctionsV2024ReturnObjDataDeployInfoResponse            `json:"deployInfo"`            /*  最近一次函数构建的状态信息  */
	Container             *CfListFunctionsV2024ReturnObjDataContainerResponse             `json:"container"`             /*  容器配置  */
	Lifecycle             *CfListFunctionsV2024ReturnObjDataLifecycleResponse             `json:"lifecycle"`             /*  生命周期配置  */
	Log                   *CfListFunctionsV2024ReturnObjDataLogResponse                   `json:"log"`                   /*  是否启用日志功能  */
	Layers                []*CfListFunctionsV2024ReturnObjDataLayersResponse              `json:"layers"`                /*  层配置  */
	Runtime               *CfListFunctionsV2024ReturnObjDataRuntimeResponse               `json:"runtime"`               /*  运行时  */
	Network               *CfListFunctionsV2024ReturnObjDataNetworkResponse               `json:"network"`               /*  网络配置  */
	CustomContainerConfig *CfListFunctionsV2024ReturnObjDataCustomContainerConfigResponse `json:"customContainerConfig"` /*  自定义镜像配置  */
	ServerlessGpuConfig   *CfListFunctionsV2024ReturnObjDataServerlessGpuConfigResponse   `json:"serverlessGpuConfig"`   /*  gpu函数配置  */
	CreateType            *int32                                                          `json:"createType"`            /*  创建类型 1:内置运行时2:自定义运行时3:自定义镜像  */
	OssMount              *CfListFunctionsV2024ReturnObjDataOssMountResponse              `json:"ossMount"`              /*  zos挂载配置  */
	FunctionName          *string                                                         `json:"functionName"`          /*  函数名  */
	Description           *string                                                         `json:"description"`           /*  说明  */
	Dns                   *CfListFunctionsV2024ReturnObjDataDnsResponse                   `json:"dns"`                   /*  DNS配置  */
	Nas                   *CfListFunctionsV2024ReturnObjDataNasResponse                   `json:"nas"`                   /*  NAS配置  */
}

type CfListFunctionsV2024ReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}

type CfListFunctionsV2024ReturnObjDataDeployInfoResponse struct {
	TaskEnd   *int32  `json:"taskEnd"`   /*  结束时间  */
	ErrMsg    *string `json:"errMsg"`    /*  错误说明  */
	Creator   *string `json:"creator"`   /*  关联的userId  */
	Id        *string `json:"id"`        /*  此次更新内部唯一id  */
	Status    *string `json:"status"`    /*  当前代码部署状态  */
	TaskBegin *int32  `json:"taskBegin"` /*  开始时间  */
}

type CfListFunctionsV2024ReturnObjDataContainerResponse struct {
	TimeZone             *string                                                              `json:"timeZone"`             /*  时区  */
	DiskSize             *int32                                                               `json:"diskSize"`             /*  磁盘规格(Mb)  */
	MaxScale             *int32                                                               `json:"maxScale"`             /*  并发实例数上限  */
	FastStart            *int32                                                               `json:"fastStart"`            /*  启动加速, 默认为0, 1(表示使用启动加速)  */
	HealthCheckConfig    *CfListFunctionsV2024ReturnObjDataContainerHealthCheckConfigResponse `json:"healthCheckConfig"`    /*  健康检查  */
	EnvironmentVariables map[string]string                                                    `json:"environmentVariables"` /*  环境变量  */
	MemorySize           *int32                                                               `json:"memorySize"`           /*  内存规格(Mb)  */
	Cpu                  *float64                                                             `json:"cpu"`                  /*  CPU规格(vCPU)  */
	RunCommand           *string                                                              `json:"runCommand"`           /*  函数服务启动命令  */
	ListenPort           *int32                                                               `json:"listenPort"`           /*  监听端口  */
}

type CfListFunctionsV2024ReturnObjDataLifecycleResponse struct {
	Initializer *CfListFunctionsV2024ReturnObjDataLifecycleInitializerResponse `json:"initializer"` /*  实例初始化  */
	PreStop     *CfListFunctionsV2024ReturnObjDataLifecyclePreStopResponse     `json:"preStop"`     /*  实例初始化  */
}

type CfListFunctionsV2024ReturnObjDataLogResponse struct {
	LogEnabled     *bool                                                `json:"logEnabled"`     /*  是否启用日志功能  */
	LogAutoConfig  *bool                                                `json:"logAutoConfig"`  /*  是否自动配置  */
	LogProjectId   *string                                              `json:"logProjectId"`   /*  日志项目id  */
	LogProjectCode *string                                              `json:"logProjectCode"` /*  日志项目Code  */
	LogUnit        *string                                              `json:"logUnit"`        /*  日志单元名称  */
	LogUnitId      *string                                              `json:"logUnitId"`      /*  日志单元id  */
	LogProject     *string                                              `json:"logProject"`     /*  日志项目名称  */
	LogUnitCode    *string                                              `json:"logUnitCode"`    /*  日志单元Code  */
	LogRuleEnabled *bool                                                `json:"logRuleEnabled"` /*  启用日志分割规则  */
	LogRule        *CfListFunctionsV2024ReturnObjDataLogLogRuleResponse `json:"logRule"`        /*  日志切割配置  */
}

type CfListFunctionsV2024ReturnObjDataLayersResponse struct {
	LayerName   *string `json:"layerName"`   /*  层名称  */
	Version     *int32  `json:"version"`     /*  版本  */
	Description *string `json:"description"` /*  描述  */
	Acl         *int32  `json:"acl"`         /*  acl 0表示自定义层，1表示官方公共层  */
}

type CfListFunctionsV2024ReturnObjDataRuntimeResponse struct {
	Runtime             *string `json:"runtime"`             /*  运行时类型  */
	HandleType          *string `json:"handleType"`          /*  请求处理程序类型  */
	ExecuteTimeout      *int32  `json:"executeTimeout"`      /*  执行超时时间  */
	Handler             *string `json:"handler"`             /*  函数执行的入口  */
	InstanceConcurrency *int32  `json:"instanceConcurrency"` /*  实例最大并发度  */
}

type CfListFunctionsV2024ReturnObjDataNetworkResponse struct {
	InternetOutAllow  *bool     `json:"InternetOutAllow"`  /*  允许函数访问公网  */
	InternetInForbid  *bool     `json:"internetInForbid"`  /*  不允许互联网公网访问函数  */
	OutVpcId          *string   `json:"outVpcId"`          /*  networkId  */
	VpcId             *int32    `json:"vpcId"`             /*  vpcId  */
	SecurityGroupName *string   `json:"securityGroupName"` /*  安全组Name  */
	SubNetId          *string   `json:"subNetId"`          /*  子网ID  */
	SubNetName        *string   `json:"subNetName"`        /*  子网name  */
	Enable            *bool     `json:"enable"`            /*  是否开启VPC  */
	FixedPublicIp     *bool     `json:"fixedPublicIp"`     /*  固定公网ip  */
	OutVpcName        *string   `json:"outVpcName"`        /*  vpcName  */
	SecurityGroupId   *string   `json:"securityGroupId"`   /*  安全组ID  */
	SubNetCidr        *string   `json:"subNetCidr"`        /*  子网CIDR  */
	AccessVpcIds      []*string `json:"accessVpcIds"`      /*  仅允许指定的VPC访问函数  */
}

type CfListFunctionsV2024ReturnObjDataCustomContainerConfigResponse struct {
	ImageDigest *string `json:"imageDigest"` /*  digest 用于指定镜像版本  */
	Image       *string `json:"image"`       /*  容器镜像地址  */
	InstanceId  *string `json:"instanceId"`  /*  crs 实例 id  */
}

type CfListFunctionsV2024ReturnObjDataServerlessGpuConfigResponse struct {
	GpuEnable         *bool   `json:"gpuEnable"`         /*  是否使用Gpu  */
	GpuEciType        *string `json:"gpuEciType"`        /*  GPU ECI 规格  */
	GpuMemorySize     *int32  `json:"gpuMemorySize"`     /*  单位是G  */
	GpuType           *string `json:"gpuType"`           /*  gpu卡型  */
	GpuProvisionCount *int32  `json:"gpuProvisionCount"` /*  配置的预留实例数量  */
}

type CfListFunctionsV2024ReturnObjDataOssMountResponse struct {
	Mounts []*CfListFunctionsV2024ReturnObjDataOssMountMountsResponse `json:"mounts"` /*  zos挂载参数  */
}

type CfListFunctionsV2024ReturnObjDataDnsResponse struct {
	NameServers []*string                                            `json:"nameServers"` /*  DNS 服务器的 IP 地址列表  */
	Searches    []*string                                            `json:"searches"`    /*  DNS 搜索域列表  */
	Options     *CfListFunctionsV2024ReturnObjDataDnsOptionsResponse `json:"options"`     /*  DNS 解析配置  */
}

type CfListFunctionsV2024ReturnObjDataNasResponse struct {
	Nas []*CfListFunctionsV2024ReturnObjDataNasNasResponse `json:"nas"` /*  nas  */
}

type CfListFunctionsV2024ReturnObjDataContainerHealthCheckConfigResponse struct {
	FailureThreshold    *int32  `json:"failureThreshold"`    /*  失败阈值  */
	GetPath             *string `json:"getPath"`             /*  检查http get path  */
	InitialDelaySeconds *int32  `json:"initialDelaySeconds"` /*  首次探测延迟时间(秒)  */
	PeriodSeconds       *int32  `json:"periodSeconds"`       /*  探测时间间隔(秒)  */
	SuccessThreshold    *int32  `json:"successThreshold"`    /*  成功阈值  */
	TimeoutSeconds      *int32  `json:"timeoutSeconds"`      /*  超时(秒)  */
}

type CfListFunctionsV2024ReturnObjDataLifecycleInitializerResponse struct {
	Handler *string `json:"handler"` /*  处理方法入口  */
	Enable  *bool   `json:"enable"`  /*  启用  */
	Timeout *int32  `json:"timeout"` /*  超时  */
}

type CfListFunctionsV2024ReturnObjDataLifecyclePreStopResponse struct {
	Handler *string `json:"handler"` /*  处理方法入口  */
	Enable  *bool   `json:"enable"`  /*  启用  */
	Timeout *int32  `json:"timeout"` /*  超时  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleResponse struct {
	RuleCode         *string                                                        `json:"ruleCode"`         /*  规则唯一编码  */
	RuleName         *string                                                        `json:"ruleName"`         /*  规则名称  */
	ExtractMode      *int32                                                         `json:"extractMode"`      /*  采集类型  */
	CollectPolicy    *string                                                        `json:"collectPolicy"`    /*  采集策略  */
	CuttingMode      *string                                                        `json:"cuttingMode"`      /*  切割模式  */
	Enable           *bool                                                          `json:"enable"`           /*  是否启用采集规则  */
	UnitCode         *string                                                        `json:"unitCode"`         /*  日志单元编码  */
	FirstLinePattern *string                                                        `json:"firstLinePattern"` /*  首行正则  */
	CustomTime       *CfListFunctionsV2024ReturnObjDataLogLogRuleCustomTimeResponse `json:"customTime"`       /*  自定义时间戳提取格式  */
	RuleConfig       *CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigResponse `json:"ruleConfig"`       /*  容器运行参数  */
	AccessType       *int32                                                         `json:"accessType"`       /*  接入类型  */
}

type CfListFunctionsV2024ReturnObjDataOssMountMountsResponse struct {
	BucketName *string `json:"bucketName"` /*  bucket名  */
	BucketPath *string `json:"bucketPath"` /*  bucket子目录  */
	MountDir   *string `json:"mountDir"`   /*  挂载本地目录  */
	ReadOnly   *bool   `json:"readOnly"`   /*  是否只读，默认false  */
	AccessUrl  *string `json:"accessUrl"`  /*  oss 访问地址  */
}

type CfListFunctionsV2024ReturnObjDataDnsOptionsResponse struct {
	Ndots *string `json:"ndots"` /*  键值对01  */
}

type CfListFunctionsV2024ReturnObjDataNasNasResponse struct {
	RemoteDir *string `json:"remoteDir"` /*  远端挂载目录  */
	SharePath *string `json:"sharePath"` /*  挂载地址  */
	LocalDir  *string `json:"localDir"`  /*  挂载本地目录  */
	SfsName   *string `json:"sfsName"`   /*  sfs 的名称  */
	SfsUID    *string `json:"sfsUID"`    /*  sfs 的 ID  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleCustomTimeResponse struct {
	Key        *string `json:"key"`        /*  key  */
	TimeFormat *string `json:"timeFormat"` /*  格式化  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigResponse struct {
	MaxPathDepth *int32                                                                  `json:"maxPathDepth"` /*  最大正则路径解析深度  */
	Delimeter    *CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigDelimeterResponse `json:"delimeter"`    /*  分隔符  */
	Regex        *CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigRegexResponse     `json:"regex"`        /*  正则切割模式  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigDelimeterResponse struct {
	Delimeter    *string                                                                               `json:"delimeter"`    /*  分隔符  */
	TypeContents []*CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigDelimeterTypeContentsResponse `json:"typeContents"` /*  类型  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigRegexResponse struct {
	RegexStr     *string                                                                           `json:"regexStr"`     /*  正则表达式  */
	TypeContents []*CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigRegexTypeContentsResponse `json:"typeContents"` /*  类型  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigDelimeterTypeContentsResponse struct {
	Key     *string `json:"key"`  /*  key  */
	RawType *string `json:"type"` /*  类型  */
}

type CfListFunctionsV2024ReturnObjDataLogLogRuleRuleConfigRegexTypeContentsResponse struct {
	Key     *string `json:"key"`  /*  key  */
	RawType *string `json:"type"` /*  类型  */
}
