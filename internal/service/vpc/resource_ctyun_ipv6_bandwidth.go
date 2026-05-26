package vpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ctyunIPv6BandwidthService{}
	_ resource.ResourceWithConfigure   = &ctyunIPv6BandwidthService{}
	_ resource.ResourceWithImportState = &ctyunIPv6BandwidthService{}
)

type ctyunIPv6BandwidthService struct {
	meta        *common.CtyunMetadata
	name        string
	orderLooper *business.OrderLooper
}

func NewCtyunIPv6BandwidthService() resource.Resource {
	return &ctyunIPv6BandwidthService{}
}

func (c *ctyunIPv6BandwidthService) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ipv6_bandwidth"
	c.name = response.TypeName
}

type CtyunIPv6BandwidthConfig struct {
	ID              types.String `tfsdk:"id"`
	RegionID        types.String `tfsdk:"region_id"`
	Name            types.String `tfsdk:"name"`
	Bandwidth       types.Int64  `tfsdk:"bandwidth"`
	CycleType       types.String `tfsdk:"cycle_type"`
	CycleCount      types.Int64  `tfsdk:"cycle_count"`
	PayVoucherPrice types.String `tfsdk:"pay_voucher_price"`
	Status          types.String `tfsdk:"status"`
	ResourceSpec    types.String `tfsdk:"resource_spec"`
	CreateTime      types.String `tfsdk:"create_time"`
	ExpireTime      types.String `tfsdk:"expire_time"`
	IpAddress       types.String `tfsdk:"ip_address"`
	Ipv6GatewayId   types.String `tfsdk:"ipv6_gateway_id"`
	Expired         types.Bool   `tfsdk:"expired"`
}

func (c *ctyunIPv6BandwidthService) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理IPv6带宽", "IPv6带宽（IPv6 Bandwidth）", "https://www.ctyun.cn/document/10026753/10037269"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "IPv6带宽ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"name": schema.StringAttribute{
				Description: "IPV6带宽实例名称",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"bandwidth": schema.Int64Attribute{
				Description: "带宽大小，取值范围1～300Mbps",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 300),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，year：按年，on_demand：按需",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-3年，支持更新",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 3),
				},
			},
			"pay_voucher_price": schema.StringAttribute{
				Optional:    true,
				Description: "代金券金额，只适用于预付费客户自动支付，若代金券支付金额传0或者控制符，则不适用代金券支付（小数会只保留2位，非负）",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "IPv6带宽状态，ACTIVE（正常） / EXPIRED（过期） / FREEZING（冻结） /CREATEING（创建中）",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_spec": schema.StringAttribute{
				Computed:    true,
				Description: "独享/共享模式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:      true,
				Description:   "创建时间",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"expire_time": schema.StringAttribute{
				Computed:      true,
				Description:   "过期时间",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ip_address": schema.StringAttribute{
				Computed:      true,
				Description:   "IP地址",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ipv6_gateway_id": schema.StringAttribute{
				Computed:      true,
				Description:   "IPv6网关",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"expired": schema.BoolAttribute{
				Computed:    true,
				Description: "是否过期",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ctyunIPv6BandwidthService) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunIPv6BandwidthConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	bandwidthId, err := c.loopCreate(ctx, &plan)
	if err != nil {
		return
	}

	plan.ID = types.StringValue(bandwidthId)

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunIPv6BandwidthService) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIPv6BandwidthConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunIPv6BandwidthService) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunIPv6BandwidthConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunIPv6BandwidthConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, &plan, &state)
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

func (c *ctyunIPv6BandwidthService) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIPv6BandwidthConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除
	masterOrderID, err := c.delete(ctx, &state)
	if err != nil {
		return
	}

	err = c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, masterOrderID)
	return
}

func (c *ctyunIPv6BandwidthService) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.orderLooper = business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
}

func (c *ctyunIPv6BandwidthService) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import %s.[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunIPv6BandwidthConfig

	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionId
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	var cycleType string
	var cycleCount int32

	cycleType, cycleCount, err = utils.CalculateMonthOnlyDiff(config.CreateTime.ValueString(), config.ExpireTime.ValueString())
	if err != nil {
		return
	}
	config.CycleType = types.StringValue(cycleType)
	if cycleCount > 0 {
		config.CycleCount = types.Int64Value(int64(cycleCount))
	} else {
		config.CycleCount = types.Int64Null()
	}

	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

// loopCreate 循环执行create
func (c *ctyunIPv6BandwidthService) loopCreate(ctx context.Context, plan *CtyunIPv6BandwidthConfig) (id string, err error) {
	clientToken := uuid.NewString()
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			id, err = c.create(ctx, clientToken, plan)
			if err != nil {
				return false
			}
			if id != "" {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = errors.New("创建时未获取到IPv6带宽ID")
	}
	return
}

// create 创建
func (c *ctyunIPv6BandwidthService) create(ctx context.Context, clientToken string, plan *CtyunIPv6BandwidthConfig) (id string, err error) {
	params := &ctvpc.CtvpcCreateIPv6BandwidthRequest{
		ClientToken:     clientToken,
		RegionID:        plan.RegionID.ValueString(),
		Bandwidth:       int32(plan.Bandwidth.ValueInt64()),
		CycleType:       plan.CycleType.ValueString(),
		CycleCount:      int32(plan.CycleCount.ValueInt64()),
		Name:            plan.Name.ValueString(),
		PayVoucherPrice: plan.PayVoucherPrice.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateIPv6BandwidthApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	id = *resp.ReturnObj.BandwidthID
	return
}

func (c *ctyunIPv6BandwidthService) getAndMerge(ctx context.Context, plan *CtyunIPv6BandwidthConfig) (err error) {
	ID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcNewIPv6BandwidthListRequest{
		RegionID:    regionID,
		BandwidthID: &ID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewIPv6BandwidthListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Objs) == 0 {
		err = common.ResourceNotExistError
		return
	}
	bandwidth := resp.ReturnObj.Objs[0]

	plan.ID = types.StringValue(ID)
	plan.Name = types.StringValue(utils.SecString(bandwidth.Name))
	plan.Bandwidth = types.Int64Value(int64(bandwidth.Bandwidth))
	plan.Status = types.StringValue(utils.SecString(bandwidth.Status))
	plan.ResourceSpec = types.StringValue(utils.SecString(bandwidth.ResourceSpec))
	plan.CreateTime = types.StringValue(utils.SecString(bandwidth.CreatedTime))
	plan.ExpireTime = types.StringValue(utils.SecString(bandwidth.ExpiredTime))
	plan.IpAddress = types.StringValue(utils.SecString(bandwidth.IpAddress))
	plan.Ipv6GatewayId = types.StringValue(utils.SecString(bandwidth.Ipv6GatewayID))
	plan.Expired = types.BoolValue(bandwidth.Expired)

	return
}

func (c *ctyunIPv6BandwidthService) update(ctx context.Context, plan, state *CtyunIPv6BandwidthConfig) (err error) {
	ID, regionID := state.ID.ValueString(), state.RegionID.ValueString()

	if !plan.Name.Equal(state.Name) {
		params := &ctvpc.CtvpcUpdateIPv6BandwidthAttributeRequest{
			RegionID:    regionID,
			BandwidthID: ID,
			Name:        plan.Name.ValueString(),
		}
		var resp *ctvpc.CtvpcUpdateIPv6BandwidthAttributeResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcUpdateIPv6BandwidthAttributeApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}

	if !plan.Bandwidth.Equal(state.Bandwidth) {
		clientToken := uuid.NewString()
		params := &ctvpc.CtvpcModifyIPv6BandwidthSpecRequest{
			ClientToken:     clientToken,
			RegionID:        regionID,
			BandwidthID:     ID,
			Bandwidth:       int32(plan.Bandwidth.ValueInt64()),
			PayVoucherPrice: plan.PayVoucherPrice.ValueStringPointer(),
		}
		var resp *ctvpc.CtvpcModifyIPv6BandwidthSpecResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcModifyIPv6BandwidthSpecApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return
		}

		masterOrderId := resp.ReturnObj.MasterOrderID

		helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
		_, err = helper.OrderLoop(ctx, c.meta.Credential, *masterOrderId)
		if err != nil {
			return
		}
	}

	return
}

func (c *ctyunIPv6BandwidthService) delete(ctx context.Context, plan *CtyunIPv6BandwidthConfig) (masterOrderID string, err error) {
	clientToken := uuid.NewString()
	ID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteIPv6BandwidthRequest{
		RegionID:    regionID,
		BandwidthID: ID,
		ClientToken: clientToken,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteIPv6BandwidthApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	masterOrderID = utils.SecString(resp.ReturnObj.MasterOrderID)
	return
}
