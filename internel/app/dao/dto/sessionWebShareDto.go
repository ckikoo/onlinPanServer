package dto

import (
	"encoding/json"
	"onlineCLoud/internel/app/ginx"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionWebShareDto struct {
	ShareId     string
	ShareUserId string
	Expire      string
	FileId      string
}

const sessionKey = "webShare_Key"

func (dto *SessionWebShareDto) SetSession(c *gin.Context) {
	session := sessions.Default(c)

	dtoJson, _ := json.Marshal(dto)
	session.Set(sessionKey+dto.ShareId, dtoJson)

	session.Save()
}
func GetSession(c *gin.Context, shareId string) *SessionWebShareDto {
	session := sessions.Default(c)

	var jsonV interface{}
	if jsonV = session.Get(sessionKey + shareId); jsonV == nil {
		ginx.ResOk(c)
		return nil
	}
	dto := new(SessionWebShareDto)

	buff := jsonV.([]uint8)
	byteSlice := make([]byte, len(buff))
	for i, v := range buff {
		byteSlice[i] = byte(v)
	}

	if err := json.Unmarshal(byteSlice, dto); err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return nil
	}

	return dto
}
