package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"
)

// 验证 RocketMQ disk_size 输入
// 要求：必须为 100 的整数倍
const (
	RocketmqDiskSizeInputError = "disk_size 输入错误，必须为 100 的整数倍！"
)

type validatorRocketmqDiskSize struct{}

func (v validatorRocketmqDiskSize) Description(ctx context.Context) string {
	return RocketmqDiskSizeInputError
}

func (v validatorRocketmqDiskSize) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorRocketmqDiskSize) ValidateInt32(ctx context.Context, request validator.Int32Request, response *validator.Int32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.String() == "" {
		return
	}
	value := request.ConfigValue.ValueInt32()
	if value%100 != 0 {
		response.Diagnostics.AddError(RocketmqDiskSizeInputError, strconv.Itoa(int(value)))
	}
}

func RocketmqDiskSize() validator.Int32 {
	return &validatorRocketmqDiskSize{}
}
