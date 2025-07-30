package entity

type Token struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM4ODM4NDAsImlhdCI6MTc1Mzg4Mjk0MCwic3ViIjoiNmQ2NWRiM2YtNTQ2NC00ODU4LWEzMmYtOTdhYWI2ZWQ1MWE1In0.264J7ppNEIKhDlEUQQk6qfFRoR-w5BlUXHzdnh4RDzUJQWt8_X7Qs-xlBpLzvCEY9D1ymcYbsP6uCzwYTYsb7A"`
	RefreshToken string `json:"refresh_token" example:"nZiZiGmwDPkRFc21izNL5Nwp94TEy+qdYwBtKWE+7CM="`
}
