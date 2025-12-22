package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
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
	_ resource.Resource              = &CtyunMongodbRestartDb{}
	_ resource.ResourceWithConfigure = &CtyunMongodbRestartDb{}
)

func NewCtyunMongodbRestartDb() resource.Resource {
	return &CtyunMongodbRestartDb{}
}

type CtyunMongodbRestartDb struct {
	meta *common.CtyunMetadata
}

type CtyunMongodbRestartDbConfig struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	RegionID   types.String `tfsdk:"region_id"`
	ProjectID  types.String `tfsdk:"project_id"`
}

func (c *CtyunMongodbRestartDb) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_restart_db"
}

func (c *CtyunMongodbRestartDb) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 instance_id:ip_RestartDb_name",
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
		},
	}
}

func (c *CtyunMongodbRestartDb) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbRestartDb) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMongodbRestartDbConfig
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *CtyunMongodbRestartDb) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMongodbRestartDbConfig
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CtyunMongodbRestartDb) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (c *CtyunMongodbRestartDb) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (c *CtyunMongodbRestartDb) checkBeforeCreate(ctx context.Context, state *CtyunMongodbRestartDbConfig, loopCount ...int) (err error) {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
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

func (c *CtyunMongodbRestartDb) create(ctx context.Context, plan *CtyunMongodbRestartDbConfig) (err error) {
	request := &mongodb.MongodbRestartDbRequest{
		ProdInstId: plan.InstanceID.ValueString(),
	}
	headers := &mongodb.MongodbRestartDbRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbRestartDbApi.Do(ctx, c.meta.Credential, request, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	}
	plan.ID = plan.InstanceID

	return
}
