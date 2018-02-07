package handlers

import (
	"github.com/valyala/fasthttp"
	"../services"
	"../models"
	"../cache"
	"../database"
	"encoding/json"
	"time"
	"github.com/satori/go.uuid"
)

const CookieAccess = "ws_auth"
const CookieRefresh = "ws_auth_refresh"


func CreateToken(ctx *fasthttp.RequestCtx) {

	redirectURL := string(ctx.QueryArgs().Peek("redirect"))
	if len(redirectURL) == 0 {
		err := serv.NewBadRequest("Redirect path is missed")
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	oauth := models.OAuth{}
	json.Unmarshal(ctx.PostBody(), &oauth)
	err := database.Auth(&oauth)
	if err != nil {
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	code := cache.CreateCode(&oauth)
	if code != nil {
		err = &serv.ErrorCode {
			Code: fasthttp.StatusInternalServerError,
			Message: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
		}
		err.WriteAsJsonResponseTo(ctx)
	}
	ctx.Redirect(redirectURL + "?code=" + code.String(), fasthttp.StatusFound)
}

func GetToken(ctx *fasthttp.RequestCtx) {

	accessCode := string(ctx.QueryArgs().Peek("code"))
	uuidCode, err := uuid.FromString(string(accessCode))

	if err != nil {
		err := serv.NewBadRequest("Access code is invalid")
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	accessToken, refreshToken := cache.ActivateCode(&uuidCode)
	if accessToken == nil || refreshToken == nil {
		err := serv.ErrorCode {
			Code: fasthttp.StatusForbidden,
			Message: "Access code is not valid \nTry register again",
		}
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	tokens := models.Tokens{
		AccessToken: *accessToken,
		RefreshToken: *refreshToken,
	}
	tokens.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
}

func RefreshToken(ctx *fasthttp.RequestCtx) {

	oldRefreshToken := string(ctx.QueryArgs().Peek("refresh_token"))
	if len(oldRefreshToken) == 0 {
		err := serv.NewBadRequest("Refresh token is missed")
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	accessToken, refreshToken := cache.RefreshToken(oldRefreshToken)
	if accessToken == nil || refreshToken == nil {
		err := serv.ErrorCode {
			Code: fasthttp.StatusForbidden,
			Message: "Refresh token is not valid \nTry register again",
		}
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	tokens := models.Tokens{
		AccessToken: *accessToken,
		RefreshToken: *refreshToken,
	}
	tokens.WriteAsJsonResponseTo(ctx, fasthttp.StatusOK)
}



// Вообще, эта функция может просто вызвать GetToken(),
// который даже находится в этом же пакете, но по легенде
// Auth-server находится на отдельной машине
func ApplicationActivate(ctx *fasthttp.RequestCtx) {

	// Парсим код доступа из query-параметров,
	// необходимый для получения токенов
	code := ctx.QueryArgs().Peek("code")
	if len(code) == 0 {
		err := serv.NewBadRequest("Access code in query arguments is missed")
		err.WriteAsJsonResponseTo(ctx)
		return
	}


	// Отправляем запрос с кодом доступа на auth-server
	statusCode, body, err := fasthttp.Get(
		nil, serv.GetConfig().Href + "/oauth/access?code=" + string(code))

	if err != nil {
		errorCode := serv.NewServerError(err)
		errorCode.WriteAsJsonResponseTo(ctx)
		return
	}

	if statusCode != fasthttp.StatusOK {
		ctx.SetStatusCode(statusCode)
		ctx.Write(body)
		return
	}

	// Если ответ auth-server'a вернулся без ошибок,
	// то парсим тело ответа в поисках токенов
	var tokens models.Tokens
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		errorCode := serv.NewServerError(err)
		errorCode.WriteAsJsonResponseTo(ctx)
		return
	}

	// И проставляем токены в качестве cookie в браузeр пользователя
	setTokenCookie(tokens, ctx)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// Аналогично, но уже по отношению к RefreshToken()
func ApplicationRefresh(ctx *fasthttp.RequestCtx) {

	// Вытаскиваем из кук refresh-token
	refreshToken := ctx.Request.Header.Cookie(CookieRefresh)
	if len(refreshToken) == 0 {
		err := serv.NewBadRequest("Refresh token in cookies not found")
		err.WriteAsJsonResponseTo(ctx)
		return
	}

	// Отправляем запрос с кодом доступа на auth-server
	statusCode, body, err := fasthttp.Get(
		nil, serv.GetConfig().Href + "/oauth/refresh?refresh_token=" + string(refreshToken))

	if err != nil {
		errorCode := serv.NewServerError(err)
		errorCode.WriteAsJsonResponseTo(ctx)
		return
	}

	if statusCode != fasthttp.StatusOK {
		ctx.SetStatusCode(statusCode)
		ctx.Write(body)
		return
	}

	var tokens models.Tokens
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		errorCode := serv.NewServerError(err)
		errorCode.WriteAsJsonResponseTo(ctx)
		return
	}
	setTokenCookie(tokens, ctx)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func setTokenCookie(tokens models.Tokens, ctx *fasthttp.RequestCtx) {
	var accessCookie, refreshCookie fasthttp.Cookie

	accessCookie.SetKey(CookieAccess)
	accessCookie.SetValue(tokens.AccessToken)
	accessCookie.SetExpire(time.Now().AddDate(0, 1, 0))
	refreshCookie.SetKey(CookieRefresh)
	refreshCookie.SetValue(tokens.RefreshToken)
	refreshCookie.SetExpire(time.Now().AddDate(0, 1, 0))

	ctx.Response.Header.SetCookie(&accessCookie)
	ctx.Response.Header.SetCookie(&refreshCookie)
}
