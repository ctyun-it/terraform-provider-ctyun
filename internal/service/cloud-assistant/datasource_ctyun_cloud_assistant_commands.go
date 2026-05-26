package cloud_assistant

import (
	"context"
	"errors"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunCloudAssistantCommands{}
	_ datasource.DataSourceWithConfigure = &CtyunCloudAssistantCommands{}
)

type CtyunCloudAssistantCommands struct {
	meta *common.CtyunMetadata
}

func NewCtyunCloudAssistantCommands() datasource.DataSource {
	return &CtyunCloudAssistantCommands{}
}

func (c *CtyunCloudAssistantCommands) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunCloudAssistantCommands) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cloud_assistant_commands"
}

func (c *CtyunCloudAssistantCommands) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询云助手命令列表", "云助手", "https://www.ctyun.cn/document/10026730/10764038"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "当前页码，默认为 1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "每页行数，最大100，默认为 10",
			},
			"is_public": schema.BoolAttribute{
				Optional:    true,
				Description: "是否为公共市场命令",
			},
			"filters": schema.ListNestedAttribute{
				Optional:    true,
				Description: "过滤条件，支持 commandID, commandName, commandType",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "过滤条件的字段名，支持 commandID, commandName, commandType",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "过滤字段对应的值",
						},
					},
				},
			},
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "命令总个数",
			},
			"commands": schema.ListNestedAttribute{
				Computed:    true,
				Description: "命令列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"command_id": schema.StringAttribute{
							Computed:    true,
							Description: "命令 ID",
						},
						"command_name": schema.StringAttribute{
							Computed:    true,
							Description: "命令名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "命令描述",
						},
						"command_type": schema.StringAttribute{
							Computed:    true,
							Description: "命令类型",
						},
						"command_content": schema.StringAttribute{
							Computed:    true,
							Description: "命令内容",
						},
						"working_directory": schema.StringAttribute{
							Computed:    true,
							Description: "运行目录",
						},
						"timeout": schema.Int32Attribute{
							Computed:    true,
							Description: "超时时间",
						},
						"is_public": schema.BoolAttribute{
							Computed:    true,
							Description: "是否为公共市场命令",
						},
						"enabled_parameter": schema.BoolAttribute{
							Computed:    true,
							Description: "是否启用自定义参数",
						},
						"default_parameter": schema.StringAttribute{
							Computed:    true,
							Description: "自定义参数默认值",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunCloudAssistantCommands) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunCloudAssistantCommandsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region_id不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)

	req := &ctecs.CtecsGetCommandsRequest{
		RegionID: regionId,
		PageNo:   config.PageNo.ValueInt32(),
		PageSize: config.PageSize.ValueInt32(),
	}

	if !config.IsPublic.IsNull() && !config.IsPublic.IsUnknown() {
		isPublic := config.IsPublic.ValueBool()
		req.IsPublic = &isPublic
	}

	if len(config.Filters) > 0 {
		filters := make([]*ctecs.CtecsGetCommandsFiltersRequest, len(config.Filters))
		for i, f := range config.Filters {
			filters[i] = &ctecs.CtecsGetCommandsFiltersRequest{
				Key:   f.Key.ValueString(),
				Value: f.Value.ValueString(),
			}
		}
		req.Filters = filters
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsGetCommandsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	}

	if resp == nil || resp.ReturnObj == nil {
		err = errors.New("查询云助手命令列表失败，接口返回异常")
		return
	}

	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)

	var commands []CommandListModel
	for _, cmd := range resp.ReturnObj.Commands {
		var model CommandListModel
		model.CommandID = types.StringValue(cmd.CommandID)
		model.CommandName = types.StringValue(cmd.CommandName)
		model.Description = types.StringValue(cmd.Description)
		model.CommandType = types.StringValue(cmd.CommandType)
		model.CommandContent = types.StringValue(cmd.CommandContent)
		model.WorkingDirectory = types.StringValue(cmd.WorkingDirectory)
		model.Timeout = types.Int32Value(cmd.Timeout)
		model.DefaultParameter = types.StringValue(cmd.DefaultParameter)
		model.CreateTime = types.StringValue(cmd.CreateTime)
		model.UpdateTime = types.StringValue(cmd.UpdateTime)

		if cmd.IsPublic != nil {
			model.IsPublic = types.BoolValue(*cmd.IsPublic)
		}
		if cmd.EnabledParameter != nil {
			model.EnabledParameter = types.BoolValue(*cmd.EnabledParameter)
		}

		commands = append(commands, model)
	}

	config.Commands = commands
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

type CommandListModel struct {
	CommandID        types.String `tfsdk:"command_id"`
	CommandName      types.String `tfsdk:"command_name"`
	Description      types.String `tfsdk:"description"`
	CommandType      types.String `tfsdk:"command_type"`
	CommandContent   types.String `tfsdk:"command_content"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
	Timeout          types.Int32  `tfsdk:"timeout"`
	IsPublic         types.Bool   `tfsdk:"is_public"`
	EnabledParameter types.Bool   `tfsdk:"enabled_parameter"`
	DefaultParameter types.String `tfsdk:"default_parameter"`
	CreateTime       types.String `tfsdk:"create_time"`
	UpdateTime       types.String `tfsdk:"update_time"`
}

type FilterModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type CtyunCloudAssistantCommandsConfig struct {
	RegionID   types.String       `tfsdk:"region_id"`
	PageNo     types.Int32        `tfsdk:"page_no"`
	PageSize   types.Int32        `tfsdk:"page_size"`
	IsPublic   types.Bool         `tfsdk:"is_public"`
	Filters    []FilterModel      `tfsdk:"filters"`
	TotalCount types.Int32        `tfsdk:"total_count"`
	Commands   []CommandListModel `tfsdk:"commands"`
}
