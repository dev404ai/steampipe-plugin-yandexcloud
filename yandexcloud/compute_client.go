package yandexcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type ZoneID string
type PlatformID string
type ServiceAccountID string
type Hostname string

type Instance struct {
	Id                 string                       `json:"id"`
	Name               string                       `json:"name"`
	ZoneId             ZoneID                       `json:"zoneId"`
	Status             Status                       `json:"status"`
	FolderId           string                       `json:"folderId"`
	Description        string                       `json:"description"`
	Labels             map[string]string            `json:"labels"`
	PlatformId         PlatformID                   `json:"platformId"`
	CreatedAt          string                       `json:"createdAt"`
	Resources          map[ResourceType]interface{} `json:"resources"`
	Metadata           map[MetadataKey]string       `json:"metadata"`
	MetadataOptions    map[string]interface{}       `json:"metadataOptions"`
	BootDisk           map[string]interface{}       `json:"bootDisk"`
	SecondaryDisks     []map[string]interface{}     `json:"secondaryDisks"`
	NetworkInterfaces  []map[string]interface{}     `json:"networkInterfaces"`
	FQDN               string                       `json:"fqdn"`
	ServiceAccountId   ServiceAccountID             `json:"serviceAccountId"`
	Hostname           Hostname                     `json:"hostname"`
	DeletionProtection bool                         `json:"deletionProtection"`
}

type InstanceID string
type FolderID string
type RetryCount int
type TimeoutSec int64
type Status string
type ResourceType string
type MetadataKey string
type Filter string
type PageToken string
type PageSize int64

const (
	StatusRunning      Status = "RUNNING"
	StatusStopped      Status = "STOPPED"
	StatusProvisioning Status = "PROVISIONING"
	StatusStarting     Status = "STARTING"
	StatusStopping     Status = "STOPPING"
	StatusUnknown      Status = "UNKNOWN"
	// ... other statuses as needed
)

type ComputeClient interface {
	ListInstances(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Instance, PageToken, error)
	GetInstance(ctx context.Context, instanceID InstanceID, timeout TimeoutSec, retry RetryCount) (*Instance, error)
	ListSnapshots(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Snapshot, PageToken, error)
	GetSnapshot(ctx context.Context, snapshotID SnapshotID, timeout TimeoutSec, retry RetryCount) (*Snapshot, error)
	ListImages(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Image, PageToken, error)
	GetImage(ctx context.Context, imageID ImageID, timeout TimeoutSec, retry RetryCount) (*Image, error)
	ListDisks(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Disk, PageToken, error)
	GetDisk(ctx context.Context, diskID DiskID, timeout TimeoutSec, retry RetryCount) (*Disk, error)
	ListFilesystems(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Filesystem, PageToken, error)
	GetFilesystem(ctx context.Context, filesystemID FilesystemID, timeout TimeoutSec, retry RetryCount) (*Filesystem, error)
	ListPlacementGroups(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*PlacementGroup, PageToken, error)
	GetPlacementGroup(ctx context.Context, placementGroupID PlacementGroupID, timeout TimeoutSec, retry RetryCount) (*PlacementGroup, error)
	ListHostGroups(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*HostGroup, PageToken, error)
	GetHostGroup(ctx context.Context, hostGroupID HostGroupID, timeout TimeoutSec, retry RetryCount) (*HostGroup, error)
	ListGPUClusters(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*GPUCluster, PageToken, error)
	GetGPUCluster(ctx context.Context, gpuClusterID GPUClusterID, timeout TimeoutSec, retry RetryCount) (*GPUCluster, error)
	ListDiskPlacementGroups(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*DiskPlacementGroup, PageToken, error)
	GetDiskPlacementGroup(ctx context.Context, diskPlacementGroupID DiskPlacementGroupID, timeout TimeoutSec, retry RetryCount) (*DiskPlacementGroup, error)
	ListSnapshotSchedules(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*SnapshotSchedule, PageToken, error)
	GetSnapshotSchedule(ctx context.Context, snapshotScheduleID SnapshotScheduleID, timeout TimeoutSec, retry RetryCount) (*SnapshotSchedule, error)
	ListReservedInstancePools(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*ReservedInstancePool, PageToken, error)
	GetReservedInstancePool(ctx context.Context, reservedInstancePoolID ReservedInstancePoolID, timeout TimeoutSec, retry RetryCount) (*ReservedInstancePool, error)
	ListZones(ctx context.Context, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Zone, PageToken, error)
	ListDiskTypes(ctx context.Context, zoneID string, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*DiskType, PageToken, error)
	ListHostTypes(ctx context.Context, zoneID string, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*HostType, PageToken, error)
	ListOperations(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Operation, PageToken, error)
	GetOperation(ctx context.Context, operationID OperationID, timeout TimeoutSec, retry RetryCount) (*Operation, error)
	Token() string
}

type ListInstancesResponse struct {
	Instances     []*Instance `json:"instances"`
	NextPageToken string      `json:"nextPageToken"`
}

type GetInstanceResponse struct {
	Instance *Instance `json:"instance"`
}

type ListSnapshotsResponse struct {
	Snapshots     []*Snapshot `json:"snapshots"`
	NextPageToken string      `json:"nextPageToken"`
}

type GetSnapshotResponse struct {
	Snapshot *Snapshot `json:"snapshot"`
}

type Snapshot struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	FolderId     string            `json:"folderId"`
	ZoneId       string            `json:"zoneId"`
	Status       string            `json:"status"`
	CreatedAt    string            `json:"createdAt"`
	SourceDiskId string            `json:"sourceDiskId"`
	Size         string            `json:"size"`
	Labels       map[string]string `json:"labels"`
}

type SnapshotID string

type Image struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	Family      string            `json:"family"`
	ProductIds  []string          `json:"productIds"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	MinDiskSize string            `json:"minDiskSize"`
	Size        string            `json:"size"`
	OsType      string            `json:"osType"`
	OsVersion   string            `json:"osVersion"`
	Labels      map[string]string `json:"labels"`
}

type ImageID string

type ListImagesResponse struct {
	Images        []*Image `json:"images"`
	NextPageToken string   `json:"nextPageToken"`
}

type GetImageResponse struct {
	Image *Image `json:"image"`
}

type Disk struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	FolderId         string            `json:"folderId"`
	ZoneId           string            `json:"zoneId"`
	TypeId           string            `json:"typeId"`
	Size             string            `json:"size"`
	Status           string            `json:"status"`
	CreatedAt        string            `json:"createdAt"`
	SourceImageId    string            `json:"sourceImageId"`
	SourceSnapshotId string            `json:"sourceSnapshotId"`
	BlockSize        string            `json:"blockSize"`
	Labels           map[string]string `json:"labels"`
}

type DiskID string

type ListDisksResponse struct {
	Disks         []*Disk `json:"disks"`
	NextPageToken string  `json:"nextPageToken"`
}

type GetDiskResponse struct {
	Disk *Disk `json:"disk"`
}

type Filesystem struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	ZoneId      string            `json:"zoneId"`
	TypeId      string            `json:"typeId"`
	Size        string            `json:"size"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
}

type FilesystemID string

type ListFilesystemsResponse struct {
	Filesystems   []*Filesystem `json:"filesystems"`
	NextPageToken string        `json:"nextPageToken"`
}

type GetFilesystemResponse struct {
	Filesystem *Filesystem `json:"filesystem"`
}

type PlacementGroup struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	ZoneId      string            `json:"zoneId"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
}

type PlacementGroupID string

type ListPlacementGroupsResponse struct {
	PlacementGroups []*PlacementGroup `json:"placementGroups"`
	NextPageToken   string            `json:"nextPageToken"`
}

type GetPlacementGroupResponse struct {
	PlacementGroup *PlacementGroup `json:"placementGroup"`
}

type HostGroup struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	ZoneId      string            `json:"zoneId"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
}

type HostGroupID string

type ListHostGroupsResponse struct {
	HostGroups    []*HostGroup `json:"hostGroups"`
	NextPageToken string       `json:"nextPageToken"`
}

type GetHostGroupResponse struct {
	HostGroup *HostGroup `json:"hostGroup"`
}

type GPUCluster struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	ZoneId      string            `json:"zoneId"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
}

type GPUClusterID string

type ListGPUClustersResponse struct {
	GPUClusters   []*GPUCluster `json:"gpuClusters"`
	NextPageToken string        `json:"nextPageToken"`
}

type GetGPUClusterResponse struct {
	GPUCluster *GPUCluster `json:"gpuCluster"`
}

type DiskPlacementGroup struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	ZoneId      string            `json:"zoneId"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
}

type DiskPlacementGroupID string

type ListDiskPlacementGroupsResponse struct {
	DiskPlacementGroups []*DiskPlacementGroup `json:"diskPlacementGroups"`
	NextPageToken       string                `json:"nextPageToken"`
}

type GetDiskPlacementGroupResponse struct {
	DiskPlacementGroup *DiskPlacementGroup `json:"diskPlacementGroup"`
}

type SnapshotSchedule struct {
	Id             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	FolderId       string                 `json:"folderId"`
	Status         string                 `json:"status"`
	CreatedAt      string                 `json:"createdAt"`
	SchedulePolicy map[string]interface{} `json:"schedulePolicy"`
	Labels         map[string]string      `json:"labels"`
}

type SnapshotScheduleID string

type ListSnapshotSchedulesResponse struct {
	SnapshotSchedules []*SnapshotSchedule `json:"snapshotSchedules"`
	NextPageToken     string              `json:"nextPageToken"`
}

type GetSnapshotScheduleResponse struct {
	SnapshotSchedule *SnapshotSchedule `json:"snapshotSchedule"`
}

type ReservedInstancePool struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FolderId    string            `json:"folderId"`
	ZoneId      string            `json:"zoneId"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
}

type ReservedInstancePoolID string

type ListReservedInstancePoolsResponse struct {
	ReservedInstancePools []*ReservedInstancePool `json:"reservedInstancePools"`
	NextPageToken         string                  `json:"nextPageToken"`
}

type GetReservedInstancePoolResponse struct {
	ReservedInstancePool *ReservedInstancePool `json:"reservedInstancePool"`
}

type Zone struct {
	Id       string `json:"id"`
	RegionId string `json:"regionId"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type ListZonesResponse struct {
	Zones         []*Zone `json:"zones"`
	NextPageToken string  `json:"nextPageToken"`
}

type DiskType struct {
	Id          string `json:"id"`
	ZoneId      string `json:"zoneId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DiskTypeID string

type ListDiskTypesResponse struct {
	DiskTypes     []*DiskType `json:"diskTypes"`
	NextPageToken string      `json:"nextPageToken"`
}

type HostType struct {
	Id          string `json:"id"`
	ZoneId      string `json:"zoneId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type HostTypeID string

type ListHostTypesResponse struct {
	HostTypes     []*HostType `json:"hostTypes"`
	NextPageToken string      `json:"nextPageToken"`
}

type Operation struct {
	Id          string                 `json:"id"`
	Description string                 `json:"description"`
	CreatedAt   string                 `json:"createdAt"`
	Done        bool                   `json:"done"`
	Status      string                 `json:"status"`
	FolderId    string                 `json:"folderId"`
	Error       map[string]interface{} `json:"error"`
	Response    map[string]interface{} `json:"response"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type OperationID string

type ListOperationsResponse struct {
	Operations    []*Operation `json:"operations"`
	NextPageToken string       `json:"nextPageToken"`
}

type GetOperationResponse struct {
	Operation *Operation `json:"operation"`
}

type yandexComputeClient struct {
	iamToken string
	http     *http.Client
	config   *Config
}

func NewComputeClient(token string, timeoutSec int64, config *Config) ComputeClient {
	return &yandexComputeClient{
		iamToken: token,
		http:     GetHTTPClient(timeoutSec),
		config:   config,
	}
}

func (c *yandexComputeClient) Token() string { return c.iamToken }

func (c *yandexComputeClient) apiGet(ctx context.Context, urlStr string, out interface{}, timeoutSec TimeoutSec, retryCount RetryCount) error {
	LogInfo(ctx, "Compute apiGet: %s", urlStr)
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	if retryCount < 1 {
		retryCount = 1
	}
	reqFactory := func() *http.Request {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
		req.Header.Set("Authorization", "Bearer "+c.iamToken)
		var ua *string
		if c.config.UserAgent != nil {
			s := string(*c.config.UserAgent)
			ua = &s
		}
		var eo *string
		if c.config.EndpointOverride != nil {
			s := string(*c.config.EndpointOverride)
			eo = &s
		}
		ApplyRequestOptions(req, ua, eo)
		LogDebug(ctx, "HTTP request URL: %s", req.URL.String())
		return req
	}
	resp, err := DoWithRetry(ctx, c.http, reqFactory, int(retryCount), int64(timeoutSec))
	if err != nil {
		LogError(ctx, "Compute GET request failed: %v", err)
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := HandleHTTPError(resp); err != nil {
		LogError(ctx, "Compute GET HTTP error: %v", err)
		return err
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		LogError(ctx, "Failed to decode compute response: %v", err)
		return err
	}
	LogInfo(ctx, "Compute apiGet success: %s", urlStr)
	return nil
}

func (c *yandexComputeClient) ListInstances(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Instance, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/instances"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListInstancesResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Instances, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetInstance(ctx context.Context, id InstanceID, timeout TimeoutSec, retry RetryCount) (*Instance, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/instances/%s?view=FULL", id)
	var respBody GetInstanceResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.Instance, nil
}

func (c *yandexComputeClient) ListSnapshots(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Snapshot, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/snapshots"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListSnapshotsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Snapshots, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetSnapshot(ctx context.Context, id SnapshotID, timeout TimeoutSec, retry RetryCount) (*Snapshot, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/snapshots/%s", id)
	var respBody GetSnapshotResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.Snapshot, nil
}

func (c *yandexComputeClient) ListImages(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Image, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/images"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListImagesResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Images, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetImage(ctx context.Context, id ImageID, timeout TimeoutSec, retry RetryCount) (*Image, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/images/%s", id)
	var respBody GetImageResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.Image, nil
}

func (c *yandexComputeClient) ListDisks(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Disk, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/disks"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListDisksResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Disks, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetDisk(ctx context.Context, id DiskID, timeout TimeoutSec, retry RetryCount) (*Disk, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/disks/%s", id)
	var respBody GetDiskResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.Disk, nil
}

func (c *yandexComputeClient) ListFilesystems(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Filesystem, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/filesystems"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListFilesystemsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Filesystems, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetFilesystem(ctx context.Context, id FilesystemID, timeout TimeoutSec, retry RetryCount) (*Filesystem, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/filesystems/%s", id)
	var respBody GetFilesystemResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.Filesystem, nil
}

func (c *yandexComputeClient) ListPlacementGroups(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*PlacementGroup, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/placementGroups"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListPlacementGroupsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.PlacementGroups, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetPlacementGroup(ctx context.Context, id PlacementGroupID, timeout TimeoutSec, retry RetryCount) (*PlacementGroup, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/placementGroups/%s", id)
	var respBody GetPlacementGroupResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.PlacementGroup, nil
}

func (c *yandexComputeClient) ListHostGroups(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*HostGroup, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/hostGroups"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListHostGroupsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.HostGroups, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetHostGroup(ctx context.Context, id HostGroupID, timeout TimeoutSec, retry RetryCount) (*HostGroup, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/hostGroups/%s", id)
	var respBody GetHostGroupResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.HostGroup, nil
}

func (c *yandexComputeClient) ListGPUClusters(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*GPUCluster, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/gpuClusters"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListGPUClustersResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.GPUClusters, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetGPUCluster(ctx context.Context, id GPUClusterID, timeout TimeoutSec, retry RetryCount) (*GPUCluster, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/gpuClusters/%s", id)
	var respBody GetGPUClusterResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.GPUCluster, nil
}

func (c *yandexComputeClient) ListDiskPlacementGroups(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*DiskPlacementGroup, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/diskPlacementGroups"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListDiskPlacementGroupsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.DiskPlacementGroups, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetDiskPlacementGroup(ctx context.Context, id DiskPlacementGroupID, timeout TimeoutSec, retry RetryCount) (*DiskPlacementGroup, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/diskPlacementGroups/%s", id)
	var respBody GetDiskPlacementGroupResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.DiskPlacementGroup, nil
}

func (c *yandexComputeClient) ListSnapshotSchedules(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*SnapshotSchedule, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/snapshotSchedules"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListSnapshotSchedulesResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.SnapshotSchedules, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetSnapshotSchedule(ctx context.Context, id SnapshotScheduleID, timeout TimeoutSec, retry RetryCount) (*SnapshotSchedule, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/snapshotSchedules/%s", id)
	var respBody GetSnapshotScheduleResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.SnapshotSchedule, nil
}

func (c *yandexComputeClient) ListReservedInstancePools(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*ReservedInstancePool, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/reservedInstancePools"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListReservedInstancePoolsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.ReservedInstancePools, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetReservedInstancePool(ctx context.Context, id ReservedInstancePoolID, timeout TimeoutSec, retry RetryCount) (*ReservedInstancePool, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/reservedInstancePools/%s", id)
	var respBody GetReservedInstancePoolResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.ReservedInstancePool, nil
}

func (c *yandexComputeClient) ListZones(ctx context.Context, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Zone, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/zones"
	params := url.Values{}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListZonesResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Zones, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) ListDiskTypes(ctx context.Context, zoneID string, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*DiskType, PageToken, error) {
	endpoint := "https://compute.api.cloud.yandex.net/compute/v1/diskTypes"
	params := url.Values{}
	if zoneID != "" {
		params.Set("zoneId", zoneID)
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListDiskTypesResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.DiskTypes, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) ListHostTypes(ctx context.Context, zoneID string, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*HostType, PageToken, error) {
	endpoint := "https://compute.api.cloud.yandex.net/compute/v1/hostTypes"
	params := url.Values{}
	if zoneID != "" {
		params.Set("zoneId", zoneID)
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListHostTypesResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.HostTypes, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) ListOperations(ctx context.Context, folderID FolderID, filter Filter, pageToken PageToken, pageSize PageSize, timeout TimeoutSec, retry RetryCount) ([]*Operation, PageToken, error) {
	const endpoint = "https://compute.api.cloud.yandex.net/compute/v1/operations"
	params := url.Values{}
	params.Set("folderId", string(folderID))
	if filter != "" {
		params.Set("filter", string(filter))
	}
	if pageToken != "" {
		params.Set("pageToken", string(pageToken))
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(int64(pageSize), 10))
	}
	var respBody ListOperationsResponse
	err := c.apiGet(ctx, fmt.Sprintf("%s?%s", endpoint, params.Encode()), &respBody, timeout, retry)
	if err != nil {
		return nil, "", err
	}
	return respBody.Operations, PageToken(respBody.NextPageToken), nil
}

func (c *yandexComputeClient) GetOperation(ctx context.Context, id OperationID, timeout TimeoutSec, retry RetryCount) (*Operation, error) {
	urlStr := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/operations/%s", id)
	var respBody GetOperationResponse
	if err := c.apiGet(ctx, urlStr, &respBody, timeout, retry); err != nil {
		return nil, err
	}
	return respBody.Operation, nil
}
