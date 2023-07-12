//go:build local

package web

func init() {
	allowedOrigins = append(allowedOrigins, "http://localhost:3000")
}
