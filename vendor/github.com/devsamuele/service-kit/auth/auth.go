package auth

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
)

// const (
// 	ScopeRContact  = "r_contact"
// 	ScopeRWContact = "rw_contact"
// )

type ctxKey int

// Key ...
const Key ctxKey = 1

// Claims ...
type Claims struct {
	jwt.StandardClaims
	TenantID string   `json:"tenant_id"`
	Scopes   []string `json:"scopes"`
}

// Authorize ...
func Authorize(r *http.Request, scopes ...string) (bool, error) {

	headerScopesString := r.Header.Get("elite-scopes")
	if headerScopesString == "" {
		// scoes not inheader -> error to handle
		return false, nil
	}

	headerScopes := make([]string, 0)
	if err := json.Unmarshal([]byte(headerScopesString), &headerScopes); err != nil {
		return false, err
	}

	for _, authScope := range scopes {
		foundScope := false
		for i := 0; i < len(headerScopes) && !foundScope; i++ {
			if authScope == headerScopes[i] {
				foundScope = true
			}
		}
		if !foundScope {
			return false, nil
		}
	}
	return true, nil
}

// Keys ...
type Keys map[string]*rsa.PrivateKey

// PublicKeyLookup ...
type PublicKeyLookup func(kid, jku string) (*rsa.PublicKey, error)
type LoadKeys func() (map[string]*rsa.PrivateKey, error)

// Auth ...
type Auth struct {
	clientOnly bool
	algorithm  string
	parser     *jwt.Parser
	lookup     PublicKeyLookup
	keys       Keys
}

// New ...
func New(algorithm string) (*Auth, error) {

	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknow algorithm %v", algorithm)
	}

	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	auth := Auth{
		clientOnly: true,
		algorithm:  algorithm,
		parser:     &parser,
		lookup:     defaultLookup,
	}

	return &auth, nil
}

func (a *Auth) SetLookup(lookup PublicKeyLookup) error {
	if lookup == nil {
		return errors.New("nil lookup func")
	}

	a.lookup = lookup
	return nil
}

func (a *Auth) LoadPrivateKeys(loadKeysFunc LoadKeys) error {
	keys, err := loadKeysFunc()
	if err != nil {
		return err
	}
	a.keys = keys
	a.clientOnly = false
	return nil
}

func defaultLookup(kid, jku string) (*rsa.PublicKey, error) {

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, jku, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jwks, err := jwk.Parse(bodyB)
	if err != nil {
		return nil, err
	}

	publicKey, ok := jwks.LookupKeyID(kid)
	if !ok {
		return nil, errors.New("kid not found")
	}

	var pk rsa.PublicKey

	if err := publicKey.Raw(&pk); err != nil {
		return nil, err
	}

	return &pk, nil
}

// //AddKey ...
// func (a *Auth) AddKey(privateKey *rsa.PrivateKey, kid string) {
// 	a.keys[kid] = privateKey
// }

// // RemoveKey ...
// func (a *Auth) RemoveKey(kid string) {
// 	delete(a.keys, kid)
// }

// GenerateToken ...
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {

	if a.clientOnly {
		return "", errors.New("auth: in client mode only")
	}

	method := jwt.GetSigningMethod(a.algorithm)

	token := jwt.NewWithClaims(method, claims)
	token.Header[kid] = kid

	privateKey, ok := a.keys[kid]
	if !ok {
		return "", errors.New("kid lookup failed")
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return tokenString, nil
}

// ValidateToken ...
func (a *Auth) ValidateToken(tokenString string) (Claims, error) {

	var claims Claims

	keyFunc := func(t *jwt.Token) (interface{}, error) {

		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id (kid) in token header")
		}

		kidID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id (kid) must be string")
		}

		jku, ok := t.Header["jku"]
		if !ok {
			return nil, errors.New("missing jku in token header")
		}

		publicKeyURL, ok := jku.(string)
		if !ok {
			return nil, errors.New("jku must be string")
		}

		return a.lookup(kidID, publicKeyURL)
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return Claims{}, errors.Wrap(err, "parsing token")
	}

	if !token.Valid {
		return Claims{}, errors.Wrap(err, "invalid token")
	}

	return claims, nil
}
