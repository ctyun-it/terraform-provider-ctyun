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
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunCcseNamespace{}
	_ resource.ResourceWithConfigure   = &ctyunCcseNamespace{}
	_ resource.ResourceWithImportState = &ctyunCcseNamespace{}
)

type ctyunCcseNamespace struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseNamespace() resource.Resource {
	return &ctyunCcseNamespace{}
}

func (c *ctyunCcseNamespace) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_namespace"
}

type CtyunCcseNamespaceConfig struct {
	ID           types.String `tfsdk:"id"`
	ClusterID    types.String `tfsdk:"cluster_id"`
	RegionID     types.String `tfsdk:"region_id"`
	ValuesYaml   types.String `tfsdk:"values_yaml"`
	ActualConfig types.String `tfsdk:"actual_config"`
	Namespace    types.String `tfsdk:"namespace"`
}

func (c *ctyunCcseNamespace) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10083472/10102631`,
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
					stringvalidator.UTF8LengthBetween(32, 32),
				},
			},
			"values_yaml": schema.StringAttribute{
				Optional:    true,
				Description: "配置参数(YAML格式)，支持更新",
				Validators: []validator.String{
					validator2.Yaml("apiVersion", "kind", "metadata.name"),
				},
			},
			"actual_config": schema.StringAttribute{
				Computed:    true,
				Description: "实际配置(YAML格式)",
			},
			"namespace": schema.StringAttribute{
				Computed:    true,
				Description: "命名空间",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ctyunCcseNamespace) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCcseNamespaceConfig
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
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunCcseNamespace) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseNamespaceConfig
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

func (c *ctyunCcseNamespace) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunCcseNamespaceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunCcseNamespaceConfig
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

func (c *ctyunCcseNamespace) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseNamespaceConfig
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
}

func (c *ctyunCcseNamespace) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunCcseNamespace) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [namespace],[clusterID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunCcseNamespaceConfig
	var namespace, clusterID, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &namespace, &clusterID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &namespace, &clusterID, &regionID)
		if err != nil {
			return
		}
	}

	if namespace == "" {
		err = fmt.Errorf("namespace不能为空")
		return
	}
	if clusterID == "" {
		err = fmt.Errorf("clusterID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	cfg.Namespace = types.StringValue(namespace)
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
func (c *ctyunCcseNamespace) checkBeforeCreate(ctx context.Context, plan CtyunCcseNamespaceConfig) (err error) {
	_, err = utils.ParseYamlValue(plan.ValuesYaml.ValueString(), "metadata.name")
	return
}

// create 创建
func (c *ctyunCcseNamespace) create(ctx context.Context, plan CtyunCcseNamespaceConfig) (err error) {
	params := &ccse2.CcseCreateNamespaceV2P2Request{
		ClusterName:         plan.ClusterID.ValueString(),
		RegionId:            plan.RegionID.ValueString(),
		TextPlainDataString: plan.ValuesYaml.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCcseApis.CcseCreateNamespaceV2P2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunCcseNamespace) getAndMerge(ctx context.Context, plan *CtyunCcseNamespaceConfig) (err error) {
	if plan.Namespace.ValueString() == "" {
		var namespace interface{}
		namespace, err = utils.ParseYamlValue(plan.ValuesYaml.ValueString(), "metadata.name")
		if err != nil {
			return
		}
		s, _ := namespace.(string)
		plan.Namespace = types.StringValue(s)
	}
	config, err := c.getNamespace(ctx, *plan)
	if err != nil {
		return
	}
	plan.ActualConfig = types.StringValue(config)
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.Namespace.ValueString(), plan.ClusterID.ValueString(), plan.RegionID.ValueString()))
	return
}

// update 更新
func (c *ctyunCcseNamespace) update(ctx context.Context, plan, state *CtyunCcseNamespaceConfig) (err error) {
	params := &ccse2.CcseUpdateNamespaceV2P2Request{
		ClusterName:         state.ClusterID.ValueString(),
		NamespaceName:       state.Namespace.ValueString(),
		RegionId:            state.RegionID.ValueString(),
		TextPlainDataString: plan.ValuesYaml.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseUpdateNamespaceV2P2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	state.ValuesYaml = plan.ValuesYaml
	return
}

// delete 删除
func (c *ctyunCcseNamespace) delete(ctx context.Context, plan CtyunCcseNamespaceConfig) (err error) {
	params := &ccse2.CcseDeleteNamespaceV2P2Request{
		ClusterName:   plan.ClusterID.ValueString(),
		NamespaceName: plan.Namespace.ValueString(),
		RegionId:      plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteNamespaceV2P2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// checkAfterDelete 删除后检查
func (c *ctyunCcseNamespace) checkAfterDelete(ctx context.Context, plan CtyunCcseNamespaceConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*5, 30)
	retryer.Start(
		func(currentTime int) bool {
			_, err = c.getNamespace(ctx, plan)
			if err != nil {
				if errors.Is(err, common.ResourceNotExistError) {
					err = nil
					executeSuccessFlag = true
				}
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("命名空间删除超时")
	}
	return
}

// getNamespace 查询命名空间
func (c *ctyunCcseNamespace) getNamespace(ctx context.Context, plan CtyunCcseNamespaceConfig) (script string, err error) {
	params := &ccse2.CcseGetNamespaceV2P2Request{
		ClusterName:   plan.ClusterID.ValueString(),
		NamespaceName: plan.Namespace.ValueString(),
		RegionId:      plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetNamespaceV2P2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		if strings.Contains(resp.Message, "不存在") || strings.Contains(resp.Message, "resource example not found") {
			err = common.ResourceNotExistError
		} else {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
		}
		return
	}
	script = resp.ReturnObj
	return
}
