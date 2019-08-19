package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

var api *cloudflare.API

func main() {
	var app = cli.NewApp()
	app.Name = "cfdyndns"
	app.Version = "1.0.0"
	app.Before = initializeAPI
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-email",
			Usage:  "Cloudflare api E-Mail",
			Value:  "",
			EnvVar: "CF_API_EMAIL",
		},
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "Cloudflare api key",
			Value:  "",
			EnvVar: "CF_API_KEY",
		},
		cli.StringFlag{
			Name:   "zone",
			Usage:  "Cloudflare zone name",
			Value:  "",
			EnvVar: "CF_API_KEY",
		},
		cli.StringSliceFlag{
			Name: "names",
			Usage: "Subdomain names to configure",
			Value: &cli.StringSlice{"home", "*.home"},
			EnvVar: "SUBDOMAIN_NAMES",
		},
	}
	app.Action = updateDNS

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func updateDNS(c *cli.Context) error {
	if err := checkFlags(c, "zone"); err != nil {
		return err
	}
	zone := c.String("zone")
	names := c.StringSlice("names")
	zoneId, err := api.ZoneIDByName(zone)
	if err != nil {
		return err
	}

	ip, err := getExternaIP()
	if err != nil {
		return err
	}

	log.Printf("Current external IP: %s\n", ip)

	for _, name := range names {
		subdomain := strings.Join([]string{name, zone}, ".")

		changed, err := dnsCreateOrUpdate(c, zoneId, subdomain, ip)
		if err != nil {
			return err
		}

		if changed {
			log.Printf("Updated zone record for '%s' successfully.", subdomain)
		} else {
			log.Printf("Not updating zone record for '%s' - the IP hasn't changed.", subdomain)
		}
	}

	return nil
}




func checkFlags(c *cli.Context, flags ...string) error {
	var missingFlags []string
	for _, flag := range flags {
		if ! c.IsSet(flag) {
			missingFlags = append(missingFlags, flag)
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("error: the following flags were empty or not provided: %s", strings.Join(missingFlags, ","))
	}

	return nil
}