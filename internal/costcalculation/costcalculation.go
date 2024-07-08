package costcalculation

import "tcost/internal/terraform"
import "tcost/internal/types"

type Estimate struct {
	ResourceName string
	DailyCost    float64
	MonthlyCost  float64
}

func CalculateCost(resource terraform.Resource, pricingData types.PricingData) Estimate {
	dailyCost := pricingData.PricePerUnit * 24
	monthlyCost := dailyCost * 30 // Simplified calculation

	return Estimate{
		ResourceName: resource.Name,
		DailyCost:    dailyCost,
		MonthlyCost:  monthlyCost,
	}
}
