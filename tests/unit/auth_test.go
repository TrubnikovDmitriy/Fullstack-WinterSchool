package unit

import (
	"testing"
	db "../../database"
	. "../../tests"
	"../../cache"
	"../../models"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"github.com/dgrijalva/jwt-go"
	"../../handlers"
)


func TestAuthHappyPath(t *testing.T) {
	person := CreateNewPerson()
	oauth := GetNewOAuth(person, 2)

	err := db.Auth(oauth)
	if err != nil {
		t.Fatalf("Can't auth common person\n%s", err)
	}

	if oauth.PersonID != person.ID {
		t.Errorf("Getted wrong ID\n" +
			"PersonID =\t%s \nAuthPersonID = \t%s\n", person.ID, oauth.PersonID)
	}
	if oauth.LastName != person.LastName {
		t.Errorf("Getted wrong last name\n" +
			"Original =\t%s \nGetted = \t%s\n", person.LastName, oauth.LastName)
	}
	if oauth.FirstName != person.FirstName {
		t.Errorf("Getted wrong first name\n" +
			"Original =\t%s \nGetted = \t%s\n", person.FirstName, oauth.FirstName)
	}

}

func TestAuthFail(t *testing.T) {

	personNotAuth := GetNewPerson()
	notAuth := GetNewOAuth(personNotAuth, 2)

	err := db.Auth(notAuth)
	if err == nil {
		t.Fatalf("Can auth random user\n(ID = %s)\n(email = %s)\n",
			notAuth.PersonID, notAuth.Email)
	}
	if err.Code != fasthttp.StatusNotFound {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestCacheCreateCode(t *testing.T) {

	oauth := CreateNewOAuth(2)
	code := cache.CreateCode(oauth)

	if code == nil || *code == uuid.Nil {
		t.Fatalf("Can't creates auth code for \npersonID = %s", oauth.PersonID)
	}
}

func TestCacheActivateCode(t *testing.T) {

	oauth := CreateNewOAuth(2)
	code := cache.CreateCode(oauth)

	accessToken, refreshToken := cache.ActivateCode(code)
	if accessToken == nil || refreshToken == nil {
		t.Fatalf("Can't activate code %s", code.String())
	}
	if *accessToken == "" || *refreshToken == "" {
		t.Fatalf("Can't refresh token (perosn_id:%s)", oauth.PersonID)
	}

	claims := handlers.ParseAccessToken(accessToken)
	if claims == nil {
		t.Fatalf("Can't parse access token %s", *accessToken)
	}

	compare(oauth, claims, t)
}

func TestCacheRefreshToken(t *testing.T) {

	oauth := CreateNewOAuth(2)
	code := cache.CreateCode(oauth)
	_, refreshToken := cache.ActivateCode(code)
	accessToken, refreshToken := cache.RefreshToken(*refreshToken)

	if accessToken == nil || refreshToken == nil {
		t.Fatalf("Can't refresh token %s", refreshToken)
	}
	if *accessToken == "" || *refreshToken == "" {
		t.Fatalf("Can't refresh token (perosn_id:%s)", oauth.PersonID)
	}

	claims := handlers.ParseAccessToken(accessToken)
	if claims == nil {
		t.Fatalf("Can't parse access token %s", *accessToken)
	}

	compare(oauth, claims, t)
}

func TestCacheDoubleRefreshDiffToken(t *testing.T) {

	oauth := CreateNewOAuth(2)
	code := cache.CreateCode(oauth)
	_, refreshToken := cache.ActivateCode(code)
	_, refreshToken = cache.RefreshToken(*refreshToken)
	accessToken, refreshToken := cache.RefreshToken(*refreshToken)

	if accessToken == nil || refreshToken == nil {
		t.Fatalf("Can't refresh token %s", refreshToken)
	}
	if *accessToken == "" || *refreshToken == "" {
		t.Fatalf("Can't refresh token (perosn_id:%s)", oauth.PersonID)
	}

	claims := handlers.ParseAccessToken(accessToken)
	if claims == nil {
		t.Fatalf("Can't parse access token %s", *accessToken)
	}

	compare(oauth, claims, t)
}

func TestCacheDoubleRefreshSameToken(t *testing.T) {

	oauth := CreateNewOAuth(2)
	code := cache.CreateCode(oauth)
	_, refreshToken := cache.ActivateCode(code)
	_, _ = cache.RefreshToken(*refreshToken)
	accessToken, refreshToken := cache.RefreshToken(*refreshToken)

	if accessToken == nil && refreshToken == nil {
		return
	} else {
		t.Fatalf("Twice activated the same refreshString %s", code.String())
	}
	if *accessToken == "" || *refreshToken == "" {
		t.Fatalf("Can't refresh token (perosn_id:%s)", oauth.PersonID)
	}
}

func compare(oauth *models.OAuth, claims *jwt.MapClaims, t* testing.T) {

	if oauth.PersonID.String() != (*claims)["person_id"] {
		t.Errorf("IDs do not match")
	}
	if oauth.LastName != (*claims)["last_name"] {
		t.Errorf("Last names do not match")
	}
	if oauth.FirstName != (*claims)["first_name"] {
		t.Errorf("First names do not match")
	}
	if oauth.AppID.String() != (*claims)["app_id"] {
		t.Errorf("App IDs do not match")
	}
	if oauth.Scope != int((*claims)["scope"].(float64)) {
		t.Errorf("Scopes do not match")
	}
	if oauth.Staff != (*claims)["staff"] {
		t.Errorf("App IDs do not match")
	}
	if handlers.IsExpire(claims) {
		t.Errorf("Token is expired")
	}
}