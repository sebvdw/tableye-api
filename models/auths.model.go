package models

type SignInResponse struct {
	AccessToken string          `json:"access_token"`
	User        UserResponse    `json:"user"`
	Dealer      *DealerResponse `json:"dealer,omitempty"`
	Casino      *CasinoResponse `json:"casino,omitempty"`
}
