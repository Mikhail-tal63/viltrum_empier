package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mikhail-tal63/viltrum_empier/config"
	utils "github.com/Mikhail-tal63/viltrum_empier/utils/JWT"
	jsonResponse "github.com/Mikhail-tal63/viltrum_empier/utils/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonResponse.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing authorization header"))
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			jsonResponse.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization format"))
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := []byte(config.ENVs.JWTSecret)
		userID, err := utils.VerifyJWT(secret, tokenString)
		if err != nil {
			jsonResponse.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
func GetUserID(ctx context.Context) (primitive.ObjectID, error) {
	userID, ok := ctx.Value(UserIDKey).(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
