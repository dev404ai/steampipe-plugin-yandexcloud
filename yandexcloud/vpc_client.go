package yandexcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type VPCNetwork struct {
	Id          string            `json:"id"`
	FolderId    string            `json:"folderId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   string            `json:"createdAt"`
}

type VPCNetworkID string

type ListVPCNetworksResponse struct {
	Networks      []*VPCNetwork `json:"networks"`
	NextPageToken string        `json:"nextPageToken"`
}

type GetVPCNetworkResponse struct {
	Network *VPCNetwork `json:"network"`
}

type VPCSubnet struct {
	Id          string            `json:"id"`
	FolderId    string            `json:"folderId"`
	NetworkId   string            `json:"networkId"`
	ZoneId      string            `json:"zoneId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   string            `json:"createdAt"`
	CidrBlocks  []string          `json:"v4CidrBlocks"`
}

type VPCSubnetID string

type ListVPCSubnetsResponse struct {
	Subnets       []*VPCSubnet `json:"subnets"`
	NextPageToken string       `json:"nextPageToken"`
}

type GetVPCSubnetResponse struct {
	Subnet *VPCSubnet `json:"subnet"`
}

type VPCRouteTable struct {
	Id           string                   `json:"id"`
	FolderId     string                   `json:"folderId"`
	NetworkId    string                   `json:"networkId"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Labels       map[string]string        `json:"labels"`
	CreatedAt    string                   `json:"createdAt"`
	StaticRoutes []map[string]interface{} `json:"staticRoutes"`
}

type VPCRouteTableID string

type ListVPCRouteTablesResponse struct {
	RouteTables   []*VPCRouteTable `json:"routeTables"`
	NextPageToken string           `json:"nextPageToken"`
}

type GetVPCRouteTableResponse struct {
	RouteTable *VPCRouteTable `json:"routeTable"`
}

type VPCSecurityGroup struct {
	Id           string                   `json:"id"`
	FolderId     string                   `json:"folderId"`
	NetworkId    string                   `json:"networkId"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Labels       map[string]string        `json:"labels"`
	CreatedAt    string                   `json:"createdAt"`
	Rules        []map[string]interface{} `json:"rules"`
	IngressRules []map[string]interface{} `json:"ingressRules"`
	EgressRules  []map[string]interface{} `json:"egressRules"`
}

type VPCSecurityGroupID string

type ListVPCSecurityGroupsResponse struct {
	SecurityGroups []*VPCSecurityGroup `json:"securityGroups"`
	NextPageToken  string              `json:"nextPageToken"`
}

type GetVPCSecurityGroupResponse struct {
	SecurityGroup *VPCSecurityGroup `json:"securityGroup"`
}

// --- VPC Address types ---
type VPCAddress struct {
	Id                  string               `json:"id"`
	FolderId            string               `json:"folderId"`
	CreatedAt           string               `json:"createdAt"`
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	Labels              map[string]string    `json:"labels"`
	ExternalIpv4Address *ExternalIpv4Address `json:"externalIpv4Address,omitempty"`
	Reserved            bool                 `json:"reserved"`
	Used                bool                 `json:"used"`
	Type                string               `json:"type"`
	IpVersion           string               `json:"ipVersion"`
	DeletionProtection  bool                 `json:"deletionProtection"`
	DnsRecords          []DnsRecord          `json:"dnsRecords"`
}

type ExternalIpv4Address struct {
	Address      string               `json:"address"`
	ZoneId       string               `json:"zoneId"`
	Requirements *AddressRequirements `json:"requirements,omitempty"`
}

type AddressRequirements struct {
	DdosProtectionProvider string `json:"ddosProtectionProvider"`
	OutgoingSmtpCapability string `json:"outgoingSmtpCapability"`
}

type DnsRecord struct {
	Fqdn      string `json:"fqdn"`
	DnsZoneId string `json:"dnsZoneId"`
	Ttl       string `json:"ttl"`
	Ptr       bool   `json:"ptr"`
}

type VPCAddressID string

type ListVPCAddressesResponse struct {
	Addresses     []*VPCAddress `json:"addresses"`
	NextPageToken string        `json:"nextPageToken"`
}

type GetVPCAddressResponse struct {
	Address *VPCAddress `json:"address"`
}

// --- VPC Gateway types ---
type VPCGateway struct {
	Id                  string                 `json:"id"`
	FolderId            string                 `json:"folderId"`
	CreatedAt           string                 `json:"createdAt"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Labels              map[string]string      `json:"labels"`
	SharedEgressGateway map[string]interface{} `json:"sharedEgressGateway"`
}

type VPCGatewayID string

type ListVPCGatewaysResponse struct {
	Gateways      []*VPCGateway `json:"gateways"`
	NextPageToken string        `json:"nextPageToken"`
}

type GetVPCGatewayResponse struct {
	Gateway *VPCGateway `json:"gateway"`
}

// --- VPC Operation types ---
type VPCOperation struct {
	Id          string                 `json:"id"`
	Description string                 `json:"description"`
	CreatedAt   string                 `json:"createdAt"`
	CreatedBy   string                 `json:"createdBy"`
	ModifiedAt  string                 `json:"modifiedAt"`
	Done        bool                   `json:"done"`
	Metadata    map[string]interface{} `json:"metadata"`
	Error       map[string]interface{} `json:"error"`
	Response    map[string]interface{} `json:"response"`
}

type VPCOperationID string

type ListVPCOperationsResponse struct {
	Operations    []*VPCOperation `json:"operations"`
	NextPageToken string          `json:"nextPageToken"`
}

type GetVPCOperationResponse struct {
	Operation *VPCOperation `json:"operation"`
}

type VPCClient interface {
	ListVPCNetworks(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCNetwork, string, error)
	GetVPCNetwork(ctx context.Context, networkID VPCNetworkID) (*VPCNetwork, error)
	ListVPCSubnets(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCSubnet, string, error)
	GetVPCSubnet(ctx context.Context, subnetID VPCSubnetID) (*VPCSubnet, error)
	ListVPCRouteTables(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCRouteTable, string, error)
	GetVPCRouteTable(ctx context.Context, routeTableID VPCRouteTableID) (*VPCRouteTable, error)
	ListVPCSecurityGroups(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCSecurityGroup, string, error)
	GetVPCSecurityGroup(ctx context.Context, securityGroupID VPCSecurityGroupID) (*VPCSecurityGroup, error)
	ListVPCAddresses(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCAddress, string, error)
	GetVPCAddress(ctx context.Context, addressID VPCAddressID) (*VPCAddress, error)
	ListVPCGateways(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCGateway, string, error)
	GetVPCGateway(ctx context.Context, gatewayID VPCGatewayID) (*VPCGateway, error)
	ListVPCOperations(ctx context.Context, pageToken string, pageSize int64) ([]*VPCOperation, string, error)
	GetVPCOperation(ctx context.Context, operationID VPCOperationID) (*VPCOperation, error)
}

type yandexVPCClient struct {
	token  string
	http   *http.Client
	config *Config
}

func NewVPCClient(token string, timeoutSec int64, config *Config) VPCClient {
	return &yandexVPCClient{
		token:  token,
		http:   GetHTTPClient(timeoutSec),
		config: config,
	}
}

func (c *yandexVPCClient) apiGet(ctx context.Context, urlStr string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		LogError(ctx, "VPC apiGet: failed to create request: %v", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		LogError(ctx, "VPC apiGet: request failed: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		LogError(ctx, "VPC API error: %s", resp.Status)
		return fmt.Errorf("VPC API error: %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(io.Reader(io.MultiReader(io.NopCloser(io.Reader(io.MultiReader())))))
	if err := json.Unmarshal(body, out); err != nil {
		LogError(ctx, "VPC apiGet: failed to decode response: %v", err)
		return err
	}
	LogInfo(ctx, "VPC apiGet success: %s", urlStr)
	return nil
}

func (c *yandexVPCClient) ListVPCNetworks(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCNetwork, string, error) {
	const endpoint = "https://vpc.api.cloud.yandex.net/vpc/v1/networks"
	params := url.Values{}
	params.Set("folderId", folderID)
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCNetworksResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.Networks, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCNetwork(ctx context.Context, networkID VPCNetworkID) (*VPCNetwork, error) {
	urlStr := fmt.Sprintf("https://vpc.api.cloud.yandex.net/vpc/v1/networks/%s", networkID)
	var respBody GetVPCNetworkResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.Network, nil
}

func (c *yandexVPCClient) ListVPCSubnets(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCSubnet, string, error) {
	const endpoint = "https://vpc.api.cloud.yandex.net/vpc/v1/subnets"
	params := url.Values{}
	params.Set("folderId", folderID)
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCSubnetsResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.Subnets, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCSubnet(ctx context.Context, subnetID VPCSubnetID) (*VPCSubnet, error) {
	urlStr := fmt.Sprintf("https://vpc.api.cloud.yandex.net/vpc/v1/subnets/%s", subnetID)
	var respBody GetVPCSubnetResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.Subnet, nil
}

func (c *yandexVPCClient) ListVPCRouteTables(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCRouteTable, string, error) {
	const endpoint = "https://vpc.api.cloud.yandex.net/vpc/v1/routeTables"
	params := url.Values{}
	params.Set("folderId", folderID)
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCRouteTablesResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.RouteTables, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCRouteTable(ctx context.Context, routeTableID VPCRouteTableID) (*VPCRouteTable, error) {
	urlStr := fmt.Sprintf("https://vpc.api.cloud.yandex.net/vpc/v1/routeTables/%s", routeTableID)
	var respBody GetVPCRouteTableResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.RouteTable, nil
}

func (c *yandexVPCClient) ListVPCSecurityGroups(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCSecurityGroup, string, error) {
	const endpoint = "https://vpc.api.cloud.yandex.net/vpc/v1/securityGroups"
	params := url.Values{}
	params.Set("folderId", folderID)
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCSecurityGroupsResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.SecurityGroups, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCSecurityGroup(ctx context.Context, securityGroupID VPCSecurityGroupID) (*VPCSecurityGroup, error) {
	urlStr := fmt.Sprintf("https://vpc.api.cloud.yandex.net/vpc/v1/securityGroups/%s", securityGroupID)
	var respBody GetVPCSecurityGroupResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.SecurityGroup, nil
}

func (c *yandexVPCClient) ListVPCAddresses(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCAddress, string, error) {
	const endpoint = "https://vpc.api.cloud.yandex.net/vpc/v1/addresses"
	params := url.Values{}
	params.Set("folderId", folderID)
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCAddressesResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.Addresses, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCAddress(ctx context.Context, addressID VPCAddressID) (*VPCAddress, error) {
	urlStr := fmt.Sprintf("https://vpc.api.cloud.yandex.net/vpc/v1/addresses/%s", addressID)
	var respBody GetVPCAddressResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.Address, nil
}

func (c *yandexVPCClient) ListVPCGateways(ctx context.Context, folderID string, pageToken string, pageSize int64) ([]*VPCGateway, string, error) {
	const endpoint = "https://vpc.api.cloud.yandex.net/vpc/v1/gateways"
	params := url.Values{}
	params.Set("folderId", folderID)
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCGatewaysResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.Gateways, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCGateway(ctx context.Context, gatewayID VPCGatewayID) (*VPCGateway, error) {
	urlStr := fmt.Sprintf("https://vpc.api.cloud.yandex.net/vpc/v1/gateways/%s", gatewayID)
	var respBody GetVPCGatewayResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.Gateway, nil
}

func (c *yandexVPCClient) ListVPCOperations(ctx context.Context, pageToken string, pageSize int64) ([]*VPCOperation, string, error) {
	const endpoint = "https://operation.api.cloud.yandex.net/operations"
	params := url.Values{}
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	var respBody ListVPCOperationsResponse
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	err := c.apiGet(ctx, urlStr, &respBody)
	if err != nil {
		return nil, "", err
	}
	return respBody.Operations, respBody.NextPageToken, nil
}

func (c *yandexVPCClient) GetVPCOperation(ctx context.Context, operationID VPCOperationID) (*VPCOperation, error) {
	urlStr := fmt.Sprintf("https://operation.api.cloud.yandex.net/operations/%s", operationID)
	var respBody GetVPCOperationResponse
	if err := c.apiGet(ctx, urlStr, &respBody); err != nil {
		return nil, err
	}
	return respBody.Operation, nil
}
