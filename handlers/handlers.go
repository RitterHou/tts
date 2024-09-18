package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tts/utils"
)

func GetVoiceList(c *gin.Context) {
	locale := c.Query("l")
	voices, err := utils.VoiceList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if locale != "" {
		filteredVoices := make([]interface{}, 0)
		for _, voice := range voices {
			if strings.Contains(voice.(map[string]interface{})["Locale"].(string), locale) {
				filteredVoices = append(filteredVoices, voice)
			}
		}
		voices = filteredVoices
	}

	_, detail := c.GetQuery("d")
	if detail {
		c.JSON(http.StatusOK, gin.H{"voices": voices})
	} else {
		voiceSimpleList := make([]map[string]string, 0)
		for _, voice := range voices {
			localName := voice.(map[string]interface{})["LocalName"].(string)
			shortName := voice.(map[string]interface{})["ShortName"].(string)
			voiceSimpleList = append(voiceSimpleList, map[string]string{
				"LocalName": localName,
				"ShortName": shortName,
			})
		}
		c.JSON(http.StatusOK, gin.H{"voices": voiceSimpleList})
	}

}

func SynthesizeVoice(c *gin.Context) {
	text := c.Query("t")
	voiceName := c.DefaultQuery("v", "zh-CN-XiaochenMultilingualNeural")
	rate := c.DefaultQuery("r", "0")
	pitch := c.DefaultQuery("p", "0")
	outputFormat := c.DefaultQuery("o", "audio-24khz-48kbitrate-mono-mp3")

	voice, err := utils.GetVoice(text, voiceName, rate, pitch, outputFormat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "audio/mpeg", voice)
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "TTS",
	})
}

func SynthesizeVoicePost(c *gin.Context) {
	var text, voiceName, rate, pitch, outputFormat string
	if c.GetHeader("Content-Type") == "application/json" {
		request := struct {
			Text         string `json:"t"`
			VoiceName    string `json:"v"`
			Rate         string `json:"r"`
			Pitch        string `json:"p"`
			OutputFormat string `json:"o"`
		}{
			Text:         "",
			VoiceName:    "zh-CN-XiaochenMultilingualNeural",
			Rate:         "0",
			Pitch:        "0",
			OutputFormat: "audio-24khz-48kbitrate-mono-mp3",
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		text, voiceName, rate, pitch, outputFormat = request.Text, request.VoiceName, request.Rate, request.Pitch, request.OutputFormat
	} else {
		text = c.PostForm("t")
		voiceName = c.DefaultPostForm("v", "zh-CN-XiaochenMultilingualNeural")
		rate = c.DefaultPostForm("r", "0")
		pitch = c.DefaultPostForm("p", "0")
		outputFormat = c.DefaultPostForm("o", "audio-24khz-48kbitrate-mono-mp3")
	}

	voice, err := utils.GetVoice(text, voiceName, rate, pitch, outputFormat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "audio/mpeg", voice)
}
