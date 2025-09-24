package business

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctzos"
)

type ZosService struct {
	meta *common.CtyunMetadata
}

func NewZosService(meta *common.CtyunMetadata) *ZosService {
	return &ZosService{meta: meta}
}

func (c ZosService) GetZosBucketInfo(ctx context.Context, bucket, regionID string) (bucketRes ctzos.ZosGetBucketInfoReturnObjResponse, err error) {
	params := &ctzos.ZosGetBucketInfoRequest{
		Bucket:   bucket,
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetBucketInfoApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	bucketRes = *resp.ReturnObj
	return
}

func (c ZosService) BuildS3Client(ctx context.Context, regionID string) (*s3.S3, error) {
	ak, sk, err := c.getAccessKey(ctx, regionID)
	if err != nil {
		return nil, err
	}
	endpoint, err := c.getEndpoint(ctx, regionID)
	if err != nil {
		return nil, err
	}

	conf := &aws.Config{
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(ak, sk, ``),
		Region:           aws.String("us-east-1"),
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{Config: *conf}))
	svc := s3.New(sess)
	return svc, nil
}

// getAccessKey 获取zosAK和zosSK
func (c ZosService) getAccessKey(ctx context.Context, regionID string) (ak, sk string, err error) {
	params := &ctzos.ZosGetKeysRequest{
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetKeysApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if len(resp.ReturnObj) == 0 {
		err = common.InvalidReturnObjError
		return
	}
	for _, r := range resp.ReturnObj {
		if r.RegionID == regionID || r.RegionID == "public" {
			ak, sk = r.AccessKey, r.SecretKey
			return
		}
	}
	err = common.InvalidReturnObjResultsError
	return
}

// getEndpoint 获取zosEndpoint
func (c ZosService) getEndpoint(ctx context.Context, regionID string) (endpoint string, err error) {
	params := &ctzos.ZosGetEndpointRequest{
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetEndpointApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.InternetEndpoint) == 0 {
		err = common.InvalidReturnObjResultsError
		return
	}
	endpoint = resp.ReturnObj.InternetEndpoint[0]
	return
}
