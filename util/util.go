package util

import (
	"strings"

	"github.com/huin/goupnp"
)

func DeviceUUID(dev *goupnp.RootDevice) string {
	return strings.TrimPrefix(dev.Device.UDN, "uuid:")
}
