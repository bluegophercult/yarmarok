//go:build local

package web

func init() {
	localRun = true
	allowedOrigins = append(allowedOrigins, "http://localhost:3000")
}
