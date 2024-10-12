package middleware

import (
	"bank-backend/pkg"
	"bank-backend/utils"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		var (
			lvState1       = utils.LogEventStateDecodeRequest
			lfState1Status = "state_1_decode_request_status"

			lf = []slog.Attr{
				pkg.LogEventName("middleware"),
			}
		)
		/*------------------------------------
		| Step 1 : Decode request
		* ----------------------------------*/

		lf = append(lf, pkg.LogEventState(lvState1))

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			lf = append(lf, pkg.LogStatusFailed(lfState1Status))
			pkg.LogWarnWithContext(c.Context(), "missing authorization header", errors.New("missing authorization header"), lf)
			return c.Status(http.StatusUnauthorized).JSON(utils.StandardResponse{
				Message: "missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(pkg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			lf = append(lf, pkg.LogStatusFailed(lfState1Status))
			pkg.LogWarnWithContext(c.Context(), "invalid jwt token ", err, lf)
			return c.Status(http.StatusUnauthorized).JSON(utils.StandardResponse{
				Message: "invalid jwt token",
			})
		}

		c.Locals("user", token)
		return c.Next()
	}
}

func RoleBasedMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var (
			lvState1       = utils.LogEventStateValidateToken
			lfState1Status = "state_1_validate_role_status"

			lf = []slog.Attr{
				pkg.LogEventName("middleware"),
			}
		)
		/*------------------------------------
		| Step 1 : validate role
		* ----------------------------------*/
		lf = append(lf, pkg.LogEventState(lvState1))
		user, ok := c.Locals("user").(*jwt.Token)
		if !ok {
			lf = append(lf, pkg.LogStatusFailed(lfState1Status))
			pkg.LogWarnWithContext(c.Context(), "invalid token or missing jwt token", errors.New("invalid token missing jwt token"), lf)
			return c.Status(http.StatusUnauthorized).JSON(utils.StandardResponse{
				Message: "invalid token or missing jwt token",
			})
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			lf = append(lf, pkg.LogStatusFailed(lfState1Status))
			pkg.LogWarnWithContext(c.Context(), "invalid token token claims", errors.New("invalid token claims"), lf)
			return c.Status(http.StatusUnauthorized).JSON(utils.StandardResponse{
				Message: "invalid token token claims",
			})
		}

		userRole, ok := claims["phone_number"].(string)
		if !ok {
			lf = append(lf, pkg.LogStatusFailed(lfState1Status))
			pkg.LogWarnWithContext(c.Context(), "phone_numbre claims missing on token", errors.New("phone_number claims missing on token"), lf)
			return c.Status(http.StatusForbidden).JSON(utils.StandardResponse{
				Message: "phone_number claims missing on token",
			})
		}
		// Store the user role in context
		c.Locals("user-phone", userRole)

		if userRole != "" {
			return c.Next()
		}
		return c.Status(http.StatusForbidden).JSON(utils.StandardResponse{
			Message: "access denied",
		})
	}
}
