package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != 32 {
		return nil, fmt.Errorf("invalid key size: must be exactly 32 bytes")
	}

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create symmetric key: %w", err)
	}

	return &PasetoMaker{
		symmetricKey: key,
	}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, role string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	// Set claims (dados) no token
	token.SetString("username", username)
	token.SetString("role", role)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	// Criptografa o token usando a chave simétrica
	encryptedToken := token.V4Encrypt(maker.symmetricKey, nil)

	return encryptedToken, nil
}

// VerifyToken checks if the token is valid and returns the claims
func (maker *PasetoMaker) VerifyToken(token string) (map[string]interface{}, error) {
	parser := paseto.NewParser()

	// Descriptografa o token
	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, nil)
	if err != nil {
		return nil, err
	}

	// Verifica se o token expirou
	expiratedAt, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("failed to get token expiration: %w", err)
	}

	if time.Now().After(expiratedAt) {
		return nil, fmt.Errorf("token has expired")
	}

	// Recupera as claims (dados) do token
	username, err := parsedToken.GetString("username")
	if err != nil {
		return nil, fmt.Errorf("failed to get username: %w", err)
	}

	role, err := parsedToken.GetString("role")
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("failed to get token isued: %w", err)
	}

	claims := make(map[string]interface{})
	claims["username"] = username
	claims["role"] = role
	claims["issued_at"] = issuedAt
	claims["expired_at"] = expiratedAt

	return claims, nil
}
