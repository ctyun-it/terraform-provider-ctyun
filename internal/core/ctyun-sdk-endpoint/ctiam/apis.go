package ctiam

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// Apis api的接口
type Apis struct {
	ServiceListApi                        *ServiceListApi
	AuthorityListApi                      *AuthorityListApi
	PolicyCreateApi                       *PolicyCreateApi
	PolicyUpdateApi                       *PolicyUpdateApi
	PolicyDeleteApi                       *PolicyDeleteApi
	PolicyGetApi                          *PolicyGetApi
	PolicyAttachUserGroupApi              *PolicyAttachUserGroupApi
	PolicyAttachUserApi                   *PolicyAttachUserApi
	PolicyInvalidUserGroupApi             *PolicyInvalidUserGroupApi
	UserCreateApi                         *UserCreateApi
	UserGetApi                            *UserGetApi
	UserUpdateApi                         *UserUpdateApi
	UserInvalidApi                        *UserInvalidApi
	UserResetPasswordApi                  *UserResetPasswordApi
	UserGroupCreateApi                    *UserGroupCreateApi
	UserGroupUpdateApi                    *UserGroupUpdateApi
	UserGroupInvalidApi                   *UserGroupInvalidApi
	UserGroupGetApi                       *UserGroupGetApi
	UserAssociationGroupApi               *UserAssociationGroupApi
	UserAttachUserGroupApi                *UserAttachUserGroupApi
	UserRemoveUserGroupApi                *UserRemoveUserGroupApi
	UserGroupQueryApi                     *UserGroupQueryApi
	IdpCreateApi                          *IdpCreateApi
	IdpDeleteApi                          *IdpDeleteApi
	IdpUpdateApi                          *IdpUpdateApi
	IdpListApi                            *IdpListApi
	PrivilegeGetApi                       *PrivilegeGetApi
	EnterpriseProjectAssignmentToGroupApi *EnterpriseProjectAssignmentToGroupApi
	EnterpriseProjectCreateApi            *EnterpriseProjectCreateApi
	EnterpriseProjectUpdateApi            *EnterpriseProjectUpdateApi
	EnterpriseProjectGetApi               *EnterpriseProjectGetApi
	EnterpriseProjectStatusUpdateApi      *EnterpriseProjectStatusUpdateApi
	EnterpriseProjectSetGroupPolicyApi    *EnterpriseProjectSetGroupPolicyApi
	EnterpriseProjectGetPolicyApi         *EnterpriseProjectGetPolicyApi
	EnterpriseProjectRemoveGroupApi       *EnterpriseProjectRemoveGroupApi
}

// NewApis 构建
func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtiamTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentTest, EndpointCtiamTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndpointCtiamProd)
	return &Apis{
		ServiceListApi:                        NewServiceListApi(client),
		AuthorityListApi:                      NewAuthorityListApi(client),
		PolicyCreateApi:                       NewPolicyCreateApi(client),
		PolicyUpdateApi:                       NewPolicyUpdateApi(client),
		PolicyDeleteApi:                       NewPolicyDeleteApi(client),
		PolicyGetApi:                          NewPolicyGetApi(client),
		PolicyAttachUserGroupApi:              NewPolicyAttachUserGroupApi(client),
		PolicyAttachUserApi:                   NewPolicyAttachUserApi(client),
		PolicyInvalidUserGroupApi:             NewPolicyInvalidUserGroupApi(client),
		UserCreateApi:                         NewUserCreateApi(client),
		UserGetApi:                            NewUserGetApi(client),
		UserUpdateApi:                         NewUserUpdateApi(client),
		UserInvalidApi:                        NewUserInvalidApi(client),
		UserResetPasswordApi:                  NewUserResetPasswordApi(client),
		UserGroupCreateApi:                    NewUserGroupCreateApi(client),
		UserGroupUpdateApi:                    NewUserGroupUpdateApi(client),
		UserGroupInvalidApi:                   NewUserGroupInvalidApi(client),
		UserGroupGetApi:                       NewUserGroupGetApi(client),
		UserAssociationGroupApi:               NewUserAssociationGroupApi(client),
		UserAttachUserGroupApi:                NewUserAttachUserGroupApi(client),
		UserRemoveUserGroupApi:                NewUserRemoveUserGroupApi(client),
		UserGroupQueryApi:                     NewUserGroupQueryApi(client),
		IdpCreateApi:                          NewIdpCreateApi(client),
		IdpDeleteApi:                          NewIdpDeleteApi(client),
		IdpUpdateApi:                          NewIdpUpdateApi(client),
		IdpListApi:                            NewIdpListApi(client),
		PrivilegeGetApi:                       NewPrivilegeGetApi(client),
		EnterpriseProjectAssignmentToGroupApi: NewEpAssignmentToGroupApi(client),
		EnterpriseProjectCreateApi:            NewEnterpriseProjectCreateApi(client),
		EnterpriseProjectUpdateApi:            NewEnterpriseProjectUpdateApi(client),
		EnterpriseProjectGetApi:               NewEnterpriseProjectGetApi(client),
		EnterpriseProjectStatusUpdateApi:      NewEnterpriseProjectStatusUpdateApi(client),
		EnterpriseProjectSetGroupPolicyApi:    NewEnterpriseProjectSetGroupPolicyApi(client),
		EnterpriseProjectGetPolicyApi:         NewEnterpriseProjectGetPolicyApi(client),
		EnterpriseProjectRemoveGroupApi:       NewEnterpriseProjectRemoveGroupApi(client),
	}
}
