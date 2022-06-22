package wrap

import (
	"net/netip"

	"github.com/gin-gonic/gin"
)

func GetAddr(ctx *gin.Context) (ap netip.AddrPort, err error) {
	var addr netip.Addr

	if ap, err = netip.ParseAddrPort(ctx.Request.RemoteAddr); err != nil {
		return ap, err
	}

	if addr, err = netip.ParseAddr(ctx.ClientIP()); err != nil {
		return ap, err
	}

	return netip.AddrPortFrom(addr, ap.Port()), nil
}

func GetAddrString(ctx *gin.Context) string {
	ap, _ := GetAddr(ctx)
	return ap.String()
}
