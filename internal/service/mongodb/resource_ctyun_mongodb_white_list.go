package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunMongodbWhiteList{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbWhiteList{}
	_ resource.ResourceWithImportState = &CtyunMongodbWhiteList{}
)

func NewCtyunMongodbWhiteList() resource.Resource {
	return &CtyunMongodbWhiteList{}
}

type CtyunMongodbWhiteList struct {
	meta *common.CtyunMetadata
}

type CtyunMongodbWhiteListConfig struct {
	ID            types.String `tfsdk:"id"`
	InstanceID    types.String `tfsdk:"instance_id"`
	RegionID      types.String `tfsdk:"region_id"`
	ProjectID     types.String `tfsdk:"project_id"`
	IpList        types.Set    `tfsdk:"ip_list"`
	GroupName     types.String `tfsdk:"group_name"`      // 白名单分组名
	IpType        types.String `tfsdk:"ip_type"`         // 白名单类型
	WhiteListType types.String `tfsdk:"white_list_type"` // 白名单分组类型
	WhiteListId   types.Int32  `tfsdk:"white_list_id"`   // 白名单分组Id

}

func (c *CtyunMongodbWhiteList) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_white_list"
}

func (c *CtyunMongodbWhiteList) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 instance_id,ip_whitelist_name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MongoDB实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"group_name": schema.StringAttribute{
				Required:    true,
				Description: "白名单分组名 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"ip_type": schema.StringAttribute{
				Required:    true,
				Description: "ip类型 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"white_list_type": schema.StringAttribute{
				Required:    true,
				Description: "白名单列表类型  支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"white_list_id": schema.Int32Attribute{
				Computed:    true,
				Description: "白名单分组Id",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"ip_list": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "IP列表 支持更新",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (c *CtyunMongodbWhiteList) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbWhiteList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbWhiteListConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBefore(ctx, plan)
	if err != nil {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunMongodbWhiteList) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMongodbWhiteListConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunMongodbWhiteList) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbWhiteListConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.checkBefore(ctx, plan)
	if err != nil {
		return
	}
	err = c.update(ctx, &plan)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunMongodbWhiteList) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMongodbWhiteListConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err = c.checkBefore(ctx, state)
	if err != nil {
		return
	}
	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunMongodbWhiteList) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunMongodbWhiteListConfig
	var instanceID, groupName, regionId string
	err = terraform_extend.Split(req.ID, &instanceID, &groupName, &regionId)
	if err != nil {
		return
	}
	cfg.InstanceID = types.StringValue(instanceID)
	cfg.GroupName = types.StringValue(groupName)
	cfg.RegionID = types.StringValue(regionId)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, cfg)...)
}

func (c *CtyunMongodbWhiteList) checkBefore(ctx context.Context, state CtyunMongodbWhiteListConfig, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return
	}

	listHeader := &mongodb.MongodbGetListHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		listHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}

	result := retryer.Start(
		func(currentTime int) bool {

			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: state.InstanceID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeader.ProjectID = state.ProjectID.ValueStringPointer()
			}
			detailResp, err3 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err3 != nil {
				err = err3
				return false
			} else if detailResp.StatusCode != 800 {
				if detailResp.Message != nil && strings.Contains(*detailResp.Message, "实例未正常运行,请稍后再试") {
					return true // 继续重试
				}
				err = fmt.Errorf("API return error. Message: %s", *detailResp.Message)
				return false
			} else if detailResp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			if detailResp.ReturnObj.ProdRunningStatus == business.MongodbRunningStatusStarted {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，实例仍未运行成功！")
	}
	return
}

func (c *CtyunMongodbWhiteList) create(ctx context.Context, plan *CtyunMongodbWhiteListConfig) (err error) {
	// 将types.Set转换为字符串切片，然后转换为JSON数组字符串
	var ipList []string
	diag := plan.IpList.ElementsAs(ctx, &ipList, false)
	if diag.HasError() {
		return
	}
	ipListBytes, err := json.Marshal(ipList)
	if err != nil {
		return fmt.Errorf("failed to marshal ip_list to JSON: %w", err)
	}
	ipListValue := string(ipListBytes)

	// 创建白名单分组
	createReq := &mongodb.MongodbCreateIpWhitelistRequest{
		ProdInstId:    plan.InstanceID.ValueString(),
		GroupName:     plan.GroupName.ValueString(),
		IpType:        plan.IpType.ValueString(),
		IpList:        ipListValue,
		WhiteListType: plan.WhiteListType.ValueString(),
	}

	headers := &mongodb.MongodbCreateIpWhitelistRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "创建MongoDB白名单分组", map[string]interface{}{
		"instance_id": plan.InstanceID.ValueString(),
	})

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbCreateIpWhitelistApi.Do(ctx, c.meta.Credential, createReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.InstanceID.ValueString(), plan.GroupName.ValueString()))

	return
}
func (c *CtyunMongodbWhiteList) getAndMerge(ctx context.Context, plan *CtyunMongodbWhiteListConfig) (err error) {
	// 查询白名单列表
	describeReq := &mongodb.MongodbDescribeIpWhitelistRequest{
		ProdInstId: plan.InstanceID.ValueString(),
	}

	headers := &mongodb.MongodbDescribeIpWhitelistRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbDescribeIpWhitelistApi.Do(ctx, c.meta.Credential, describeReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}

	// 查找对应的白名单分组
	var found bool
	for _, group := range resp.ReturnObj.WhitelistGroup {
		if group.GroupName == plan.GroupName.ValueString() {
			// 更新状态
			plan.ID = types.StringValue(fmt.Sprintf("%d,%s", group.Id, plan.InstanceID.ValueString()))
			plan.WhiteListId = types.Int32Value(group.Id)

			// 设置IP列表，API返回的是逗号分隔的字符串
			var ipList []string
			if group.IpList != "" {
				ipList = strings.Split(group.IpList, ",")
			}

			ipListSet, diags := types.SetValueFrom(ctx, types.StringType, ipList)
			if diags.HasError() {
				return fmt.Errorf("failed to set ip_list value")
			}
			plan.IpList = ipListSet
			plan.IpType = types.StringValue(group.IpType)
			plan.WhiteListType = types.StringValue(fmt.Sprintf("%d", group.WhiteListType))

			found = true
		}
	}
	if !found {
		return fmt.Errorf("API return error. Message: mongodb white list not found")
	}
	return
}
func (c *CtyunMongodbWhiteList) update(ctx context.Context, plan *CtyunMongodbWhiteListConfig) (err error) {
	// 更新白名单分组
	// 将types.Set转换为字符串切片，然后转换为JSON数组字符串
	var ipList []string
	diag := plan.IpList.ElementsAs(ctx, &ipList, false)
	if diag.HasError() {
		return
	}
	ipListBytes, err := json.Marshal(ipList)
	if err != nil {
		return fmt.Errorf("failed to marshal ip_list to JSON: %w", err)
	}
	ipListValue := string(ipListBytes)

	updateReq := &mongodb.MongodbUpdateIpWhitelistRequest{
		ProdInstId:    plan.InstanceID.ValueString(),
		GroupName:     plan.GroupName.ValueString(),
		IpType:        plan.IpType.ValueString(),
		IpList:        ipListValue,
		WhiteListType: plan.WhiteListType.ValueString(),
		WhiteListId:   fmt.Sprintf("%d", plan.WhiteListId.ValueInt32()),
	}

	headers := &mongodb.MongodbUpdateIpWhitelistRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbUpdateIpWhitelistApi.Do(ctx, c.meta.Credential, updateReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	return
}
func (c *CtyunMongodbWhiteList) delete(ctx context.Context, state *CtyunMongodbWhiteListConfig) (err error) {
	// 删除白名单分组
	deleteReq := &mongodb.MongodbDeleteIpWhitelistRequest{
		ProdInstId:  state.InstanceID.ValueString(),
		WhiteListId: fmt.Sprintf("%d", state.WhiteListId.ValueInt32()),
	}

	headers := &mongodb.MongodbDeleteIpWhitelistRequestHeaders{
		RegionID: state.RegionID.ValueString(),
	}
	if !state.ProjectID.IsNull() {
		headers.ProjectID = state.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "删除MongoDB白名单分组", map[string]interface{}{
		"instance_id":       state.InstanceID.ValueString(),
		"ip_whitelist_name": state.GroupName.ValueString(),
	})

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbDeleteIpWhitelistApi.Do(ctx, c.meta.Credential, deleteReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	return
}
