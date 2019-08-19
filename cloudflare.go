package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/urfave/cli"
)

func initializeAPI(c *cli.Context) error {
	if err := checkFlags(c, "api-email", "api-key"); err != nil {
		return err
	}

	apiEmail := c.String("api-email")
	apiKey := c.String("api-key")

	// Be aware the following code sets the global package `api` variable
	var err error
	api, err = cloudflare.New(apiKey, apiEmail)
	if err != nil {
		return err
	}

	return nil
}

func findRecord(zoneID string, subdomain string) (*cloudflare.DNSRecord, error){
	rr := cloudflare.DNSRecord{
			Name: subdomain,
	}
	records, err := api.DNSRecords(zoneID, rr)
	if err != nil {
		return nil, fmt.Errorf("error fetching DNS records: %s", err)
	}

	if len(records) > 0 {
		return &records[0], nil
	}

	return nil, nil
}

func dnsCreateOrUpdate(c *cli.Context, zoneId string, subdomain string, ip string) (changed bool, err error) {
	rr, err := findRecord(zoneId, subdomain)
	if err != nil {
		return false, err
	}

	var resp *cloudflare.DNSRecordResponse
	if rr == nil {
		rr = &cloudflare.DNSRecord{
			Name: subdomain,
		}
		rr.Type = "A"
		rr.TTL = 120
		rr.Content = ip

		resp, err = api.CreateDNSRecord(zoneId, *rr)
		if err != nil {
			return false, err
		}

		if resp.Success {
			return true, nil
		}
	}

	if rr.Content == ip {
		return false, nil
	}

	rr.Content = ip
	err = api.UpdateDNSRecord(zoneId, rr.ID, *rr)
	if err == nil {
		return true, nil
	}

	return false, err
}
