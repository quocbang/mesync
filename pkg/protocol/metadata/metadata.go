package metadata

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// metadata header list
const (
	AuthenticationKey = "Auth-Key"
	RequestID         = "Request-ID"
	UserID            = "User-ID"
	FactoryIDs        = "Factory-IDs"
)

// MetaContext defines the functions needed for composing metadata
type MetaContext interface {
	GetAuthenticationKey() string
	GetRequestID() string
	GetUserID() string
	GetFactoryIDs() []string
}

// ContextWith returns context with specific metadata
func ContextWith(ctx MetaContext) context.Context {
	md := metadata.Pairs(
		AuthenticationKey, ctx.GetAuthenticationKey(),
		RequestID, ctx.GetRequestID(),
		UserID, ctx.GetUserID(),
	)
	md.Append(FactoryIDs, ctx.GetFactoryIDs()...)

	return metadata.NewOutgoingContext(context.Background(), md)
}

// OutgoingContext content.
type OutgoingContext struct {
	authenticationKey string
	requestID         string
	userID            string
	factoryIDs        []string
}

// NewOutgoingContext creates metadata outgoing context.
func NewOutgoingContext(authKey, requestID, userID string, factoryIDs []string) OutgoingContext {
	return OutgoingContext{
		authenticationKey: authKey,
		requestID:         requestID,
		userID:            userID,
		factoryIDs:        factoryIDs,
	}
}

// GetAuthenticationKey receives authentication key.
func (c OutgoingContext) GetAuthenticationKey() string {
	return c.authenticationKey
}

// GetRequestID receives request ID.
func (c OutgoingContext) GetRequestID() string {
	return c.requestID
}

// GetUserID receives user ID.
func (c OutgoingContext) GetUserID() string {
	return c.userID
}

// GetFactoryIDs receives factory ID.
func (c OutgoingContext) GetFactoryIDs() []string {
	return c.factoryIDs
}
