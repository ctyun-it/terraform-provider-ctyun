package cloud_assistant

import (
	"context"
	"errors"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunCloudAssistantInvocationResults{}
	_ datasource.DataSourceWithConfigure = &CtyunCloudAssistantInvocationResults{}
)

type CtyunCloudAssistantInvocationResults struct {
	meta *common.CtyunMetadata
}

func NewCtyunCloudAssistantInvocationResults() datasource.DataSource {
	return &CtyunCloudAssistantInvocationResults{}
}

func (c *CtyunCloudAssistantInvocationResults) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunCloudAssistantInvocationResults) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_cloud_assistant_invocation_results"
}

func (c *CtyunCloudAssistantInvocationResults) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询云助手命令执行结果", "云助手", "https://www.ctyun.cn/document/10026730/10764038"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID",
			},
			"command_id": schema.StringAttribute{
				Optional:    true,
				Description: "命令 ID",
			},
			"invoked_id": schema.StringAttribute{
				Optional:    true,
				Description: "命令执行 ID",
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
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "命令总个数",
			},
			"results": schema.ListNestedAttribute{
				Computed:    true,
				Description: "命令执行结果列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"invoked_id": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行 ID",
						},
						"command_id": schema.StringAttribute{
							Computed:    true,
							Description: "命令 ID",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "云主机 ID",
						},
						"invocation_status": schema.StringAttribute{
							Computed:    true,
							Description: "命令运行状态：Pending/Running/Success/Failed",
						},
						"output": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行输出信息",
						},
						"exit_code": schema.Int32Attribute{
							Computed:    true,
							Description: "命令退出码",
						},
						"error_info": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行失败原因",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行创建时间",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "命令执行完成时间",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunCloudAssistantInvocationResults) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunCloudAssistantInvocationResultsConfig
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

	req := &ctecs.CtecsCloudAssistantDescribeInvocationResultsRequest{
		RegionID: regionId,
		PageNo:   config.PageNo.ValueInt32(),
		PageSize: config.PageSize.ValueInt32(),
	}

	if !config.CommandID.IsNull() && !config.CommandID.IsUnknown() && config.CommandID.ValueString() != "" {
		req.CommandID = config.CommandID.ValueString()
	}
	if !config.InvokedID.IsNull() && !config.InvokedID.IsUnknown() && config.InvokedID.ValueString() != "" {
		req.InvokedID = config.InvokedID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsCloudAssistantDescribeInvocationResultsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询云助手执行结果失败(invoked_id=%s, command_id=%s)，接口返回nil，请联系研发确认问题原因！", config.InvokedID.ValueString(), config.CommandID.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	var results []InvocationResultDatasourceModel

	for _, resultItem := range resp.ReturnObj.Results {
		results = append(results, InvocationResultDatasourceModel{
			InvokedID:        types.StringValue(resultItem.InvokedID),
			CommandID:        types.StringValue(resultItem.CommandID),
			InstanceID:       types.StringValue(resultItem.InstanceID),
			InvocationStatus: types.StringValue(resultItem.InvocationStatus),
			Output:           types.StringValue(resultItem.Output),
			ExitCode:         types.Int32Value(resultItem.ExitCode),
			ErrorInfo:        types.StringValue(resultItem.ErrorInfo),
			CreateTime:       types.StringValue(resultItem.CreateTime),
			UpdateTime:       types.StringValue(resultItem.UpdateTime),
		})
	}
	config.Results = results

	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

type InvocationResultDatasourceModel struct {
	InvokedID        types.String `tfsdk:"invoked_id"`
	CommandID        types.String `tfsdk:"command_id"`
	InstanceID       types.String `tfsdk:"instance_id"`
	InvocationStatus types.String `tfsdk:"invocation_status"`
	Output           types.String `tfsdk:"output"`
	ExitCode         types.Int32  `tfsdk:"exit_code"`
	ErrorInfo        types.String `tfsdk:"error_info"`
	CreateTime       types.String `tfsdk:"create_time"`
	UpdateTime       types.String `tfsdk:"update_time"`
}

type CtyunCloudAssistantInvocationResultsConfig struct {
	RegionID   types.String                      `tfsdk:"region_id"`
	CommandID  types.String                      `tfsdk:"command_id"`
	InvokedID  types.String                      `tfsdk:"invoked_id"`
	PageNo     types.Int32                       `tfsdk:"page_no"`
	PageSize   types.Int32                       `tfsdk:"page_size"`
	TotalCount types.Int32                       `tfsdk:"total_count"`
	Results    []InvocationResultDatasourceModel `tfsdk:"results"`
}
