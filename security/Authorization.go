package security

import (
	"GoMin/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func RoleBaseAuthMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS512 {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
				}
				return []byte(config.AppConfig.Jwt.Secret), nil
			})

			if err != nil || !token.Valid {
				fmt.Println("err: ", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid claims")
			}

			if requiredRole == "SYSTEM" {
				sub, ok := claims["sub"].(string)
				if !ok || sub != "system" {
					return echo.NewHTTPError(http.StatusForbidden, "Unauthorized role")
				}
			} else {
				ok, err2, done := checkExp(claims)
				if done {
					return err2
				}

				authRole, ok := claims["auth"].(string)
				if !ok || authRole != requiredRole {
					return echo.NewHTTPError(http.StatusForbidden, "Unauthorized role")
				}

				bid, _ := claims["bid"].(float64)
				c.Set("user", claims)
				c.Set("bid", int(bid))
			}

			return next(c)
		}
	}
}

func checkExp(claims jwt.MapClaims) (bool, error, bool) {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid exp claim"), true
	}

	expirationTime := time.Unix(int64(exp), 0)
	if time.Now().After(expirationTime) {
		return false, echo.NewHTTPError(http.StatusUnauthorized, "Token is expired"), true
	}
	return ok, nil, false
}
