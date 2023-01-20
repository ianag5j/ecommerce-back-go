package jwk

import (
	"io"
	"net/http"
	"os"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func ValidateToken(at string) (jwt.Token, error) {
	resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + "pem")
	if err != nil {
		return jwt.New(), err
	}
	defer resp.Body.Close()

	pem, _ := io.ReadAll(resp.Body)
	jwkey, err := jwk.ParseKey(pem, jwk.WithPEM(true))

	if err != nil {
		return jwt.New(), err
	}

	return jwt.Parse([]byte(at), jwt.WithKey(jwa.RS256, jwkey), jwt.WithValidate(true))
}
