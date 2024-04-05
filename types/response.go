package types

type BakalariErrorResponse struct {
	Error string
}

type BakalariLoginResponse struct {
	Access_Token string
	Refresh_Token string
	Token_Type string
	Expires_In int
}
