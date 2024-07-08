package awspricing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tcost/internal/terraform"
	"tcost/internal/types"
)

func GetPricing(resource terraform.Resource) (types.PricingData, error) {
	// Example URL, replace with actual AWS Pricing API endpoint
	url := fmt.Sprintf("https://api.pricing.us-east-1.amazonaws.com/pricing?serviceCode=%s&filters[0].field=instanceType&filters[0].value=%s",
		resource.Type, resource.Configuration["instanceType"])

	resp, err := http.Get(url)
	if err != nil {
		return types.PricingData{}, err
	}
	defer resp.Body.Close()

	var pricingData types.PricingData
	err = json.NewDecoder(resp.Body).Decode(&pricingData)
	if err != nil {
		return types.PricingData{}, err
	}

	return pricingData, nil
}
