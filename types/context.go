//nolint:revive // types package is intentional for shared domain types
package types

type ContextType string

const AuthorizationKey ContextType = "authorization"

const UserContextKey ContextType = "user_claims"
