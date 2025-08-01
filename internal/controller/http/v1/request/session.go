package request

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"rN3aqj1enSeFhz7lMzgAtvUZWRz4GZ8qDEy0yXUG4hQ="`
}
