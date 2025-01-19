package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symmetrickey []byte
}

func NewPasetoMaker(symmetrickey string)(Maker,error){
	if len(symmetrickey)!=chacha20poly1305.KeySize{
		return nil,fmt.Errorf("invalid key size : must be exactly %d characters",chacha20poly1305.KeySize)
	}
	maker:=&PasetoMaker{
		paseto: paseto.NewV2(),
		symmetrickey: []byte(symmetrickey),
	}
	return maker,nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker)CreateToken(username string,duration time.Duration)(string,error){
	payload,err :=NewPayload(username,duration)
	if err!=nil{
		return"",err
	}
	return maker.paseto.Encrypt(maker.symmetrickey,payload,nil)
}

// VerifyToken checks if the token is valid orn not
func (maker *PasetoMaker)VerifyToken(token string)(*Payload,error){
	payload := &Payload{}
	err := maker.paseto.Decrypt(token,maker.symmetrickey,payload,nil)
	if err!=nil{
		return nil,err
	}
	err = payload.Valid()
	if err!=nil{
		return nil,ErrInvalidToken
	}
	return payload,nil
}