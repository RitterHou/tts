package utils

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	data, err := GetVoice("你好，世界，今天是", "zh-CN-XiaochenNeural", "0", "0", "audio-24khz-48kbitrate-mono-mp3")
	if err != nil {
		panic(err)
	}
	f, _ := os.Create("voice.mp3")
	i, _ := f.Write(data)
	fmt.Println(i)
}
