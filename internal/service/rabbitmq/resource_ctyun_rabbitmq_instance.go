package rabbitmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"

	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"

	"time"
)

var (
	_ resource.Resource                = &ctyunRabbitmqInstance{}
	_ resource.ResourceWithConfigure   = &ctyunRabbitmqInstance{}
	_ resource.ResourceWithImportState = &ctyunRabbitmqInstance{}
)

type ctyunRabbitmqInstance struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunRabbitmqInstance() resource.Resource {
	return &ctyunRabbitmqInstance{}
}

func (c *ctyunRabbitmqInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_instance"
}

type CtyunRabbitmqInstanceConfig struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	MasterOrderID   types.String `tfsdk:"master_order_id"`
	ProjectID       types.String `tfsdk:"project_id"`
	RegionID        types.String `tfsdk:"region_id"`
	ZoneList        types.Set    `tfsdk:"zone_list"`
	InstanceName    types.String `tfsdk:"instance_name"`
	SpecName        types.String `tfsdk:"spec_name"`
	DiskType        types.String `tfsdk:"disk_type"`
	DiskSize        types.Int32  `tfsdk:"disk_size"`
	NodeNum         types.Int32  `tfsdk:"node_num"`
	VpcID           types.String `tfsdk:"vpc_id"`
	SubnetID        types.String `tfsdk:"subnet_id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`
	CycleType       types.String `tfsdk:"cycle_type"`
	CycleCount      types.Int32  `tfsdk:"cycle_count"`

	zoneList []string
}

func (c *ctyunRabbitmqInstance) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10000118/10001967`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "名称",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "主订单号",
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
				Description: "实例的规格类型，建议使用ctyun_rabbitmq_specs查看，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"node_num": schema.Int32Attribute{
				Required:    true,
				Description: "节点数。支持1、3、5、7、9，支持更新",
				Validators: []validator.Int32{
					int32validator.OneOf(1, 3, 5, 7, 9),
				},
			},
			"disk_type": schema.StringAttribute{
				Required:    true,
				Description: "磁盘类型，通常支持SAS、SSD、FAST-SSD",
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
				},
			},
			"zone_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "实例所在可用区信息，只能传一个或三个可用区，可通过ctyun_regions查看",
				PlanModifiers: []planmodifier.Set{
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
				Description: "虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Required:    true,
				Description: "安全组ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SecurityGroupValidate(),
				},
			},

			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf("month", "on_demand"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month时才生效，当cycle_type=month，支持传递1、2、3、4、5、6、12、24、36",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					int32validator.OneOf(1, 2, 3, 5, 6, 7, 12, 24, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunRabbitmqInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRabbitmqInstanceConfig
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

func (c *ctyunRabbitmqInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRabbitmqInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunRabbitmqInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunRabbitmqInstanceConfig
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
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRabbitmqInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	instance, err := c.getByID(ctx, state)
	if err != nil {
		return
	}
	// 如果状态不是已退订状态，则执行退订
	if instance.Status != business.RabbitMqStatusUnsubscribed {
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
	// 销毁
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDestroy(ctx, state)
	if err != nil {
		return
	}
	response.Diagnostics.AddWarning("删除RabbitMq集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

func (c *ctyunRabbitmqInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionID]
func (c *ctyunRabbitmqInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunRabbitmqInstanceConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunRabbitmqInstance) checkBeforeCreate(ctx context.Context, plan *CtyunRabbitmqInstanceConfig) (err error) {
	regionID, projectID := plan.RegionID.ValueString(), plan.ProjectID.ValueString()
	vpc, subnetID, sgID := plan.VpcID.ValueString(), plan.SubnetID.ValueString(), plan.SecurityGroupID.ValueString()
	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpc, regionID, projectID)
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
func (c *ctyunRabbitmqInstance) checkSpecParams(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
	nodeNum := plan.NodeNum.ValueInt32()
	specName := plan.SpecName.ValueString()
	diskType := plan.DiskType.ValueString()

	if strings.HasSuffix(specName, "single") && nodeNum != 1 {
		return fmt.Errorf("单机版实例节点数必须为1")
	} else if strings.HasSuffix(specName, "cluster") && nodeNum < 3 {
		return fmt.Errorf("集群版实例节点数必须大于等于3")
	}
	// 组装请求体
	params := &amqp.AmqpProdDetailRequest{
		RegionId: plan.RegionID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpProdDetailApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var skuRes amqp.AmqpProdDetailReturnObjDataSeriesSkuResItemResponse
	var skuDisk amqp.AmqpProdDetailReturnObjDataSeriesSkuDiskItemResponse
	for _, s := range resp.ReturnObj.Data.Series {
		for _, p := range s.Sku {
			if p.ProdName == "集群版" && plan.NodeNum.ValueInt32() >= 3 {
				skuRes = p.ResItem
				skuDisk = p.DiskItem
				break
			} else if p.ProdName == "单机版" && plan.NodeNum.ValueInt32() == 1 {
				skuRes = p.ResItem
				skuDisk = p.DiskItem
				break
			}
		}
	}

	var specAvailable bool
	for _, r := range skuRes.ResItems {
		for _, s := range r.Spec {
			if s.SpecName == specName {
				specAvailable = true
				break
			}
		}
		if specAvailable {
			break
		}
	}
	if !specAvailable {
		return fmt.Errorf("本资源池不支持 %s", specName)
	}

	var diskAvailable bool
	for _, d := range skuDisk.ResItems {
		if d == diskType {
			diskAvailable = true
			break
		}
	}
	if !diskAvailable {
		return fmt.Errorf("本资源池不支持 %s", diskType)
	}

	return
}

// checkZoneList 检查zoneList
func (c *ctyunRabbitmqInstance) checkZoneList(ctx context.Context, plan *CtyunRabbitmqInstanceConfig) (err error) {
	zones, err := business.NewRegionService(c.meta).GetZonesByRegionID(ctx, plan.RegionID.ValueString())
	if err != nil {
		return err
	}
	z := map[string]bool{}
	for _, az := range zones {
		z[az] = true
	}

	var zoneList []string
	var str []types.String
	plan.ZoneList.ElementsAs(ctx, &str, true)
	for _, s := range str {
		zoneList = append(zoneList, s.ValueString())
	}
	plan.zoneList = zoneList
	return
}

// create 创建
func (c *ctyunRabbitmqInstance) create(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (masterOrderID string, err error) {
	switch plan.CycleType.ValueString() {
	case business.OrderCycleTypeMonth:
		return c.createPrePayOrder(ctx, plan)
	case business.OrderCycleTypeOnDemand:
		return c.createPostPayOrder(ctx, plan)
	}
	return
}

// createPrePayOrder 创建包年包月
func (c *ctyunRabbitmqInstance) createPrePayOrder(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (masterOrderID string, err error) {
	params := &amqp.AmqpInstancesCreatePrePayOrderRequest{
		CycleCnt:        plan.CycleCount.ValueInt32(),
		RegionId:        plan.RegionID.ValueString(),
		ClusterName:     plan.InstanceName.ValueString(),
		ProjectId:       plan.ProjectID.ValueString(),
		SpecName:        plan.SpecName.ValueString(),
		NodeNum:         plan.NodeNum.ValueInt32(),
		DiskType:        plan.DiskType.ValueString(),
		DiskSize:        plan.DiskSize.ValueInt32(),
		VpcId:           plan.VpcID.ValueString(),
		SubnetId:        plan.SubnetID.ValueString(),
		SecurityGroupId: plan.SecurityGroupID.ValueString(),
		ZoneList:        plan.zoneList,
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesCreatePrePayOrderApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.Data.NewOrderId
	return
}

// createPostPayOrder 创建按需
func (c *ctyunRabbitmqInstance) createPostPayOrder(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (masterOrderID string, err error) {
	params := &amqp.AmqpInstancesCreatePostPayOrderRequest{
		RegionId:        plan.RegionID.ValueString(),
		ClusterName:     plan.InstanceName.ValueString(),
		ProjectId:       plan.ProjectID.ValueString(),
		SpecName:        plan.SpecName.ValueString(),
		NodeNum:         plan.NodeNum.ValueInt32(),
		DiskType:        plan.DiskType.ValueString(),
		DiskSize:        plan.DiskSize.ValueInt32(),
		VpcId:           plan.VpcID.ValueString(),
		SubnetId:        plan.SubnetID.ValueString(),
		SecurityGroupId: plan.SecurityGroupID.ValueString(),
		ZoneList:        plan.zoneList,
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesCreatePostPayOrderApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = resp.ReturnObj.Data.NewOrderId
	return
}

// getAndMerge 从远端查询
func (c *ctyunRabbitmqInstance) getAndMerge(ctx context.Context, plan *CtyunRabbitmqInstanceConfig) (err error) {
	instance, err := c.getByID(ctx, *plan)
	if err != nil {
		return
	}
	plan.InstanceName = types.StringValue(instance.ClusterName)
	plan.Name = plan.InstanceName
	if plan.ZoneList.IsNull() {
		plan.ZoneList = types.SetNull(types.StringType)
	}

	plan.DiskSize = types.Int32Value(utils.StringToInt32Must(instance.Space) / instance.NodeCount)
	plan.NodeNum = types.Int32Value(instance.NodeCount)
	plan.SpecName = types.StringValue(instance.Prod)
	return
}

func (c *ctyunRabbitmqInstance) checkBeforeUpdate(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	instance, err := c.getByID(ctx, state)
	if err != nil {
		return
	}
	if instance.Status != 1 {
		return fmt.Errorf("请在实例处于运行中状态时再进行更新操作")
	}

	return nil
}

// update 更新
func (c *ctyunRabbitmqInstance) update(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
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
func (c *ctyunRabbitmqInstance) updateDiskSize(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
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
func (c *ctyunRabbitmqInstance) diskExtend(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	params := &amqp.AmqpInstancesDiskExtendRequest{
		RegionId:       state.RegionID.ValueString(),
		ProdInstId:     state.ID.ValueString(),
		DiskExtendSize: plan.DiskSize.ValueInt32(),
		AutoPay:        true,
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesDiskExtendApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterUpdateDiskSize 检查磁盘大小是否变更成功
func (c *ctyunRabbitmqInstance) checkAfterUpdateDiskSize(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *amqp.AmqpInstancesQueryDetailResponseReturnObjData
			instance, err = c.getByID(ctx, state)
			if err != nil {
				return false
			}
			if instance.Status != 1 || utils.StringToInt32Must(instance.Space) != plan.DiskSize.ValueInt32()*instance.NodeCount {
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
		err = fmt.Errorf("磁盘变配时间过长")
	}
	return
}

// updateNodeNum 更新节点数量
func (c *ctyunRabbitmqInstance) updateNodeNum(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
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
func (c *ctyunRabbitmqInstance) nodeExtend(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	params := &amqp.AmqpInstancesNodeExtendRequest{
		RegionId:      state.RegionID.ValueString(),
		ProdInstId:    state.ID.ValueString(),
		ExtendNodeNum: plan.NodeNum.ValueInt32(),
		AutoPay:       true,
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesNodeExtendApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterUpdateNodeNum 检查节点数量是否变更成功
func (c *ctyunRabbitmqInstance) checkAfterUpdateNodeNum(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *amqp.AmqpInstancesQueryDetailResponseReturnObjData
			instance, err = c.getByID(ctx, state)
			if err != nil {
				return false
			}
			if instance.Status != 1 || instance.NodeCount != plan.NodeNum.ValueInt32() {
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
		err = fmt.Errorf("节点数量变配时间过长")
	}
	return
}

// updateSpec 更新规格
func (c *ctyunRabbitmqInstance) updateSpec(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
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
func (c *ctyunRabbitmqInstance) specExtend(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	params := &amqp.AmqpInstancesSpecExtendRequest{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.ID.ValueString(),
		SpecName:   plan.SpecName.ValueString(),
		AutoPay:    true,
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesSpecExtendApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterUpdateSpec 检查规格是否变更成功
func (c *ctyunRabbitmqInstance) checkAfterUpdateSpec(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	var successCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *amqp.AmqpInstancesQueryDetailResponseReturnObjData
			instance, err = c.getByID(ctx, state)
			if err != nil {
				return false
			}
			if instance.Status != 1 || instance.Prod != plan.SpecName.ValueString() {
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
		err = fmt.Errorf("规格变配时间过长")
	}
	return
}

// updateName 更新实例名称
func (c *ctyunRabbitmqInstance) updateName(ctx context.Context, plan, state CtyunRabbitmqInstanceConfig) (err error) {
	if plan.InstanceName.Equal(state.InstanceName) {
		return
	}
	params := &amqp.AmqpInstancesInstanceNameRequest{
		RegionId:     state.RegionID.ValueString(),
		ProdInstId:   state.ID.ValueString(),
		InstanceName: plan.InstanceName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesInstanceNameApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	return
}

// unsubscribe 退订
func (c *ctyunRabbitmqInstance) unsubscribe(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
	params := &amqp.AmqpInstancesUnsubscribeInstRequest{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesUnsubscribeInstApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// destroy 销毁
func (c *ctyunRabbitmqInstance) destroy(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
	params := &amqp.AmqpInstanceDeleteRequest{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstanceDeleteApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunRabbitmqInstance) checkAfterCreate(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (id string, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *amqp.AmqpInstancesQueryResponseReturnObjData
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			if instance == nil || instance.Status != 1 || instance.ProdInstId == "" {
				return true
			}
			time.Sleep(30 * time.Second)
			id = instance.ProdInstId
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("创建时间过长")
	}
	return
}

// checkAfterUnsubscribe 退订后检查
func (c *ctyunRabbitmqInstance) checkAfterUnsubscribe(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *amqp.AmqpInstancesQueryResponseReturnObjData
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			if instance != nil && instance.Status != business.RabbitMqStatusUnsubscribed {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("退订时间过长")
	}
	return
}

// checkAfterDestroy 销毁后检查
func (c *ctyunRabbitmqInstance) checkAfterDestroy(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var instance *amqp.AmqpInstancesQueryResponseReturnObjData
			instance, err = c.getByName(ctx, plan)
			if err != nil {
				return false
			}
			if instance != nil {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("销毁时间过长")
	}
	return
}

// getByName 根据名称查询集群
func (c *ctyunRabbitmqInstance) getByName(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (instance *amqp.AmqpInstancesQueryResponseReturnObjData, err error) {
	params := &amqp.AmqpInstancesQueryRequest{
		RegionId: plan.RegionID.ValueString(),
		PageNum:  1,
		PageSize: 100,
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesQueryApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, r := range resp.ReturnObj.Data {
		if r.ClusterName == plan.InstanceName.ValueString() {
			instance = r
			return
		}
	}
	return
}

// getByID 根据ID查询集群
func (c *ctyunRabbitmqInstance) getByID(ctx context.Context, plan CtyunRabbitmqInstanceConfig) (instance *amqp.AmqpInstancesQueryDetailResponseReturnObjData, err error) {
	params := &amqp.AmqpInstancesQueryDetailRequest{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkAmqpApis.AmqpInstancesQueryDetailApi.Do(ctx, c.meta.Credential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = resp.ReturnObj.Data
	return
}

// parseSpec 从规格名称解析cpu和mem
func (c *ctyunRabbitmqInstance) parseSpec(s string) (u, m int, err error) {
	re := regexp.MustCompile(`(\d+)u(\d+)g`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		err = fmt.Errorf("invalid format: %s", s)
		return
	}

	if _, err = fmt.Sscanf(matches[1], "%d", &u); err != nil {
		return
	}
	if _, err = fmt.Sscanf(matches[2], "%d", &m); err != nil {
		return
	}
	return
}
