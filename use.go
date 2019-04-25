package NetWorkRate

import (
	"fmt"
)

const KB = 1024

//将值转换为固定长度的字符串，目前暂时只支持字符串和float64
func VolumeGenerate(length int, data interface{}) string {
	tmp := ""
	dc, ok := data.(string)
	if ok {
		tmp = dc
	} else {
		dc, ok := data.(float64)
		if ok {
			tmp = fmt.Sprintf("%.4f", dc)
		} else {
			return ""
		}
	}

	d := length - len(tmp)
	if d <= 0 {
		return tmp
	}
	temp := ""
	for i := 0; i < d; i++ {
		temp += " "
	}
	return tmp + temp
}

//打印速率
func LinesPrint(length int, rate *IORate) {
	lines := ""
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, rate.Name))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, rate.RecvPacketsRate/KB))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, rate.SentPacketsRate/KB))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, rate.RecvBytesRate/KB))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, rate.SentBytesRate/KB))
	fmt.Println(lines)
}

//第一行的说明文字
func FirstLine(length int) {
	lines := ""
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, "interface"))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, "rxpck/s"))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, "txpck/s"))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, "rxkB/s"))
	lines = fmt.Sprintf("%s%s", lines, VolumeGenerate(length, "txkB/s"))
	fmt.Println(lines)
}
