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
	return ctx
}

func releaseContext(ctx *Context) {
	ctx.Context = nil
	ctx.Conn = nil
	ctx.Hub = nil
	contextPool.Put(ctx)
}

var DefaultConnSliceLen = 16
var connSlicePool = sync.Pool{
	// Hàm New được gọi khi Pool trống.
	// Thường trả về nil hoặc slice rỗng, chúng ta sẽ cấp phát
	// khi lấy ra (Get) để kiểm soát dung lượng (capacity).
	New: func() any {
		slice := make([]*Connection, 0, DefaultConnSliceLen)
		return &slice // luôn lưu *slice trong Pool
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
