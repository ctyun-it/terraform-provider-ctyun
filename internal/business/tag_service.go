package business

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type TagsService struct {
	meta *common.CtyunMetadata
}

func NewTagsService(meta *common.CtyunMetadata) *TagsService {
	return &TagsService{meta: meta}
}
func (v TagsService) BindTags(ctx context.Context, regionId, resourceType, resourceID string, tags *types.Set) (err error) {

	// 最大轮询100次直到没有 "errorCode":"Openapi.Nat.BindLabelFailed"
	retryer, err := NewRetryer(time.Second*10, 61)
	if err != nil {
		return nil
	}
	result := retryer.Start(func(currentTime int) bool {
		err = v.Bind(ctx, regionId, resourceType, resourceID, tags)
		if err != nil {
			// 检查错误是否包含BindLabelFailed，如果是则继续重试
			if strings.Contains(err.Error(), "bind label failed") {
				return true // 继续重试
			}
			return false // 其他错误直接退出
		}
		return false // 成功则退出
	})
	if result.ReturnReason == ReachMaxLoopTime {
		err = errors.New("轮询间隔已达10分钟，标签绑定仍未成功！")
		return
	}
	tagsSet, err := v.QueryAll(ctx, regionId, resourceType, resourceID)
	if err != nil {
		return
	}
	tags = &tagsSet
	return
}
func (v TagsService) Bind(ctx context.Context, regionId, resourceType, resourceID string, tags *types.Set) (err error) {
	var planTags []Tag
	diag := tags.ElementsAs(ctx, &planTags, true)
	if diag.HasError() {
		return
	}
	//循环遍历 tags 中的标签 和值  分别添加到params 调用Do 方法
	for _, tag := range planTags {
		params := &ctvpc.CtvpcResourceBindLabelRequest{
			RegionID:     regionId,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			LabelKey:     tag.LabelKey.ValueString(),
			LabelValue:   tag.LabelValue.ValueString(),
		}
		var resp *ctvpc.CtvpcResourceBindLabelResponse
		resp, err = v.meta.Apis.SdkCtVpcApis.CtvpcResourceBindLabelApi.Do(ctx, v.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf("API return error. Message: %s", *resp.Message)
			return
		}
	}
	return
}

func (v TagsService) Unbind(ctx context.Context, regionId, resourceType, resourceID string, tags *types.Set) (err error) {

	//err = v.Query(ctx, regionId, resourceType, resourceID, tags)
	//if err != nil {
	//	return
	//}
	var planTags []Tag
	diags := tags.ElementsAs(ctx, &planTags, true)
	if diags.HasError() {
		err = fmt.Errorf("failed to convert tags to slice: %v", diags.Errors())
		return
	}
	//循环遍历 tags 中的标签 和值  分别添加到params 调用Do 方法
	for _, tag := range planTags {
		params := &ctvpc.CtvpcResourceUnbindLabelRequest{
			RegionID:     regionId,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			LabelID:      tag.LabelID.ValueString(),
		}

		var resp *ctvpc.CtvpcResourceUnbindLabelResponse
		resp, err = v.meta.Apis.SdkCtVpcApis.CtvpcResourceUnbindLabelApi.Do(ctx, v.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf("API return error. Message: %s", *resp.Message)
			return
		}
	}
	return
}
func (v TagsService) UpdateBind(ctx context.Context, regionId, resourceType, resourceID string, tags *types.Set) (err error) {

	unBindTags, err := v.QueryAll(ctx, regionId, resourceType, resourceID)
	if err != nil {
		return
	}
	// 执行解绑操作
	if !unBindTags.IsNull() {
		err = v.Unbind(ctx, regionId, resourceType, resourceID, &unBindTags)
		if err != nil {
			return
		}
	}

	// 执行绑定操作
	if !tags.IsNull() {
		err = v.Bind(ctx, regionId, resourceType, resourceID, tags)
		if err != nil {
			return
		}
	}
	tagsSet, err := v.QueryAll(ctx, regionId, resourceType, resourceID)
	if err != nil {
		return
	}
	tags = &tagsSet
	return
}

//func (v TagsService) Query(ctx context.Context, regionId, resourceType, resourceID string, tags *types.Set) (err error) {
//	allTagsList, err := v.QueryAll(ctx, regionId, resourceType, resourceID)
//	if err != nil {
//		return
//	}
//	var planTags []Tag
//	diags := tags.ElementsAs(ctx, &planTags, true)
//	if diags.HasError() {
//		err = fmt.Errorf(diags[0].Detail())
//		return
//	}
//	var allTags []Tag
//	diags = allTagsList.ElementsAs(ctx, &allTags, true)
//	if diags.HasError() {
//		err = fmt.Errorf(diags[0].Detail())
//		return
//	}
//
//	// 遍历planTags，为每个标签找到对应的LabelID
//	for i, planTag := range planTags {
//		for _, result := range allTags {
//			if planTag.LabelKey.ValueString() == result.LabelKey.ValueString() &&
//				planTag.LabelValue.ValueString() == result.LabelValue.ValueString() {
//				// 正确地将LabelID赋值给planTag
//				planTags[i].LabelID = result.LabelID
//				break
//			}
//		}
//	}
//
//	// 将更新后的planTags转换为types.List
//	tagsSet, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(Tag{}), planTags)
//	if diags.HasError() {
//		err = fmt.Errorf(diags[0].Detail())
//		return
//	}
//	tags = &tagsSet // 注意：这里需要解引用tags指针
//	return
//}

func (v TagsService) QueryAll(ctx context.Context, regionId string, resourceType string, resourceID string) (tags types.Set, err error) {
	params := &ctvpc.CtvpcQueryLabelsByResourceRequest{
		RegionID:     regionId,
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}

	resp, err := v.meta.Apis.SdkCtVpcApis.CtvpcQueryLabelsByResourceApi.Do(ctx, v.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var resultTags []Tag
	for _, apiTag := range resp.ReturnObj.Results {
		tag := Tag{
			LabelID:    types.StringPointerValue(apiTag.LabelID),
			LabelKey:   types.StringPointerValue(apiTag.LabelKey),
			LabelValue: types.StringPointerValue(apiTag.LabelValue),
		}
		resultTags = append(resultTags, tag)
	}
	tags, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(Tag{}), resultTags)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return
	}
	return
}

type Tag struct {
	// 标签键（必填）
	// 约束：1~32字符，不能换行或以空格开头/结尾，同一镜像的标签键不可重复
	LabelID types.String `tfsdk:"id"`
	// 约束：1~32字符，不能换行或以空格开头/结尾，同一镜像的标签键不可重复
	LabelKey types.String `tfsdk:"key"`

	// 标签值（必填）
	// 约束：1~32字符，不能换行或以空格开头/结尾
	LabelValue types.String `tfsdk:"value"`
}
