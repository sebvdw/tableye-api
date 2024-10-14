package models

// SignInResponse represents the response structure for user sign-in
type SignInResponse struct {
	AccessToken string          `json:"access_token"`
	User        UserResponse    `json:"user"`
	Dealer      *DealerResponse `json:"dealer,omitempty"`
	Casino      *CasinoResponse `json:"casino,omitempty"`
}
