package ctxlog

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// CtxLogKeyType is a single type we use so no colisions in the context
type CtxLogKeyType int

const (
	// CtxLogKey defines a context key that can hold a slice of context keys
	CtxLogKey CtxLogKeyType = iota
)

type CtxLogTags map[string]interface{}

// AddRpcTagsToContext adds the given log tag mappings (logTagsToAdd) to the
// given context, creating a new one if necessary. Returns the resulting
// context with the new log tag mappings.
func AddTagsToContext(ctx context.Context, logTagsToAdd CtxLogTags) context.Context {
	currTags, ok := TagsFromContext(ctx)
	if !ok {
		currTags = make(CtxLogTags)
	}
	for key, tag := range logTagsToAdd {
		currTags[key] = tag
	}

	return context.WithValue(ctx, CtxLogKey, currTags)
}

// TagsFromContext returns the tags being passed along with the given context.
func TagsFromContext(ctx context.Context) (CtxLogTags, bool) {
	logTags, ok := ctx.Value(CtxLogKey).(CtxLogTags)
	if ok {
		ret := make(CtxLogTags)
		for k, v := range logTags {
			ret[k] = v
		}
		return ret, true
	}
	return nil, false
}

func RandStringB64(numTriads int) string {
	buf, err := RandBytes(numTriads * 3)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(buf)
}

func RandBytes(length int) ([]byte, error) {
	var n int
	var err error
	buf := make([]byte, length)
	if n, err = rand.Read(buf); err != nil {
		return nil, err
	}
	// rand.Read uses io.ReadFull internally, so this check should never fail.
	if n != length {
		return nil, fmt.Errorf("RandBytes got too few bytes, %d < %d", n, length)
	}
	return buf, nil
}

func WithLogTag(ctx context.Context, k string) context.Context {
	return WithLogTagWithValue(ctx, k, RandStringB64(3))
}

func WithLogTagWithValue(ctx context.Context, k, v string) context.Context {

	// Don't overwrite existing tag?
	if tags, ok := TagsFromContext(ctx); ok {
		if _, found := tags[k]; found {
			return ctx
		}
	}

	newTags := make(CtxLogTags)
	newTags[k] = v
	return AddTagsToContext(ctx, newTags)
}
