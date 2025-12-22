package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunKafkaConsumerGroup{}
	_ resource.ResourceWithConfigure   = &ctyunKafkaConsumerGroup{}
	_ resource.ResourceWithImportState = &ctyunKafkaConsumerGroup{}
)

type ctyunKafkaConsumerGroup struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunKafkaConsumerGroup() resource.Resource {
	return &ctyunKafkaConsumerGroup{}
}

func (c *ctyunKafkaConsumerGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_consumer_group"
}

type CtyunKafkaConsumerGroupConfig struct {
	ID            types.Int32  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	RegionID      types.String `tfsdk:"region_id"`
	InstanceId    types.String `tfsdk:"instance_id"`
	Description   types.String `tfsdk:"description"`
	Ctime         types.String `tfsdk:"create_time"`
	State         types.String `tfsdk:"state"`
	CoordinatorId types.Int32  `tfsdk:"coordinator_id"`

	// 重置消费点配置对象
	ResetConfig *CtyunKafkaConsumerGroupResetConfig `tfsdk:"reset_config"` /*  重置时间点毫秒时间戳，type=1时必填。  */
}

// 重置消费点配置结构体
type CtyunKafkaConsumerGroupResetConfig struct {
	TopicName          types.String                                                      `tfsdk:"topic_name"`           /*  主题名称。  */
	Type               types.Int32                                                       `tfsdk:"type"`                 /*  类型，<li>0：重置到latest。 <li>1：按时间重置。<li>2：重置到earliest。<li>3：按位点重置，此类型参数partitionShiftList为必填。  */
	PartitionShiftList []*ctgkafka.CtgkafkaConsumerGroupResetV3PartitionShiftListRequest `tfsdk:"partition_shift_list"` /*  位点重置列表，当type为3时必填。  */
	Time               types.Int64                                                       `tfsdk:"timestamp"`            /*  重置时间点毫秒时间戳，type=1时必填。  */
}

func (c *ctyunKafkaConsumerGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10145103`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				PlanModifiers: []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "消费组名称，规则如下：以英文字母、数字、下划线开头，且只能由英文字母、数字、句点、中划线、下划线组成 长度3-64。 名称不可重复。 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.-]*$`),
						"必须以英文字母、数字、下划线开头，只能包含英文字母、数字、句点、中划线、下划线",
					),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID。支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "消费组描述，规则如下：\n不能以+,-,@,= 特殊字符开头。\n长度不能大于200。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 200),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^+\-@=].*$`),
						"不能以+,-,@,=特殊字符开头",
					),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "消费组状态",
			},
			"coordinator_id": schema.Int32Attribute{
				Computed:    true,
				Description: "协调器编号",
			},
			// 重置消费点配置对象
			"reset_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "重置消费点配置",
				Attributes: map[string]schema.Attribute{
					//重置消费点请求参数
					"topic_name": schema.StringAttribute{
						Required:    true,
						Description: "主题名称 支持更新",
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
					},
					"type": schema.Int32Attribute{
						Required:    true,
						Description: "类型，0：重置到latest；1：按时间重置；2：重置到earliest；3：按位点重置，此类型参数partitionShiftList为必填 支持更新",
						Validators: []validator.Int32{
							int32validator.Between(0, 3),
						},
					},
					"timestamp": schema.Int64Attribute{
						Optional:    true,
						Description: "重置时间点毫秒时间戳，type=1时必填 支持更新",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"partition_shift_list": schema.ListNestedAttribute{
						Optional:    true,
						Description: "位点重置列表，当type为3时必填 支持更新",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"partition": schema.Int32Attribute{
									Optional:    true,
									Description: "主题分区号 支持更新",
									Validators: []validator.Int32{
										int32validator.AtLeast(1),
									},
								},
								"shift_by": schema.Int64Attribute{
									Optional:    true,
									Description: "主题分区消费位点向左或向右移动的相对位置，例如当前offset是1000，当shiftBy=-10重置后offset=990，当shiftBy=10重置后offset=1010。支持更新",
									Validators: []validator.Int64{
										int64validator.AtLeast(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunKafkaConsumerGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunKafkaConsumerGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	// 创建后检查
	id, err := c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.Int32Value(id)

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunKafkaConsumerGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaConsumerGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaConsumerGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunKafkaConsumerGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunKafkaConsumerGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.reset(ctx, plan, state)
	if err != nil {
		return
	}
	state.ResetConfig = plan.ResetConfig
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaConsumerGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaConsumerGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 销毁
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunKafkaConsumerGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

func (c *ctyunKafkaConsumerGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceId],[groupName],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunKafkaConsumerGroupConfig
	var instanceId, regionID, groupName string

	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceId, &groupName)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceId, &groupName, &regionID)
		if err != nil {
			return
		}
	}

	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if groupName == "" {
		err = fmt.Errorf("groupName不能为空")
		return
	}

	cfg.RegionID = types.StringValue(regionID)
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.Name = types.StringValue(groupName)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunKafkaConsumerGroup) create(ctx context.Context, plan CtyunKafkaConsumerGroupConfig) (err error) {
	params := &ctgkafka.CtgkafkaConsumerGroupCreateV3Request{
		RegionId:    plan.RegionID.ValueString(),
		GroupName:   plan.Name.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		Description: plan.Description.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaConsumerGroupCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
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

// update 更新
func (c *ctyunKafkaConsumerGroup) update(ctx context.Context, plan, state CtyunKafkaConsumerGroupConfig) (err error) {

	params := &ctgkafka.CtgkafkaConsumerGroupUpdateRequest{
		RegionId:    state.RegionID.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		GroupName:   plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaConsumerGroupUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if resp.ReturnObj.Data != "update success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
	}
	return
}

// reset 重置消费组消费点
func (c *ctyunKafkaConsumerGroup) reset(ctx context.Context, plan, state CtyunKafkaConsumerGroupConfig) (err error) {
	// 检查 ResetConfig 是否存在
	if plan.ResetConfig == nil {
		return
	}
	// 检查必要参数是否存在
	if plan.ResetConfig.Type.IsNull() || plan.ResetConfig.TopicName.IsNull() {
		return
	}

	params := &ctgkafka.CtgkafkaConsumerGroupResetV3Request{
		RegionId:   state.RegionID.ValueString(),
		ProdInstId: state.InstanceId.ValueString(),
		GroupName:  plan.Name.ValueString(),
		TopicName:  plan.ResetConfig.TopicName.ValueString(),
		RawType:    plan.ResetConfig.Type.ValueInt32(),
		Time:       plan.ResetConfig.Time.ValueInt64(),
	}

	// 判空处理PartitionShiftList
	if plan.ResetConfig.PartitionShiftList != nil && len(plan.ResetConfig.PartitionShiftList) > 0 {
		params.PartitionShiftList = make([]*ctgkafka.CtgkafkaConsumerGroupResetV3PartitionShiftListRequest, 0, len(plan.ResetConfig.PartitionShiftList))
		for _, item := range plan.ResetConfig.PartitionShiftList {
			if item != nil {
				params.PartitionShiftList = append(params.PartitionShiftList, &ctgkafka.CtgkafkaConsumerGroupResetV3PartitionShiftListRequest{
					Partition: item.Partition,
					ShiftBy:   item.ShiftBy,
				})
			}
		}
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaConsumerGroupResetV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if resp.ReturnObj.Data != "reset success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
	}
	return
}

// destroy 销毁
func (c *ctyunKafkaConsumerGroup) destroy(ctx context.Context, plan CtyunKafkaConsumerGroupConfig) (err error) {
	params := &ctgkafka.CtgkafkaConsumerGroupDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		GroupName:  plan.Name.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaConsumerGroupDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
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

// checkAfterCreate 创建后检查
func (c *ctyunKafkaConsumerGroup) checkAfterCreate(ctx context.Context, plan CtyunKafkaConsumerGroupConfig) (id int32, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			err = c.getAndMerge(ctx, &plan)
			if err != nil {
				return false
			}
			if plan.ID.IsNull() {
				return true
			}
			id = plan.ID.ValueInt32()
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

// getAndMerge 从远端查询
func (c *ctyunKafkaConsumerGroup) getAndMerge(ctx context.Context, plan *CtyunKafkaConsumerGroupConfig) (err error) {
	params := &ctgkafka.CtgkafkaConsumerGroupQueryV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		GroupName:  plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaConsumerGroupQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Data) > 0 {
		plan.ID = types.Int32Value(resp.ReturnObj.Data[0].Id)
		plan.Ctime = types.StringValue(utils.ConvertToUTCZ(utils.Layout2, resp.ReturnObj.Data[0].Ctime))
		plan.State = types.StringValue(resp.ReturnObj.Data[0].State)
		plan.CoordinatorId = types.Int32Value(resp.ReturnObj.Data[0].CoordinatorId)
		plan.Description = types.StringValue(resp.ReturnObj.Data[0].Description)
		plan.Name = types.StringValue(resp.ReturnObj.Data[0].Name)
		if plan.ID.IsNull() {
			err = common.InvalidReturnObjResultsError
		}
	}
	return
}
