package yandexcloud

// BillingAccount describes the structure of a Yandex Cloud billing account
// Used in the REST client, tables, and tests
type BillingAccount struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	CreatedAt   string            `json:"createdAt"`
	CountryCode string            `json:"countryCode"`
	Balance     string            `json:"balance"`
	Currency    string            `json:"currency"`
	Active      bool              `json:"active"`
	Labels      map[string]string `json:"labels"`
}
