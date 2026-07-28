package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfCreateFunctionV2024Api
/* 创建函数 */
type CfCreateFunctionV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfCreateFunctionV2024Api(client *core.CtyunClient) *CfCreateFunctionV2024Api {
	return &CfCreateFunctionV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/functions",
			ContentType:  "application/json",
		},
	}
}

func (a *CfCreateFunctionV2024Api) Do(ctx context.Context, credential core.Credential, req *CfCreateFunctionV2024Request) (*CfCreateFunctionV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfCreateFunctionV2024Request
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
	var resp CfCreateFunctionV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfCreateFunctionV2024Request struct {
	RegionId              string                                             `json:"regionId"`                        /*  资源池id  */
	Container             *CfCreateFunctionV2024ContainerRequest             `json:"container,omitempty"`             /*  容器配置  */
	Lifecycle             *CfCreateFunctionV2024LifecycleRequest             `json:"lifecycle,omitempty"`             /*  生命周期配置  */
	Log                   *CfCreateFunctionV2024LogRequest                   `json:"log,omitempty"`                   /*  是否启用日志功能  */
	Layers                []*CfCreateFunctionV2024LayersRequest              `json:"layers,omitempty"`                /*  层配置  */
	Runtime               *CfCreateFunctionV2024RuntimeRequest               `json:"runtime,omitempty"`               /*  运行时  */
	Code                  *CfCreateFunctionV2024CodeRequest                  `json:"code,omitempty"`                  /*  代码配置  */
	Network               *CfCreateFunctionV2024NetworkRequest               `json:"network,omitempty"`               /*  网络配置  */
	CustomContainerConfig *CfCreateFunctionV2024CustomContainerConfigRequest `json:"customContainerConfig,omitempty"` /*  自定义镜像配置  */
	ServerlessGpuConfig   *CfCreateFunctionV2024ServerlessGpuConfigRequest   `json:"serverlessGpuConfig,omitempty"`   /*  gpu 函数配置  */
	CreateType            int32                                              `json:"createType"`                      /*  创建函数的类型   1:标准运行时   2:自定义运行时   3:自定义镜像  */
	OssMount              *CfCreateFunctionV2024OssMountRequest              `json:"ossMount,omitempty"`              /*  zos挂载配置  */
	FunctionName          string                                             `json:"functionName"`                    /*  函数名  */
	Description           *string                                            `json:"description,omitempty"`           /*  说明  */
	Role                  *string                                            `json:"role,omitempty"`                  /*  角色  */
	Dns                   *CfCreateFunctionV2024DnsRequest                   `json:"dns,omitempty"`                   /*  DNS配置  */
	Nas                   *CfCreateFunctionV2024NasRequest                   `json:"nas,omitempty"`                   /*  NAS配置  */
}

type CfCreateFunctionV2024ContainerRequest struct {
	TimeZone             string                                                  `json:"timeZone"`                       /*  时区  */
	DiskSize             int32                                                   `json:"diskSize"`                       /*  磁盘规格(Mb)  */
	MaxScale             *int32                                                  `json:"maxScale,omitempty"`             /*  并发实例数上限  */
	FastStart            *int32                                                  `json:"fastStart,omitempty"`            /*  启动加速, 默认为0, 1(表示使用启动加速)  */
	HealthCheckConfig    *CfCreateFunctionV2024ContainerHealthCheckConfigRequest `json:"healthCheckConfig,omitempty"`    /*  健康检查  */
	EnvironmentVariables map[string]string                                       `json:"environmentVariables,omitempty"` /*  环境变量(map<string>)  */
	MemorySize           int32                                                   `json:"memorySize"`                     /*  内存规格(Mb)  */
	Cpu                  float64                                                 `json:"cpu"`                            /*  CPU规格(vCPU)  */
	RunCommand           *string                                                 `json:"runCommand,omitempty"`           /*  函数服务启动命令  */
	ListenPort           int32                                                   `json:"listenPort"`                     /*  监听端口(如使用标准类型函数请传递8080)  */
}

type CfCreateFunctionV2024LifecycleRequest struct {
	Initializer *CfCreateFunctionV2024LifecycleInitializerRequest `json:"initializer,omitempty"` /*  实例初始化  */
	PreStop     *CfCreateFunctionV2024LifecyclePreStopRequest     `json:"preStop,omitempty"`     /*  实例初始化  */
}

type CfCreateFunctionV2024LogRequest struct {
	LogEnabled     *bool                                   `json:"logEnabled,omitempty"`     /*  是否启用日志功能  */
	LogAutoConfig  *bool                                   `json:"logAutoConfig,omitempty"`  /*  是否自动配置  */
	LogProjectId   *string                                 `json:"logProjectId,omitempty"`   /*  日志项目id  */
	LogProjectCode *string                                 `json:"logProjectCode,omitempty"` /*  日志项目Code  */
	LogUnit        *string                                 `json:"logUnit,omitempty"`        /*  日志单元名称  */
	LogUnitId      *string                                 `json:"logUnitId,omitempty"`      /*  日志单元id  */
	LogProject     *string                                 `json:"logProject,omitempty"`     /*  日志项目名称  */
	LogUnitCode    *string                                 `json:"logUnitCode,omitempty"`    /*  日志单元Code  */
	LogRuleEnabled *bool                                   `json:"logRuleEnabled,omitempty"` /*  启用日志分割规则  */
	LogRule        *CfCreateFunctionV2024LogLogRuleRequest `json:"logRule,omitempty"`        /*  日志切割配置  */
}

type CfCreateFunctionV2024LayersRequest struct {
	LayerName *string `json:"layerName,omitempty"` /*  层名称  */
	Version   *int32  `json:"version,omitempty"`   /*  版本  */
	Acl       *int32  `json:"acl,omitempty"`       /*  0表示自定义层，1表示官方公共层  */
}

type CfCreateFunctionV2024RuntimeRequest struct {
	Runtime             *string `json:"runtime,omitempty"`             /*  运行时类型  */
	HandleType          *string `json:"handleType,omitempty"`          /*  请求处理程序类型（标准运行时必填）  */
	ExecuteTimeout      *int32  `json:"executeTimeout,omitempty"`      /*  执行超时时间  */
	Handler             *string `json:"handler,omitempty"`             /*  函数执行的入口 （标准运行时必填）  */
	InstanceConcurrency *int32  `json:"instanceConcurrency,omitempty"` /*  实例最大并发度  */
}

type CfCreateFunctionV2024CodeRequest struct {
	OssBucketName *string `json:"ossBucketName,omitempty"` /*  zos桶名  */
	OssObjectName *string `json:"ossObjectName,omitempty"` /*  zos对象路径名  */
	ZipFile       *string `json:"zipFile,omitempty"`       /*  函数代码 ZIP 包的Base 64编码  */
}

type CfCreateFunctionV2024NetworkRequest struct {
	InternetOutAllow  *bool     `json:"InternetOutAllow,omitempty"`  /*  允许函数访问公网  */
	InternetInForbid  *bool     `json:"internetInForbid,omitempty"`  /*  不允许互联网公网访问函数  */
	OutVpcId          *string   `json:"outVpcId,omitempty"`          /*  networkId  */
	VpcId             *int32    `json:"vpcId,omitempty"`             /*  vpcId  */
	SecurityGroupName *string   `json:"securityGroupName,omitempty"` /*  安全组Name  */
	SubNetId          *string   `json:"subNetId,omitempty"`          /*  子网ID  */
	SubNetName        *string   `json:"subNetName,omitempty"`        /*  子网name  */
	Enable            *bool     `json:"enable,omitempty"`            /*  是否开启VPC  */
	FixedPublicIp     *bool     `json:"fixedPublicIp,omitempty"`     /*  固定公网ip  */
	OutVpcName        *string   `json:"outVpcName,omitempty"`        /*  vpcName  */
	SecurityGroupId   *string   `json:"securityGroupId,omitempty"`   /*  安全组ID  */
	SubNetCidr        *string   `json:"subNetCidr,omitempty"`        /*  子网CIDR  */
	AccessVpcIds      []*string `json:"accessVpcIds,omitempty"`      /*  仅允许指定的VPC访问函数  */
}

type CfCreateFunctionV2024CustomContainerConfigRequest struct {
	ImageDigest *string `json:"imageDigest,omitempty"` /*  digest 用于指定镜像版本  */
	Image       *string `json:"image,omitempty"`       /*  容器镜像地址  */
	InstanceId  *string `json:"instanceId,omitempty"`  /*  crs 实例 id  */
}

type CfCreateFunctionV2024ServerlessGpuConfigRequest struct {
	GpuEnable         *bool   `json:"gpuEnable,omitempty"`         /*  是否使用Gpu  */
	GpuEciType        *string `json:"gpuEciType,omitempty"`        /*  GPU ECI 规格  */
	GpuMemorySize     *int32  `json:"gpuMemorySize,omitempty"`     /*  单位是G  */
	GpuType           *string `json:"gpuType,omitempty"`           /*  gpu卡型  */
	GpuProvisionCount *int32  `json:"gpuProvisionCount,omitempty"` /*  配置的预留实例数量  */
}

type CfCreateFunctionV2024OssMountRequest struct {
	Mounts []*CfCreateFunctionV2024OssMountMountsRequest `json:"mounts,omitempty"` /*  zos挂载参数  */
}

type CfCreateFunctionV2024DnsRequest struct {
	NameServers []*string                               `json:"nameServers,omitempty"` /*  DNS 服务器的 IP 地址列表  */
	Searches    []*string                               `json:"searches,omitempty"`    /*  DNS 搜索域列表  */
	Options     *CfCreateFunctionV2024DnsOptionsRequest `json:"options,omitempty"`     /*  DNS 解析配置  */
}

type CfCreateFunctionV2024NasRequest struct {
	Nas []*CfCreateFunctionV2024NasNasRequest `json:"nas,omitempty"` /*  nas配置  */
}

type CfCreateFunctionV2024ContainerHealthCheckConfigRequest struct {
	FailureThreshold    int32  `json:"failureThreshold"`    /*  失败阈值  */
	GetPath             string `json:"getPath"`             /*  检查http get path  */
	InitialDelaySeconds int32  `json:"initialDelaySeconds"` /*  首次探测延迟时间(秒)  */
	PeriodSeconds       int32  `json:"periodSeconds"`       /*  探测时间间隔(秒)  */
	SuccessThreshold    int32  `json:"successThreshold"`    /*  成功阈值  */
	TimeoutSeconds      int32  `json:"timeoutSeconds"`      /*  超时(秒)  */
}

type CfCreateFunctionV2024LifecycleInitializerRequest struct {
	Handler string `json:"handler"` /*  处理方法入口  */
	Enable  bool   `json:"enable"`  /*  启用  */
	Timeout int32  `json:"timeout"` /*  超时  */
}

type CfCreateFunctionV2024LifecyclePreStopRequest struct {
	Handler string `json:"handler"` /*  处理方法入口  */
	Enable  bool   `json:"enable"`  /*  启用  */
	Timeout int32  `json:"timeout"` /*  超时  */
}

type CfCreateFunctionV2024LogLogRuleRequest struct {
	RuleCode         *string                                           `json:"ruleCode,omitempty"`         /*  规则唯一编码  */
	RuleName         *string                                           `json:"ruleName,omitempty"`         /*  规则名称  */
	ExtractMode      *int32                                            `json:"extractMode,omitempty"`      /*  采集类型  */
	CollectPolicy    *string                                           `json:"collectPolicy,omitempty"`    /*  采集策略  */
	CuttingMode      *string                                           `json:"cuttingMode,omitempty"`      /*  切割模式  */
	Enable           *bool                                             `json:"enable,omitempty"`           /*  是否启用采集规则  */
	UnitCode         *string                                           `json:"unitCode,omitempty"`         /*  日志单元编码  */
	FirstLinePattern *string                                           `json:"firstLinePattern,omitempty"` /*  首行正则  */
	CustomTime       *CfCreateFunctionV2024LogLogRuleCustomTimeRequest `json:"customTime,omitempty"`       /*  自定义时间戳提取格式  */
	RuleConfig       *CfCreateFunctionV2024LogLogRuleRuleConfigRequest `json:"ruleConfig,omitempty"`       /*  容器运行参数  */
	AccessType       *int32                                            `json:"accessType,omitempty"`       /*  接入类型  */
}

type CfCreateFunctionV2024OssMountMountsRequest struct {
	BucketName *string `json:"bucketName,omitempty"` /*  bucket名  */
	BucketPath *string `json:"bucketPath,omitempty"` /*  bucket子目录  */
	MountDir   *string `json:"mountDir,omitempty"`   /*  挂载本地目录  */
	ReadOnly   *bool   `json:"readOnly,omitempty"`   /*  是否只读，默认false  */
}

type CfCreateFunctionV2024DnsOptionsRequest struct {
	Ndots *string `json:"ndots,omitempty"` /*  键值对01  */
}

type CfCreateFunctionV2024NasNasRequest struct {
	RemoteDir *string `json:"remoteDir,omitempty"` /*  远端挂载目录  */
	SharePath *string `json:"sharePath,omitempty"` /*  挂载地址  */
	LocalDir  *string `json:"localDir,omitempty"`  /*  挂载本地目录  */
	SfsName   *string `json:"sfsName,omitempty"`   /*  sfs的名称  */
	SfsUID    *string `json:"sfsUID,omitempty"`    /*  sfs的ID  */
}

type CfCreateFunctionV2024LogLogRuleCustomTimeRequest struct {
	Key        *string `json:"key,omitempty"`        /*  key  */
	TimeFormat *string `json:"timeFormat,omitempty"` /*  格式化  */
}

type CfCreateFunctionV2024LogLogRuleRuleConfigRequest struct {
	MaxPathDepth *int32                                                     `json:"maxPathDepth,omitempty"` /*  最大正则路径解析深度  */
	Delimeter    *CfCreateFunctionV2024LogLogRuleRuleConfigDelimeterRequest `json:"delimeter,omitempty"`    /*  分隔符  */
	Regex        *CfCreateFunctionV2024LogLogRuleRuleConfigRegexRequest     `json:"regex,omitempty"`        /*  正则切割模式  */
}

type CfCreateFunctionV2024LogLogRuleRuleConfigDelimeterRequest struct {
	Delimeter    *string                                                                  `json:"delimeter,omitempty"`    /*  分隔符  */
	TypeContents []*CfCreateFunctionV2024LogLogRuleRuleConfigDelimeterTypeContentsRequest `json:"typeContents,omitempty"` /*  类型  */
}

type CfCreateFunctionV2024LogLogRuleRuleConfigRegexRequest struct {
	RegexStr     *string                                                              `json:"regexStr,omitempty"`     /*  正则表达式  */
	TypeContents []*CfCreateFunctionV2024LogLogRuleRuleConfigRegexTypeContentsRequest `json:"typeContents,omitempty"` /*  类型  */
}

type CfCreateFunctionV2024LogLogRuleRuleConfigDelimeterTypeContentsRequest struct {
	Key     *string `json:"key,omitempty"`  /*  key  */
	RawType *string `json:"type,omitempty"` /*  类型  */
}

type CfCreateFunctionV2024LogLogRuleRuleConfigRegexTypeContentsRequest struct {
	Key     *string `json:"key,omitempty"`  /*  key  */
	RawType *string `json:"type,omitempty"` /*  类型  */
}

type CfCreateFunctionV2024Response struct {
	StatusCode *int32                                  `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                                 `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                 `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfCreateFunctionV2024ReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfCreateFunctionV2024ReturnObjResponse struct {
	FunctionId            *string                                                      `json:"functionId"`            /*  函数 id  */
	DeployInfo            *CfCreateFunctionV2024ReturnObjDeployInfoResponse            `json:"deployInfo"`            /*  最近一次函数构建的状态信息（可以为空）  */
	Container             *CfCreateFunctionV2024ReturnObjContainerResponse             `json:"container"`             /*  容器配置  */
	Lifecycle             *CfCreateFunctionV2024ReturnObjLifecycleResponse             `json:"lifecycle"`             /*  生命周期配置  */
	Log                   *CfCreateFunctionV2024ReturnObjLogResponse                   `json:"log"`                   /*  是否启用日志功能  */
	Layers                []*CfCreateFunctionV2024ReturnObjLayersResponse              `json:"layers"`                /*  层配置  */
	Runtime               *CfCreateFunctionV2024ReturnObjRuntimeResponse               `json:"runtime"`               /*  运行时  */
	Network               *CfCreateFunctionV2024ReturnObjNetworkResponse               `json:"network"`               /*  网络配置  */
	CustomContainerConfig *CfCreateFunctionV2024ReturnObjCustomContainerConfigResponse `json:"customContainerConfig"` /*  自定义镜像配置  */
	ServerlessGpuConfig   *CfCreateFunctionV2024ReturnObjServerlessGpuConfigResponse   `json:"serverlessGpuConfig"`   /*  gpu函数配置  */
	CreateType            *int32                                                       `json:"createType"`            /*  创建类型 1:内置运行时2:自定义运行时3:自定义镜像  */
	OssMount              *CfCreateFunctionV2024ReturnObjOssMountResponse              `json:"ossMount"`              /*  zos挂载配置  */
	FunctionName          *string                                                      `json:"functionName"`          /*  函数名  */
	Description           *string                                                      `json:"description"`           /*  说明  */
	Role                  *string                                                      `json:"role"`                  /*  角色  */
	Dns                   *CfCreateFunctionV2024ReturnObjDnsResponse                   `json:"dns"`                   /*  DNS配置  */
	Nas                   *CfCreateFunctionV2024ReturnObjNasResponse                   `json:"nas"`                   /*  NAS配置  */
}

type CfCreateFunctionV2024ReturnObjDeployInfoResponse struct {
	TaskEnd   *int32  `json:"taskEnd"`   /*  结束时间  */
	ErrMsg    *string `json:"errMsg"`    /*  错误说明  */
	Creator   *string `json:"creator"`   /*  关联的userId  */
	Id        *string `json:"id"`        /*  此次更新内部唯一id  */
	Status    *string `json:"status"`    /*  当前代码部署状态  */
	TaskBegin *int32  `json:"taskBegin"` /*  开始时间  */
}

type CfCreateFunctionV2024ReturnObjContainerResponse struct {
	TimeZone             *string                                                           `json:"timeZone"`             /*  时区  */
	DiskSize             *int32                                                            `json:"diskSize"`             /*  磁盘规格(Mb)  */
	MaxScale             *int32                                                            `json:"maxScale"`             /*  并发实例数上限  */
	FastStart            *int32                                                            `json:"fastStart"`            /*  启动加速, 默认为0, 1(表示使用启动加速)  */
	HealthCheckConfig    *CfCreateFunctionV2024ReturnObjContainerHealthCheckConfigResponse `json:"healthCheckConfig"`    /*  健康检查  */
	EnvironmentVariables map[string]string                                                 `json:"environmentVariables"` /*  环境变量  */
	MemorySize           *int32                                                            `json:"memorySize"`           /*  内存规格(Mb)  */
	Cpu                  *float64                                                          `json:"cpu"`                  /*  CPU规格(vCPU)  */
	RunCommand           *string                                                           `json:"runCommand"`           /*  函数服务启动命令  */
	ListenPort           *int32                                                            `json:"listenPort"`           /*  监听端口  */
	Image                *string                                                           `json:"image"`                /*  基础镜像地址  */
}

type CfCreateFunctionV2024ReturnObjLifecycleResponse struct {
	Initializer *CfCreateFunctionV2024ReturnObjLifecycleInitializerResponse `json:"initializer"` /*  实例初始化  */
	PreStop     *CfCreateFunctionV2024ReturnObjLifecyclePreStopResponse     `json:"preStop"`     /*  实例初始化  */
}

type CfCreateFunctionV2024ReturnObjLogResponse struct {
	LogEnabled     *bool                                             `json:"logEnabled"`     /*  是否启用日志功能  */
	LogAutoConfig  *bool                                             `json:"logAutoConfig"`  /*  是否自动配置  */
	LogProjectId   *string                                           `json:"logProjectId"`   /*  日志项目id  */
	LogProjectCode *string                                           `json:"logProjectCode"` /*  日志项目Code  */
	LogUnit        *string                                           `json:"logUnit"`        /*  日志单元名称  */
	LogUnitId      *string                                           `json:"logUnitId"`      /*  日志单元id  */
	LogProject     *string                                           `json:"logProject"`     /*  日志项目名称  */
	LogUnitCode    *string                                           `json:"logUnitCode"`    /*  日志单元Code  */
	LogRuleEnabled *bool                                             `json:"logRuleEnabled"` /*  启用日志分割规则  */
	LogRule        *CfCreateFunctionV2024ReturnObjLogLogRuleResponse `json:"logRule"`        /*  日志切割配置  */
}

type CfCreateFunctionV2024ReturnObjLayersResponse struct {
	LayerName   *string `json:"layerName"`   /*  层名称  */
	Version     *int32  `json:"version"`     /*  版本  */
	Description *string `json:"description"` /*  描述  */
	Acl         *int32  `json:"acl"`         /*  0表示自定义层，1表示官方公共层  */
}

type CfCreateFunctionV2024ReturnObjRuntimeResponse struct {
	Runtime             *string `json:"runtime"`             /*  运行时类型  */
	HandleType          *string `json:"handleType"`          /*  请求处理程序类型  */
	ExecuteTimeout      *int32  `json:"executeTimeout"`      /*  执行超时时间  */
	Handler             *string `json:"handler"`             /*  函数执行的入口  */
	InstanceConcurrency *int32  `json:"instanceConcurrency"` /*  实例最大并发度  */
}

type CfCreateFunctionV2024ReturnObjNetworkResponse struct {
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

type CfCreateFunctionV2024ReturnObjCustomContainerConfigResponse struct {
	ImageDigest *string `json:"imageDigest"` /*  digest 用于指定镜像版本  */
	Image       *string `json:"image"`       /*  容器镜像地址  */
	InstanceId  *string `json:"instanceId"`  /*  crs 实例 id  */
}

type CfCreateFunctionV2024ReturnObjServerlessGpuConfigResponse struct {
	GpuEnable         *bool   `json:"gpuEnable"`         /*  是否使用Gpu  */
	GpuEciType        *string `json:"gpuEciType"`        /*  GPU ECI 规格  */
	GpuMemorySize     *int32  `json:"gpuMemorySize"`     /*  单位是G  */
	GpuType           *string `json:"gpuType"`           /*  gpu卡型  */
	GpuProvisionCount *int32  `json:"gpuProvisionCount"` /*  配置的预留实例数量  */
}

type CfCreateFunctionV2024ReturnObjOssMountResponse struct {
	Mounts []*CfCreateFunctionV2024ReturnObjOssMountMountsResponse `json:"mounts"` /*  zos挂载参数  */
}

type CfCreateFunctionV2024ReturnObjDnsResponse struct {
	NameServers []*string                                         `json:"nameServers"` /*  DNS 服务器的 IP 地址列表  */
	Searches    []*string                                         `json:"searches"`    /*  DNS 搜索域列表  */
	Options     *CfCreateFunctionV2024ReturnObjDnsOptionsResponse `json:"options"`     /*  DNS 解析配置  */
}

type CfCreateFunctionV2024ReturnObjNasResponse struct {
	Nas []*CfCreateFunctionV2024ReturnObjNasNasResponse `json:"nas"` /*  nas  */
}

type CfCreateFunctionV2024ReturnObjContainerHealthCheckConfigResponse struct {
	FailureThreshold    *int32  `json:"failureThreshold"`    /*  失败阈值  */
	GetPath             *string `json:"getPath"`             /*  检查http get path  */
	InitialDelaySeconds *int32  `json:"initialDelaySeconds"` /*  首次探测延迟时间(秒)  */
	PeriodSeconds       *int32  `json:"periodSeconds"`       /*  探测时间间隔(秒)  */
	SuccessThreshold    *int32  `json:"successThreshold"`    /*  成功阈值  */
	TimeoutSeconds      *int32  `json:"timeoutSeconds"`      /*  超时(秒)  */
}

type CfCreateFunctionV2024ReturnObjLifecycleInitializerResponse struct {
	Handler *string `json:"handler"` /*  处理方法入口  */
	Enable  *bool   `json:"enable"`  /*  启用  */
	Timeout *int32  `json:"timeout"` /*  超时  */
}

type CfCreateFunctionV2024ReturnObjLifecyclePreStopResponse struct {
	Handler *string `json:"handler"` /*  处理方法入口  */
	Enable  *bool   `json:"enable"`  /*  启用  */
	Timeout *int32  `json:"timeout"` /*  超时  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleResponse struct {
	RuleCode         *string                                                     `json:"ruleCode"`         /*  规则唯一编码  */
	RuleName         *string                                                     `json:"ruleName"`         /*  规则名称  */
	ExtractMode      *int32                                                      `json:"extractMode"`      /*  采集类型  */
	CollectPolicy    *string                                                     `json:"collectPolicy"`    /*  采集策略  */
	CuttingMode      *string                                                     `json:"cuttingMode"`      /*  切割模式  */
	Enable           *bool                                                       `json:"enable"`           /*  是否启用采集规则  */
	UnitCode         *string                                                     `json:"unitCode"`         /*  日志单元编码  */
	FirstLinePattern *string                                                     `json:"firstLinePattern"` /*  首行正则  */
	CustomTime       *CfCreateFunctionV2024ReturnObjLogLogRuleCustomTimeResponse `json:"customTime"`       /*  自定义时间戳提取格式  */
	RuleConfig       *CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigResponse `json:"ruleConfig"`       /*  容器运行参数  */
	AccessType       *int32                                                      `json:"accessType"`       /*  接入类型  */
}

type CfCreateFunctionV2024ReturnObjOssMountMountsResponse struct {
	BucketName *string `json:"bucketName"` /*  bucket名  */
	BucketPath *string `json:"bucketPath"` /*  bucket子目录  */
	MountDir   *string `json:"mountDir"`   /*  挂载本地目录  */
	ReadOnly   *bool   `json:"readOnly"`   /*  是否只读，默认false  */
	AccessUrl  *string `json:"accessUrl"`  /*  oss 访问地址  */
}

type CfCreateFunctionV2024ReturnObjDnsOptionsResponse struct {
	Ndots *string `json:"ndots"` /*  键值对01  */
}

type CfCreateFunctionV2024ReturnObjNasNasResponse struct {
	RemoteDir *string `json:"remoteDir"` /*  远端挂载目录  */
	SharePath *string `json:"sharePath"` /*  挂载地址  */
	LocalDir  *string `json:"localDir"`  /*  挂载本地目录  */
	SfsName   *string `json:"sfsName"`   /*  sfs 的名称  */
	SfsUID    *string `json:"sfsUID"`    /*  sfs 的 ID  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleCustomTimeResponse struct {
	Key        *string `json:"key"`        /*  key  */
	TimeFormat *string `json:"timeFormat"` /*  格式化  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigResponse struct {
	MaxPathDepth *int32                                                               `json:"maxPathDepth"` /*  最大正则路径解析深度  */
	Delimeter    *CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigDelimeterResponse `json:"delimeter"`    /*  分隔符  */
	Regex        *CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigRegexResponse     `json:"regex"`        /*  正则切割模式  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigDelimeterResponse struct {
	Delimeter    *string                                                                            `json:"delimeter"`    /*  分隔符  */
	TypeContents []*CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigDelimeterTypeContentsResponse `json:"typeContents"` /*  类型  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigRegexResponse struct {
	RegexStr     *string                                                                        `json:"regexStr"`     /*  正则表达式  */
	TypeContents []*CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigRegexTypeContentsResponse `json:"typeContents"` /*  类型  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigDelimeterTypeContentsResponse struct {
	Key     *string `json:"key"`  /*  key  */
	RawType *string `json:"type"` /*  类型  */
}

type CfCreateFunctionV2024ReturnObjLogLogRuleRuleConfigRegexTypeContentsResponse struct {
	Key     *string `json:"key"`  /*  key  */
	RawType *string `json:"type"` /*  类型  */
}
