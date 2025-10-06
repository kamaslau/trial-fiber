package middlewares

import (
	"log"
	"net/http"
	"strings"

	"app/src/internal/models"
	"app/src/internal/utils/drivers"
	jwtutils "app/src/internal/utils/jwt"

	"github.com/gofiber/fiber/v3"
)

// TryToken Try to parse userID from token, save to local context if found, for future authentication
func TryToken(c fiber.Ctx) error {
	token := ""

	// Parse the raw JWT string (empty if not present) from cookie
	if token = c.Cookies("token"); token != "" {
		log.Printf("🔍 JWT Token parsed from Cookies: %s", token)
	}

	// Parse the raw JWT string (empty if not present) from Authorization Header
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// log.Println("🔍 Authorization header is not found")
			return c.Next()
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			// log.Println("🔍 Authorization header is not in Bearer format")
			return c.Next()
		}

		// 提取token
		token = strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			// log.Println("🔍 Empty token found in Authorization header")
			return c.Next()
		}
		log.Printf("🔍 JWT Token parsed from Authorization header: %s", token)
	}

	// Parse and persist certain info
	if payload, err := jwtutils.VerifyJWT(token); err != nil {
		log.Printf("🛑 Failed to verify JWT: %v", err)
	} else {
		// log.Printf("🔍 JWT Token stored:\n UserID: %s\n ExpiresAt: %v\n TokenID: %s", payload.UserID, payload.ExpiresAt, payload.TokenID)

		// check if token is revoked/blocked; check handler token->DeleteOne for reference
		var revoked = false
		if value, err := drivers.CacheGet("token-revoke-user:" + payload.UserID); err != nil {
			log.Printf("Failed to fetch cache: %v", err)
		} else {
			log.Printf("Got cache: %v", value)

			if value == "" || value == payload.TokenID {
				revoked = true
			} else {
				log.Print("Token is not revoked")
			}
		}
		// log.Printf("revoked: %v", revoked)

		// Persist payload
		if !revoked {
			c.Locals("userID", payload.UserID)
			c.Locals("expiresAt", payload.ExpiresAt)
			c.Locals("tokenID", payload.TokenID)
		}
	}

	return c.Next()
}

// ACRequirement
type ACRequirement struct {
	Roles       []string // 角色需求，留空则不限
	Permissions []string // 权限需求，留空则不限
}

// NeedsRoles 检查是否有有效的角色要求（已设置且非空）
func (ac *ACRequirement) NeedsRoles() bool {
	return len(ac.Roles) > 0
}

// NeedsPermissions 检查是否有有效的权限要求（已设置且非空）
func (ac *ACRequirement) NeedsPermissions() bool {
	return len(ac.Permissions) > 0
}

// NeedsAC 检查是否有任何访问控制要求
func (ac *ACRequirement) NeedsAC() bool {
	return len(ac.Roles) > 0 || len(ac.Permissions) > 0
}

// Access Control Policy
// OR (default) match Role or Permission
// AND both Role and Permission should meet
type ACPolicy string

// Try AC Try to get roles and permissions for certain user from userId, save to local context for future usage
// Do NOT refactor this to a real middleware for using with routers, unless only entity level is all that is need to be controled
func TryAC(c fiber.Ctx, r ACRequirement, p ...ACPolicy) error {
	log.Printf("🔒 Checking access control...")

	// No requirement, skip
	if !r.NeedsAC() {
		return nil
	}

	// Get local userID
	localUserID := ""
	if value, ok := c.Locals("userID").(string); !ok {
		return fiber.NewError(http.StatusUnauthorized)
	} else {
		localUserID = value
	}
	log.Printf("Local userID: %s", localUserID)

	// Decide policy
	var policy ACPolicy
	if len(p) > 0 {
		policy = p[0]
	} else {
		policy = "OR"
	}

	// For "AND" policy, both Roles and Permissions should be explicitly constrained
	if policy == "AND" && !r.NeedsRoles() && !r.NeedsPermissions() {
		return fiber.NewError(http.StatusForbidden, "Need both role and permission to match, but got ACRequirement xxx")
	}

	// Check Roles Needed
	if r.NeedsRoles() {
		// log.Printf("NeedsRoles: %v", r.Roles)

		// Get user roles and permissions from local context
		userRoles := []string{}
		if value, ok := c.Locals("roles").([]string); ok {
			userRoles = value
			log.Printf("Local User Roles: %v", userRoles)
		}

		// Get user roles from database using proper association
		var user models.User
		if err := drivers.DBClient.Preload("Roles").Where("id = ?", localUserID).First(&user).Error; err != nil {
			log.Printf("Failed to get user roles from database: %v", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		// Extract role names from the user's roles
		for _, role := range user.Roles {
			userRoles = append(userRoles, role.Name)
		}
		// log.Printf("Fetched User Roles: %v", userRoles)
		c.Locals("roles", userRoles)

		if len(userRoles) == 0 {
			return fiber.NewError(http.StatusForbidden)
		}

		if !intersect(r.Roles, userRoles) {
			// log.Printf("NeedsRoles: not intersect")
			if policy == "AND" {
				return fiber.NewError(http.StatusForbidden)
			}
		}
	}

	// Check Permissions Needed
	if r.NeedsPermissions() {
		// log.Printf("NeedsPermissions: %v", r.Permissions)

		// Get user permissions from local context
		userPermissions := []string{}
		if value, ok := c.Locals("permissions").([]string); ok {
			userPermissions = value
			log.Printf("Local User Permissions: %v", userPermissions)
		}

		// Get user roles from database using proper association
		var user models.User
		if err := drivers.DBClient.Preload("Roles.Permissions").Where("id = ?", localUserID).First(&user).Error; err != nil {
			log.Printf("Failed to get user roles and permissions from database: %v", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		// Extract permission names from the user's roles
		permSet := make(map[string]struct{})
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				permSet[perm.Name] = struct{}{}
			}
		}
		for permName := range permSet {
			userPermissions = append(userPermissions, permName)
		}
		// log.Printf("Fetched User Permissions: %v", userPermissions)
		c.Locals("permissions", userPermissions)

		if len(userPermissions) == 0 {
			return fiber.NewError(http.StatusForbidden)
		}

		if !intersect(r.Permissions, userPermissions) {
			// log.Printf("NeedsPermissions: not intersect")
			return fiber.NewError(http.StatusForbidden)
		}
	}

	// Continue, if no breaks should apply
	log.Printf("✅ Access control passed")
	return nil
}

func intersect(a, b []string) bool {
	// Maximize efficiency by making b the smaller one, as indexing target
	if len(a) < len(b) {
		a, b = b, a
	}

	set := make(map[string]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	for _, v := range a {
		if _, ok := set[v]; ok {
			return true
		}
	}
	return false
}
