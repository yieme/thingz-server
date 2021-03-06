package types

type UtilizationResponse struct {

	// Timestamp of when the metric was captured
	Timestamp int64 `json:"ts"`

	// Period
	Period *Period `json:"period"`

	// Criteria
	Criteria string `json:"criteria"`

	// Method
	Method string `json:"method"`

	// Resources pf metrics
	Resources *ResourceUtilizationList `json:"resources"`
}
