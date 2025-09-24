package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunMysqlWhiteLists{}
	_ datasource.DataSourceWithConfigure = &CtyunMysqlWhiteLists{}
)

type CtyunMysqlWhiteLists struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMysqlWhiteLists) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10033813/10133794**`,
		Attributes: map[string]schema.Attribute{
			"prod_inst_id": schema.StringAttribute{
				Required:    true,
				Description: "Mysql实例id",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id",
			},
			"white_lists": schema.ListNestedAttribute{
				Computed:    true,
				Description: "白名单列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_name": schema.StringAttribute{
							Computed:    true,
							Description: "白名单分组名",
						},
						"group_white_list_count": schema.Int32Attribute{
							Computed:    true,
							Description: "白名单分组组内数量",
						},
						"prod_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "Mysql实例ID",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"group_white_list": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "白名单IP列表",
						},
						"access_machine_type": schema.StringAttribute{
							Computed:    true,
							Description: "访问类型",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMysqlWhiteLists) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlAccessWhiteListsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}

	params := mysql.TeledbGetAccessWhiteListRequest{
		OuterProdInstID: config.ProdInstID.ValueString(),
	}
	header := mysql.TeledbGetAccessWhiteListRequestHeader{
		InstID:   config.ProdInstID.ValueStringPointer(),
		RegionID: regionId,
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetAccessWhiteList.Do(ctx, c.meta.Credential, &params, &header)
	if err != nil {
		return
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
	var whiteList []CtyunMysqlAccessWhiteListModel
	for _, whiteInfo := range resp.ReturnObj {
		var whiteListInfo CtyunMysqlAccessWhiteListModel
		whiteListInfo.GroupName = types.StringValue(whiteInfo.GroupName)
		whiteListInfo.GroupWhiteListCount = types.Int32Value(whiteInfo.GroupWhiteListCount)
		whiteListInfo.ProdInstID = types.StringValue(whiteInfo.OuterProdInstID)
		whiteListInfo.CreatedTime = types.StringValue(fmt.Sprintf("%d", whiteInfo.CreateTime))
		whiteListInfo.UpdatedTime = types.StringValue(fmt.Sprintf("%d", whiteInfo.UpdateTime))
		whiteListInfo.AccessMachineType = types.StringValue(whiteInfo.AccessMachineType)
		groupWhiteList, diags := types.SetValueFrom(ctx, types.StringType, whiteInfo.WhiteList)
		if diags.HasError() {
			return
		}
		whiteListInfo.GroupWhiteList = groupWhiteList
		whiteList = append(whiteList, whiteListInfo)
	}
	config.WhiteLists = whiteList
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func NewCtyunMysqlWhiteLists() *CtyunMysqlWhiteLists {
	return &CtyunMysqlWhiteLists{}
}
func (c *CtyunMysqlWhiteLists) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMysqlWhiteLists) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_white_lists"
}

type CtyunMysqlAccessWhiteListModel struct {
	GroupName           types.String `tfsdk:"group_name"`
	GroupWhiteListCount types.Int32  `tfsdk:"group_white_list_count"`
	ProdInstID          types.String `tfsdk:"prod_inst_id"`
	CreatedTime         types.String `tfsdk:"created_time"`
	UpdatedTime         types.String `tfsdk:"updated_time"`
	GroupWhiteList      types.Set    `tfsdk:"group_white_list"`
	AccessMachineType   types.String `tfsdk:"access_machine_type"`
}

type CtyunMysqlAccessWhiteListsConfig struct {
	ProdInstID types.String                     `tfsdk:"prod_inst_id"`
	ProjectID  types.String                     `tfsdk:"project_id"`
	RegionID   types.String                     `tfsdk:"region_id"`
	WhiteLists []CtyunMysqlAccessWhiteListModel `tfsdk:"white_lists"`
}
