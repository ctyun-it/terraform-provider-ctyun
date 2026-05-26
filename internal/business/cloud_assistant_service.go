package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/common"

type CloudAssistantService struct {
	meta *common.CtyunMetadata
}

func NewCloudAssistantService(meta *common.CtyunMetadata) *CloudAssistantService {
	return &CloudAssistantService{meta: meta}
}
