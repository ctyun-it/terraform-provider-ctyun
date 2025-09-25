package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbUpdatePortApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdatePortApi(client *ctyunsdk.CtyunClient) *MongodbUpdatePortApi {
	return &MongodbUpdatePortApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/modifyPort",
		},
	}
}

func (this *MongodbUpdatePortApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdatePortRequest, headers *MongodbUpdatePortRequestHeader) (updatedResponse *MongodbUpdatePortResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}
	if headers.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	if req.ProdInstId == "" {
		err = errors.New("missing required field: InstID")
		return
	}
	if req.NewPort == "" {
		err = errors.New("missing required field: Port")
		return
	}
	builder.AddHeader("regionId", headers.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	updatedResponse = &MongodbUpdatePortResponse{}
	err = resp.Parse(updatedResponse)
	if err != nil {
		return
	}
	return updatedResponse, nil
}

type MongodbUpdatePortRequest struct {
	ProdInstId string `json:"prodInstId"`
	NewPort    string `json:"newPort"`
}

type MongodbUpdatePortRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
	RegionID  string  `json:"regionId"`
}

type MongodbUpdatePortResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
