package main

import (
	"errors"
	"fs_backend/apierrors"
	"fs_backend/models"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (apiCfg ApiConfig) generateJwtToken(data models.JWTData) (string, error) {
	var secretKey = []byte(apiCfg.jwtSecret)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":        data.Issuer,
			"remoteAddr": data.RemoteAddr,
			"tokenId":    data.TokenId,
			"accountId":  data.AccountId,
			"name":       data.Name,
			"username":   data.Username,
			"isOwner":    data.IsOwner,
		},
	)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

func (apifn *ApiConfig) parseJwtFromHeader(req *http.Request) (models.JWTData, error) {
	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		return models.JWTData{}, apierrors.NoJWTToken{}
	}
	vals := strings.Split(tokenString, " ")
	if len(vals) != 2 {
		return models.JWTData{}, apierrors.InvalidToken{}
	}
	tokenString = vals[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(apifn.jwtSecret)
		return secret, nil
	})
	if err != nil {
		log.Default().Panicln(err.Error())
		return models.JWTData{}, apierrors.InvalidToken{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		issuer := claims["iss"]
		RemoteAddr := claims["remoteAddr"]
		tokenId := claims["tokenId"]
		accountId := claims["accountId"]
		name := claims["name"]
		username := claims["username"]
		isOwner := claims["isOwner"]
		if issuer == nil || RemoteAddr == nil || tokenId == nil || accountId == nil ||
			name == nil || username == nil || isOwner == nil {
			return models.JWTData{}, apierrors.InvalidToken{}
		}
		data := models.JWTData{
			Issuer:     claims["iss"].(string),
			RemoteAddr: claims["remoteAddr"].(string),
			TokenId:    claims["tokenId"].(string),
			AccountId:  claims["accountId"].(string),
			Name:       claims["name"].(string),
			Username:   claims["username"].(string),
			IsOwner:    claims["isOwner"].(bool),
		}
		return data, nil
	} else {
		return models.JWTData{}, err
	}
}

func (apifn *ApiConfig) authMiddleware(handler func(http.ResponseWriter, *http.Request, models.JWTData)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		claims, err := apifn.parseJwtFromHeader(req)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.NoJWTToken{}) {
				ErrorResponseWriter(res, apierrors.ResErrTokenNotFound, http.StatusBadRequest)
				return
			}
			if errors.Is(err, apierrors.InvalidToken{}) {
				ErrorResponseWriter(res, apierrors.ResErrTokenInvalid, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		handler(res, req, claims)
	}
}
