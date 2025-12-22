package business

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctiam"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
)

type IamService struct {
	meta *common.CtyunMetadata
}

func NewIamService(meta *common.CtyunMetadata) *IamService {
	return &IamService{meta: meta}
}

func (c IamService) QueryAkList(ctx context.Context, userID string) (aks []*ctiam.CtiamQueryAkReturnObjAccessKeyUserListAccessKeyListResponse, err error) {
	params := &ctiam.CtiamQueryAkRequest{
		UserIdList: []string{userID},
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamQueryAkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, r := range resp.ReturnObj.AccessKeyUserList {
		if utils.SecString(r.UserId) == userID {
			aks = r.AccessKeyList
			return
		}
	}
	err = fmt.Errorf("not found userID %s", userID)
	return
}

// DecryptSK 解密SK
func (c IamService) DecryptSK(secretSK, ak string) (decrypted string, err error) {
	decodedTxt, err := hex.DecodeString(secretSK)
	if err != nil {
		return
	}

	decodedKey, err := hex.DecodeString(ak)
	if err != nil {
		return
	}

	decrypted, err = utils.Decrypt(decodedTxt, decodedKey)
	if err != nil {
		return
	}
	return
}

func (c IamService) QueryUserList(ctx context.Context, pageNum, pageSize int32) (users *ctiam.CtiamQueryUsersReturnObjResponse, err error) {
	params := &ctiam.CtiamQueryUsersRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamQueryUsersApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	users = resp.ReturnObj
	return
}
