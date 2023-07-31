package web

// Middleware is used to process HTTP requests before they reach the final Handler.
type Middleware = func(Handler) Handler

// WrapMiddlewares chains a series of Middleware applied from top to bottom order.
func WrapMiddlewares(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}
