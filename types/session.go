package types

import "fmt"

type Session struct {
    UserId int64
    SessionId string
}

func (session *Session) ToCookie(stay_logged bool) string {
    age := ""
    if stay_logged {
        age = "max-age=31536000;"
    }

    return fmt.Sprintf("session=%s;%ssamesite=strict;secure;HTTPonly", session.SessionId, age)
}
