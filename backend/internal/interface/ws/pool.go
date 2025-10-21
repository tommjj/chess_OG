package ws

import "sync"

var contextPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

func acquireContext() *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Context = nil
	ctx.Conn = nil
	ctx.Hub = nil
	ctx.Payload = nil
	return ctx
}

func releaseContext(ctx *Context) {
	ctx.Context = nil
	ctx.Conn = nil
	ctx.Hub = nil
	ctx.Payload = nil
	contextPool.Put(ctx)
}

// DefaultConnSliceLen is the default length of the connection slice pool
var DefaultConnSliceLen = 16

// connSlicePool is a pool of []*Connection slices to reduce allocations
var connSlicePool = sync.Pool{
	New: func() any {
		slice := make([]*Connection, 0, DefaultConnSliceLen)
		return &slice
	},
}

func acquireConnSlice() *[]*Connection {
	conns := connSlicePool.Get().(*[]*Connection)
	*conns = (*conns)[:0]
	return conns
}

func releaseConnSlice(conns *[]*Connection) {
	*conns = (*conns)[:0]
	connSlicePool.Put(conns)
}
