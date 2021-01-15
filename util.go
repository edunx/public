package public

import (
	"github.com/edunx/lua"
	"math"
)

func CheckUserData(L *lua.LState, idx int) UserData {
	v := L.ToUserData(idx)

	switch ud := v.Value.(type) {
	case UserData:
		return ud
	default:
		L.RaiseError("must common.Userdata , got %v", ud)
		return nil
	}
}

func CheckTransport(ud *lua.LUserData) Transport {
	switch tp := ud.Value.(type) {
	case Transport:
		return tp
	default:
		return nil
	}
}

func CheckTransportByTable(key string, tb *lua.LTable) Transport {
	v := tb.RawGetString(key)

	switch ud := v.(type) {
	case *lua.LUserData:
		if tp := CheckTransport(ud); tp == nil {
			Out.Err("%s not transport userdata", key)
			return nil
		} else {
			return tp
		}
	default:
		Out.Err("%s must userdata , got %v", key, ud)
		return nil
	}
}

func Round(val float64, precision int64) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(precision))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
