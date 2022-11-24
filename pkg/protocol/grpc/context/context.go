package context

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// Keys for application flows
type key int

const (
	// AuthKey for client-server validation key.
	AuthKey key = iota + 1
	// RequestID records each incoming requestID.
	RequestID
	// FactoryIDs determine database settings.
	FactoryIDs
	// UserID to record who make requests.
	UserID
)

// GetAuthKey returns auth-key from header.
func GetAuthKey(ctx context.Context) string {
	return getStringContext(ctx, AuthKey)
}

// GetRequestID returns RequestID.
func GetRequestID(ctx context.Context) string {
	return getStringContext(ctx, RequestID)
}

// GetFactoryIDs returns FactoryIDs.
func GetFactoryIDs(ctx context.Context) []string {
	return getStringsContext(ctx, FactoryIDs)
}

// GetUserID returns UserID.
func GetUserID(ctx context.Context) string {
	return getStringContext(ctx, UserID)
}

type fieldOptions struct {
	useDefault func() string
	multiple   bool
}

// FieldOption definition.
type FieldOption func(*fieldOptions)

func parseFieldOptions(opts []FieldOption) fieldOptions {
	var o fieldOptions
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithDefault set default value of parsed field.
func WithDefault(f func() string) FieldOption {
	return func(o *fieldOptions) {
		o.useDefault = f
	}
}

// WithMultiple means the field is multiple values.
func WithMultiple() FieldOption {
	return func(o *fieldOptions) {
		o.multiple = true
	}
}

// Parser for context fields from MD
type Parser struct {
	ctx context.Context
	md  metadata.MD
}

// NewParser create a new Parser.
func NewParser(ctx context.Context, md metadata.MD) *Parser {
	return &Parser{ctx: ctx, md: md}
}

// Parse parses metadata key-value pairs to context.
// you could set a default value of the desired key by
// assigning WithDefault(f func() string) inside the third parameter (opts).
func (p *Parser) Parse(mdKey string, key interface{}, opts ...FieldOption) *Parser {
	o := parseFieldOptions(opts)

	if v := p.md.Get(mdKey); len(v) > 0 {
		if o.multiple {
			p.ctx = context.WithValue(p.ctx, key, v)
		} else {
			p.ctx = context.WithValue(p.ctx, key, v[0])
		}
		return p
	}

	if o.useDefault != nil {
		p.ctx = context.WithValue(p.ctx, key, o.useDefault())
	}
	return p
}

// Done return context of parsed metadata.
func (p *Parser) Done() context.Context {
	return p.ctx
}

// getStringContext return context key's string value.
func getStringContext(ctx context.Context, k interface{}) string {
	s, _ := ctx.Value(k).(string)
	return s
}

// getStringsContext return context key's string value.
func getStringsContext(ctx context.Context, k interface{}) []string {
	ss, ok := ctx.Value(k).([]string)
	if ok {
		return ss
	}

	s, ok := ctx.Value(k).(string)
	if ok {
		return []string{s}
	}
	return []string{}
}
