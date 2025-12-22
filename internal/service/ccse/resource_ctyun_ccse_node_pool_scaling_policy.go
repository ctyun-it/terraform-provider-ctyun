package ccse

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
//	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
//	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
//	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"strings"
//)
//
//var (
//	_ resource.Resource                = &ctyunCcseNodePoolScalingPolicy{}
//	_ resource.ResourceWithConfigure   = &ctyunCcseNodePoolScalingPolicy{}
//	_ resource.ResourceWithImportState = &ctyunCcseNodePoolScalingPolicy{}
//)
//
//type ctyunCcseNodePoolScalingPolicy struct {
//	meta *common.CtyunMetadata
//}
//
//func NewCtyunCcseScalingNodePoolPolicy() resource.Resource {
//	return &ctyunCcseNodePoolScalingPolicy{}
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
//	response.TypeName = request.ProviderTypeName + "_ccse_node_pool_scaling_policy"
//}
//
//type CtyunCcseScalingNodePoolPolicyConfig struct {
//	ID           types.String `tfsdk:"id"`
//	ClusterID    types.String `tfsdk:"cluster_id"`
//	RegionID     types.String `tfsdk:"region_id"`
//	ValuesYaml   types.String `tfsdk:"values_yaml"`
//	ActualConfig types.String `tfsdk:"actual_config"`
//	NodePoolName types.String `tfsdk:"node_pool_name"`
//	Name         types.String `tfsdk:"name"`
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
//	response.Schema = schema.Schema{
//		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10083472/10269202`,
//		Attributes: map[string]schema.Attribute{
//			"id": schema.StringAttribute{
//				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
//				Computed:      true,
//				Description:   "ID",
//			},
//			"region_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//				Validators: []validator.String{
//					stringvalidator.UTF8LengthAtLeast(1),
//				},
//				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
//			},
//			"cluster_id": schema.StringAttribute{
//				Required:    true,
//				Description: "集群ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//				Validators: []validator.String{
//					stringvalidator.UTF8LengthBetween(32, 32),
//				},
//			},
//			"values_yaml": schema.StringAttribute{
//				Optional:    true,
//				Description: "配置参数(YAML格式)，支持更新",
//				Validators: []validator.String{
//					validator2.Yaml("apiVersion", "kind", "metadata.name"),
//				},
//			},
//			"actual_config": schema.StringAttribute{
//				Computed:    true,
//				Description: "实际配置(YAML格式)",
//			},
//			"name": schema.StringAttribute{
//				Computed:    true,
//				Description: "策略名称，为配置参数中的metadata.name",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"node_pool_name": schema.StringAttribute{
//				Computed:    true,
//				Description: "节点池名称",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//		},
//	}
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var plan CtyunCcseScalingNodePoolPolicyConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	err = c.checkBeforeCreate(ctx, &plan)
//	if err != nil {
//		return
//	}
//	// 创建
//	err = c.create(ctx, plan)
//	if err != nil {
//		return
//	}
//	// 反查信息
//	err = c.getAndMerge(ctx, &plan)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunCcseScalingNodePoolPolicyConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 查询远端
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		if errors.Is(err, common.ResourceNotExistError) {
//			err = nil
//			response.State.RemoveResource(ctx)
//		}
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	// tf文件中的
//	var plan CtyunCcseScalingNodePoolPolicyConfig
//	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// state中的
//	var state CtyunCcseScalingNodePoolPolicyConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 更新
//	err = c.update(ctx, &plan, &state)
//	if err != nil {
//		return
//	}
//
//	// 查询远端信息
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//
//	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			response.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunCcseScalingNodePoolPolicyConfig
//	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//	// 删除
//	err = c.delete(ctx, state)
//	if err != nil {
//		return
//	}
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//	meta := request.ProviderData.(*common.CtyunMetadata)
//	c.meta = meta
//}
//
//func (c *ctyunCcseNodePoolScalingPolicy) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			title := "导入失败：" + err.Error()
//			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[clusterID],[region_id]"
//			response.Diagnostics.AddError(title, detail)
//		}
//	}()
//	var cfg CtyunCcseScalingNodePoolPolicyConfig
//	var nodePoolName, clusterID, regionID string
//	// 根据分隔符数量判断是否输入了regionID
//	if strings.Count(request.ID, common.ImportSeparator) < 2 {
//		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
//		err = terraform_extend.Split(request.ID, &nodePoolName, &clusterID)
//		if err != nil {
//			return
//		}
//	} else {
//		err = terraform_extend.Split(request.ID, &nodePoolName, &clusterID, &regionID)
//		if err != nil {
//			return
//		}
//	}
//
//	if nodePoolName == "" {
//		err = fmt.Errorf("nodePoolName不能为空")
//		return
//	}
//	if clusterID == "" {
//		err = fmt.Errorf("clusterID不能为空")
//		return
//	}
//	if regionID == "" {
//		err = fmt.Errorf("regionID不能为空")
//		return
//	}
//
//	cfg.NodePoolName = types.StringValue(nodePoolName)
//	cfg.Name = types.StringValue(fmt.Sprintf("%s-%s", nodePoolName, clusterID))
//	cfg.RegionID = types.StringValue(regionID)
//	cfg.ClusterID = types.StringValue(clusterID)
//	// 查询远端
//	err = c.getAndMerge(ctx, &cfg)
//	if err != nil {
//		return
//	}
//	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
//}
//
//// checkBeforeCreate 创建前检查
//func (c *ctyunCcseNodePoolScalingPolicy) checkBeforeCreate(ctx context.Context, plan *CtyunCcseScalingNodePoolPolicyConfig) (err error) {
//	config, err := utils.ParseYamlValue(plan.ValuesYaml.ValueString(), "metadata.name")
//	if err != nil {
//		return
//	}
//	s, _ := config.(string)
//	ss := strings.SplitN(s, "-", 2)
//	if len(ss) != 2 {
//		err = fmt.Errorf("invalid metadata.name")
//		return
//	}
//	if ss[1] != plan.ClusterID.ValueString() {
//		err = fmt.Errorf("metadata.name must be ${NodePoolName}-${ClusterID}")
//		return
//	}
//	plan.Name = types.StringValue(s)
//	plan.NodePoolName = types.StringValue(ss[0])
//	return
//}
//
//// create 创建
//func (c *ctyunCcseNodePoolScalingPolicy) create(ctx context.Context, plan CtyunCcseScalingNodePoolPolicyConfig) (err error) {
//	params := &ccse.CcseCreateClusterAutoscalerPolicyRequest{
//		ClusterId:           plan.ClusterID.ValueString(),
//		RegionId:            plan.RegionID.ValueString(),
//		TextPlainDataString: plan.ValuesYaml.ValueString(),
//	}
//
//	resp, err := c.meta.Apis.SdkCcseApis.CcseCreateClusterAutoscalerPolicyApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	}
//	return
//}
//
//// getAndMerge 从远端查询
//func (c *ctyunCcseNodePoolScalingPolicy) getAndMerge(ctx context.Context, plan *CtyunCcseScalingNodePoolPolicyConfig) (err error) {
//	if plan.NodePoolName.ValueString() == "" {
//		var config interface{}
//		config, err = utils.ParseYamlValue(plan.ValuesYaml.ValueString(), "metadata.name")
//		if err != nil {
//			return
//		}
//		s, _ := config.(string)
//		ss := strings.SplitN(s, "-", 2)
//		if len(ss) != 2 {
//			err = fmt.Errorf("invalid metadata.name")
//			return
//		}
//		plan.Name = types.StringValue(s)
//		plan.NodePoolName = types.StringValue(ss[0])
//	}
//	config, err := c.getScaling(ctx, *plan)
//	if err != nil {
//		return
//	}
//	plan.ActualConfig = types.StringValue(config)
//	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.NodePoolName.ValueString(), plan.ClusterID.ValueString(), plan.RegionID.ValueString()))
//	return
//}
//
//// update 更新
//func (c *ctyunCcseNodePoolScalingPolicy) update(ctx context.Context, plan, state *CtyunCcseScalingNodePoolPolicyConfig) (err error) {
//	params := &ccse.CcseUpdateClusterAutoscalerPolicyRequest{
//		ClusterId:           state.ClusterID.ValueString(),
//		Name:                state.Name.ValueString(),
//		RegionId:            state.RegionID.ValueString(),
//		TextPlainDataString: plan.ValuesYaml.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCcseApis.CcseUpdateClusterAutoscalerPolicyApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	}
//	state.ValuesYaml = plan.ValuesYaml
//	return
//}
//
//// delete 删除
//func (c *ctyunCcseNodePoolScalingPolicy) delete(ctx context.Context, plan CtyunCcseScalingNodePoolPolicyConfig) (err error) {
//	params := &ccse.CcseDeleteClusterAutoscalerPolicyRequest{
//		ClusterId: plan.ClusterID.ValueString(),
//		Name:      plan.Name.ValueString(),
//		RegionId:  plan.RegionID.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteClusterAutoscalerPolicyApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		return
//	}
//	return
//}
//
//// getScaling 查询弹性伸缩策略
//func (c *ctyunCcseNodePoolScalingPolicy) getScaling(ctx context.Context, plan CtyunCcseScalingNodePoolPolicyConfig) (script string, err error) {
//	params := &ccse.CcseGetClusterAutoscalerPolicyRequest{
//		ClusterId: plan.ClusterID.ValueString(),
//		Name:      plan.Name.ValueString(),
//		RegionId:  plan.RegionID.ValueString(),
//	}
//	resp, err := c.meta.Apis.SdkCcseApis.CcseGetClusterAutoscalerPolicyApi.Do(ctx, c.meta.SdkCredential, params)
//	if err != nil {
//		return
//	} else if resp.StatusCode != common.NormalStatusCode {
//		if strings.Contains(resp.Message, "not found") {
//			err = common.ResourceNotExistError
//		} else {
//			err = fmt.Errorf("API return error. Message: %s", resp.Message)
//		}
//		return
//	}
//	script = resp.ReturnObj
//	return
//}
