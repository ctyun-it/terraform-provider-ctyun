package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunMysqlWhiteList{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlWhiteList{}
	_ resource.ResourceWithImportState = &CtyunMysqlWhiteList{}
)

type CtyunMysqlWhiteList struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlWhiteList) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10133794`,
		Attributes: map[string]schema.Attribute{
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
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"prod_inst_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql实例id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"group_name": schema.StringAttribute{
				Required:    true,
				Description: "白名单分组名（必须以小写字母开头，且必须以小写字母或数字结尾，可包含数字或下划线，不含其他特殊字符)",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z]([a-z0-9_]*[a-z0-9])?$"), "白名单分组名不符合要求"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_white_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "白名单ip列表，举例：['192.168.0.1', '192.168.0.*'],指定IP地址192.168.0.1：表示允许192.168.0.1的IP地址访问实例。 指定IP地址192.168.0.*：表示允许从192.168.0.1到192.168.0.255的IP地址访问实例。",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"group_white_list_count": schema.Int32Attribute{
				Computed:    true,
				Description: "白名单分组组内数量",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间",
			},
			"access_machine_type": schema.StringAttribute{
				Computed:    true,
				Description: "访问类型",
			},
		},
	}
}

func (c *CtyunMysqlWhiteList) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlWhiteListConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkStatus(ctx, plan)
	if err != nil {
		return
	}
	// 开始创建
	err = c.CreateMysqlAccessWhiteList(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlAccessWhiteList(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlWhiteList) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlWhiteListConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlAccessWhiteList(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlWhiteList) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlWhiteListConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlWhiteListConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkStatus(ctx, plan)
	if err != nil {
		return
	}
	err = c.updateMysqlWhiteList(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlAccessWhiteList(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlWhiteList) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunMysqlWhiteListConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkStatus(ctx, state)
	if err != nil {
		return
	}
	params := &mysql.TeledbDeleteAccessWhiteListRequest{
		OuterProdInstID: state.ProdInstID.ValueString(),
		GroupName:       state.GroupName.ValueString(),
	}
	header := &mysql.TeledbDeleteAccessWhiteListRequestHeader{
		InstID:   state.ProdInstID.ValueStringPointer(),
		RegionID: state.RegionID.ValueString(),
	}
	if !state.ProjectID.IsNull() && !state.ProjectID.IsUnknown() {
		header.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbDeleteAccessWhiteList.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("删除mysql白名单过程中，response返回为空, 请稍后再试！")
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

func (c *CtyunMysqlWhiteList) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_white_list"
}
func NewCtyunMysqlWhiteList() resource.Resource {
	return &CtyunMysqlWhiteList{}
}

func (c *CtyunMysqlWhiteList) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// todo
}

func (c *CtyunMysqlWhiteList) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlWhiteList) CreateMysqlAccessWhiteList(ctx context.Context, config *CtyunMysqlWhiteListConfig) (err error) {
	params := &mysql.TeledbCreateAccessWhiteListRequest{
		OuterProdInstID: config.ProdInstID.ValueString(),
		GroupName:       config.GroupName.ValueString(),
	}
	var groupWhiteList []string
	diags := config.GroupWhiteList.ElementsAs(ctx, &groupWhiteList, false)
	if diags.HasError() {
		return
	}
	params.GroupWhiteList = groupWhiteList
	header := &mysql.TeledbCreateAccessWhiteListRequestHeader{
		InstID:   config.ProdInstID.ValueStringPointer(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateAccessWhiteList.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("创建mysql白名单过程中，response返回为空, 请稍后再试！")
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}

	return

}

func (c *CtyunMysqlWhiteList) getAndMergeMysqlAccessWhiteList(ctx context.Context, config *CtyunMysqlWhiteListConfig) (err error) {
	params := mysql.TeledbGetAccessWhiteListRequest{
		OuterProdInstID: config.ProdInstID.ValueString(),
	}
	header := mysql.TeledbGetAccessWhiteListRequestHeader{
		InstID:   config.ProdInstID.ValueStringPointer(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetAccessWhiteList.Do(ctx, c.meta.Credential, &params, &header)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("查询mysql白名单过程中，response返回为空, 请稍后再试！")
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, whileListInfo := range resp.ReturnObj {
		if whileListInfo.GroupName == config.GroupName.ValueString() {
			config.GroupWhiteListCount = types.Int32Value(whileListInfo.GroupWhiteListCount)
			config.AccessMachineType = types.StringValue(whileListInfo.AccessMachineType)
			config.GroupName = types.StringValue(whileListInfo.GroupName)
			config.CreatedTime = types.StringValue(fmt.Sprintf("%d", whileListInfo.CreateTime))
			config.UpdatedTime = types.StringValue(fmt.Sprintf("%d", whileListInfo.UpdateTime))
			groupWhiteList, diags := types.SetValueFrom(ctx, types.StringType, whileListInfo.WhiteList)
			if diags.HasError() {
				return
			}
			config.GroupWhiteList = groupWhiteList
		}
	}
	return
}

func (c *CtyunMysqlWhiteList) updateMysqlWhiteList(ctx context.Context, state *CtyunMysqlWhiteListConfig, plan *CtyunMysqlWhiteListConfig) (err error) {
	params := &mysql.TeledbUpdateAccessWhiteListRequest{
		OuterProdInstID: state.ProdInstID.ValueString(),
		GroupName:       state.GroupName.ValueString(),
	}
	var groupWhiteList []string
	diags := plan.GroupWhiteList.ElementsAs(ctx, &groupWhiteList, false)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return
	}
	params.GroupWhiteList = groupWhiteList
	header := &mysql.TeledbUpdateAccessWhiteListRequestHeader{
		InstID:   state.ProdInstID.ValueStringPointer(),
		RegionID: state.RegionID.ValueString(),
	}
	if !state.ProjectID.IsNull() && !state.ProjectID.IsUnknown() {
		header.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateAccessWhiteList.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("更新mysql白名单过程中，response返回为空, 请稍后再试！")
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}

	return
}

type CtyunMysqlWhiteListConfig struct {
	ProjectID           types.String `tfsdk:"project_id"`
	RegionID            types.String `tfsdk:"region_id"`
	ProdInstID          types.String `tfsdk:"prod_inst_id"`
	GroupName           types.String `tfsdk:"group_name"`
	GroupWhiteList      types.Set    `tfsdk:"group_white_list"`
	GroupWhiteListCount types.Int32  `tfsdk:"group_white_list_count"`
	CreatedTime         types.String `tfsdk:"created_time"`
	UpdatedTime         types.String `tfsdk:"updated_time"`
	AccessMachineType   types.String `tfsdk:"access_machine_type"` // 访问类型
}

// checkStatus 数据库状态为running
func (c *CtyunMysqlWhiteList) checkStatus(ctx context.Context, state CtyunMysqlWhiteListConfig) (err error) {
	retryer, err := business.NewRetryer(time.Second*30, 10)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			var resp *mysql.DetailRespReturnObj
			resp, err = c.getDetail(ctx, state)
			if err != nil {
				return false
			}
			// 资源不正常，继续轮询
			if resp.ProdRunningStatus != 0 {
				return true
			}
			return false
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("资源状态异常")
	}
	return
}

func (c *CtyunMysqlWhiteList) getDetail(ctx context.Context, state CtyunMysqlWhiteListConfig) (instance *mysql.DetailRespReturnObj, err error) {
	// 获取实例详情
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: state.ProdInstID.ValueString(),
	}
	detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   state.ProdInstID.ValueString(),
		RegionID: state.RegionID.ValueString(),
	}
	if state.ProjectID.ValueString() != "" {
		detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = resp.ReturnObj
	return
}
