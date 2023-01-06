package logs

import (
	"context"
	"sync"
)

type KeyVal struct {
	key interface{}
	val interface{}
}

type LogField struct {
	kvs       []KeyVal
	fieldLock sync.Mutex
}

func (l *LogField) AddField(key, val interface{}) {
	l.fieldLock.Lock()
	defer l.fieldLock.Unlock()
	l.kvs = append(l.kvs, KeyVal{key: key, val: val})
}

type kvsIdKey struct{}

func WithFieldContext(ctx context.Context) context.Context {
	fields := &LogField{}
	return context.WithValue(ctx, kvsIdKey{}, fields)
}

func AddField(ctx context.Context, key string, val interface{}) {
	field := getFields(ctx)
	if field == nil {
		return
	}
	field.AddField(key, val)
}

func getFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}
	return field
}
