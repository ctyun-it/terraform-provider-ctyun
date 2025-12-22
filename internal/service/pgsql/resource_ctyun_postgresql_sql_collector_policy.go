package pgsql

//
//import (
//	"context"
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
//	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
//	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
//	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
//	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"github.com/hashicorp/terraform-plugin-log/tflog"
//	"strings"
//)
//
//var (
//	_ resource.Resource                = &CtyunPostgresqlCollectorPolicy{}
//	_ resource.ResourceWithConfigure   = &CtyunPostgresqlCollectorPolicy{}
//	_ resource.ResourceWithImportState = &CtyunPostgresqlCollectorPolicy{}
//)
//
//func NewCtyunPostgresqlCollectorPolicy() resource.Resource {
//	return &CtyunPostgresqlCollectorPolicy{}
//}
//
//type CtyunPostgresqlCollectorPolicy struct {
//	meta *common.CtyunMetadata
//}
//
//type CtyunPostgresqlCollectorPolicyConfig struct {
//	ID                 types.String `tfsdk:"id"`
//	InstanceID         types.String `tfsdk:"instance_id"`
//	RegionID           types.String `tfsdk:"region_id"`
//	ProjectID          types.String `tfsdk:"project_id"`
//	SqlCollectorStatus types.String `tfsdk:"sql_collector_status"`
//	LogInterval        types.Int32  `tfsdk:"log_interval"`
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
//	resp.TypeName = req.ProviderTypeName + "_postgresql_sql_collector_policy"
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
//	resp.Schema = schema.Schema{
//		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10161917`,
//		Attributes: map[string]schema.Attribute{
//			"id": schema.StringAttribute{
//				Computed:    true,
//				Description: "资源唯一标识，与实例ID相同",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"instance_id": schema.StringAttribute{
//				Required:    true,
//				Description: "PostgreSQL实例ID",
//				Validators: []validator.String{
//					stringvalidator.LengthAtLeast(1),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//			},
//			"region_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
//				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//				Validators: []validator.String{
//					stringvalidator.UTF8LengthAtLeast(1),
//				},
//			},
//			// 项目相关
//			"project_id": schema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.RequiresReplace(),
//				},
//				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
//				Validators: []validator.String{
//					validator2.Project(),
//				},
//			},
//			"sql_collector_status": schema.StringAttribute{
//				Required:    true,
//				Description: "SQL审计状态，支持更新。取值：enable | disabled",
//				Validators: []validator.String{
//					stringvalidator.OneOf("enable", "disabled"),
//				},
//			},
//			"log_interval": schema.Int32Attribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "定时收集SQL日志的间隔，支持更新。单位：分钟，默认5分钟，取值：5，10，30，60",
//				Default:     int32default.StaticInt32(5),
//				Validators: []validator.Int32{
//					int32validator.OneOf(5, 10, 30, 60),
//				},
//			},
//		},
//	}
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
//	if req.ProviderData == nil {
//		return
//	}
//	meta := req.ProviderData.(*common.CtyunMetadata)
//	c.meta = meta
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			resp.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var plan CtyunPostgresqlCollectorPolicyConfig
//	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	err = c.create(ctx, &plan)
//	if err != nil {
//		return
//	}
//
//	plan.ID = plan.InstanceID
//
//	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			resp.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunPostgresqlCollectorPolicyConfig
//	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//	err = c.getAndMerge(ctx, &state)
//	if err != nil {
//		return
//	}
//	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			resp.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var plan CtyunPostgresqlCollectorPolicyConfig
//	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	err = c.update(ctx, &plan)
//	if err != nil {
//		return
//	}
//	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			resp.Diagnostics.AddError(err.Error(), err.Error())
//		}
//	}()
//	var state CtyunPostgresqlCollectorPolicyConfig
//	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	err = c.delete(ctx, &state)
//	if err != nil {
//		return
//	}
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
//	var err error
//	defer func() {
//		if err != nil {
//			title := "导入失败：" + err.Error()
//			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[projectID],[regionID]"
//			response.Diagnostics.AddError(title, detail)
//		}
//	}()
//	var config CtyunPostgresqlCollectorPolicyConfig
//	var instanceId, regionId, projectId string
//	if strings.Count(request.ID, common.ImportSeparator) < 1 {
//		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
//		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
//		instanceId = request.ID
//	} else {
//		err = terraform_extend.Split(request.ID, &instanceId, &projectId, &regionId)
//		if err != nil {
//			return
//		}
//	}
//	if instanceId == "" {
//		err = fmt.Errorf("ID不能为空")
//		return
//	}
//	if regionId == "" {
//		err = fmt.Errorf("regionID不能为空")
//		return
//	}
//	config.InstanceID = types.StringValue(instanceId)
//	config.RegionID = types.StringValue(regionId)
//	config.ProjectID = types.StringValue(projectId)
//	err = c.getAndMerge(ctx, &config)
//	if err != nil {
//		return
//	}
//	response.Diagnostics.Append(response.State.Set(ctx, config)...)
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) create(ctx context.Context, plan *CtyunPostgresqlCollectorPolicyConfig) (err error) {
//	request := &pgsql.PostgresqlCollectorPolicyRequest{
//		ProdInstId:         plan.InstanceID.ValueString(),
//		SqlCollectorStatus: plan.SqlCollectorStatus.ValueString(),
//	}
//
//	if !plan.LogInterval.IsNull() && !plan.LogInterval.IsUnknown() {
//		logInterval := plan.LogInterval.ValueInt32()
//		request.LogInterval = &logInterval
//	}
//
//	header := &pgsql.PostgresqlCollectorPolicyRequestHeader{
//		RegionId: plan.RegionID.ValueStringPointer(),
//	}
//	if !plan.ProjectID.IsNull() {
//		header.ProjectId = plan.ProjectID.ValueStringPointer()
//	}
//
//	resp, err := c.meta.Apis.SdkCtPgsqlApis.PostgresqlCollectorPolicyApi.Do(ctx, c.meta.Credential, request, header)
//	if err != nil {
//		return err
//	} else if resp.StatusCode != common.NormalStatusCode {
//		return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
//	}
//
//	return nil
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) getAndMerge(ctx context.Context, state *CtyunPostgresqlCollectorPolicyConfig) (err error) {
//	// 使用查询接口获取当前状态
//	request := &pgsql.PostgresqlGetCollectorPolicyRequest{
//		ProdInstId: state.InstanceID.ValueString(),
//	}
//
//	header := &pgsql.PostgresqlGetCollectorPolicyRequestHeader{
//		RegionId: state.RegionID.ValueStringPointer(),
//	}
//	if !state.ProjectID.IsNull() {
//		header.ProjectId = state.ProjectID.ValueStringPointer()
//	}
//
//	tflog.Info(ctx, "查询PostgreSQL SQL审计策略", map[string]interface{}{
//		"instance_id": state.InstanceID.ValueString(),
//	})
//
//	resp, err := c.meta.Apis.SdkCtPgsqlApis.PostgresqlGetCollectorPolicyApi.Do(ctx, c.meta.Credential, request, header)
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode != common.NormalStatusCode {
//		return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
//	}
//
//	if resp.ReturnObj == nil {
//		return fmt.Errorf("API return empty response object")
//	}
//
//	// 更新状态
//	state.SqlCollectorStatus = types.StringValue(resp.ReturnObj.SqlCollectorStatus)
//	state.LogInterval = types.Int32Value(resp.ReturnObj.LogInterval)
//	state.ID = types.StringValue(resp.ReturnObj.ProdInstId)
//
//	return nil
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) update(ctx context.Context, plan *CtyunPostgresqlCollectorPolicyConfig) (err error) {
//	// 更新操作与创建操作相同，都是调用设置策略接口
//	request := &pgsql.PostgresqlCollectorPolicyRequest{
//		ProdInstId:         plan.InstanceID.ValueString(),
//		SqlCollectorStatus: plan.SqlCollectorStatus.ValueString(),
//	}
//
//	if !plan.LogInterval.IsNull() && !plan.LogInterval.IsUnknown() {
//		logInterval := plan.LogInterval.ValueInt32()
//		request.LogInterval = &logInterval
//	}
//
//	header := &pgsql.PostgresqlCollectorPolicyRequestHeader{
//		RegionId: plan.RegionID.ValueStringPointer(),
//	}
//	if !plan.ProjectID.IsNull() {
//		header.ProjectId = plan.ProjectID.ValueStringPointer()
//	}
//
//	tflog.Info(ctx, "更新PostgreSQL SQL审计策略", map[string]interface{}{
//		"instance_id": plan.InstanceID.ValueString(),
//	})
//
//	resp, err := c.meta.Apis.SdkCtPgsqlApis.PostgresqlCollectorPolicyApi.Do(ctx, c.meta.Credential, request, header)
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode != common.NormalStatusCode {
//		return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
//	}
//
//	return nil
//}
//
//func (c *CtyunPostgresqlCollectorPolicy) delete(ctx context.Context, state *CtyunPostgresqlCollectorPolicyConfig) (err error) {
//	// 删除操作实际上是将SQL审计状态设置为disabled
//	state.SqlCollectorStatus = types.StringValue("disabled")
//
//	request := &pgsql.PostgresqlCollectorPolicyRequest{
//		ProdInstId:         state.InstanceID.ValueString(),
//		SqlCollectorStatus: state.SqlCollectorStatus.ValueString(),
//	}
//
//	header := &pgsql.PostgresqlCollectorPolicyRequestHeader{
//		RegionId: state.RegionID.ValueStringPointer(),
//	}
//	if !state.ProjectID.IsNull() {
//		header.ProjectId = state.ProjectID.ValueStringPointer()
//	}
//
//	tflog.Info(ctx, "禁用PostgreSQL SQL审计策略", map[string]interface{}{
//		"instance_id": state.InstanceID.ValueString(),
//	})
//
//	resp, err := c.meta.Apis.SdkCtPgsqlApis.PostgresqlCollectorPolicyApi.Do(ctx, c.meta.Credential, request, header)
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode != common.NormalStatusCode {
//		return fmt.Errorf("API return error. Status Code: %d, Message: %s, Error: %s", resp.StatusCode, resp.Message, resp.Error)
//	}
//
//	return nil
//}
