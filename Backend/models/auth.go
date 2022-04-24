package models

type AuthSignin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64 //Access Token expies
	RtExpires    int64 //refresh token expires
}

type AuthUsedToken struct {
	Token        string
	RefreshToken string
}
