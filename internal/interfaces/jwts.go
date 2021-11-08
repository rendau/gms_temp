package interfaces

type Jwts interface {
	JwtCreate(sub string, expSeconds int64, payload map[string]interface{}) (string, error)
}
