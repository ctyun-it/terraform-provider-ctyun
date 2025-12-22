package acl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunPrefixLists{}
	_ datasource.DataSourceWithConfigure = &CtyunPrefixLists{}
)

type CtyunPrefixLists struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrefixLists() datasource.DataSource {
	return &CtyunPrefixLists{}
}
func (c *CtyunPrefixLists) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPrefixLists) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_prefix_lists"
}

func (c *CtyunPrefixLists) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10298321",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "前缀列表ID，精确查询",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "查询内容（支持模糊查询）",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "分页页码，默认为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认为10，取值范围：1~50",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"prefix_lists": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "前缀列表ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "前缀列表名称",
						},
						"limit": schema.Int32Attribute{
							Computed:    true,
							Description: "最大条目限制",
						},
						"address_type": schema.StringAttribute{
							Computed:    true,
							Description: "地址类型（ipv4/ipv6）",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "前缀列表描述",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
						"prefix_list_rules": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"prefix_list_rule_id": schema.StringAttribute{
										Computed:    true,
										Description: "前缀列表规则ID",
									},
									"cidr": schema.StringAttribute{
										Computed:    true,
										Description: "CIDR地址块",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "规则描述",
									},
								},
							},
							Description: "前缀列表规则",
						},
					},
				},
				Description: "前缀列表",
			},
		},
	}
}

func (c *CtyunPrefixLists) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPrefixListsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)
	params := &ctvpc.CtvpcPrefixlistQueryRequest{
		RegionID: config.RegionID.ValueString(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.QueryContent.IsNull() {
		params.QueryContent = config.QueryContent.ValueStringPointer()
	}
	if !config.ID.IsNull() {
		params.PrefixListID = config.ID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcPrefixlistQueryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("prefix list查询失败，接口返回nil，请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var prefixLists []CtyunPrefixListModel
	for _, prefixItem := range resp.ReturnObj.PrefixList {
		var prefix CtyunPrefixListModel
		prefix.ID = types.StringValue(*prefixItem.PrefixListID)
		prefix.Name = types.StringValue(*prefixItem.Name)
		prefix.Limit = types.Int32Value(prefixItem.Limit)
		prefix.AddressType = types.StringValue(business.PrefixAddressTyperRevMap[prefixItem.AddressType])
		prefix.Description = types.StringValue(*prefixItem.Description)
		prefix.CreateTime = types.StringValue(utils.ConvertToUTCZ(utils.Layout1, utils.SecString(prefixItem.CreatedAt)))
		prefix.UpdateTime = types.StringValue(utils.ConvertToUTCZ(utils.Layout1, utils.SecString(prefixItem.UpdatedAt)))
		for _, rule := range prefixItem.PrefixListRules {
			var ruleModel CtyunPrefixRuleModel
			ruleModel.PrefixListRuleID = types.StringValue(*rule.PrefixListRuleID)
			ruleModel.Cidr = types.StringValue(*rule.Cidr)
			ruleModel.Description = types.StringValue(*rule.Description)
			prefix.PrefixListRules = append(prefix.PrefixListRules, ruleModel)
		}
		prefixLists = append(prefixLists, prefix)
	}
	config.PrefixLists = prefixLists
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	return

}

type CtyunPrefixRuleModel struct {
	PrefixListRuleID types.String `tfsdk:"prefix_list_rule_id"`
	Cidr             types.String `tfsdk:"cidr"`
	Description      types.String `tfsdk:"description"`
}

type CtyunPrefixListModel struct {
	ID              types.String           `tfsdk:"id"`
	Name            types.String           `tfsdk:"name"`
	Limit           types.Int32            `tfsdk:"limit"`
	AddressType     types.String           `tfsdk:"address_type"`
	Description     types.String           `tfsdk:"description"`
	CreateTime      types.String           `tfsdk:"create_time"`
	UpdateTime      types.String           `tfsdk:"update_time"`
	PrefixListRules []CtyunPrefixRuleModel `tfsdk:"prefix_list_rules"`
}

type CtyunPrefixListsConfig struct {
	RegionID     types.String           `tfsdk:"region_id"`
	ID           types.String           `tfsdk:"id"`
	QueryContent types.String           `tfsdk:"query_content"`
	PageNo       types.Int32            `tfsdk:"page_no"`
	PageSize     types.Int32            `tfsdk:"page_size"`
	PrefixLists  []CtyunPrefixListModel `tfsdk:"prefix_lists"`
}
