package wxapkg

import (
	"fmt"
	"os"
)

const (
	// wechat 3.8.* version
	wxappletPath = "/Users/%s/Library/Containers/com.tencent.xinWeChat/Data/.wxapplet/packages/"
)

func GetWXAppletPath() string {
	return fmt.Sprintf(wxappletPath, os.Getenv("USER"))
}
