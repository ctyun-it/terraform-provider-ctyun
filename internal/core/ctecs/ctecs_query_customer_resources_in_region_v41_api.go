package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryCustomerResourcesInRegionV41Api
/* 根据regionID查询用户已有资源
 */type CtecsQueryCustomerResourcesInRegionV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryCustomerResourcesInRegionV41Api(client *core.CtyunClient) *CtecsQueryCustomerResourcesInRegionV41Api {
	return &CtecsQueryCustomerResourcesInRegionV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/region/customer-resources",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryCustomerResourcesInRegionV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryCustomerResourcesInRegionV41Request) (*CtecsQueryCustomerResourcesInRegionV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryCustomerResourcesInRegionV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryCustomerResourcesInRegionV41Request struct {
	RegionID string /*  资源池ID  */
}

type CtecsQueryCustomerResourcesInRegionV41Response struct {
	StatusCode  int32                                                    `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                                   `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                                   `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQueryCustomerResourcesInRegionV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResponse struct {
	Resources *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesResponse `json:"resources"` /*  资源信息  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesResponse struct {
	VM              *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVMResponse              `json:"VM"`              /*  云主机  */
	Volume          *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVolumeResponse          `json:"Volume"`          /*  磁盘  */
	VOLUME_SNAPSHOT *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVOLUME_SNAPSHOTResponse `json:"VOLUME_SNAPSHOT"` /*  磁盘快照  */
	VPC             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVPCResponse             `json:"VPC"`             /*  VPC  */
	Public_IP       *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesPublic_IPResponse       `json:"Public_IP"`       /*  公网IP  */
	BMS             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesBMSResponse             `json:"BMS"`             /*  物理机  */
	NAT             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesNATResponse             `json:"NAT"`             /*  NAT  */
	Disk_Backup     *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesDisk_BackupResponse     `json:"Disk_Backup"`     /*  磁盘备份  */
	Vm_Group        *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVm_GroupResponse        `json:"Vm_Group"`        /*  云主机组  */
	SNAPSHOT        *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesSNAPSHOTResponse        `json:"SNAPSHOT"`        /*  云主机快照  */
	ACLLIST         *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesACLLISTResponse         `json:"ACLLIST"`         /*  ACL  */
	IP_POOL         *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIP_POOLResponse         `json:"IP_POOL"`         /*  共享带宽  */
	IMAGE           *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIMAGEResponse           `json:"IMAGE"`           /*  私有镜像  */
	LB_LISTENER     *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLB_LISTENERResponse     `json:"LB_LISTENER"`     /*  负载均衡监听器  */
	LOADBALANCER    *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLOADBALANCERResponse    `json:"LOADBALANCER"`    /*  负载均衡  */
	OS_Backup       *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesOS_BackupResponse       `json:"OS_Backup"`       /*  操作系统备份  */
	CBR             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBRResponse             `json:"CBR"`             /*  云主机备份  */
	CERT            *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCERTResponse            `json:"CERT"`            /*  负载均衡证书  */
	CBR_VBS         *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBR_VBSResponse         `json:"CBR_VBS"`         /*  磁盘存储备份  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVMResponse struct {
	Vm_shutd_count       int32 `json:"vm_shutd_count,omitempty"`       /*  已关机云主机数量  */
	Expire_count         int32 `json:"expire_count,omitempty"`         /*  过期云主机数量  */
	Expire_running_count int32 `json:"expire_running_count,omitempty"` /*  已过期的运行中云主机数量  */
	Expire_shutd_count   int32 `json:"expire_shutd_count,omitempty"`   /*  已过期的关机云主机数量  */
	Vm_running_count     int32 `json:"vm_running_count,omitempty"`     /*  运行中云主机数量  */
	Total_count          int32 `json:"total_count,omitempty"`          /*  云主机总数  */
	Cpu_count            int32 `json:"cpu_count,omitempty"`            /*  CPU总数  */
	Memory_count         int32 `json:"memory_count,omitempty"`         /*  总内存大小  */
	Detail_total_count   int32 `json:"detail_total_count,omitempty"`   /*  云主机总数  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVolumeResponse struct {
	Vo_root_count      int32 `json:"vo_root_count,omitempty"`      /*  系统盘数量  */
	Vo_disk_count      int32 `json:"vo_disk_count,omitempty"`      /*  数据盘数量  */
	Total_count        int32 `json:"total_count,omitempty"`        /*  磁盘总数  */
	Detail_total_count int32 `json:"detail_total_count,omitempty"` /*  磁盘总数  */
	Total_size         int32 `json:"total_size,omitempty"`         /*  磁盘总大小  */
	Vo_disk_size       int32 `json:"vo_disk_size,omitempty"`       /*  数据盘大小  */
	Vo_root_size       int32 `json:"vo_root_size,omitempty"`       /*  系统盘大小  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVOLUME_SNAPSHOTResponse struct {
	Total_count        int32 `json:"total_count,omitempty"`        /*  磁盘快照总数  */
	Detail_total_count int32 `json:"detail_total_count,omitempty"` /*  磁盘快照总数  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVPCResponse struct {
	Total_count        int32 `json:"total_count,omitempty"`        /*  VPC总数  */
	Detail_total_count int32 `json:"detail_total_count,omitempty"` /*  VPC总数  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesPublic_IPResponse struct {
	Total_count        int32 `json:"total_count,omitempty"`        /*  公网IP总数  */
	Detail_total_count int32 `json:"detail_total_count,omitempty"` /*  公网IP总数  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesBMSResponse struct {
	Total_count          int32 `json:"total_count,omitempty"`          /*  物理机总数  */
	Detail_total_count   int32 `json:"detail_total_count,omitempty"`   /*  物理机总数  */
	Memory_count         int32 `json:"memory_count,omitempty"`         /*  固定为0  */
	Cpu_count            int32 `json:"cpu_count,omitempty"`            /*  固定为0  */
	Bm_shutd_count       int32 `json:"bm_shutd_count,omitempty"`       /*  固定为0  */
	Expire_running_count int32 `json:"expire_running_count,omitempty"` /*  固定为0  */
	Bm_running_count     int32 `json:"bm_running_count,omitempty"`     /*  固定为0  */
	Expire_count         int32 `json:"expire_count,omitempty"`         /*  固定为0  */
	Expire_shutd_count   int32 `json:"expire_shutd_count,omitempty"`   /*  固定为0  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesNATResponse struct {
	Total_count        int32                                                                      `json:"total_count,omitempty"`        /*  nat总数  */
	Detail_total_count int32                                                                      `json:"detail_total_count,omitempty"` /*  nat总数  */
	Detail             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesNATDetailResponse `json:"detail"`                       /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesDisk_BackupResponse struct {
	Total_count        int32                                                                              `json:"total_count,omitempty"`        /*  磁盘备份总数  */
	Detail_total_count int32                                                                              `json:"detail_total_count,omitempty"` /*  磁盘备份总数  */
	Detail             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesDisk_BackupDetailResponse `json:"detail"`                       /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVm_GroupResponse struct {
	Total_count int32                                                                           `json:"total_count,omitempty"` /*  云主机组总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVm_GroupDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesSNAPSHOTResponse struct {
	Total_count int32                                                                           `json:"total_count,omitempty"` /*  云主机快照总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesSNAPSHOTDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesACLLISTResponse struct {
	Total_count int32                                                                          `json:"total_count,omitempty"` /*  ACL总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesACLLISTDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIP_POOLResponse struct {
	Total_count int32                                                                          `json:"total_count,omitempty"` /*  共享带宽总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIP_POOLDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIMAGEResponse struct {
	Total_count int32                                                                        `json:"total_count,omitempty"` /*  私有镜像总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIMAGEDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLB_LISTENERResponse struct {
	Total_count int32                                                                              `json:"total_count,omitempty"` /*  负载均衡监听器总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLB_LISTENERDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLOADBALANCERResponse struct {
	Total_count int32                                                                               `json:"total_count,omitempty"` /*  负载均衡总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLOADBALANCERDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesOS_BackupResponse struct {
	Total_count        int32 `json:"total_count,omitempty"`        /*  固定为0  */
	Detail_total_count int32 `json:"detail_total_count,omitempty"` /*  固定为0  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBRResponse struct {
	Total_count        int32                                                                      `json:"total_count,omitempty"`        /*  云主机备份总数  */
	Detail_total_count int32                                                                      `json:"detail_total_count,omitempty"` /*  云主机备份总数  */
	Total_size         int32                                                                      `json:"total_size,omitempty"`         /*  云主机备份总大小  */
	Detail             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBRDetailResponse `json:"detail"`                       /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCERTResponse struct {
	Total_count int32                                                                       `json:"total_count,omitempty"` /*  负载均衡证书总数  */
	Detail      *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCERTDetailResponse `json:"detail"`                /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBR_VBSResponse struct {
	Total_count        int32                                                                          `json:"total_count,omitempty"`        /*  磁盘存储备份总数  */
	Detail_total_count int32                                                                          `json:"detail_total_count,omitempty"` /*  磁盘存储备份总数  */
	Detail             *CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBR_VBSDetailResponse `json:"detail"`                       /*  对应资源池id下的数量  */
}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesNATDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesDisk_BackupDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesVm_GroupDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesSNAPSHOTDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesACLLISTDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIP_POOLDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesIMAGEDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLB_LISTENERDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesLOADBALANCERDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBRDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCERTDetailResponse struct{}

type CtecsQueryCustomerResourcesInRegionV41ReturnObjResourcesCBR_VBSDetailResponse struct{}
