package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"
)

// 验证并行文件sfs_size输入
// 要求：步长512
const (
	SfsSizeInputError = "sfs_size输入错误，步长为512！"
)

type validatorSfsSize struct{}

func (v validatorSfsSize) Description(ctx context.Context) string {
	return SfsSizeInputError
}

func (v validatorSfsSize) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorSfsSize) ValidateInt32(ctx context.Context, request validator.Int32Request, response *validator.Int32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.String() == "" {
		return
	}
	value := request.ConfigValue.ValueInt32()
	if value%512 != 0 {
		response.Diagnostics.AddError(SfsSizeInputError, strconv.Itoa(int(value)))
	}
}

func SfsSize() validator.Int32 {
	return &validatorSfsSize{}
}
