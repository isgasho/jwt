package jwt

// Sign signs and generates a new token based on the algorithm and a secret key.
// The claims is the payload, the actual body of the token, should
// contain information about a specific authorized client.
// Note that the payload part is not encrypted,
// therefore it should NOT contain any private information.
// See the `Verify` function to decode and verify the result token.
//
// Example Code to pass only standard Claims:
//
//  token, err := jwt.Sign(jwt.HS256, []byte("secret"), jwt.Claims{...})
//
// Example Code to pass custom and expiration Claims manually:
//
//  now := time.Now()
//  token, err := jwt.Sign(jwt.HS256, []byte("secret"), map[string]interface{}{
//    "iat": now.Unix(),
//    "exp": now.Add(15 * time.Minute).Unix(),
//    "foo": "bar",
//  })
//
// Example Code for custom and standard claims using a SignOption:
//
//  token, err := jwt.Sign(jwt.HS256, []byte("secret"), jwt.Map{"foo":"bar"}, jwt.MaxAge(15 * time.Minute))
//  OR
//  token, err := jwt.Sign(jwt.HS256, []byte("secret"), jwt.Map{"foo":"bar"}, jwt.Claims {Expiry: ...})
//
// Example Code for custom type as Claims + standard Claims:
//
//  type User struct { Username string `json:"username"` }
//  token, err := jwt.Sign(jwt.HS256, []byte("secret"), User{Username: "kataras"}, jwt.MaxAge(15 * time.Minute))
func Sign(alg Alg, key PrivateKey, claims interface{}, opts ...SignOption) ([]byte, error) {
	if len(opts) > 0 {
		var standardClaims Claims
		for _, opt := range opts {
			opt.ApplyClaims(&standardClaims)
		}

		claims = Merge(claims, standardClaims)
	}

	return encodeToken(alg, key, claims)
}

// SignOption is just a helper which sets the standard claims at the `Sign` function.
//
// Available SignOptions:
// - MaxAge(time.Duration)
// - Claims{}
type SignOption interface {
	// ApplyClaims should apply standard claims.
	// Accepts the destination claims.
	ApplyClaims(*Claims)
}

// SignOptionFunc completes the `SignOption`. It's a helper to pass a `SignOption` as a function.
type SignOptionFunc func(*Claims)

// ApplyClaims completes the `SignOption` interface.
func (f SignOptionFunc) ApplyClaims(c *Claims) {
	f(c)
}
