package public

import (
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/global"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/model"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiLogin struct{}

func (a *ApiLogin) Login(c *gin.Context) {
	var l model.RequestLogin
	if err := model.RequestShouldBindJSON(c, &l); err != nil {
		return
	}
	if global.Store.Verify(l.CaptchaId, l.Captcha, true) {
		u := &model.SysUser{Username: l.Username, Password: l.Password}
		if err, user := servicePublic.ServiceLogin.Login(u); err != nil {
			global.GqaLogger.Error(l.Username + utils.GqaI18n("LoginFailed"))
			model.ResponseErrorMessage(utils.GqaI18n("LoginFailed"), c)
			if err = servicePublic.ServiceLogin.LogLogin(l.Username, c, "yesNo_no", utils.GqaI18n("LoginFailed")); err != nil {
				global.GqaLogger.Error(utils.GqaI18n("LoginLogError"), zap.Any("err", err))
			}
		} else {
			a.createToken(*user, c)
		}
	} else {
		model.ResponseErrorMessage(utils.GqaI18n("CaptchaError"), c)
		if err := servicePublic.ServiceLogin.LogLogin(l.Username, c, "yesNo_no", utils.GqaI18n("CaptchaError")); err != nil {
			global.GqaLogger.Error(utils.GqaI18n("LoginLogError"), zap.Any("err", err))
		}
	}
}

func (a *ApiLogin) createToken(user model.SysUser, c *gin.Context) {
	ss := utils.CreateToken(user.Username)
	if ss == "" {
		global.GqaLogger.Error(utils.GqaI18n("JwtConfigError"))
		model.ResponseErrorMessage(utils.GqaI18n("JwtConfigError"), c)
		return
	}
	if err := servicePublic.ServiceLogin.LogLogin(user.Username, c, "yesNo_yes", utils.GqaI18n("LoginSuccess")); err != nil {
		global.GqaLogger.Error(utils.GqaI18n("LoginLogError"), zap.Any("err", err))
	}
	if err := servicePublic.ServiceLogin.SaveOnline(user.Username, ss); err != nil {
		global.GqaLogger.Error(utils.GqaI18n("UserOnlineFailed"), zap.Any("err", err))
	}
	model.ResponseSuccessMessageData(model.ResponseLogin{
		Avatar:   user.Avatar,
		Username: user.Username,
		Nickname: user.Nickname,
		RealName: user.RealName,
		Token:    ss,
	}, utils.GqaI18n("LoginSuccess"), c)
}
