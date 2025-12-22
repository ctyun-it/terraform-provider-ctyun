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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource              = &CtyunMysqlAudit{}
	_ resource.ResourceWithConfigure = &CtyunMysqlAudit{}
)

type CtyunMysqlAudit struct {
	meta         *common.CtyunMetadata
	mysqlService *business.MysqlService
}

func (c *CtyunMysqlAudit) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_audit"
}
func NewCtyunMysqlAudit() resource.Resource {
	return &CtyunMysqlAudit{}
}

func (c *CtyunMysqlAudit) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.mysqlService = business.NewMysqlService(c.meta)
}

func (c *CtyunMysqlAudit) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10133568",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql实例Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
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
			"audit_switch": schema.BoolAttribute{
				Required:    true,
				Description: "sql审计开关状态。false：关闭状态；true：开启状态。",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *CtyunMysqlAudit) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlAuditConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开启/关闭sql审计
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，获取mysql详情
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAudit) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunMysqlAudit) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMysqlAudit) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlAudit) create(ctx context.Context, config *CtyunMysqlAuditConfig) error {
	// 进行审计之前，确定数据库是否running
	err := c.startedLoop(ctx, config)
	if err != nil {
		return err
	}
	params := &mysql.TeledbStartAuditRequest{
		OuterProdInstId: config.InstID.ValueString(),
		AuditSwitch:     config.AuditSwitch.ValueBool(),
	}
	header := &mysql.TeledbStartAuditRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbStartAuditApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("开始/关闭数据库（id=%s）SQL审计失败，接口返回nil", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	status, err := c.getMysqlAuditStatus(ctx, config)
	if err != nil || status == nil {
		return err
	}

	if *status != config.AuditSwitch.ValueBool() {
		err = fmt.Errorf("开始/关闭数据库（id=%s）SQL审计失败，当前状态为：%t", config.InstID.ValueString(), *status)
		return nil
	}
	return nil
}

func (c *CtyunMysqlAudit) getAndMerge(ctx context.Context, config *CtyunMysqlAuditConfig) error {
	return nil
}

func (c *CtyunMysqlAudit) getMysqlAuditStatus(ctx context.Context, config *CtyunMysqlAuditConfig) (*bool, error) {
	params := &mysql.TeledbGetAuditStatusRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetAuditStatusRequestHeader{
		RegionID: config.RegionID.ValueString(),
		InstID:   config.InstID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	// 先确保mysql实例状态running

	err := c.startedLoop(ctx, config)
	if err != nil {
		return nil, err
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetAuditStatusApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取数据库（id=%s）SQL审计状态失败，接口返回nil", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return &resp.ReturnObj.AuditLogSwitch, nil
}

// startedLoop 等待实例处于启动状态
func (c *CtyunMysqlAudit) startedLoop(ctx context.Context, state *CtyunMysqlAuditConfig) (err error) {
	retryer, err := business.NewRetryer(time.Second*30, 60)
	if err != nil {
		return
	}
	var cnt int
	result := retryer.Start(
		func(currentTime int) bool {
			var instance *mysql.DetailRespReturnObj
			instance, err = c.mysqlService.GetDetailByID(
				ctx,
				state.InstID.ValueString(),
				state.ProjectID.ValueString(),
				state.RegionID.ValueString(),
			)
			if err != nil {
				return false
			}
			runningStatus := instance.ProdRunningStatus
			orderStatus := instance.ProdOrderStatus
			// 若变配前，发现数据库已冻结，将其恢复
			if orderStatus == business.MysqlOrderStatusPause {
				err = fmt.Errorf("数据库id=%s处于冻结状态，请恢复后在再请求！", state.InstID.ValueString())
				if err != nil {
					return false
				}
			}
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
				cnt++
				if cnt > 1 {
					return false
				}
			}
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("订单处于暂停状态，不可进行变更操作")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStopping || runningStatus == business.MysqlRunningStatusStopped {
				err = errors.New("主机处于关机状态，不可进行变更操作")
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
	}
	return
}

type CtyunMysqlAuditConfig struct {
	InstID      types.String `tfsdk:"instance_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	RegionID    types.String `tfsdk:"region_id"`
	AuditSwitch types.Bool   `tfsdk:"audit_switch"`
}
