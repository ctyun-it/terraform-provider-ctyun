package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type ctyunScalingEcsProtection struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
	imageService  *business.ImageService
}

func NewCtyunScalingEcsProtection() resource.Resource {
	return &ctyunScalingEcsProtection{}
}

func (c *ctyunScalingEcsProtection) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_ecs_protection"
}

func (c *ctyunScalingEcsProtection) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)
	c.imageService = business.NewImageService(c.meta)
}

func (c *ctyunScalingEcsProtection) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：//www.ctyun.cn/document/10027725/10216534`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "伸缩组ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"instance_id_list": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "需要开启伸缩保护的的云主机uuid列表。伸缩组内云主机清单可以根据data.ctyun_scaling_ecs_list获取。支持更新。",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"protect_status": schema.BoolAttribute{
				Required:    true,
				Description: "开始保护或者停止保护伸缩组内的一台或者多台云主机。false=关闭云主机保护，true=开启云主机保护",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunScalingEcsProtection) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunScalingEcsProtectionConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 确认instance_id_list机器都在伸缩组内
	isValid, err := c.preCheckBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	if isValid {
		err = c.updateScalingEcsProtection(ctx, &plan)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingEcsProtection) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *ctyunScalingEcsProtection) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *ctyunScalingEcsProtection) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

// 确认instance_id_list机器都在伸缩组内
func (c *ctyunScalingEcsProtection) preCheckBeforeCreate(ctx context.Context, config *CtyunScalingEcsProtectionConfig) (bool, error) {

	instanceList, err := c.getInstanceListByGroupID(ctx, config)
	if err != nil {
		return false, err
	}
	var instanceUUIDList []string
	diags := config.InstanceIDList.ElementsAs(ctx, &instanceUUIDList, true)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return false, err
	}
	// 确认待添加云主机未在伸缩组内
	removeIntersection := c.FindIntersection(instanceUUIDList, instanceList)
	if len(removeIntersection) != len(instanceUUIDList) {
		err = fmt.Errorf("待删除的云主机中有部分未加入伸缩组，符合删除条件列表为：%s", strings.Join(removeIntersection, ", "))
		return false, err
	}
	return true, nil
}

func (c *ctyunScalingEcsProtection) FindIntersection(instanceUUIDList []string, scalingInstanceList []*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse) []string {
	// 使用map记录第一个数组的元素
	set := make(map[string]bool)
	for _, item := range instanceUUIDList {
		set[item] = true
	}

	var intersection []string
	// 遍历第二个数组，检查元素是否在第一个数组中存在
	for _, item := range scalingInstanceList {
		instanceID := item.InstanceID
		if set[instanceID] {
			intersection = append(intersection, instanceID)
			set[instanceID] = false // 标记已添加，避免重复
		}
	}
	return intersection
}

func (c *ctyunScalingEcsProtection) getInstanceListByGroupID(ctx context.Context, config *CtyunScalingEcsProtectionConfig) ([]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse, error) {
	var pageSize, pageNo int32
	pageSize = 100
	pageNo = 1
	resp, err := c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount := resp.ReturnObj.TotalCount
	totalPageNo := pageNo
	if totalCount > pageSize {
		totalPageNo = totalCount/pageSize + 1
	}
	var instances []*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse
	for pageNo <= totalPageNo {
		instanceList := resp.ReturnObj.InstanceList
		for _, instance := range instanceList {
			instances = append(instances, instance)
		}
		pageNo++
		if pageNo > totalPageNo {
			break
		}
		resp, err = c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return instances, nil
}

func (c *ctyunScalingEcsProtection) requestInstanceListByGroup(ctx context.Context, config *CtyunScalingEcsProtectionConfig, pageNo, pageSize int32) (*scaling.ScalingGroupQueryInstanceListResponse, error) {
	params := &scaling.ScalingGroupQueryInstanceListRequest{
		RegionID: config.RegionID.ValueString(),
		GroupID:  config.GroupID.ValueInt64(),
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupQueryInstanceListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询group id为%d下的云主机列表失败，接口范围nil。请联系研发，或稍后重试！", config.GroupID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *ctyunScalingEcsProtection) updateScalingEcsProtection(ctx context.Context, plan *CtyunScalingEcsProtectionConfig) error {

	if !plan.ProtectStatus.IsNull() {
		var instanceUUIDs []string
		diags := plan.InstanceIDList.ElementsAs(ctx, &instanceUUIDs, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}

		instanceIds, err := c.getInstanceAssocIdByUUID(ctx, plan, instanceUUIDs)
		if err != nil {
			return err
		}
		// 关闭云主机保护
		if !plan.ProtectStatus.ValueBool() {
			err = c.disableProtectEcs(ctx, plan, instanceIds, instanceUUIDs)
			if err != nil {
				return err
			}
		} else if plan.ProtectStatus.ValueBool() {
			// 开启云主机保护
			err = c.enableProtectEcs(ctx, plan, instanceIds, instanceUUIDs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ctyunScalingEcsProtection) getInstanceAssocIdByUUID(ctx context.Context, state *CtyunScalingEcsProtectionConfig, UUIDs []string) ([]int32, error) {
	instanceMap, err := c.getInstanceMap(ctx, state)
	if err != nil {
		return nil, err
	}
	var instanceIdList []int32
	for _, uuid := range UUIDs {
		ecsInfo := instanceMap[uuid]
		if ecsInfo == nil {
			continue
		}
		instanceIdList = append(instanceIdList, ecsInfo.Id)
	}
	return instanceIdList, nil
}

func (c *ctyunScalingEcsProtection) getInstanceMap(ctx context.Context, config *CtyunScalingEcsProtectionConfig) (map[string]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse, error) {
	var pageNo, pageSize int32
	pageNo = 1
	pageSize = 100
	pageEndNo := pageNo
	ecsListResp, err := c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	totalCount := ecsListResp.ReturnObj.TotalCount

	if totalCount > pageSize {
		pageEndNo = totalCount / pageSize
	}
	// 先获取所有ecs列表，并设置成map
	instanceMap := make(map[string]*scaling.ScalingGroupQueryInstanceListReturnObjInstanceListResponse)
	for pageNo <= pageEndNo {
		ecsList := ecsListResp.ReturnObj.InstanceList
		for _, ecs := range ecsList {
			instanceId := ecs.InstanceID
			instanceMap[instanceId] = ecs
		}

		pageNo++
		if pageNo > pageEndNo {
			break
		}
		ecsListResp, err = c.requestInstanceListByGroup(ctx, config, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
	}
	return instanceMap, nil
}

func (c *ctyunScalingEcsProtection) disableProtectEcs(ctx context.Context, state *CtyunScalingEcsProtectionConfig, instanceIDs []int32, instanceUUIDs []string) error {

	params := &scaling.ScalingGroupProtectDisableRequest{
		RegionID:       state.RegionID.ValueString(),
		GroupID:        state.GroupID.ValueInt64(),
		InstanceIDList: instanceIDs,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupProtectDisableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("关闭云主机保护失败，接口返回nil，ecs列表：%s", strings.Join(instanceUUIDs, ", "))
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunScalingEcsProtection) enableProtectEcs(ctx context.Context, state *CtyunScalingEcsProtectionConfig, instanceIds []int32, instanceUUIDs []string) error {
	params := scaling.ScalingGroupProtectEnableRequest{
		RegionID:       state.RegionID.ValueString(),
		GroupID:        state.GroupID.ValueInt64(),
		InstanceIDList: instanceIds,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingGroupProtectEnableApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("开启云主机保护失败，接口返回nil。ecs列表：%s", strings.Join(instanceUUIDs, ", "))
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

type CtyunScalingEcsProtectionConfig struct {
	RegionID       types.String `tfsdk:"region_id"`
	GroupID        types.Int64  `tfsdk:"group_id"`
	InstanceIDList types.Set    `tfsdk:"instance_id_list"`
	ProtectStatus  types.Bool   `tfsdk:"protect_status"`
}
