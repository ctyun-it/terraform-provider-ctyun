package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	_ resource.Resource                = &ctyunRedisBackup{}
	_ resource.ResourceWithConfigure   = &ctyunRedisBackup{}
	_ resource.ResourceWithImportState = &ctyunRedisBackup{}
)

type ctyunRedisBackup struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisBackup() resource.Resource {
	return &ctyunRedisBackup{}
}

func (c *ctyunRedisBackup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_backup"
}

type CtyunRedisBackupConfig struct {
	ID           types.String `tfsdk:"id"`
	InstanceId   types.String `tfsdk:"instance_id"`
	RegionId     types.String `tfsdk:"region_id"`
	Remark       types.String `tfsdk:"remark"`
	Name         types.String `tfsdk:"name"`
	CreateTime   types.String `tfsdk:"create_time"`
	Status       types.String `tfsdk:"status"`
	Type         types.Int32  `tfsdk:"type"`
	IpType       types.String `tfsdk:"ip_type"`
	DownloadUrls types.Map    `tfsdk:"download_urls"`
}

func (c *ctyunRedisBackup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10142282`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "资源唯一标识符",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
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
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "备注信息，不超过128个字符",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "备份名，格式为YYYYMMDDHHMMSS",
			},
			"create_time": schema.StringAttribute{
				Computed:      true,
				Description:   "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "节点状态: success(成功), processing(进行中), fail(失败)",
			},
			"type": schema.Int32Attribute{
				Computed:    true,
				Description: "备份类型: 0(手动备份), 1(自动备份)",
			},
			"ip_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "获取备份文件下载链接的输入参数：网络类型，可选值：publicIp(公网IP)、privateIp(私网IP)，默认为privateIp 支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("publicIp", "privateIp"),
				},
			},
			"download_urls": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "备份文件下载链接，key为Redis节点名，value为备份文件下载URL",
			},
		},
	}
}

func (c *ctyunRedisBackup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisBackupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建备份
	name, err := c.create(ctx, plan)
	if err != nil {
		return
	}

	// 设置备份名
	plan.Name = types.StringValue(name)

	err = c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}

	// 查询备份信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisBackup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisBackupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	err = c.getBackupRdbDownLoadUrl(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisBackup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunRedisBackup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisBackupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除备份
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunRedisBackup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRedisBackup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceId],[restoreName],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()

	var cfg CtyunRedisBackupConfig

	var instanceId, regionId, restoreName string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceId, &restoreName)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceId, &restoreName, &regionId)
		if err != nil {
			return
		}
	}

	if instanceId == "" {
		err = fmt.Errorf("实例ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if restoreName == "" {
		err = fmt.Errorf("名称不能为空")
		return
	}
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.RegionId = types.StringValue(regionId)
	cfg.Name = types.StringValue(restoreName)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建备份
func (c *ctyunRedisBackup) create(ctx context.Context, plan CtyunRedisBackupConfig) (name string, err error) {
	params := &ctgdcs2.Dcs2CreateBackupRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Remark:     plan.Remark.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2CreateBackupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return resp.ReturnObj.RestoreName, nil
}

// checkAfterCreate 创建后检查
func (c *ctyunRedisBackup) checkAfterCreate(ctx context.Context, plan CtyunRedisBackupConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var status string
			status, err = c.getBackupTasks(ctx, &plan)
			if err != nil {
				return false
			}
			if status != "success" {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("创建时间过长")
	}
	return
}

// destroy 删除备份
func (c *ctyunRedisBackup) destroy(ctx context.Context, plan CtyunRedisBackupConfig) (err error) {
	params := &ctgdcs2.Dcs2DeleteBackupRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		RestoreName: plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DeleteBackupApi.Do(ctx, c.meta.SdkCredential, params)
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

// getAndMerge 从远端查询备份信息
func (c *ctyunRedisBackup) getAndMerge(ctx context.Context, plan *CtyunRedisBackupConfig) (err error) {
	params := &ctgdcs2.Dcs2DescribeBackupsRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		RestoreName: plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeBackupsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Rows == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 查找匹配的备份信息
	var backupData *ctgdcs2.Dcs2DescribeBackupsReturnObjRowsResponse
	restoreName := plan.Name.ValueString()

	for _, backup := range resp.ReturnObj.Rows {
		if backup.RestoreName == restoreName {
			backupData = backup
			break
		}
	}

	if backupData == nil {
		err = fmt.Errorf("backup %s not found", restoreName)
		return
	}

	// 设置备份信息
	plan.Name = types.StringValue(backupData.RestoreName)
	plan.CreateTime = types.StringValue(utils.FromBJTimeToUTCZ(backupData.CreateTime))
	plan.Status = types.StringValue(backupData.Status)
	plan.Type = types.Int32Value(backupData.RawType)
	if backupData.Remark != "" {
		plan.Remark = types.StringValue(backupData.Remark)
	} else {
		plan.Remark = types.StringNull()
	}

	// 如果IpType未设置，设置默认值
	if plan.IpType.IsNull() || plan.IpType.IsUnknown() {
		plan.IpType = types.StringValue("privateIp")
	}
	// 如果DownloadUrls未设置，设置为空map
	if plan.DownloadUrls.IsNull() || plan.DownloadUrls.IsUnknown() {
		emptyMap := types.MapNull(types.StringType)
		plan.DownloadUrls = emptyMap
	}

	// 设置ID
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.InstanceId.ValueString(), plan.RegionId.ValueString(), backupData.RestoreName))

	return
}

// getBackupTasks 查询备份任务执行情况
func (c *ctyunRedisBackup) getBackupTasks(ctx context.Context, plan *CtyunRedisBackupConfig) (status string, err error) {
	params := &ctgdcs2.Dcs2DescribeBackupTasksRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		RestoreName: plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeBackupTasksApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 查找匹配的备份信息
	var backupTaskData *ctgdcs2.Dcs2DescribeBackupTasksReturnObjResponse

	backupTaskData = resp.ReturnObj
	if backupTaskData == nil {
		err = fmt.Errorf("backup task %s not found", plan.Name.ValueString())
		return
	}
	// 设置备份信息
	plan.Name = types.StringValue(backupTaskData.RestoreName)
	plan.CreateTime = types.StringValue(backupTaskData.CreateTime)
	plan.Status = types.StringValue(backupTaskData.Status)

	return backupTaskData.Status, nil
}

// getBackupRdbDownLoadUrl 查询备份文件下载链接
func (c *ctyunRedisBackup) getBackupRdbDownLoadUrl(ctx context.Context, plan *CtyunRedisBackupConfig) (err error) {
	// 如果IpType未设置，使用默认值
	ipType := "privateIp"
	if !plan.IpType.IsNull() && !plan.IpType.IsUnknown() {
		ipType = plan.IpType.ValueString()
	}
	params := &ctgdcs2.Dcs2GetRdbDownLoadUrlRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		RestoreName: plan.Name.ValueString(),
		IpType:      ipType,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2GetRdbDownLoadUrlApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回的下载链接
	// ReturnObj格式: map[string]string，key为Redis节点名，value为备份文件下载URL
	downloadUrls := make(map[string]string)
	if returnObjMap, ok := resp.ReturnObj.(map[string]interface{}); ok {
		// 遍历所有节点的下载链接
		for nodeName, url := range returnObjMap {
			if urlStr, ok := url.(string); ok {
				downloadUrls[nodeName] = urlStr
			}
		}

		// 将下载链接设置到plan中
		urlsMap, diag := types.MapValueFrom(ctx, types.StringType, downloadUrls)
		if diag.HasError() {
			err = fmt.Errorf("failed to convert download URLs to map: %v", diag.Errors())
			return
		}
		plan.DownloadUrls = urlsMap
	} else {
		err = fmt.Errorf("unexpected return object type: %T", resp.ReturnObj)
		return
	}

	// 设置IpType
	plan.IpType = types.StringValue(ipType)

	return
}
