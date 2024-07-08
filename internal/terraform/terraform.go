package terraform

type Plan struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	Configuration map[string]interface{} `json:"configuration"`
}

func ParsePlan(plan Plan) []Resource {
	return plan.Resources
}
