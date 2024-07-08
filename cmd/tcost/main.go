package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"tcost/internal/awspricing"
	"tcost/internal/costcalculation"
	"tcost/internal/prefetch"
	"tcost/internal/terraform"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tcost",
		Short: "tcost calculates AWS resource costs based on a Terraform plan",
	}

	var prefetchCmd = &cobra.Command{
		Use:   "prefetch",
		Short: "Prefetch AWS service attributes",
		Run: func(cmd *cobra.Command, args []string) {
			prefetch.PrefetchServiceAttributes()
		},
	}

	var calculateCmd = &cobra.Command{
		Use:   "calculate [terraform-plan.json]",
		Short: "Calculate costs based on a Terraform plan",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			planFile := args[0]
			planData, err := ioutil.ReadFile(planFile)
			if err != nil {
				log.Fatalf("Failed to read Terraform plan file: %v", err)
			}

			var plan terraform.Plan
			err = json.Unmarshal(planData, &plan)
			if err != nil {
				log.Fatalf("Failed to parse Terraform plan JSON: %v", err)
			}

			resources := terraform.ParsePlan(plan)
			estimates := []costcalculation.Estimate{}

			for _, resource := range resources {
				pricingData, err := awspricing.GetPricing(resource)
				if err != nil {
					log.Printf("Failed to get pricing for resource %v: %v", resource, err)
					continue
				}
				estimate := costcalculation.CalculateCost(resource, pricingData)
				estimates = append(estimates, estimate)
			}

			for _, estimate := range estimates {
				fmt.Printf("Resource: %s\n", estimate.ResourceName)
				fmt.Printf("Daily Cost: $%.2f\n", estimate.DailyCost)
				fmt.Printf("Monthly Cost: $%.2f\n", estimate.MonthlyCost)
				fmt.Println("----------")
			}
		},
	}

	rootCmd.AddCommand(prefetchCmd)
	rootCmd.AddCommand(calculateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
