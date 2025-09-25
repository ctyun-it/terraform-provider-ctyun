package ccse

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource                = &ctyunCcsePlugin{}
	_ resource.ResourceWithConfigure   = &ctyunCcsePlugin{}
	_ resource.ResourceWithImportState = &ctyunCcsePlugin{}
)

type ctyunCcsePlugin struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcsePlugin() resource.Resource {
	return &ctyunCcsePlugin{}
}

func (c *ctyunCcsePlugin) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_plugin"
}

type CtyunCcsePluginConfig struct {
	ID           types.String `tfsdk:"id"`
	ClusterID    types.String `tfsdk:"cluster_id"`
	RegionID     types.String `tfsdk:"region_id"`
	ChartName    types.String `tfsdk:"chart_name"`
	ChartVersion types.String `tfsdk:"chart_version"`
	ValuesYaml   types.String `tfsdk:"values_yaml"`
	ValuesJson   types.String `tfsdk:"values_json"`
	Namespace    types.String `tfsdk:"namespace"`
}

func (c *ctyunCcsePlugin) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10083472/10102631**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "集群ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"chart_name": schema.StringAttribute{
				Required:    true,
				Description: "插件名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"chart_version": schema.StringAttribute{
				Required:    true,
				Description: "插件版本号，可通过ctyun_ccse_plugin_market查询，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values_yaml": schema.StringAttribute{
				Optional:    true,
				Description: "插件配置参数(YAML格式)，与values_json二选一。",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("values_json")),
					validator2.Yaml(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"values_json": schema.StringAttribute{
				Optional:    true,
				Description: "插件配置参数(JSON格式)，与values_yaml二选一。",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("values_yaml")),
					validator2.Json(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Computed:    true,
				Description: "命名空间",
			},
		},
	}
}

func (c *ctyunCcsePlugin) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCcsePluginConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	// 创建后检查
	err = c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunCcsePlugin) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcsePluginConfig
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

func (c *ctyunCcsePlugin) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunCcsePluginConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunCcsePluginConfig
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

func (c *ctyunCcsePlugin) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcsePluginConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDelete(ctx, state)
	if err != nil {
		return
	}
	//response.State.RemoveResource(ctx)
}

func (c *ctyunCcsePlugin) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [pluginName],[clusterID],[regionID]
func (c *ctyunCcsePlugin) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunCcsePluginConfig
	var chartName, clusterID, regionID string
	err = terraform_extend.Split(request.ID, &chartName, &clusterID, &regionID)
	if err != nil {
		return
	}
	cfg.ChartName = types.StringValue(chartName)
	cfg.RegionID = types.StringValue(regionID)
	cfg.ClusterID = types.StringValue(clusterID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunCcsePlugin) checkBeforeCreate(ctx context.Context, plan CtyunCcsePluginConfig) (err error) {
	p, err := c.getByChartName(ctx, plan)
	if err != nil {
		if errors.Is(err, common.InvalidReturnObjResultsError) {
			err = nil
		}
		return
	}
	if p.Status == "deployed" {
		err = fmt.Errorf("插件 %s 不可重复安装", plan.ChartName.ValueString())
		return
	}

	return
}

// checkAfterCreate 创建后检查
func (c *ctyunCcsePlugin) checkAfterCreate(ctx context.Context, plan CtyunCcsePluginConfig) (err error) {
	var executeSuccessFlag bool
	var failedCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 30)
	retryer.Start(
		func(currentTime int) bool {
			var plugin *ccse2.CcseListPluginInstancesReturnObjRecordsResponse
			plugin, err = c.getByChartName(ctx, plan)
			if err != nil {
				if errors.Is(err, common.InvalidReturnObjResultsError) {
					return true
				}
				return false
			}
			if plugin.Status == "failed" {
				failedCnt++
			}
			if failedCnt > 1 {
				err = fmt.Errorf("安装失败，请登录控制台查看原因")
				return false
			}
			if plugin.Status != "deployed" {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("插件安装超时，请登录控制台查看原因")
	}
	return
}

// create 创建
func (c *ctyunCcsePlugin) create(ctx context.Context, plan CtyunCcsePluginConfig) (err error) {
	params := &ccse2.CcseDeployPluginInstanceRequest{
		ClusterId:    plan.ClusterID.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		ChartName:    plan.ChartName.ValueString(),
		ChartVersion: plan.ChartVersion.ValueString(),
		Values:       plan.ValuesYaml.ValueString(),
		ValuesJson:   plan.ValuesJson.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCcseApis.CcseDeployPluginInstanceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// getAndMerge 从远端查询
func (c *ctyunCcsePlugin) getAndMerge(ctx context.Context, plan *CtyunCcsePluginConfig) (err error) {
	plugin, err := c.getByChartName(ctx, *plan)
	if err != nil {
		return
	}
	plan.Namespace = types.StringValue(plugin.Namespace)
	plan.ChartName = types.StringValue(plugin.ChartName)
	plan.ChartVersion = types.StringValue(plugin.ChartVersion)
	plan.ClusterID = types.StringValue(plugin.ClusterId)
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plugin.ChartName, plugin.ClusterId, plan.RegionID.ValueString()))

	return
}

// update 更新
func (c *ctyunCcsePlugin) update(ctx context.Context, plan, state *CtyunCcsePluginConfig) (err error) {
	err = c.updateChartVersion(ctx, *plan, *state)
	if err != nil {
		return
	}
	return
}

// updateChartVersion 更新chart_version
func (c *ctyunCcsePlugin) updateChartVersion(ctx context.Context, plan, state CtyunCcsePluginConfig) (err error) {
	if plan.ChartVersion.Equal(state.ChartVersion) {
		return
	}
	params := &ccse2.CcseUpgradePluginInstanceRequest{
		ClusterId:    plan.ClusterID.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		ChartName:    plan.ChartName.ValueString(),
		ChartVersion: plan.ChartVersion.ValueString(),
		Values:       state.ValuesYaml.ValueString(),
		ValuesJson:   state.ValuesJson.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCcseApis.CcseUpgradePluginInstanceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return c.checkAfterChartVersion(ctx, plan)
}

// checkAfterChartVersion 升级插件版本后检查
func (c *ctyunCcsePlugin) checkAfterChartVersion(ctx context.Context, plan CtyunCcsePluginConfig) (err error) {
	var executeSuccessFlag bool
	var failedCnt int
	retryer, _ := business.NewRetryer(time.Second*10, 30)
	retryer.Start(
		func(currentTime int) bool {
			var plugin *ccse2.CcseListPluginInstancesReturnObjRecordsResponse
			plugin, err = c.getByChartName(ctx, plan)
			if err != nil {
				return false
			}

			if plugin.Status == "failed" {
				failedCnt++
			}
			if failedCnt > 1 {
				err = fmt.Errorf("升级插件版本失败")
				return false
			}

			if plugin.Status != "deployed" || plugin.ChartVersion != plan.ChartVersion.ValueString() {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("插件升级超时")
	}
	return
}

// delete 删除
func (c *ctyunCcsePlugin) delete(ctx context.Context, plan CtyunCcsePluginConfig) (err error) {
	params := &ccse2.CcseDeletePluginInstanceRequest{
		ClusterId:    plan.ClusterID.ValueString(),
		RegionId:     plan.RegionID.ValueString(),
		InstanceName: plan.ChartName.ValueString(),
		Namespace:    plan.Namespace.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseDeletePluginInstanceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// checkAfterDelete 删除后检查
func (c *ctyunCcsePlugin) checkAfterDelete(ctx context.Context, plan CtyunCcsePluginConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 30)
	retryer.Start(
		func(currentTime int) bool {
			var plugin *ccse2.CcseListPluginInstancesReturnObjRecordsResponse
			plugin, err = c.getByChartName(ctx, plan)
			if err != nil {
				if errors.Is(err, common.InvalidReturnObjResultsError) {
					err = nil
					executeSuccessFlag = true
				}
				return false
			}
			if plugin.Status != "uninstalled" {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("插件卸载超时")
	}
	return
}

// getByChartName通过插件名称查询
func (c *ctyunCcsePlugin) getByChartName(ctx context.Context, plan CtyunCcsePluginConfig) (plugin *ccse2.CcseListPluginInstancesReturnObjRecordsResponse, err error) {
	params := &ccse2.CcseListPluginInstancesRequest{
		ClusterId: plan.ClusterID.ValueString(),
		RegionId:  plan.RegionID.ValueString(),
		ChartName: plan.ChartName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseListPluginInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if len(resp.ReturnObj.Records) == 0 {
		err = common.InvalidReturnObjResultsError
		return
	}
	plugin = resp.ReturnObj.Records[0]
	return
}
