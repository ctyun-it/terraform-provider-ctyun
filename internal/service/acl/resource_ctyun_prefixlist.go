package acl

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunPrefix{}
	_ resource.ResourceWithConfigure   = &CtyunPrefix{}
	_ resource.ResourceWithImportState = &CtyunPrefix{}
)

type CtyunPrefix struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunPrefix) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_prefix_list"
}

func (c *CtyunPrefix) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunPrefix() resource.Resource {
	return &CtyunPrefix{}
}

func (c *CtyunPrefix) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPrefixConfig
	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunPrefix) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10298321",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "前缀列表名称，支持更新。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
					validator2.AclName(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "前缀列表描述，支持更新。支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:\"{},./;'[\\]·！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
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
			"limit": schema.Int32Attribute{
				Required:    true,
				Description: "前缀列表支持的最大条目容量，创建后将无法修改,限制1-200条，具体以账户配额为准,不能小于前缀列表规则个数",
				Validators: []validator.Int32{
					int32validator.Between(1, 200),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"address_type": schema.StringAttribute{
				Required:    true,
				Description: "地址类型，取值范围：ipv4，ipv6",
				Validators: []validator.String{
					stringvalidator.OneOf(business.PrefixAddressTypeIpv4, business.PrefixAddressTypeIpv6),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prefix_list_rules": schema.SetNestedAttribute{
				Required:    true,
				Description: "前缀规则列表",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cidr": schema.StringAttribute{
							Required:    true,
							Description: "CIDR格式的IP地址段。例如：192.168.0.0/16",
							Validators: []validator.String{
								validator2.Cidr(),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "前缀规则描述。支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:\"{},./;'[\\]·！@#￥%……&*（） —— -+={}\\《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128",
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 128),
								validator2.Desc(),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "前缀列表id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunPrefix) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPrefixConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
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
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrefix) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPrefixConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrefix) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPrefixConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPrefixConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrefix) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPrefixConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunPrefix) getAndMerge(ctx context.Context, config *CtyunPrefixConfig) error {
	params := &ctvpc.CtvpcPrefixlistShowRequest{
		RegionID:     config.RegionID.ValueString(),
		PrefixListID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistShowApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取prefix失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	returnObj := resp.ReturnObj
	config.Name = types.StringValue(*returnObj.Name)
	config.Limit = types.Int32Value(returnObj.Limit)
	config.AddressType = types.StringValue(business.PrefixAddressTyperRevMap[returnObj.AddressType])
	config.Description = types.StringValue(*returnObj.Description)
	config.CreateTime = types.StringValue(utils.ConvertToUTCZ(utils.Layout1, utils.SecString(returnObj.CreatedAt)))
	config.UpdateTime = types.StringValue(utils.ConvertToUTCZ(utils.Layout1, utils.SecString(returnObj.UpdatedAt)))
	rules := returnObj.PrefixListRules
	var prefixRules []CtyunPrefixModel
	for _, ruleItem := range rules {
		var rule CtyunPrefixModel
		rule.Cidr = types.StringValue(*ruleItem.Cidr)
		rule.Description = types.StringValue(*ruleItem.Description)
		prefixRules = append(prefixRules, rule)
	}
	rulesTmp, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(CtyunPrefixModel{}), prefixRules)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.PrefixListRules = rulesTmp
	return nil
}

func (c *CtyunPrefix) create(ctx context.Context, config *CtyunPrefixConfig) error {
	params := &ctvpc.CtvpcPrefixlistCreateRequest{
		RegionID:        config.RegionID.ValueString(),
		Name:            config.Name.ValueString(),
		Limit:           config.Limit.ValueInt32(),
		AddressType:     business.PrefixAddressTypeMap[config.AddressType.ValueString()],
		PrefixListRules: nil,
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.Description = config.Description.ValueStringPointer()
	}
	var prefixListRules []*ctvpc.CtvpcPrefixlistCreatePrefixListRulesRequest
	var prefixes []CtyunPrefixModel
	diags := config.PrefixListRules.ElementsAs(ctx, &prefixes, true)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	for _, rule := range prefixes {
		var prefix ctvpc.CtvpcPrefixlistCreatePrefixListRulesRequest
		prefix.Cidr = rule.Cidr.ValueString()
		prefix.Description = rule.Description.ValueStringPointer()
		prefixListRules = append(prefixListRules, &prefix)
	}
	params.PrefixListRules = prefixListRules
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistCreateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建prefix失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.PrefixListID)
	return nil
}

func (c *CtyunPrefix) update(ctx context.Context, state *CtyunPrefixConfig, plan *CtyunPrefixConfig) error {
	params := &ctvpc.CtvpcPrefixlistUpdateRequest{
		RegionID:     state.RegionID.ValueString(),
		PrefixListID: state.ID.ValueString(),
		Name:         plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && !plan.Description.Equal(state.Description) {
		params.Description = plan.Description.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新prefix失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunPrefix) delete(ctx context.Context, config CtyunPrefixConfig) error {
	params := &ctvpc.CtvpcPrefixlistDeleteRequest{
		RegionID:     config.RegionID.ValueString(),
		PrefixListID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistDeleteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除prefix失败（prefixlist id=%s），接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

type CtyunPrefixModel struct {
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
}
type CtyunPrefixConfig struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	RegionID        types.String `tfsdk:"region_id"`
	Limit           types.Int32  `tfsdk:"limit"`
	AddressType     types.String `tfsdk:"address_type"`
	PrefixListRules types.Set    `tfsdk:"prefix_list_rules"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdateTime      types.String `tfsdk:"update_time"`
}
