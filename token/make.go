package token

import "time"

type Maker interface {
	//通过用户openid和有效时间生成token
	CreateToken(openid string, duration time.Duration) (string, error)

	//验证token是否有效
	VerifyToken(token string) (*Payload, error)
}
