package cloudflare

import (
	"context"
	"fmt"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Client struct {
	api    *cloudflare.API
	zoneID string
}

type DNSRecord struct {
	ID       string
	Type     string
	Name     string
	Content  string
	TTL      int
	Priority *uint16
	Proxied  *bool
}

func NewClient(token, email string) (*Client, error) {
	var api *cloudflare.API
	var err error

	if email != "" {
		// Using API Key (legacy)
		api, err = cloudflare.New(token, email)
	} else {
		// Using API Token
		api, err = cloudflare.NewWithAPIToken(token)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare client: %w", err)
	}

	return &Client{api: api}, nil
}

func (c *Client) SetZone(ctx context.Context, domain string) error {
	// Find zone by name
	zoneID, err := c.api.ZoneIDByName(domain)
	if err != nil {
		return fmt.Errorf("failed to find zone %s: %w", domain, err)
	}

	c.zoneID = zoneID
	return nil
}

func (c *Client) ListZones(ctx context.Context) ([]cloudflare.Zone, error) {
	zones, err := c.api.ListZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}

	return zones, nil
}

func (c *Client) ListDNSRecords(ctx context.Context) ([]DNSRecord, error) {
	if c.zoneID == "" {
		return nil, fmt.Errorf("zone not set")
	}

	rc := cloudflare.ZoneIdentifier(c.zoneID)
	records, _, err := c.api.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	var dnsRecords []DNSRecord
	for _, record := range records {
		dnsRecord := DNSRecord{
			ID:       record.ID,
			Type:     record.Type,
			Name:     record.Name,
			Content:  record.Content,
			TTL:      record.TTL,
			Priority: record.Priority,
			Proxied:  record.Proxied,
		}
		dnsRecords = append(dnsRecords, dnsRecord)
	}

	return dnsRecords, nil
}

func (c *Client) FindDNSRecord(ctx context.Context, name, content string, recordType string) ([]DNSRecord, error) {
	if c.zoneID == "" {
		return nil, fmt.Errorf("zone not set")
	}

	params := cloudflare.ListDNSRecordsParams{}
	if name != "" {
		params.Name = name
	}
	if recordType != "" {
		params.Type = recordType
	}
	if content != "" {
		params.Content = content
	}

	rc := cloudflare.ZoneIdentifier(c.zoneID)
	records, _, err := c.api.ListDNSRecords(ctx, rc, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find DNS records: %w", err)
	}

	var dnsRecords []DNSRecord
	for _, record := range records {
		dnsRecord := DNSRecord{
			ID:       record.ID,
			Type:     record.Type,
			Name:     record.Name,
			Content:  record.Content,
			TTL:      record.TTL,
			Priority: record.Priority,
			Proxied:  record.Proxied,
		}
		dnsRecords = append(dnsRecords, dnsRecord)
	}

	return dnsRecords, nil
}

func (c *Client) AddDNSRecord(ctx context.Context, recordType, name, content string, ttl int, priority *uint16, proxied bool) (*DNSRecord, error) {
	if c.zoneID == "" {
		return nil, fmt.Errorf("zone not set")
	}

	params := cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    name,
		Content: content,
		TTL:     ttl,
	}

	// Set proxied status (only for A, AAAA, CNAME)
	recordTypeUpper := strings.ToUpper(recordType)
	if recordTypeUpper == "A" || recordTypeUpper == "AAAA" || recordTypeUpper == "CNAME" {
		params.Proxied = &proxied
	}

	// Set priority for MX and SRV records
	if priority != nil && (recordTypeUpper == "MX" || recordTypeUpper == "SRV") {
		params.Priority = priority
	}

	rc := cloudflare.ZoneIdentifier(c.zoneID)
	record, err := c.api.CreateDNSRecord(ctx, rc, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS record: %w", err)
	}

	dnsRecord := &DNSRecord{
		ID:       record.ID,
		Type:     record.Type,
		Name:     record.Name,
		Content:  record.Content,
		TTL:      record.TTL,
		Priority: record.Priority,
		Proxied:  record.Proxied,
	}

	return dnsRecord, nil
}

func (c *Client) UpdateDNSRecord(ctx context.Context, recordID, recordType, name, content string, ttl int, priority *uint16, proxied bool) (*DNSRecord, error) {
	if c.zoneID == "" {
		return nil, fmt.Errorf("zone not set")
	}

	params := cloudflare.UpdateDNSRecordParams{
		Type:    recordType,
		Name:    name,
		Content: content,
		TTL:     ttl,
	}

	// Set proxied status
	recordTypeUpper := strings.ToUpper(recordType)
	if recordTypeUpper == "A" || recordTypeUpper == "AAAA" || recordTypeUpper == "CNAME" {
		params.Proxied = &proxied
	}

	// Set priority for MX and SRV records
	if priority != nil && (recordTypeUpper == "MX" || recordTypeUpper == "SRV") {
		params.Priority = priority
	}

	rc := cloudflare.ZoneIdentifier(c.zoneID)
	params.ID = recordID
	record, err := c.api.UpdateDNSRecord(ctx, rc, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update DNS record: %w", err)
	}

	dnsRecord := &DNSRecord{
		ID:       record.ID,
		Type:     record.Type,
		Name:     record.Name,
		Content:  record.Content,
		TTL:      record.TTL,
		Priority: record.Priority,
		Proxied:  record.Proxied,
	}

	return dnsRecord, nil
}

func (c *Client) DeleteDNSRecord(ctx context.Context, recordID string) error {
	if c.zoneID == "" {
		return fmt.Errorf("zone not set")
	}

	rc := cloudflare.ZoneIdentifier(c.zoneID)
	err := c.api.DeleteDNSRecord(ctx, rc, recordID)
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}

	return nil
}
