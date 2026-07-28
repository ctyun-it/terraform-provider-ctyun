package rocketmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/rocketmq"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ctyunRocketmqInstance{}
	_ resource.ResourceWithConfigure   = &ctyunRocketmqInstance{}
	_ resource.ResourceWithImportState = &ctyunRocketmqInstance{}
)

type ctyunRocketmqInstance struct {
	meta        *common.CtyunMetadata
	name        string
	vpcService  *business.VpcService
	sgService   *business.SecurityGroupService
	orderLooper *business.OrderLooper
}

func NewCtyunRocketmqInstance() resource.Resource {
	return &ctyunRocketmqInstance{}
}

func (c *ctyunRocketmqInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_instance"
	c.name = response.TypeName
}

type CtyunRocketmqInstanceConfig struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	MasterOrderID       types.String `tfsdk:"master_order_id"`
	ProjectID           types.String `tfsdk:"project_id"`
	RegionID            types.String `tfsdk:"region_id"`
	ZoneList            types.Set    `tfsdk:"zone_list"`
	InstanceName        types.String `tfsdk:"instance_name"`
	SpecName            types.String `tfsdk:"spec_name"`
	AutoRenew           types.Bool   `tfsdk:"auto_renew"`
	AutoRenewCycleCount types.Int32  `tfsdk:"auto_renew_cycle_count"`
	NodeNum             types.Int32  `tfsdk:"node_num"`
	DiskType            types.String `tfsdk:"disk_type"`
	DiskSize            types.Int32  `tfsdk:"disk_size"`
	//AzInfo              types.String `tfsdk:"az_info"`
	VpcID           types.String `tfsdk:"vpc_id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`
	SubnetID        types.String `tfsdk:"subnet_id"`
	CycleType       types.String `tfsdk:"cycle_type"`
	CycleCount      types.Int32  `tfsdk:"cycle_count"`
	CreateTime      types.String `tfsdk:"create_time"`
	ExpireTime      types.String `tfsdk:"expire_time"`

	//zoneList []string
}

func (c *ctyunRocketmqInstance) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理 RocketMQ 实例", "分布式消息服务RocketMQ", "https://www.ctyun.cn/document/10000114/10123760"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "实例 ID",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "实例名称",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					explanmodifier.Project(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"master_order_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "主订单号",
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"instance_name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"spec_name": schema.StringAttribute{
				Required:    true,
				Description: "实例的规格类型，建议使用ctyun_rocketmq_specs查看，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "是否自动续订，仅在 cycle_type 为 month 时生效。默认不自动续订",
			},
			"auto_renew_cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "自动续订周期时长（单位：月），仅在 auto_renew 为 true 时必填。取值范围：1,2,3,4,5,6,12,24,36",
				Validators: []validator.Int32{
					int32validator.OneOf(1, 2, 3, 4, 5, 6, 12, 24, 36),
				},
			},
			"node_num": schema.Int32Attribute{
				Required:    true,
				Description: "broker 节点数，值等于代理数*2，取值范围为 [1,32]，单机版传 1",
				Validators: []validator.Int32{
					int32validator.Between(1, 32),
				},
			},
			"disk_type": schema.StringAttribute{
				Required:    true,
				Description: "存储类型，支持 SAS、SSD、FAST-SSD",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("SAS", "SSD", "FAST-SSD"),
				},
			},
			"disk_size": schema.Int32Attribute{
				Required:    true,
				Description: "单个节点的磁盘存储空间，单位为GB，必须为100的倍数，实例总存储空间为diskSize * nodeNum，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(100, 10000),
					validator2.RocketmqDiskSize(),
				},
			},

			"zone_list": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "实例所在可用区信息，只能传一个或三个可用区，可通过ctyun_regions查看",
				PlanModifiers: []planmodifier.Set{
					explanmodifier.NullIgnoreSet(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(3),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "私有云 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SecurityGroupValidate(),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，on_demand：按需，支持更新。当此值为month时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeMonth, business.OrderCycleTypeOnDemand),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，该参数在 cycle_type 为 month 时才生效，当 cycle_type=month，支持传递 1、2、3、4、5、6、12、24、36，从按需变为包周期时支持更新",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					int32validator.OneOf(1, 2, 3, 4, 5, 6, 12, 24, 36),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，UTC 格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "到期时间，为 UTC 格式，按需时为空",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ctyunRocketmqInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRocketmqInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	// 创建
	masterOrderID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.MasterOrderID = types.StringValue(masterOrderID)
	// 创建后检查
	id, err := c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(id)

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRocketmqInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRocketmqInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRocketmqInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunRocketmqInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunRocketmqInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeUpdate(ctx, plan, state)
	if err != nil {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	state.CycleType, state.CycleCount = plan.CycleType, plan.CycleCount
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	if !plan.ZoneList.IsUnknown() && !plan.ZoneList.IsNull() && state.ZoneList.IsNull() {
		state.ZoneList = plan.ZoneList
		response.Diagnostics.AddWarning("zone_list的更新仅写入状态文件", "在import时，状态文件中zone_list为null，允许用模板中的值进行一次更新，该更新不触发远程调用")
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRocketmqInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRocketmqInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	instance, err := c.getByID(ctx, state)
	if err != nil {
		return
	}
	// 如果状态不是已退订状态，则执行退订
	if instance.State != business.RabbitMqStatusUnsubscribed { // Changed from string to business.RabbitMqStatusUnsubscribed
		// 退订
		err = c.unsubscribe(ctx, state)
		if err != nil {
			return
		}
		err = c.checkAfterUnsubscribe(ctx, state)
		if err != nil {
			return
		}
		time.Sleep(60 * time.Second)
	}

	response.Diagnostics.AddWarning("删除RabbitMq集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

func (c *ctyunRocketmqInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
	c.orderLooper = business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
}

func (c *ctyunRocketmqInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRocketmqInstanceConfig
	var id, regionID string
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		id = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &id, &regionID)
		if err != nil {
			return
		}
	}

	if id == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}

	if cfg.ZoneList.IsNull() {
		cfg.ZoneList = types.SetNull(types.StringType)
	}

	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunRocketmqInstance) checkBeforeCreate(ctx context.Context, plan *CtyunRocketmqInstanceConfig) (err error) {
	regionID := plan.RegionID.ValueString()
	vpc, subnetID, sgID := plan.VpcID.ValueString(), plan.SubnetID.ValueString(), plan.SecurityGroupID.ValueString()
	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpc, regionID)
	if err != nil {
		return err
	}
	_, exist := subnets[subnetID]
	if !exist {
		err = fmt.Errorf("子网不存在")
		return err
	}
	err = c.sgService.MustExistInVpc(ctx, vpc, sgID, regionID)
	if err != nil {
		return err
	}
	err = c.checkZoneList(ctx, plan)
	if err != nil {
		return err
	}
	err = c.checkSpecParams(ctx, *plan)
	if err != nil {
		return err
	}
	return nil
}

// checkSpecParams 检查规格参数
func (c *ctyunRocketmqInstance) checkSpecParams(ctx context.Context, plan CtyunRocketmqInstanceConfig) (err error) {
	nodeNum := plan.NodeNum.ValueInt32()
	specName := plan.SpecName.ValueString()

	if strings.HasSuffix(specName, "single") && nodeNum != 1 {
		return fmt.Errorf("单机版实例节点数必须为1")
	} else if strings.HasSuffix(specName, "cluster") && nodeNum < 3 {
		return fmt.Errorf("集群版实例节点数必须大于等于3")
	}
	// 组装请求体
	params := &rocketmq.RocketmqProdDetailRequest{
		RegionId: plan.RegionID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.RocketmqApis.RocketmqProdDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var specAvailable bool
	for _, s := range resp.ReturnObj.Data {
		if s.SpecName == specName {
			specAvailable = true
			break
		}
	}

	if !specAvailable {
		return fmt.Errorf("本资源池不支持 %s", specName)
	}

	return
}

// checkZoneList 检查用户输入的可用区是否在指定资源池支持的可用区列表中
func (c *ctyunRocketmqInstance) checkZoneList(ctx context.Context, plan *CtyunRocketmqInstanceConfig) (err error) {
	// 从 API 获取资源池支持的可用区列表
	zones, err := business.NewRegionService(c.meta).GetZonesByRegionID(ctx, plan.RegionID.ValueString())
	if err != nil {
		return err
	}

	// 将可用区列表转换为 map，用于 O(1) 时间复杂度的查找
	// map 结构：key=可用区名称（如 "az1"），value=true（标记存在）
	// 例如：{"az1": true, "az2": true, "az3": true}
	validZones := map[string]bool{}
	for _, az := range zones {
		validZones[az] = true
	}

	// 获取用户配置的可用区列表
	var str []types.String
	plan.ZoneList.ElementsAs(ctx, &str, true)

	// 验证每个用户输入的可用区是否在资源池支持的列表中
	for _, s := range str {
		zone := s.ValueString()
		// map[zone] 语法：查找 map 中是否存在 key=zone
		// 如果存在返回 true，不存在返回 false
		// !validZones[zone] 表示：如果该可用区不在支持的列表中，则报错
		if !validZones[zone] {
			err = fmt.Errorf("可用区 %s 不在资源池 %s 支持的可用区列表中", zone, plan.RegionID.ValueString())
			return
		}
	}

	return
}

// buildAzInfo 构建 azInfo 参数
func (c *ctyunRocketmqInstance) buildAzInfo(ctx context.Context, plan CtyunRocketmqInstanceConfig) (azInfoStr string, err error) {
	var zoneList []string
	var str []types.String
	plan.ZoneList.ElementsAs(ctx, &str, true)
	for _, s := range str {
		zoneList = append(zoneList, s.ValueString())
	}

	// 构建 JSON 数组格式的 azInfo
	azInfoArray := make([]map[string]interface{}, 0, len(zoneList))
	for _, zone := range zoneList {
		azInfo := map[string]interface{}{
			"azName": zone,
			"azId":   0, // azId 设置为 0，让后端自动分配
		}
		azInfoArray = append(azInfoArray, azInfo)
	}

	// 转换为 JSON 字符串
	azInfoJSON, err := json.Marshal(azInfoArray)
	if err != nil {
		return "", err
	}
	azInfoStr = string(azInfoJSON)

	return azInfoStr, nil
}

// create 创建
func (c *ctyunRocketmqInstance) create(ctx context.Context, plan CtyunRocketmqInstanceConfig) (masterOrderID string, err error) {
	switch plan.CycleType.ValueString() {
	case business.OrderCycleTypeMonth:
		return c.createPrePayOrder(ctx, plan)
	case business.OrderCycleTypeOnDemand:
		return c.createPostPayOrder(ctx, plan)
	}
	return
}

// createPrePayOrder 创建包年包月
func (c *ctyunRocketmqInstance) createPrePayOrder(ctx context.Context, plan CtyunRocketmqInstanceConfig) (masterOrderID string, err error) {
	azInfo, err := c.buildAzInfo(ctx, plan)
	if err != nil {
		return
	}
	params := &rocketmq.RocketmqCreatePostPayOrderRequest{
		CycleCnt:    plan.CycleCount.ValueInt32(),
		RegionId:    plan.RegionID.ValueString(),
		ClusterName: plan.InstanceName.ValueString(),
		//ProjectId:           plan.ProjectID.ValueString(),
		SpecName:            plan.SpecName.ValueString(),
		NodeNum:             plan.NodeNum.ValueInt32(),
		DiskType:            plan.DiskType.ValueString(),
		DiskSize:            plan.DiskSize.ValueInt32(),
		VpcId:               plan.VpcID.ValueString(),
		SubnetId:            plan.SubnetID.ValueString(),
		SecurityGroupId:     plan.SecurityGroupID.ValueString(),
		AzInfo:              azInfo,
		AutoRenew:           plan.AutoRenew.ValueBool(),
		AutoRenewCycleCount: plan.AutoRenewCycleCount.ValueInt32(),
	}

	resp, err := c.meta.Apis.RocketmqApis.RocketmqCreatePostPayOrderApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.Data.NewOrderId
	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, masterOrderID)
	return
}

// createPostPayOrder 创建按需
func (c *ctyunRocketmqInstance) createPostPayOrder(ctx context.Context, plan CtyunRocketmqInstanceConfig) (masterOrderID string, err error) {
	azInfo, err := c.buildAzInfo(ctx, plan)
	if err != nil {
		return
	}
	params := &rocketmq.RocketmqCreatePostPayOrderRequest{
		RegionId:    plan.RegionID.ValueString(),
		ClusterName: plan.InstanceName.ValueString(),
		//ProjectId:       plan.ProjectID.ValueString(),
		SpecName:        plan.SpecName.ValueString(),
		NodeNum:         plan.NodeNum.ValueInt32(),
		DiskType:        plan.DiskType.ValueString(),
		DiskSize:        plan.DiskSize.ValueInt32(),
		VpcId:           plan.VpcID.ValueString(),
		SubnetId:        plan.SubnetID.ValueString(),
		SecurityGroupId: plan.SecurityGroupID.ValueString(),
		AzInfo:          azInfo,
	}

	resp, err := c.meta.Apis.RocketmqApis.RocketmqCreatePostPayOrderApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.Data.NewOrderId
	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, masterOrderID)
	return
}

// getAndMerge 从远端查询
func (c *ctyunRocketmqInstance) getAndMerge(ctx context.Context, plan *CtyunRocketmqInstanceConfig) (err error) {
	instance, err := c.getByID(ctx, *plan)
	if err != nil {
		return
	}
	plan.InstanceName = types.StringValue(instance.ProdInstName)
	plan.Name = plan.InstanceName
	plan.SpecName = types.StringValue(instance.MachineSpec)
	plan.NodeNum = types.Int32Value(instance.NodeSize)
	plan.DiskType = types.StringValue(instance.DiskType)
	// 从 "200G" 转换为 200
	diskSpace := c.parseDiskSpace(instance.DiskSpace)
	plan.DiskSize = types.Int32Value(diskSpace)

	cTime := ""
	if instance.CrtTime != "" && instance.CrtTime != "null" {
		cTime = utils.ConvertToUTCZ(time.RFC3339, instance.CrtTime)
	}
	eTime := ""
	if instance.ExpTime != "" && instance.ExpTime != "null" {
		eTime = utils.ConvertToUTCZ(time.RFC3339, instance.ExpTime)
	}
	plan.CreateTime = types.StringValue(cTime)
	plan.ExpireTime = types.StringValue(eTime)

	return
}

func (c *ctyunRocketmqInstance) checkBeforeUpdate(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	instance, err := c.getByID(ctx, state)
	if err != nil {
		return
	}
	if instance.State != 1 {
		return fmt.Errorf("请在实例处于运行中状态时再进行更新操作")
	}
	if strings.Contains(plan.SpecName.ValueString(), "single") && !strings.Contains(state.SpecName.ValueString(), "single") {
		return fmt.Errorf("不支持单机版和集群版互转")
	}
	if strings.Contains(plan.SpecName.ValueString(), "cluster") && !strings.Contains(state.SpecName.ValueString(), "cluster") {
		return fmt.Errorf("不支持单机版和集群版互转")
	}
	if plan.CycleType.Equal(state.CycleType) && !plan.CycleCount.Equal(state.CycleCount) {
		return fmt.Errorf("不支持续订")
	}

	return nil
}

// update 更新
func (c *ctyunRocketmqInstance) update(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {

	err = c.updateName(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateDiskSize(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateNodeNum(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateSpec(ctx, plan, state)
	if err != nil {
		return
	}

	return
}

// updateDiskSize 更新磁盘大小
func (c *ctyunRocketmqInstance) updateDiskSize(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	if plan.DiskSize.Equal(state.DiskSize) {
		return
	}
	if plan.DiskSize.ValueInt32() > state.DiskSize.ValueInt32() {
		err = c.diskExtend(ctx, plan, state)
	} else {
		err = fmt.Errorf("目前不支持磁盘缩容")
	}
	if err != nil {
		return
	}
	return c.checkAfterUpdateDiskSize(ctx, plan, state)
}

// diskExtend 磁盘扩容
func (c *ctyunRocketmqInstance) diskExtend(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	params := &rocketmq.RocketmqDiskExtendRequest{
		RegionId:       state.RegionID.ValueString(),
		ProdInstId:     state.ID.ValueString(),
		DiskExtendSize: plan.DiskSize.ValueInt32(),
		AutoPay:        true,
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqDiskExtendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, resp.ReturnObj.Data.NewOrderId)
	return
}

// checkAfterUpdateDiskSize 检查磁盘大小是否变更成功
func (c *ctyunRocketmqInstance) checkAfterUpdateDiskSize(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *rocketmq.ProdInstInfo
			instance, err = c.getByID(ctx, state)
			if err != nil {
				return false
			}
			diskSpace := c.parseDiskSpace(instance.DiskSpace)
			if instance.State != 1 || diskSpace != plan.DiskSize.ValueInt32() {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("实例 %s(%s) 磁盘变配时间过长", plan.InstanceName.ValueString(), state.ID.ValueString())
	}
	return
}

// updateNodeNum 更新节点数量
func (c *ctyunRocketmqInstance) updateNodeNum(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	if plan.NodeNum.Equal(state.NodeNum) {
		return
	}
	if plan.NodeNum.ValueInt32() > state.NodeNum.ValueInt32() {
		err = c.nodeExtend(ctx, plan, state)
	} else {
		err = fmt.Errorf("目前不支持节点缩容")
	}
	if err != nil {
		return
	}
	return c.checkAfterUpdateNodeNum(ctx, plan, state)
}

// nodeExtend 节点扩容
func (c *ctyunRocketmqInstance) nodeExtend(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	params := &rocketmq.RocketmqNodeExtendRequest{
		RegionId:      state.RegionID.ValueString(),
		ProdInstId:    state.ID.ValueString(),
		ExtendNodeNum: plan.NodeNum.ValueInt32(),
		AutoPay:       true,
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqNodeExtendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, resp.ReturnObj.Data.NewOrderId)
	return
}

// checkAfterUpdateNodeNum 检查节点数量是否变更成功
func (c *ctyunRocketmqInstance) checkAfterUpdateNodeNum(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *rocketmq.ProdInstInfo
			instance, err = c.getByID(ctx, state)
			if err != nil {
				return false
			}
			if instance.State != 1 || instance.NodeSize != plan.NodeNum.ValueInt32() {
				return true
			}
			successCnt++
			if successCnt < 3 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("实例 %s(%s) 节点数量变配时间过长", plan.InstanceName.ValueString(), state.ID.ValueString())
	}
	return
}

// updateSpec 更新规格
func (c *ctyunRocketmqInstance) updateSpec(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	if plan.SpecName.Equal(state.SpecName) {
		return
	}
	ou, om, _ := c.parseSpec(state.SpecName.ValueString())
	u, m, _ := c.parseSpec(plan.SpecName.ValueString())
	if u <= ou && m <= om {
		err = fmt.Errorf("只支持规格扩容")
		return
	}
	err = c.specExtend(ctx, plan, state)
	if err != nil {
		return
	}
	return c.checkAfterUpdateSpec(ctx, plan, state)
}

// specExtend 规格扩容
func (c *ctyunRocketmqInstance) specExtend(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	cpuNum, memSize := c.parseSpecToCpuMem(plan.SpecName.ValueString())
	params := &rocketmq.RocketmqSpecExtendRequest{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.ID.ValueString(),
		CpuNum:     cpuNum,
		MemSize:    memSize,
		AutoPay:    true,
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqSpecExtendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, resp.ReturnObj.Data.NewOrderId)
	return
}

// checkAfterUpdateSpec 检查规格是否变更成功
func (c *ctyunRocketmqInstance) checkAfterUpdateSpec(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *rocketmq.ProdInstInfo
			instance, err = c.getByID(ctx, state)
			if err != nil {
				return false
			}
			if instance.State != 1 || instance.MachineSpec != plan.SpecName.ValueString() {
				return true
			}
			successCnt++
			if successCnt < 3 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("实例 %s(%s) 规格变配时间过长", plan.InstanceName.ValueString(), state.ID.ValueString())
	}
	return
}

// updateName 更新实例名称
func (c *ctyunRocketmqInstance) updateName(ctx context.Context, plan, state CtyunRocketmqInstanceConfig) (err error) {
	if plan.InstanceName.Equal(state.InstanceName) {
		return
	}
	params := &rocketmq.RocketmqInstanceNameV3Request{
		RegionId:     state.RegionID.ValueString(),
		ProdInstId:   state.ID.ValueString(),
		InstanceName: plan.InstanceName.ValueString(),
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqInstanceNameV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	return
}

// unsubscribe 退订
func (c *ctyunRocketmqInstance) unsubscribe(ctx context.Context, plan CtyunRocketmqInstanceConfig) (err error) {
	params := &rocketmq.RocketmqUnsubscribeInstRequest{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqUnsubscribeInstApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil ||
		len(resp.ReturnObj.Data.BatchOrderPlacementResults) == 0 ||
		len(resp.ReturnObj.Data.BatchOrderPlacementResults[0].OrderPlacedEvents) == 0 {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID := resp.ReturnObj.Data.BatchOrderPlacementResults[0].OrderPlacedEvents[0].NewOrderId
	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, masterOrderID)
	return
}

//// destroy 销毁
//func (c *ctyunRocketmqInstance) destroy(ctx context.Context, plan CtyunRocketmqInstanceConfig) (err error) {
//	params := &rocketmq.RocketmqInstanceDeleteRequest{
//		RegionId:   plan.RegionID.ValueString(),
//		ProdInstId: plan.ID.ValueString(),
//	}
//	resp, err := c.meta.Apis.RocketmqApis.RocketmqInstanceDeleteApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	}
//	return
//}

// checkAfterCreate 创建后检查
func (c *ctyunRocketmqInstance) checkAfterCreate(ctx context.Context, plan CtyunRocketmqInstanceConfig) (id string, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *rocketmq.ProdInstInfo
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			if instance == nil || instance.State != 1 || instance.ProdInstId == "" {
				return true
			}
			id = instance.ProdInstId
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("实例 %s 创建时间过长", plan.InstanceName.ValueString())
	}
	return
}

// checkAfterUnsubscribe 退订后检查
func (c *ctyunRocketmqInstance) checkAfterUnsubscribe(ctx context.Context, state CtyunRocketmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *rocketmq.ProdInstInfo
			instance, err = c.getByName(ctx, state)
			if err != nil {
				return false
			}
			if instance != nil && instance.State != business.RabbitMqStatusUnsubscribed {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("实例 %s(%s) 退订时间过长", state.InstanceName.ValueString(), state.ID.ValueString())
	}
	return
}

// getByName 根据名称查询集群
func (c *ctyunRocketmqInstance) getByName(ctx context.Context, plan CtyunRocketmqInstanceConfig) (instance *rocketmq.ProdInstInfo, err error) {
	params := &rocketmq.RocketmqInstQueryV3Request{
		RegionId: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqInstQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, r := range resp.ReturnObj.ProdInstList {
		if r.ProdInstName == plan.InstanceName.ValueString() {
			instance = r
			return
		}
	}
	return
}

// getByID 根据ID查询集群
func (c *ctyunRocketmqInstance) getByID(ctx context.Context, plan CtyunRocketmqInstanceConfig) (instance *rocketmq.ProdInstInfo, err error) {
	params := &rocketmq.RocketmqInstQueryV3Request{
		RegionId: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqInstQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, r := range resp.ReturnObj.ProdInstList {
		if r.ProdInstId == plan.ID.ValueString() {
			instance = r
			return
		}
	}
	err = common.ResourceNotExistError
	return
}

// parseDiskSpace 从磁盘空间字符串解析出数字
// 例如："200G" -> 200, "100G" -> 100
func (c *ctyunRocketmqInstance) parseDiskSpace(diskSpace string) int32 {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(diskSpace)
	if len(matches) >= 2 {
		return utils.StringToInt32Must(matches[1])
	}
	return 0
}

// parseSpecToCpuMem 从规格名称解析 CPU 核数和内存大小
// 例如：rocketmq.2u4g.cluster -> cpuNum=2, memSize=4
func (c *ctyunRocketmqInstance) parseSpecToCpuMem(specName string) (cpuNum, memSize int32) {
	re := regexp.MustCompile(`(\d+)u(\d+)g`)
	matches := re.FindStringSubmatch(specName)
	if len(matches) == 3 {
		cpuNum = utils.StringToInt32Must(matches[1])
		memSize = utils.StringToInt32Must(matches[2])
	}
	return
}

// parseSpec 从规格名称解析 cpu 和 mem（保留原有方法，返回 int 类型）
func (c *ctyunRocketmqInstance) parseSpec(s string) (u, m int, err error) {
	re := regexp.MustCompile(`(\d+)u(\d+)g`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		err = fmt.Errorf("invalid format: %s", s)
		return
	}

	u, err = strconv.Atoi(matches[1])
	if err != nil {
		return
	}
	m, err = strconv.Atoi(matches[2])
	if err != nil {
		return
	}
	return
}
