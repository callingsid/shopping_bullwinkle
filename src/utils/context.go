package utils

import "context"

type key int
var txn_uid = key(1)

type key2 int
var method = key2(1)

type ContextValue struct {
	Txn_uid string
	Method string
}

func SetContext(ctx context.Context, id string, httpMethod string) context.Context {
	//TODO sid -- put context data in struct
	return context.WithValue(ctx, "values", ContextValue{Txn_uid:id, Method: httpMethod})
}
func GetContext(ctx context.Context) (ContextValue) {
	val := ctx.Value("values").(ContextValue)
	return val
}








