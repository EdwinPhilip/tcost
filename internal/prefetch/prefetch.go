package prefetch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
)

type ServiceAttributes struct {
	ServiceCode    string   `json:"serviceCode"`
	AttributeNames []string `json:"attributeNames"`
}

func PrefetchServiceAttributes() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	pricingSvc := pricing.NewFromConfig(cfg)

	log.Println("Fetching service codes...")
	serviceCodes, err := getServiceCodes(pricingSvc)
	if err != nil {
		log.Fatalf("Failed to get service codes: %v", err)
	}
	log.Printf("Fetched %d service codes\n", len(serviceCodes))

	var allServiceAttributes []ServiceAttributes

	for i, serviceCode := range serviceCodes {
		log.Printf("Fetching attributes for service %s (%d/%d)...", serviceCode, i+1, len(serviceCodes))
		attributeNames, err := getAttributes(pricingSvc, serviceCode)
		if err != nil {
			log.Printf("Failed to get attributes for service %s: %v", serviceCode, err)
			continue
		}

		serviceAttributes := ServiceAttributes{
			ServiceCode:    serviceCode,
			AttributeNames: attributeNames,
		}
		allServiceAttributes = append(allServiceAttributes, serviceAttributes)
		log.Printf("Fetched %d attributes for service %s\n", len(attributeNames), serviceCode)
	}

	jsonData, err := json.MarshalIndent(allServiceAttributes, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("service_attributes.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	fmt.Println("Successfully fetched service attributes and saved to service_attributes.json")
}

func getServiceCodes(pricingSvc *pricing.Client) ([]string, error) {
	var serviceCodes []string
	var nextToken *string

	for {
		input := &pricing.DescribeServicesInput{
			NextToken: nextToken,
		}

		output, err := pricingSvc.DescribeServices(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		for _, service := range output.Services {
			serviceCodes = append(serviceCodes, aws.ToString(service.ServiceCode))
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return serviceCodes, nil
}

func getAttributes(pricingSvc *pricing.Client, serviceCode string) ([]string, error) {
	var attributeNames []string
	var nextToken *string

	for {
		input := &pricing.DescribeServicesInput{
			ServiceCode: aws.String(serviceCode),
			NextToken:   nextToken,
		}

		output, err := pricingSvc.DescribeServices(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		if len(output.Services) > 0 {
			for _, attributeName := range output.Services[0].AttributeNames {
				attributeNames = append(attributeNames, attributeName)
			}
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return attributeNames, nil
}
