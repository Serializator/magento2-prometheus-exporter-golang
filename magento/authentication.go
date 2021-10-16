package magento

import (
	"fmt"
	"net/http"
)

type Authenticator interface {
	apply(request *http.Request) error
}

type bearerAuthenticator struct {
	Token string
}

func (authenticator *bearerAuthenticator) apply(request *http.Request) error {
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authenticator.Token))
	return nil
}

func NewBearerAuthenticator(token string) Authenticator {
	return &bearerAuthenticator{token}
}
