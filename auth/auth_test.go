package auth

import "testing"

func TestGetToken(t *testing.T) {
	requestAuth := RequestAuth{
		ID:    "505",
		Email: "abuarquemf@gmail.com",
	}
	authService := New()
	_, err := authService.GetToken(&requestAuth)
	if err != nil {
		t.Errorf("want error nil, got %q", err)
	}
}

func TestIsValid(t *testing.T) {
	requestAuth := RequestAuth{
		ID:    "505",
		Email: "abuarquemf@gmail.com",
	}
	authService := New()
	token, err := authService.GetToken(&requestAuth)
	if err != nil {
		t.Errorf("want error nil, got %q", err)
	}
	isValid := authService.IsValid(token)
	if isValid == false {
		t.Errorf("want is Valid true")
	}
}

func TestGetClaims(t *testing.T) {
	requestAuth := RequestAuth{
		ID:    "505",
		Email: "abuarquemf@gmail.com",
	}
	authService := New()
	token, err := authService.GetToken(&requestAuth)
	if err != nil {
		t.Errorf("want error nil, got %q", err)
	}
	claims, err := GetClaims(token)
	if err != nil {
		t.Errorf("want err nil when getting claims")
	}
	id := claims["id"]
	if id != "505" {
		t.Errorf("want id 505, got %s", id)
	}
	email := claims["email"]
	if email != "abuarquemf@gmail.com" {
		t.Errorf("want email abuarquemf@gmail.com, got %s", id)
	}
}
