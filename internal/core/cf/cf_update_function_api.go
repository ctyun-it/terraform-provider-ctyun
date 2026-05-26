package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfUpdateFunctionApi
/* 修改函数 */
type CfUpdateFunctionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfUpdateFunctionApi(client *core.CtyunClient) *CfUpdateFunctionApi {
	return &CfUpdateFunctionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/functions/{functionName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfUpdateFunctionApi) Do(ctx context.Context, credential core.Credential, req *CfUpdateFunctionRequest) (*CfUpdateFunctionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("updateType", req.UpdateType)
	_, err := ctReq.WriteJson(struct {
		*CfUpdateFunctionRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		UpdateType   interface{} `json:"updateType,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
	}{
		req, nil, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfUpdateFunctionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfUpdateFunctionRequest struct {
	FunctionName          string                                        `json:"functionName"`                    /*  函数名  */
	RegionId              string                                        `json:"regionId"`                        /*  资源池id  */
	UpdateType            string                                        `json:"updateType"`                      /*  更新配置类型，枚举值:基础配置:basic；代码:code；运行时:runtime；层:layer；权限:role；日志:log；网络:network；存储NAS:nas；存储OSS:oss；生命周期:lifecycle；DNS配置:dns；异步配置:asyncConfig；异步成功目标:asyncOk；异步失败目标:asyncFail 健康检查:health; 全部或多项配置  */
	Container             *CfUpdateFunctionContainerRequest             `json:"container,omitempty"`             /*  容器配置  */
	Lifecycle             *CfUpdateFunctionLifecycleRequest             `json:"lifecycle,omitempty"`             /*  生命周期配置  */
	Log                   *CfUpdateFunctionLogRequest                   `json:"log,omitempty"`                   /*  是否启用日志功能  */
	Layers                []*CfUpdateFunctionLayersRequest              `json:"layers,omitempty"`                /*  层配置  */
	Runtime               *CfUpdateFunctionRuntimeRequest               `json:"runtime,omitempty"`               /*  运行时  */
	Code                  *CfUpdateFunctionCodeRequest                  `json:"code,omitempty"`                  /*  代码配置  */
	Network               *CfUpdateFunctionNetworkRequest               `json:"network,omitempty"`               /*  网络配置  */
	CustomContainerConfig *CfUpdateFunctionCustomContainerConfigRequest `json:"customContainerConfig,omitempty"` /*  镜像仓库配置  */
	ServerlessGpuConfig   *CfUpdateFunctionServerlessGpuConfigRequest   `json:"serverlessGpuConfig,omitempty"`   /*  gpu 函数配置  */
	CreateType            *int32                                        `json:"createType,omitempty"`            /*  创建类型 1:内置运行时2:自定义运行时3:自定义镜像  */
	OssMount              *CfUpdateFunctionOssMountRequest              `json:"ossMount,omitempty"`              /*  zos挂载配置  */
	Description           *string                                       `json:"description,omitempty"`           /*  说明  */
	Role                  *string                                       `json:"role,omitempty"`                  /*  角色  */
	Dns                   *CfUpdateFunctionDnsRequest                   `json:"dns,omitempty"`                   /*  DNS配置  */
	Nas                   *CfUpdateFunctionNasRequest                   `json:"nas,omitempty"`                   /*  NAS配置  */
}

type CfUpdateFunctionContainerRequest struct {
	TimeZone             *string                                            `json:"timeZone,omitempty"`             /*  时区  */
	DiskSize             *int32                                             `json:"diskSize,omitempty"`             /*  磁盘规格(Mb)  */
	MaxScale             *int32                                             `json:"maxScale,omitempty"`             /*  并发实例数上限  */
	FastStart            *int32                                             `json:"fastStart,omitempty"`            /*  启动加速, 默认为0, 1(表示使用启动加速)  */
	HealthCheckConfig    *CfUpdateFunctionContainerHealthCheckConfigRequest `json:"healthCheckConfig,omitempty"`    /*  健康检查  */
	EnvironmentVariables map[string]string                                  `json:"environmentVariables,omitempty"` /*  环境变量(map<string>)  */
	MemorySize           *int32                                             `json:"memorySize,omitempty"`           /*  内存规格(Mb)  */
	Cpu                  *float64                                           `json:"cpu,omitempty"`                  /*  CPU规格(vCPU)  */
	RunCommand           *string                                            `json:"runCommand,omitempty"`           /*  函数服务启动命令  */
	ListenPort           int32                                              `json:"listenPort,,omitempty"`          /*  监听端口(如使用标准类型函数请传递8080)  */
	Image                *string                                            `json:"image,omitempty"`                /*  基础镜像地址(非自定义镜像必填)  */
}

type CfUpdateFunctionLifecycleRequest struct {
	Initializer *CfUpdateFunctionLifecycleInitializerRequest `json:"initializer,omitempty"` /*  实例初始化  */
	PreStop     *CfUpdateFunctionLifecyclePreStopRequest     `json:"preStop,omitempty"`     /*  实例初始化  */
}

type CfUpdateFunctionLogRequest struct {
	LogEnabled     *bool                              `json:"logEnabled,omitempty"`     /*  是否启用日志功能  */
	LogAutoConfig  *bool                              `json:"logAutoConfig,omitempty"`  /*  是否自动配置  */
	LogProjectId   *string                            `json:"logProjectId,omitempty"`   /*  日志项目id  */
	LogProjectCode *string                            `json:"logProjectCode,omitempty"` /*  日志项目Code  */
	LogUnit        *string                            `json:"logUnit,omitempty"`        /*  日志单元名称  */
	LogUnitId      *string                            `json:"logUnitId,omitempty"`      /*  日志单元id  */
	LogProject     *string                            `json:"logProject,omitempty"`     /*  日志项目名称  */
	LogUnitCode    *string                            `json:"logUnitCode,omitempty"`    /*  日志单元Code  */
	LogRuleEnabled *bool                              `json:"logRuleEnabled,omitempty"` /*  启用日志分割规则  */
	LogRule        *CfUpdateFunctionLogLogRuleRequest `json:"logRule,omitempty"`        /*  日志切割配置  */
}

type CfUpdateFunctionLayersRequest struct {
	LayerName *string `json:"layerName,omitempty"` /*  层名称  */
	Version   *int32  `json:"version,omitempty"`   /*  版本  */
	Acl       *int32  `json:"acl,omitempty"`       /*  0表示自定义层，1表示官方公共层  */
}

type CfUpdateFunctionRuntimeRequest struct {
	Runtime             *string `json:"runtime,omitempty"`             /*  运行时类型  */
	HandleType          *string `json:"handleType,omitempty"`          /*  请求处理程序类型（标准运行时必填）  */
	ExecuteTimeout      *int32  `json:"executeTimeout,omitempty"`      /*  执行超时时间  */
	Handler             *string `json:"handler,omitempty"`             /*  函数执行的入口 （标准运行时必填）  */
	InstanceConcurrency *int32  `json:"instanceConcurrency,omitempty"` /*  实例最大并发度  */
}

type CfUpdateFunctionCodeRequest struct {
	OssBucketName *string `json:"ossBucketName,omitempty"` /*  zos桶名  */
	OssObjectName *string `json:"ossObjectName,omitempty"` /*  zos对象路径名  */
	ZipFile       *string `json:"zipFile,omitempty"`       /*  函数代码 ZIP 包的Base 64编码  */
}

type CfUpdateFunctionNetworkRequest struct {
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

type CfUpdateFunctionCustomContainerConfigRequest struct {
	ImageDigest *string `json:"imageDigest,omitempty"` /*  digest 用于指定镜像版本  */
	Image       *string `json:"image,omitempty"`       /*  容器镜像地址  */
	InstanceId  *string `json:"instanceId,omitempty"`  /*  crs 实例 id  */
}

type CfUpdateFunctionServerlessGpuConfigRequest struct {
	GpuEnable         *bool   `json:"gpuEnable,omitempty"`         /*  是否使用Gpu  */
	GpuEciType        *string `json:"gpuEciType,omitempty"`        /*  GPU ECI 规格  */
	GpuMemorySize     *int32  `json:"gpuMemorySize,omitempty"`     /*  单位是G  */
	GpuType           *string `json:"gpuType,omitempty"`           /*  gpu卡型  */
	GpuProvisionCount *int32  `json:"gpuProvisionCount,omitempty"` /*  配置的预留实例数量  */
}

type CfUpdateFunctionOssMountRequest struct {
	Mounts []*CfUpdateFunctionOssMountMountsRequest `json:"mounts,omitempty"` /*  zos挂载参数  */
}

type CfUpdateFunctionDnsRequest struct {
	NameServers []*string                          `json:"nameServers,omitempty"` /*  DNS 服务器的 IP 地址列表  */
	Searches    []*string                          `json:"searches,omitempty"`    /*  DNS 搜索域列表  */
	Options     *CfUpdateFunctionDnsOptionsRequest `json:"options,omitempty"`     /*  DNS 解析配置  */
}

type CfUpdateFunctionNasRequest struct {
	Nas []*CfUpdateFunctionNasNasRequest `json:"nas,omitempty"` /*  nas配置  */
}

type CfUpdateFunctionContainerHealthCheckConfigRequest struct {
	FailureThreshold    int32  `json:"failureThreshold"`    /*  失败阈值  */
	GetPath             string `json:"getPath"`             /*  检查http get path  */
	InitialDelaySeconds int32  `json:"initialDelaySeconds"` /*  首次探测延迟时间(秒)  */
	PeriodSeconds       int32  `json:"periodSeconds"`       /*  探测时间间隔(秒)  */
	SuccessThreshold    int32  `json:"successThreshold"`    /*  成功阈值  */
	TimeoutSeconds      int32  `json:"timeoutSeconds"`      /*  超时(秒)  */
}

type CfUpdateFunctionLifecycleInitializerRequest struct {
	Handler string `json:"handler"` /*  处理方法入口  */
	Enable  bool   `json:"enable"`  /*  启用  */
	Timeout int32  `json:"timeout"` /*  超时  */
}

type CfUpdateFunctionLifecyclePreStopRequest struct {
	Handler string `json:"handler"` /*  处理方法入口  */
	Enable  bool   `json:"enable"`  /*  启用  */
	Timeout int32  `json:"timeout"` /*  超时  */
}

type CfUpdateFunctionLogLogRuleRequest struct {
	RuleCode         *string                                      `json:"ruleCode,omitempty"`         /*  规则唯一编码  */
	RuleName         *string                                      `json:"ruleName,omitempty"`         /*  规则名称  */
	ExtractMode      *int32                                       `json:"extractMode,omitempty"`      /*  采集类型  */
	CollectPolicy    *string                                      `json:"collectPolicy,omitempty"`    /*  采集策略  */
	CuttingMode      *string                                      `json:"cuttingMode,omitempty"`      /*  切割模式  */
	Enable           *bool                                        `json:"enable,omitempty"`           /*  是否启用采集规则  */
	UnitCode         *string                                      `json:"unitCode,omitempty"`         /*  日志单元编码  */
	FirstLinePattern *string                                      `json:"firstLinePattern,omitempty"` /*  首行正则  */
	CustomTime       *CfUpdateFunctionLogLogRuleCustomTimeRequest `json:"customTime,omitempty"`       /*  自定义时间戳提取格式  */
	RuleConfig       *CfUpdateFunctionLogLogRuleRuleConfigRequest `json:"ruleConfig,omitempty"`       /*  容器运行参数  */
	AccessType       *int32                                       `json:"accessType,omitempty"`       /*  接入类型  */
}

type CfUpdateFunctionOssMountMountsRequest struct {
	BucketName *string `json:"bucketName,omitempty"` /*  bucket名  */
	BucketPath *string `json:"bucketPath,omitempty"` /*  bucket子目录  */
	MountDir   *string `json:"mountDir,omitempty"`   /*  挂载本地目录  */
	ReadOnly   *bool   `json:"readOnly,omitempty"`   /*  是否只读，默认false  */
}

type CfUpdateFunctionDnsOptionsRequest struct {
	Ndots *string `json:"ndots,omitempty"` /*  键值对01  */
}

type CfUpdateFunctionNasNasRequest struct {
	RemoteDir *string `json:"remoteDir,omitempty"` /*  远端挂载目录  */
	SharePath *string `json:"sharePath,omitempty"` /*  挂载地址  */
	LocalDir  *string `json:"localDir,omitempty"`  /*  挂载本地目录  */
	SfsName   *string `json:"sfsName,omitempty"`   /*  sfs的名称  */
	SfsUID    *string `json:"sfsUID,omitempty"`    /*  sfs的ID  */
}

type CfUpdateFunctionLogLogRuleCustomTimeRequest struct {
	Key        *string `json:"key,omitempty"`        /*  key  */
	TimeFormat *string `json:"timeFormat,omitempty"` /*  格式化  */
}

type CfUpdateFunctionLogLogRuleRuleConfigRequest struct {
	MaxPathDepth *int32                                                `json:"maxPathDepth,omitempty"` /*  最大正则路径解析深度  */
	Delimeter    *CfUpdateFunctionLogLogRuleRuleConfigDelimeterRequest `json:"delimeter,omitempty"`    /*  分隔符  */
	Regex        *CfUpdateFunctionLogLogRuleRuleConfigRegexRequest     `json:"regex,omitempty"`        /*  正则切割模式  */
}

type CfUpdateFunctionLogLogRuleRuleConfigDelimeterRequest struct {
	Delimeter    *string                                                             `json:"delimeter,omitempty"`    /*  分隔符  */
	TypeContents []*CfUpdateFunctionLogLogRuleRuleConfigDelimeterTypeContentsRequest `json:"typeContents,omitempty"` /*  类型  */
}

type CfUpdateFunctionLogLogRuleRuleConfigRegexRequest struct {
	RegexStr     *string                                                         `json:"regexStr,omitempty"`     /*  正则表达式  */
	TypeContents []*CfUpdateFunctionLogLogRuleRuleConfigRegexTypeContentsRequest `json:"typeContents,omitempty"` /*  类型  */
}

type CfUpdateFunctionLogLogRuleRuleConfigDelimeterTypeContentsRequest struct {
	Key     *string `json:"key,omitempty"`  /*  key  */
	RawType *string `json:"type,omitempty"` /*  类型  */
}

type CfUpdateFunctionLogLogRuleRuleConfigRegexTypeContentsRequest struct {
	Key     *string `json:"key,omitempty"`  /*  key  */
	RawType *string `json:"type,omitempty"` /*  类型  */
}

type CfUpdateFunctionResponse struct {
	StatusCode *int32                             `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                            `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                            `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfUpdateFunctionReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfUpdateFunctionReturnObjResponse struct {
	FunctionId            *string                                                 `json:"functionId"`            /*  函数 id  */
	DeployInfo            *CfUpdateFunctionReturnObjDeployInfoResponse            `json:"deployInfo"`            /*  最近一次函数构建的状态信息  */
	Container             *CfUpdateFunctionReturnObjContainerResponse             `json:"container"`             /*  容器配置  */
	Lifecycle             *CfUpdateFunctionReturnObjLifecycleResponse             `json:"lifecycle"`             /*  生命周期配置  */
	Log                   *CfUpdateFunctionReturnObjLogResponse                   `json:"log"`                   /*  是否启用日志功能  */
	Layers                []*CfUpdateFunctionReturnObjLayersResponse              `json:"layers"`                /*  层配置  */
	Runtime               *CfUpdateFunctionReturnObjRuntimeResponse               `json:"runtime"`               /*  运行时  */
	Network               *CfUpdateFunctionReturnObjNetworkResponse               `json:"network"`               /*  网络配置  */
	CustomContainerConfig *CfUpdateFunctionReturnObjCustomContainerConfigResponse `json:"customContainerConfig"` /*  自定义镜像配置  */
	ServerlessGpuConfig   *CfUpdateFunctionReturnObjServerlessGpuConfigResponse   `json:"serverlessGpuConfig"`   /*  gpu函数配置  */
	CreateType            *int32                                                  `json:"createType"`            /*  创建类型 1:内置运行时2:自定义运行时3:自定义镜像  */
	OssMount              *CfUpdateFunctionReturnObjOssMountResponse              `json:"ossMount"`              /*  zos挂载配置  */
	FunctionName          *string                                                 `json:"functionName"`          /*  函数名  */
	Description           *string                                                 `json:"description"`           /*  说明  */
	Role                  *string                                                 `json:"role"`                  /*  角色  */
	Dns                   *CfUpdateFunctionReturnObjDnsResponse                   `json:"dns"`                   /*  DNS配置  */
	Nas                   *CfUpdateFunctionReturnObjNasResponse                   `json:"nas"`                   /*  NAS配置  */
}

type CfUpdateFunctionReturnObjDeployInfoResponse struct {
	TaskEnd   *int32  `json:"taskEnd"`   /*  结束时间  */
	ErrMsg    *string `json:"errMsg"`    /*  错误说明  */
	Creator   *string `json:"creator"`   /*  关联的userId  */
	Id        *string `json:"id"`        /*  此次更新内部唯一id  */
	Status    *string `json:"status"`    /*  当前代码部署状态  */
	TaskBegin *int32  `json:"taskBegin"` /*  开始时间  */
}

type CfUpdateFunctionReturnObjContainerResponse struct {
	TimeZone             *string                                                      `json:"timeZone"`             /*  时区  */
	DiskSize             *int32                                                       `json:"diskSize"`             /*  磁盘规格(Mb)  */
	MaxScale             *int32                                                       `json:"maxScale"`             /*  并发实例数上限  */
	FastStart            *int32                                                       `json:"fastStart"`            /*  启动加速, 默认为0, 1(表示使用启动加速)  */
	HealthCheckConfig    *CfUpdateFunctionReturnObjContainerHealthCheckConfigResponse `json:"healthCheckConfig"`    /*  健康检查  */
	EnvironmentVariables map[string]string                                            `json:"environmentVariables"` /*  环境变量  */
	MemorySize           *int32                                                       `json:"memorySize"`           /*  内存规格(Mb)  */
	Cpu                  *float64                                                     `json:"cpu"`                  /*  CPU规格(vCPU)  */
	RunCommand           *string                                                      `json:"runCommand"`           /*  函数服务启动命令  */
	ListenPort           *int32                                                       `json:"listenPort"`           /*  监听端口  */
	Image                *string                                                      `json:"image"`                /*  基础镜像地址  */
}

type CfUpdateFunctionReturnObjLifecycleResponse struct {
	Initializer *CfUpdateFunctionReturnObjLifecycleInitializerResponse `json:"initializer"` /*  实例初始化  */
	PreStop     *CfUpdateFunctionReturnObjLifecyclePreStopResponse     `json:"preStop"`     /*  实例初始化  */
}

type CfUpdateFunctionReturnObjLogResponse struct {
	LogEnabled     *bool                                        `json:"logEnabled"`     /*  是否启用日志功能  */
	LogAutoConfig  *bool                                        `json:"logAutoConfig"`  /*  是否自动配置  */
	LogProjectId   *string                                      `json:"logProjectId"`   /*  日志项目id  */
	LogProjectCode *string                                      `json:"logProjectCode"` /*  日志项目Code  */
	LogUnit        *string                                      `json:"logUnit"`        /*  日志单元名称  */
	LogUnitId      *string                                      `json:"logUnitId"`      /*  日志单元id  */
	LogProject     *string                                      `json:"logProject"`     /*  日志项目名称  */
	LogUnitCode    *string                                      `json:"logUnitCode"`    /*  日志单元Code  */
	LogRuleEnabled *bool                                        `json:"logRuleEnabled"` /*  启用日志分割规则  */
	LogRule        *CfUpdateFunctionReturnObjLogLogRuleResponse `json:"logRule"`        /*  日志切割配置  */
}

type CfUpdateFunctionReturnObjLayersResponse struct {
	LayerName   *string `json:"layerName"`   /*  层名称  */
	Version     *int32  `json:"version"`     /*  版本  */
	Description *string `json:"description"` /*  描述  */
	Acl         *int32  `json:"acl"`         /*  0表示自定义层，1表示官方公共层  */
}

type CfUpdateFunctionReturnObjRuntimeResponse struct {
	Runtime             *string `json:"runtime"`             /*  运行时类型  */
	HandleType          *string `json:"handleType"`          /*  请求处理程序类型  */
	ExecuteTimeout      *int32  `json:"executeTimeout"`      /*  执行超时时间  */
	Handler             *string `json:"handler"`             /*  函数执行的入口  */
	InstanceConcurrency *int32  `json:"instanceConcurrency"` /*  实例最大并发度  */
}

type CfUpdateFunctionReturnObjNetworkResponse struct {
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

type CfUpdateFunctionReturnObjCustomContainerConfigResponse struct {
	ImageDigest *string `json:"imageDigest"` /*  digest 用于指定镜像版本  */
	Image       *string `json:"image"`       /*  容器镜像地址  */
	InstanceId  *string `json:"instanceId"`  /*  crs 实例 id  */
}

type CfUpdateFunctionReturnObjServerlessGpuConfigResponse struct {
	GpuEnable         *bool   `json:"gpuEnable"`         /*  是否使用Gpu  */
	GpuEciType        *string `json:"gpuEciType"`        /*  GPU ECI 规格  */
	GpuMemorySize     *int32  `json:"gpuMemorySize"`     /*  单位是G  */
	GpuType           *string `json:"gpuType"`           /*  gpu卡型  */
	GpuProvisionCount *int32  `json:"gpuProvisionCount"` /*  配置的预留实例数量  */
}

type CfUpdateFunctionReturnObjOssMountResponse struct {
	Mounts []*CfUpdateFunctionReturnObjOssMountMountsResponse `json:"mounts"` /*  zos挂载参数  */
}

type CfUpdateFunctionReturnObjDnsResponse struct {
	NameServers []*string                                    `json:"nameServers"` /*  DNS 服务器的 IP 地址列表  */
	Searches    []*string                                    `json:"searches"`    /*  DNS 搜索域列表  */
	Options     *CfUpdateFunctionReturnObjDnsOptionsResponse `json:"options"`     /*  DNS 解析配置  */
}

type CfUpdateFunctionReturnObjNasResponse struct {
	Nas []*CfUpdateFunctionReturnObjNasNasResponse `json:"nas"` /*  nas  */
}

type CfUpdateFunctionReturnObjContainerHealthCheckConfigResponse struct {
	FailureThreshold    *int32  `json:"failureThreshold"`    /*  失败阈值  */
	GetPath             *string `json:"getPath"`             /*  检查http get path  */
	InitialDelaySeconds *int32  `json:"initialDelaySeconds"` /*  首次探测延迟时间(秒)  */
	PeriodSeconds       *int32  `json:"periodSeconds"`       /*  探测时间间隔(秒)  */
	SuccessThreshold    *int32  `json:"successThreshold"`    /*  成功阈值  */
	TimeoutSeconds      *int32  `json:"timeoutSeconds"`      /*  超时(秒)  */
}

type CfUpdateFunctionReturnObjLifecycleInitializerResponse struct {
	Handler *string `json:"handler"` /*  处理方法入口  */
	Enable  *bool   `json:"enable"`  /*  启用  */
	Timeout *int32  `json:"timeout"` /*  超时  */
}

type CfUpdateFunctionReturnObjLifecyclePreStopResponse struct {
	Handler *string `json:"handler"` /*  处理方法入口  */
	Enable  *bool   `json:"enable"`  /*  启用  */
	Timeout *int32  `json:"timeout"` /*  超时  */
}

type CfUpdateFunctionReturnObjLogLogRuleResponse struct {
	RuleCode         *string                                                `json:"ruleCode"`         /*  规则唯一编码  */
	RuleName         *string                                                `json:"ruleName"`         /*  规则名称  */
	ExtractMode      *int32                                                 `json:"extractMode"`      /*  采集类型  */
	CollectPolicy    *string                                                `json:"collectPolicy"`    /*  采集策略  */
	CuttingMode      *string                                                `json:"cuttingMode"`      /*  切割模式  */
	Enable           *bool                                                  `json:"enable"`           /*  是否启用采集规则  */
	UnitCode         *string                                                `json:"unitCode"`         /*  日志单元编码  */
	FirstLinePattern *string                                                `json:"firstLinePattern"` /*  首行正则  */
	CustomTime       *CfUpdateFunctionReturnObjLogLogRuleCustomTimeResponse `json:"customTime"`       /*  自定义时间戳提取格式  */
	RuleConfig       *CfUpdateFunctionReturnObjLogLogRuleRuleConfigResponse `json:"ruleConfig"`       /*  容器运行参数  */
	AccessType       *int32                                                 `json:"accessType"`       /*  接入类型  */
}

type CfUpdateFunctionReturnObjOssMountMountsResponse struct {
	BucketName *string `json:"bucketName"` /*  bucket名  */
	BucketPath *string `json:"bucketPath"` /*  bucket子目录  */
	MountDir   *string `json:"mountDir"`   /*  挂载本地目录  */
	ReadOnly   *bool   `json:"readOnly"`   /*  是否只读，默认false  */
	AccessUrl  *string `json:"accessUrl"`  /*  oss 访问地址  */
}

type CfUpdateFunctionReturnObjDnsOptionsResponse struct {
	Ndots *string `json:"ndots"` /*  键值对01  */
}

type CfUpdateFunctionReturnObjNasNasResponse struct {
	RemoteDir *string `json:"remoteDir"` /*  远端挂载目录  */
	SharePath *string `json:"sharePath"` /*  挂载地址  */
	LocalDir  *string `json:"localDir"`  /*  挂载本地目录  */
	SfsName   *string `json:"sfsName"`   /*  sfs 的名称  */
	SfsUID    *string `json:"sfsUID"`    /*  sfs 的 ID  */
}

type CfUpdateFunctionReturnObjLogLogRuleCustomTimeResponse struct {
	Key        *string `json:"key"`        /*  key  */
	TimeFormat *string `json:"timeFormat"` /*  格式化  */
}

type CfUpdateFunctionReturnObjLogLogRuleRuleConfigResponse struct {
	MaxPathDepth *int32                                                          `json:"maxPathDepth"` /*  最大正则路径解析深度  */
	Delimeter    *CfUpdateFunctionReturnObjLogLogRuleRuleConfigDelimeterResponse `json:"delimeter"`    /*  分隔符  */
	Regex        *CfUpdateFunctionReturnObjLogLogRuleRuleConfigRegexResponse     `json:"regex"`        /*  正则切割模式  */
}

type CfUpdateFunctionReturnObjLogLogRuleRuleConfigDelimeterResponse struct {
	Delimeter    *string                                                                       `json:"delimeter"`    /*  分隔符  */
	TypeContents []*CfUpdateFunctionReturnObjLogLogRuleRuleConfigDelimeterTypeContentsResponse `json:"typeContents"` /*  类型  */
}

type CfUpdateFunctionReturnObjLogLogRuleRuleConfigRegexResponse struct {
	RegexStr     *string                                                                   `json:"regexStr"`     /*  正则表达式  */
	TypeContents []*CfUpdateFunctionReturnObjLogLogRuleRuleConfigRegexTypeContentsResponse `json:"typeContents"` /*  类型  */
}

type CfUpdateFunctionReturnObjLogLogRuleRuleConfigDelimeterTypeContentsResponse struct {
	Key     *string `json:"key"`  /*  key  */
	RawType *string `json:"type"` /*  类型  */
}

type CfUpdateFunctionReturnObjLogLogRuleRuleConfigRegexTypeContentsResponse struct {
	Key     *string `json:"key"`  /*  key  */
	RawType *string `json:"type"` /*  类型  */
}
