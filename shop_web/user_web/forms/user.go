package forms

type PassWordLoginForm struct{
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type SmsLoginForm struct{
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Code string `form:"code" json:"code" binding:"required,min=6,max=6"`
}

type RegisterForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	NickName string `form:"NickName" json:"nick_name" binding:"max=10"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code string `form:"code" json:"code" binding:"required,min=6,max=6"`
}

type UpdateUserForm struct {
	Name string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}