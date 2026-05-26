package opensearch

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/opensearch"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunSearchInstance{}
	_ resource.ResourceWithConfigure   = &CtyunSearchInstance{}
	_ resource.ResourceWithImportState = &CtyunSearchInstance{}
)

func NewCtyunSearchInstance() resource.Resource {
	return &CtyunSearchInstance{}
}

type CtyunSearchInstance struct {
	meta        *common.CtyunMetadata
	name        string
	orderLooper *business.OrderLooper
}

type CtyunSearchInstanceConfig struct {
	ID              types.String `tfsdk:"id"`
	ClusterName     types.String `tfsdk:"name"`
	RegionID        types.String `tfsdk:"region_id"`
	ZoneList        types.Set    `tfsdk:"zone_list"`
	VPCID           types.String `tfsdk:"vpc_id"`
	SubnetID        types.String `tfsdk:"subnet_id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`
	EnableIPv6      types.Bool   `tfsdk:"enable_ipv6"`
	Password        types.String `tfsdk:"password"`
	ClusterType     types.Int32  `tfsdk:"cluster_type"`
	OSType          types.String `tfsdk:"os_type"`
	EnableHTTPS     types.String `tfsdk:"enable_https"`
	CycleCount      types.Int64  `tfsdk:"cycle_count"`
	CycleType       types.String `tfsdk:"cycle_type"`
	NodeDetails     types.Set    `tfsdk:"node_details"`
	nodeDetails     []CtyunSearchNodeDetail
	Status          types.String `tfsdk:"status"`
	Version         types.String `tfsdk:"version"`
}

type CtyunSearchNodeDetail struct {
	HostNum       types.Int32  `tfsdk:"host_num"`
	StorageType   types.String `tfsdk:"storage_type"`
	StorageSpace  types.Int32  `tfsdk:"storage_space"`
	FlavorName    types.String `tfsdk:"flavor_name"`
	NodeGroupType types.String `tfsdk:"node_group_type"`
}

// NodeDetailObjectType 返回 node_details 的 ObjectType
var NodeDetailObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"host_num":        types.Int32Type,
		"storage_type":    types.StringType,
		"storage_space":   types.Int32Type,
		"flavor_name":     types.StringType,
		"node_group_type": types.StringType,
	},
}

func (c *CtyunSearchInstance) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_instance"
	c.name = resp.TypeName
}

func (c *CtyunSearchInstance) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理OpenSearch实例", "天翼云OpenSearch服务", "https://www.ctyun.cn/document/10026730/10040008"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称，由大小写字母、数字、下划线(_)或连字符(-)组成，且不以下划线(_)或连字符(-)开头，长度是1-32位",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9_]*$`), "必须由大小写字母、数字、下划线(_)或连字符(-)组成，且不以下划线(_)或连字符(-)开头"),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"zone_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "实例所在可用区信息，只能传一个或三个可用区，可通过ctyun_zones查看",
				PlanModifiers: []planmodifier.Set{
					explanmodifier.NullIgnoreSet(),
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(3),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "VPC ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable_ipv6": schema.BoolAttribute{
				Required:    true,
				Description: "开启IPv6：开启:true 关闭:false",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "组件密码，创建时必填。密码应为数字、大写字母、小写字母、特殊符号 (@$!%*#_~?) 的组合，长度在 12－26 位",
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(8, 26),
					validator2.SearchPassword(),
				},
			},
			"cluster_type": schema.Int32Attribute{
				Required:    true,
				Description: "集群类型：1：OpenSearch，2：Elasticsearch",
				Validators: []validator.Int32{
					int32validator.OneOf(opensearch.ClusterTypeOpenSearch, opensearch.ClusterTypeElasticsearch),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"os_type": schema.StringAttribute{
				Required:    true,
				Description: "操作系统类型，ctyun操作系统：CTyun、麒麟操作系统：Kylin",
				Validators: []validator.String{
					stringvalidator.OneOf(opensearch.OSTypeCTyun, opensearch.OSTypeKylin),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable_https": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "不开启 https: CLOSE,开启：OPEN；默认 CLOSE",
				Validators: []validator.String{
					stringvalidator.OneOf(opensearch.EnableHTTPSOpen, opensearch.EnableHTTPSClose),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString(opensearch.EnableHTTPSClose),
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "付费周期类型，month：包月，year：包年，on_demand：按需付费",
				Validators: []validator.String{
					stringvalidator.OneOf(
						opensearch.CycleTypeMonthlyStr,
						opensearch.CycleTypeYearlyStr,
						opensearch.CycleTypeOnDemandStr,
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，cycle_type=month时取值为1-11，cycle_type=year时取值为1-5，cycle_type=on_demand时无需填写",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(opensearch.CycleTypeMonthlyStr),
						types.StringValue(opensearch.CycleTypeYearlyStr),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(opensearch.CycleTypeOnDemandStr),
					),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"node_details": schema.SetNestedAttribute{
				Required:    true,
				Description: "节点组详情列表，支持配置数据节点(MASTER)、专属主节点(EXCLUSIVE_MASTER)、协调节点(COORDINATE)和冷数据节点(COLD)四类节点",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host_num": schema.Int32Attribute{
							Required:    true,
							Description: "节点数量，MASTER 节点最小为 3，最大为 50；EXCLUSIVE_MASTER 节点最大为 3；COORDINATE 节点最大为 32；COLD 节点最大为 50",
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
						},
						"storage_type": schema.StringAttribute{
							Required:    true,
							Description: "存储类型：SSD-genric（通用型SSD）、SAS（高IO）、SSD（超高IO）、XSSD-0、XSSD-1",
							Validators: []validator.String{
								stringvalidator.OneOf(
									opensearch.IOTypeSSDGeneric,
									opensearch.IOTypeSAS,
									opensearch.IOTypeSSD,
									opensearch.IOTypeXSSD0,
									opensearch.IOTypeXSSD1,
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"storage_space": schema.Int32Attribute{
							Required:    true,
							Description: "存储空间(GB)，MASTER 节点可选：40-6144GB；EXCLUSIVE_MASTER 节点固定 40GB；COORDINATE 节点固定 40GB；COLD 节点可选：40-6144GB",
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.RequiresReplace(),
							},
						},
						"flavor_name": schema.StringAttribute{
							Required:    true,
							Description: "实例规格名称，每个资源池可用区下可选择的机型，参考订购页面展示的机型信息，如 s3.medium.2、s3.large.4 等",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"node_group_type": schema.StringAttribute{
							Required:    true,
							Description: "节点组类型：MASTER（数据节点）/EXCLUSIVE_MASTER（专属master节点）/COORDINATE（专属协调节点）/COLD（冷数据节点）",
							Validators: []validator.String{
								stringvalidator.OneOf(
									opensearch.NodeGroupTypeMaster,
									opensearch.NodeGroupTypeExclusiveMaster,
									opensearch.NodeGroupTypeCoordinate,
									opensearch.NodeGroupTypeCold,
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "实例状态",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "版本信息",
			},
		},
	}
}

func (c *CtyunSearchInstance) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.orderLooper = business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)

}

func (c *CtyunSearchInstance) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSearchInstanceConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	// 等待创建完成并获取实例详情
	time.Sleep(5 * time.Second)
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSearchInstance) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSearchInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 查询远端确认资源是否存在
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			resp.State.RemoveResource(ctx)
		}
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunSearchInstance) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSearchInstanceConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state CtyunSearchInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 目前只支持节点扩容
	err = c.scale(ctx, &plan, state)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunSearchInstance) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSearchInstanceConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = c.unsubscribe(ctx, &state)
	if err != nil {
		return
	}
	// 等待集群状态变为已销毁
	err = c.waitForClusterDestroyed(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunSearchInstance) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例：%s 失败：%s", c.name, req.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id]", c.name)
			resp.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunSearchInstanceConfig
	config.ID = types.StringValue(req.ID)

	// 查询远端
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

// create 创建实例
func (c *CtyunSearchInstance) create(ctx context.Context, plan *CtyunSearchInstanceConfig) (err error) {
	// 验证密码是否填写
	if plan.Password.IsNull() || plan.Password.IsUnknown() {
		err = fmt.Errorf("password is null")
		return err
	}

	// 将 Set 转换为 slice
	plan.nodeDetails = []CtyunSearchNodeDetail{}
	diags := plan.NodeDetails.ElementsAs(ctx, &plan.nodeDetails, false)
	if diags.HasError() {
		err = fmt.Errorf("failed to convert node_details: %v", diags)
		return
	}

	// 构造节点详情列表
	nodeDetails := make([]opensearch.NodeDetail, len(plan.nodeDetails))
	for i, node := range plan.nodeDetails {
		nodeDetails[i] = opensearch.NodeDetail{
			HostNum:        node.HostNum.ValueInt32(),
			IOType:         node.StorageType.ValueString(),
			Volume:         node.StorageSpace.ValueInt32(),
			IaasVMSpecCode: mapIaasVmSpecCode(node.FlavorName.ValueString()),
			NodeGroupType:  node.NodeGroupType.ValueString(),
		}
	}

	// 构造创建请求
	// 将 ZoneList 转换为逗号分隔的字符串
	var zoneListStr string
	var zones []string
	plan.ZoneList.ElementsAs(ctx, &zones, false)
	zoneListStr = strings.Join(zones, ",")

	createReq := &opensearch.OrderNewRequest{
		AvailableZoneID: zoneListStr,
		ClusterName:     plan.ClusterName.ValueString(),

		RegionID:        plan.RegionID.ValueString(),
		VPCID:           plan.VPCID.ValueString(),
		SubnetID:        plan.SubnetID.ValueString(),
		SecurityGroupID: plan.SecurityGroupID.ValueString(),
		ComponentPwd:    plan.Password.ValueString(),
		ClusterType:     plan.ClusterType.ValueInt32(),
		OSType:          plan.OSType.ValueString(),
		NodeDetails:     nodeDetails,
		AutoPay:         true,
	}
	if plan.EnableIPv6.ValueBool() {
		createReq.EnableIPv6 = opensearch.EnableIPv6Open
	} else {
		createReq.EnableIPv6 = opensearch.EnableIPv6Close
	}

	// 将 cycle_type 字符串转换为 int32
	var cycleTypeInt int32
	var cycleCnt int32
	var payType int32
	switch plan.CycleType.ValueString() {
	case opensearch.CycleTypeMonthlyStr:
		cycleTypeInt = opensearch.CycleTypeMonthly
		cycleCnt = int32(plan.CycleCount.ValueInt64())
		payType = opensearch.PayTypePrepaid
	case opensearch.CycleTypeYearlyStr:
		cycleTypeInt = opensearch.CycleTypeYearly
		cycleCnt = int32(plan.CycleCount.ValueInt64())
		payType = opensearch.PayTypePrepaid
	case opensearch.CycleTypeOnDemandStr:
		//cycleTypeInt = opensearch.CycleTypeMonthly
		cycleCnt = 0
		payType = opensearch.PayTypePostpaid
	}
	createReq.PayType = payType
	if cycleCnt > 0 {
		createReq.CycleCnt = &cycleCnt
	}
	// 未赋值，为按需
	if cycleTypeInt != 0 {
		createReq.CycleType = &cycleTypeInt
	}

	// 设置可选参数
	if !plan.EnableHTTPS.IsNull() && !plan.EnableHTTPS.IsUnknown() {
		enableHTTPS := plan.EnableHTTPS.ValueString()
		createReq.EnableHTTPS = &enableHTTPS
	}

	// 调用创建 API
	resp, err := c.meta.Apis.SdkOpensearchApis.OpensearchNewClusterApi.Do(ctx, c.meta.SdkCredential, createReq)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. StatusCode: %d, Message: %s, %s", resp.StatusCode, resp.Message, resp.ReturnObj)
		return
	}

	// 等待集群状态变为运行中
	err = c.waitForClusterRunning(ctx, plan)
	if err != nil {
		return
	}

	return
}

// waitForClusterRunning 等待集群状态变为运行中
func (c *CtyunSearchInstance) waitForClusterRunning(ctx context.Context, config *CtyunSearchInstanceConfig) (err error) {
	// 最大等待时间 30 分钟，每 30 秒查询一次，共 60 次
	retryer, err := business.NewRetryer(time.Second*30, 60)
	if err != nil {
		return
	}
	result := retryer.Start(func(currentTime int) bool {
		// 使用 List 接口查询，传入 cluster name
		listReq := &opensearch.ListInstancesRequest{
			RegionID:    config.RegionID.ValueString(),
			PageIndex:   1,
			PageSize:    10,
			ClusterName: config.ClusterName.ValueString(),
			ClusterType: config.ClusterType.ValueInt32(),
		}

		resp, err := c.meta.Apis.SdkOpensearchApis.OpensearchListInstancesApi.Do(ctx, c.meta.SdkCredential, listReq)
		if err != nil {
			return false
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("查询集群状态失败，StatusCode: %d, Message: %s", resp.StatusCode, resp.Message)
			return false
		} else if len(resp.ReturnObj.Records) == 0 {
			// 集群不存在，继续轮询
			return true
		}

		// 查找匹配的集群
		var foundCluster *opensearch.ClusterRecord
		for _, record := range resp.ReturnObj.Records {
			if record.ClusterName == config.ClusterName.ValueString() {
				foundCluster = &record
				config.ID = types.StringValue(record.ClusterID)
				break
			}
		}

		if foundCluster == nil {
			// 继续轮询
			return true
		}

		// 检查集群状态类型（2 表示运行中）
		if foundCluster.ClusterStateType == 2 {
			// 更新配置中的状态
			config.Status = types.StringValue(foundCluster.ClusterState)
			config.ID = types.StringValue(foundCluster.ClusterID)
			return false
		}

		// 如果状态异常（6 表示异常），直接返回错误
		if foundCluster.ClusterStateType == 6 {
			err = fmt.Errorf("集群创建失败，当前状态：%s", foundCluster.ClusterState)
			if foundCluster.ClusterMessage != nil && *foundCluster.ClusterMessage != "" {
				err = fmt.Errorf("集群创建失败，当前状态：%s, 错误原因：%s", foundCluster.ClusterState, *foundCluster.ClusterMessage)
			}
			return false
		}

		// 继续轮询
		return true
	})

	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("等待集群运行超时，最大等待时间：30分钟")
	}
	return
}

// waitForClusterDestroyed 等待集群已销毁
func (c *CtyunSearchInstance) waitForClusterDestroyed(ctx context.Context, config *CtyunSearchInstanceConfig) (err error) {
	// 最大等待时间 30 分钟，每 30 秒查询一次，共 60 次
	retryer, err := business.NewRetryer(time.Second*30, 60)
	if err != nil {
		return
	}
	result := retryer.Start(func(currentTime int) bool {
		// 使用 List 接口查询，传入 cluster name
		listReq := &opensearch.ListInstancesRequest{
			RegionID:    config.RegionID.ValueString(),
			PageIndex:   1,
			PageSize:    10,
			ClusterName: config.ClusterName.ValueString(),
			ClusterType: config.ClusterType.ValueInt32(),
		}

		resp, err := c.meta.Apis.SdkOpensearchApis.OpensearchListInstancesApi.Do(ctx, c.meta.SdkCredential, listReq)
		if err != nil {
			return false
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("查询集群状态失败，StatusCode: %d, Message: %s", resp.StatusCode, resp.Message)
			return false
		}

		// 查找匹配的集群
		var foundCluster *opensearch.ClusterRecord
		for _, record := range resp.ReturnObj.Records {
			if record.ClusterName == config.ClusterName.ValueString() {
				foundCluster = &record
				break
			}
		}

		if foundCluster == nil {
			// 集群不存在，说明已销毁成功
			return false
		}

		// 检查集群状态类型（5 表示已销毁）
		if foundCluster.ClusterStateType == 5 {
			return false
		}

		// 继续轮询
		return true
	})

	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("等待集群销毁超时，最大等待时间：30分钟")
	}
	return
}

// getAndMerge 获取实例详情并合并到配置中
func (c *CtyunSearchInstance) getAndMerge(ctx context.Context, config *CtyunSearchInstanceConfig) (err error) {
	clusterID := config.ID.ValueString()

	getReq := &opensearch.GetInstanceRequest{
		ClusterID: clusterID,
	}

	resp, err := c.meta.Apis.SdkOpensearchApis.OpensearchGetInstanceApi.Do(ctx, c.meta.SdkCredential, getReq)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. StatusCode: %d, Message: %s", resp.StatusCode, resp.Message)
		return
	}

	// 填充数据
	instance := resp.ReturnObj
	config.ID = types.StringValue(instance.ClusterID)
	config.ClusterName = types.StringValue(instance.ClusterName)
	config.Status = types.StringValue(instance.State)
	config.RegionID = types.StringValue(instance.RegionID)
	// 将逗号分隔的可用区字符串转换为 Set
	zones := strings.Split(instance.AvailableZoneID, ",")
	var zoneList []types.String
	for _, zone := range zones {
		zoneList = append(zoneList, types.StringValue(zone))
	}
	config.ZoneList, _ = types.SetValueFrom(ctx, types.StringType, zoneList)
	config.Version = types.StringValue(instance.ClusterTypeVersion)
	config.VPCID = types.StringValue(instance.VPCID)
	config.SubnetID = types.StringValue(instance.SubnetID)
	config.SecurityGroupID = types.StringValue(instance.SecurityGroupID)
	config.EnableIPv6 = types.BoolValue(instance.EnableIpv6 == opensearch.EnableIPv6Open)
	config.OSType = types.StringValue(instance.OSType)
	config.ClusterType = types.Int32Value(instance.ClusterType)
	var nodeDetails []CtyunSearchNodeDetail

	dataHostInfos := instance.DataHostInfos
	if len(dataHostInfos) > 0 {
		nodeDetail := CtyunSearchNodeDetail{
			HostNum:       types.Int32Value(int32(len(dataHostInfos))),
			StorageType:   types.StringValue(mapIOTypeToEnglish(dataHostInfos[0].IOTypeName)),
			StorageSpace:  types.Int32Value(dataHostInfos[0].DiskVolumn),
			FlavorName:    types.StringValue(dataHostInfos[0].IaasVMSpecCode),
			NodeGroupType: types.StringValue(opensearch.NodeGroupTypeMaster),
		}
		nodeDetails = append(nodeDetails, nodeDetail)
	}

	exclusiveMasterHostInfos := instance.ExclusiveMasterHostInfos
	if len(exclusiveMasterHostInfos) > 0 {
		nodeDetail := CtyunSearchNodeDetail{
			HostNum:       types.Int32Value(int32(len(exclusiveMasterHostInfos))),
			StorageType:   types.StringValue(mapIOTypeToEnglish(exclusiveMasterHostInfos[0].IOTypeName)),
			StorageSpace:  types.Int32Value(exclusiveMasterHostInfos[0].DiskVolumn),
			FlavorName:    types.StringValue(exclusiveMasterHostInfos[0].IaasVMSpecCode),
			NodeGroupType: types.StringValue(opensearch.NodeGroupTypeExclusiveMaster),
		}
		nodeDetails = append(nodeDetails, nodeDetail)
	}

	coordinateHostInfos := instance.CoordinateHostInfos
	if len(coordinateHostInfos) > 0 {
		nodeDetail := CtyunSearchNodeDetail{
			HostNum:       types.Int32Value(int32(len(coordinateHostInfos))),
			StorageType:   types.StringValue(mapIOTypeToEnglish(coordinateHostInfos[0].IOTypeName)),
			StorageSpace:  types.Int32Value(coordinateHostInfos[0].DiskVolumn),
			FlavorName:    types.StringValue(coordinateHostInfos[0].IaasVMSpecCode),
			NodeGroupType: types.StringValue(opensearch.NodeGroupTypeCoordinate),
		}
		nodeDetails = append(nodeDetails, nodeDetail)
	}

	coldHostInfos := instance.ColdHostInfos
	if len(coldHostInfos) > 0 {
		nodeDetail := CtyunSearchNodeDetail{
			HostNum:       types.Int32Value(int32(len(coldHostInfos))),
			StorageType:   types.StringValue(mapIOTypeToEnglish(coldHostInfos[0].IOTypeName)),
			StorageSpace:  types.Int32Value(coldHostInfos[0].DiskVolumn),
			FlavorName:    types.StringValue(coldHostInfos[0].IaasVMSpecCode),
			NodeGroupType: types.StringValue(opensearch.NodeGroupTypeCold),
		}
		nodeDetails = append(nodeDetails, nodeDetail)
	}

	// 将 nodeDetails slice 转换为 Set
	config.nodeDetails = nodeDetails
	nodeDetailsSet, diags := types.SetValueFrom(ctx, NodeDetailObjectType, nodeDetails)
	if diags.HasError() {
		err = fmt.Errorf("failed to convert node_details to Set: %v", diags)
		return
	}
	config.NodeDetails = nodeDetailsSet

	// 根据 PayType 推导 cycle_type
	// PayType=1（包年包月）时，根据到期时间计算 cycle_type 和 cycle_count
	// PayType=2（按需付费）时，cycle_type = "on_demand"
	if instance.PayType == "2" {
		config.CycleType = types.StringValue(opensearch.CycleTypeOnDemandStr)
		config.CycleCount = types.Int64Null()
	} else if instance.CreateTime > 0 && instance.ClusterDueTime > 0 {
		createTimeStr := time.Unix(instance.CreateTime/1000, 0).Format(time.RFC3339)
		expireTimeStr := time.Unix(instance.ClusterDueTime/1000, 0).Format(time.RFC3339)
		cycleType, cycleCount, err := utils.CalculateMonthOnlyDiff(createTimeStr, expireTimeStr)
		if err == nil {
			config.CycleType = types.StringValue(cycleType)
			config.CycleCount = types.Int64Value(int64(cycleCount))
		}
	}

	return
}

// mapIOTypeToEnglish 将 IO类型从中文映射回英文
func mapIOTypeToEnglish(ioTypeName string) string {
	switch ioTypeName {
	case "通用型SSD":
		return opensearch.IOTypeSSDGeneric
	case "高IO":
		return opensearch.IOTypeSAS
	case "超高IO":
		return opensearch.IOTypeSSD
	default:
		// 其他情况直接返回（如 XSSD-0、XSSD-1 等）
		return ioTypeName
	}
}

// mapIaasVmSpecCode 将 API 返回的机型名称（esearch-*格式）映射回 Terraform 配置中使用的名称（s7.*、c7.*等格式）
func mapIaasVmSpecCode(apiVmSpecName string) string {
	// 创建映射表：key 是 API 返回的值，value 是 Terraform 配置的值
	vmSpecMapping := map[string]string{
		// 通用型
		"esearch-4c16g":   "s7.xlarge.4",
		"esearch-8c16g":   "s7.2xlarge.2",
		"esearch-8c32g":   "s7.2xlarge.4",
		"esearch-16c32g":  "s7.4xlarge.2",
		"esearch-16c64g":  "s7.4xlarge.4",
		"esearch-32c64g":  "s7.8xlarge.2",
		"esearch-32c128g": "s7.8xlarge.4",

		// 计算型
		"esearch-c4c16g":   "c7.xlarge.4",
		"esearch-c8c16g":   "c7.2xlarge.2",
		"esearch-c8c32g":   "c7.2xlarge.4",
		"esearch-c16c32g":  "c7.4xlarge.2",
		"esearch-c16c64g":  "c7.4xlarge.4",
		"esearch-c32c64g":  "c7.8xlarge.2",
		"esearch-c32c128g": "c7.8xlarge.4",
		"esearch-c64c128g": "c7.16xlarge.2",

		// 内存型
		"esearch-m4c32g":   "m7.xlarge.8",
		"esearch-m8c64g":   "m7.2xlarge.8",
		"esearch-m16c128g": "m7.4xlarge.8",

		// 通用增强型
		"esearch-eis4c8g":   "s8.xlarge.2",
		"esearch-eis4c16g":  "s8.xlarge.4",
		"esearch-eis8c16g":  "s8.2xlarge.2",
		"esearch-eis8c32g":  "s8.2xlarge.4",
		"esearch-eis16c32g": "s8.4xlarge.2",
		"esearch-eis16c64g": "s8.4xlarge.4",
		"esearch-eis32c64g": "s8.8xlarge.2",

		// 计算增强型
		"esearch-eic4c8g":   "c8.xlarge.2",
		"esearch-eic4c16g":  "c8.xlarge.4",
		"esearch-eic8c16g":  "c8.2xlarge.2",
		"esearch-eic8c32g":  "c8.2xlarge.4",
		"esearch-eic16c32g": "c8.4xlarge.2",
		"esearch-eic16c64g": "c8.4xlarge.4",
		"esearch-eic32c64g": "c8.8xlarge.2",

		// 内存增强型
		"esearch-eim4c32g": "m8.xlarge.8",
		"esearch-eim8c64g": "m8.2xlarge.8",

		// 海光通用型
		"esearch-h1s4c8g":   "hs1.xlarge.2",
		"esearch-h1s4c16g":  "hs1.xlarge.4",
		"esearch-h1s8c16g":  "hs1.2xlarge.2",
		"esearch-h1s8c32g":  "hs1.2xlarge.4",
		"esearch-h1s16c32g": "hs1.4xlarge.2",
		"esearch-h1s16c64g": "hs1.4xlarge.4",

		// 海光计算型
		"esearch-h1c4c8g":   "hc1.xlarge.2",
		"esearch-h1c4c16g":  "hc1.xlarge.4",
		"esearch-h1c8c16g":  "hc1.2xlarge.2",
		"esearch-h1c8c32g":  "hc1.2xlarge.4",
		"esearch-h1c16c32g": "hc1.4xlarge.2",
		"esearch-h1c16c64g": "hc1.4xlarge.4",
		"esearch-h1c32c64g": "hc1.8xlarge.2",

		// 海光内存型
		"esearch-h1m4c32g": "hm1.xlarge.8",
		"esearch-h1m8c64g": "hm1.2xlarge.8",

		// 海光 4 计算型
		"esearch-h3c4c8g":    "hc3.xlarge.2",
		"esearch-h3c4c16g":   "hc3.xlarge.4",
		"esearch-h3c8c16g":   "hc3.2xlarge.2",
		"esearch-h3c8c32g":   "hc3.2xlarge.4",
		"esearch-h3c16c32g":  "hc3.4xlarge.2",
		"esearch-h3c16c64g":  "hc3.4xlarge.4",
		"esearch-h3c32c64g":  "hc3.8xlarge.2",
		"esearch-h3c32c128g": "hc3.8xlarge.4",
		"esearch-h3c64c128g": "hc3.16xlarge.2",

		// 海光 4 内存型
		"esearch-h3m4c32g":   "hm3.xlarge.8",
		"esearch-h3m8c64g":   "hm3.2xlarge.8",
		"esearch-h3m16c128g": "hm3.4xlarge.8",

		// 鲲鹏通用型
		"esearch-k1s4c8g":   "ks1.xlarge.2",
		"esearch-k1s4c16g":  "ks1.xlarge.4",
		"esearch-k1s8c16g":  "ks1.2xlarge.2",
		"esearch-k1s8c32g":  "ks1.2xlarge.4",
		"esearch-k1s16c32g": "ks1.4xlarge.2",
		"esearch-k1s16c64g": "ks1.4xlarge.4",

		// 鲲鹏计算型
		"esearch-k1c4c8g":   "kc1.xlarge.2",
		"esearch-k1c4c16g":  "kc1.xlarge.4",
		"esearch-k1c8c16g":  "kc1.2xlarge.2",
		"esearch-k1c8c32g":  "kc1.2xlarge.4",
		"esearch-k1c16c32g": "kc1.4xlarge.2",
		"esearch-k1c16c64g": "kc1.4xlarge.4",
		"esearch-k1c32c64g": "kc1.8xlarge.2",

		// 鲲鹏内存型
		"esearch-k1m4c32g": "km1.xlarge.8",
		"esearch-k1m8c64g": "km1.2xlarge.8",

		// 飞腾通用型
		"esearch-f1s4c8g":   "fs1.xlarge.2",
		"esearch-f1s4c16g":  "fs1.xlarge.4",
		"esearch-f1s8c16g":  "fs1.2xlarge.2",
		"esearch-f1s8c32g":  "fs1.2xlarge.4",
		"esearch-f1s16c32g": "fs1.4xlarge.2",
		"esearch-f1s16c64g": "fs1.4xlarge.4",

		// 飞腾计算型
		"esearch-f1c4c8g":   "fc1.xlarge.2",
		"esearch-f1c4c16g":  "fc1.xlarge.4",
		"esearch-f1c8c16g":  "fc1.2xlarge.2",
		"esearch-f1c8c32g":  "fc1.2xlarge.4",
		"esearch-f1c16c32g": "fc1.4xlarge.2",
		"esearch-f1c16c64g": "fc1.4xlarge.4",

		// 飞腾内存型
		"esearch-f1m4c32g": "fm1.xlarge.8",
		"esearch-f1m8c64g": "fm1.2xlarge.8",
	}

	// 如果映射表中存在，返回对应的值；否则返回原值
	if tfSpec, exists := vmSpecMapping[apiVmSpecName]; exists {
		return tfSpec
	}
	return apiVmSpecName
}

// scale 扩容实例
func (c *CtyunSearchInstance) scale(ctx context.Context, plan *CtyunSearchInstanceConfig, state CtyunSearchInstanceConfig) (err error) {
	// 将 Set 转换为 slice
	plan.nodeDetails = []CtyunSearchNodeDetail{}
	diags := plan.NodeDetails.ElementsAs(ctx, &plan.nodeDetails, false)
	if diags.HasError() {
		err = fmt.Errorf("failed to convert plan node_details: %v", diags)
		return
	}

	state.nodeDetails = []CtyunSearchNodeDetail{}
	diags = state.NodeDetails.ElementsAs(ctx, &state.nodeDetails, false)
	if diags.HasError() {
		err = fmt.Errorf("failed to convert state node_details: %v", diags)
		return
	}

	// 构建 state.nodeDetails 的 map，key 为 nodeGroupType，用于 O(1) 查找
	stateNodeMap := make(map[string]int32)
	for _, nodeState := range state.nodeDetails {
		stateNodeMap[nodeState.NodeGroupType.ValueString()] = nodeState.HostNum.ValueInt32()
	}

	// 过滤出需要扩容的节点组并进行校验
	var scaleNodes []struct {
		nodeGroupType   string
		increaseHostNum int32
	}
	for _, node := range plan.nodeDetails {
		nodeGroupType := node.NodeGroupType.ValueString()
		originalHostNum := stateNodeMap[nodeGroupType]
		increaseHostNum := node.HostNum.ValueInt32() - originalHostNum

		// 跳过不需要扩容的节点组
		if increaseHostNum <= 0 {
			continue
		}

		// 校验：EXCLUSIVE_MASTER 节点组不允许扩容
		if nodeGroupType == opensearch.NodeGroupTypeExclusiveMaster {
			return fmt.Errorf("EXCLUSIVE_MASTER节点组不允许扩容")
		}

		scaleNodes = append(scaleNodes, struct {
			nodeGroupType   string
			increaseHostNum int32
		}{
			nodeGroupType:   nodeGroupType,
			increaseHostNum: increaseHostNum,
		})
	}

	// 执行扩容
	for _, scaleNode := range scaleNodes {
		scaleReq := &opensearch.ScaleInstanceRequest{
			ClusterID:       plan.ID.ValueString(),
			NodeGroupName:   scaleNode.nodeGroupType,
			IncreaseHostNum: scaleNode.increaseHostNum,
			AutoPay:         true,
		}
		var resp *opensearch.ScaleInstanceResponse
		resp, err = c.meta.Apis.SdkOpensearchApis.OpensearchScaleInstanceApi.Do(ctx, c.meta.SdkCredential, scaleReq)
		if err != nil {
			return fmt.Errorf("扩容失败：%w", err)
		} else if resp.StatusCode != 200 {
			err = fmt.Errorf("API return error. StatusCode: %d, Message: %s, %s", resp.StatusCode, resp.Message, resp.ReturnObj)
		}
		// 等待集群状态变为运行中
		err = c.waitForClusterRunning(ctx, plan)
		if err != nil {
			return
		}
	}

	return
}

// unsubscribe 退订实例
func (c *CtyunSearchInstance) unsubscribe(ctx context.Context, state *CtyunSearchInstanceConfig) (err error) {
	unsubscribeReq := &opensearch.UnsubscribeInstanceRequest{
		ClusterID: state.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkOpensearchApis.OpensearchUnsubscribeInstanceApi.Do(ctx, c.meta.SdkCredential, unsubscribeReq)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. StatusCode: %d, Message: %s", resp.StatusCode, resp.Message)
		return
	}

	return
}
