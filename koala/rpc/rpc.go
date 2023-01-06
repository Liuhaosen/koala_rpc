package rpc

import "modtest/gostudy/lesson2/ibinarytree/koala/middleware"

func BuildClientMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware
	if len(mids) == 0 {
		return handle
	}
	m := middleware.Chain(mids[0], mids...)
	return m(handle)
}
