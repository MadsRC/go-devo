package devo

type AlertCorrelationTrigger struct {
	Kind string `json:"kind"`
	ExternalOffset string `json:"externalOffset"`
	InternalPeriod string `json:"internalPeriod"`
	InternalOffset string `json:"internalOffset"`
	Period string `json:"period"`
	Threshold string `json:"threshold"`
	BackPeriod string `json:"backPeriod"`
	Absolute string `json:"absolute"`
	AggregationColumn string `json:"aggregationColumn"`
}

type AlertCorrelationContext struct {
	QuerySourceCode string `json:"querySourceCode"`
	Priority int `json:"priority"`
	CorrelationTrigger AlertCorrelationTrigger `json:"correlationTrigger`
}

type Alert struct {
	ID int `json:"id"`
	CreationDate string `json:"creationDate`
	Name	string `json:"name"`
	Message string `json:"message"`
	Description string `json:"description"`
	Subcategory string `json:"subcategory"`
	CategoryID int `json:"categoryId"`
	SubcategoryID int `json:"subcategoryId"`
	IsActive bool `json:"isActive"`
	IsAlertChain bool `json:"isAlertChain"`
	AlertCorrelationContext AlertCorrelationContext
	actionPolicyID map[string]string `json:"actionPolicyId"`
}