package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &ctyunCcseNodeAssociation{}
	_ resource.ResourceWithConfigure   = &ctyunCcseNodeAssociation{}
	_ resource.ResourceWithImportState = &ctyunCcseNodeAssociation{}
)

type ctyunCcseNodeAssociation struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseNodeAssociation() resource.Resource {
	return &ctyunCcseNodeAssociation{}
}

func (c *ctyunCcseNodeAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_node_association"
}

type CtyunCcseNodeAssociationConfig struct {
	ID                       types.String `tfsdk:"id"`
	RegionID                 types.String `tfsdk:"region_id"`
	AzName                   types.String `tfsdk:"az_name"`
	ClusterID                types.String `tfsdk:"cluster_id"`
	InstanceType             types.String `tfsdk:"instance_type"`
	InstanceID               types.String `tfsdk:"instance_id"`
	MirrorID                 types.String `tfsdk:"mirror_id"`
	VisibilityPostHostScript types.String `tfsdk:"visibility_post_host_script"`
	VisibilityHostScript     types.String `tfsdk:"visibility_host_script"`
	Password                 types.String `tfsdk:"password"`
	Name                     types.String `tfsdk:"name"`
	DefaultPoolID            types.String `tfsdk:"default_pool_id"`
	NodeType                 types.String `tfsdk:"node_type"`
	NodeStatus               types.String `tfsdk:"node_status"`
	IsSchedule               types.Bool   `tfsdk:"is_schedule"`
	IsEvict                  types.Bool   `tfsdk:"is_evict"`
}

func (c *ctyunCcseNodeAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10083472/10318452**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "集群ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "主机ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.RegexMatches(regexp.MustCompile("^ss-[a-z0-9]{28}$"), "不符合裸金属id规范"),
						validator2.UUID(),
					),
				},
			},
			"visibility_post_host_script": schema.StringAttribute{
				Optional:    true,
				Description: "部署后执行自定义脚本，base64编码",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"visibility_host_script": schema.StringAttribute{
				Optional:    true,
				Description: "部署前执行自定义脚本，base64编码",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"instance_type": schema.StringAttribute{
				Required:    true,
				Description: "实例类型，支持ecs（云主机）、ebm（裸金属）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.CcseSlaveInstanceTypeEcs, business.CcseSlaveInstanceTypeEbm),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mirror_id": schema.StringAttribute{
				Required:    true,
				Description: "镜像id，可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "用户密码，需要满足以下规则：长度在8～30个字符；必须包含大写字母、小写字母、数字以及特殊符号中的三项；特殊符号可选：()`~!@#$%^&*_-+=|{}[]:;'<>,.?/\\且不能以斜线号/开头",
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.All(
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEcs),
							),
							validator2.EcsPassword(),
						),

						stringvalidator.All(
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEbm),
							),
							validator2.EbmPassword(),
						),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Sensitive: true,
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "纳管后的节点名称，通常和主机名称相同",
			},
			"default_pool_id": schema.StringAttribute{
				Computed:    true,
				Description: "集群默认的节点池ID，纳管的节点都属于此节点池",
			},
			"node_type": schema.StringAttribute{
				Computed:    true,
				Description: "节点类型，master或slave",
			},
			"node_status": schema.StringAttribute{
				Computed:    true,
				Description: "节点状态，normal：健康，abnormal：异常，expulsion：驱逐中。",
			},
			"is_schedule": schema.BoolAttribute{
				Computed:    true,
				Description: "是否调度",
			},
			"is_evict": schema.BoolAttribute{
				Computed:    true,
				Description: "是否驱逐",
			},
		},
	}
}

func (c *ctyunCcseNodeAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCcseNodeAssociationConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	err = c.association(ctx, plan)
	if err != nil {
		return
	}
	// 创建后检查
	name, err := c.checkAfterAssociation(ctx, plan)
	if err != nil {
		return
	}
	plan.Name = types.StringValue(name)

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunCcseNodeAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseNodeAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunCcseNodeAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *ctyunCcseNodeAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseNodeAssociationConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.dissociation(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDissociation(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunCcseNodeAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [name],[clusterID],[regionID]
func (c *ctyunCcseNodeAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunCcseNodeAssociationConfig
	var name, clusterID, regionID string
	err = terraform_extend.Split(request.ID, &name, &clusterID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ClusterID = types.StringValue(clusterID)
	cfg.Name = types.StringValue(name)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// association 纳管节点池
func (c *ctyunCcseNodeAssociation) association(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (err error) {
	params := &ccse2.CcseAttachClusterNodesRequest{
		ClusterId: plan.ClusterID.ValueString(),
		RegionId:  plan.RegionID.ValueString(),
		Instances: []*ccse2.CcseAttachClusterNodesInstancesRequest{
			{
				InstanceId: plan.InstanceID.ValueString(),
				AzName:     plan.AzName.ValueString(),
			},
		},
		VmType:                   plan.InstanceType.ValueString(),
		Runtime:                  "CONTAINERD",
		ImageUuid:                plan.MirrorID.ValueString(),
		VisibilityPostHostScript: plan.VisibilityPostHostScript.ValueString(),
		VisibilityHostScript:     plan.VisibilityHostScript.ValueString(),
		LoginType:                "password",
		Password:                 plan.Password.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCcseApis.CcseAttachClusterNodesApi.Do(ctx, c.meta.SdkCredential, params)
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

// checkAfterAssociation 纳管后检查
func (c *ctyunCcseNodeAssociation) checkAfterAssociation(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (name string, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var node *ccse2.CcseListClusterNodesReturnObjResponse
			node, err = c.getNodeByInstanceID(ctx, plan)
			if err != nil {
				return false
			}
			if node == nil || node.NodeStatus != "normal" {
				return true
			}
			name = node.NodeName
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("纳管时间过长")
		return
	}
	return
}

// getNodeByInstanceID 根据主机id查询节点
func (c *ctyunCcseNodeAssociation) getNodeByInstanceID(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (node *ccse2.CcseListClusterNodesReturnObjResponse, err error) {
	params := &ccse2.CcseListClusterNodesRequest{
		ClusterId: plan.ClusterID.ValueString(),
		RegionId:  plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseListClusterNodesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, n := range resp.ReturnObj {
		if n.EcsId == plan.InstanceID.ValueString() {
			node = n
			return
		}
	}
	return
}

// getNodeDetailByName 根据节点名称查询节点详情
func (c *ctyunCcseNodeAssociation) getNodeDetailByName(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (node *ccse2.CcseGetNodeDetailReturnObjResponse, err error) {
	params := &ccse2.CcseGetNodeDetailRequest{
		ClusterId: plan.ClusterID.ValueString(),
		RegionId:  plan.RegionID.ValueString(),
		NodeName:  plan.Name.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetNodeDetailApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	node = resp.ReturnObj
	return
}

// getCustomPoolID 查询节点池ID
func (c *ctyunCcseNodeAssociation) getCustomPoolID(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (poolID string, err error) {
	params := &ccse2.CcseListNodePoolsRequest{
		RegionId:     plan.RegionID.ValueString(),
		ClusterId:    plan.ClusterID.ValueString(),
		NodePoolName: "CustomPool",
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCcseApis.CcseListNodePoolsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Records) == 0 {
		err = common.InvalidReturnObjError
		return
	}
	poolID = resp.ReturnObj.Records[0].Id
	return
}

// getAndMerge 从远端查询
func (c *ctyunCcseNodeAssociation) getAndMerge(ctx context.Context, plan *CtyunCcseNodeAssociationConfig) (err error) {
	poolID, err := c.getCustomPoolID(ctx, *plan)
	if err != nil {
		return
	}
	plan.DefaultPoolID = types.StringValue(poolID)
	node, err := c.getNodeDetailByName(ctx, *plan)
	if err != nil {
		return
	}
	plan.IsSchedule = types.BoolValue(map[int32]bool{1: true, 0: false}[node.IsSchedule])
	plan.IsEvict = types.BoolValue(map[int32]bool{1: true, 0: false}[node.IsEvict])
	plan.NodeType = types.StringValue(map[int32]string{1: "master", 0: "slave"}[node.NodeType])
	plan.NodeStatus = types.StringValue(node.NodeStatus)
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.Name.ValueString(), plan.ClusterID.ValueString(), plan.RegionID.ValueString()))
	return
}

// dissociation 移除节点
func (c *ctyunCcseNodeAssociation) dissociation(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (err error) {
	params := &ccse2.CcseRemoveNodeV2Request{
		ClusterId:  plan.ClusterID.ValueString(),
		NodePoolId: plan.DefaultPoolID.ValueString(),
		RegionId:   plan.RegionID.ValueString(),
		Nodes:      []string{plan.Name.ValueString()},
		LoginType:  "password",
		Password:   plan.Password.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseRemoveNodeV2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// checkAfterDissociation 移除后检查
func (c *ctyunCcseNodeAssociation) checkAfterDissociation(ctx context.Context, plan CtyunCcseNodeAssociationConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var node *ccse2.CcseListClusterNodesReturnObjResponse
			node, err = c.getNodeByInstanceID(ctx, plan)
			if err != nil {
				return false
			}
			if node != nil {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("移除节点时间过长")
		return
	}
	return
}
