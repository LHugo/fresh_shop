package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_web/user_web/forms"
	"shop_web/user_web/global"
	"shop_web/user_web/global/response"
	"shop_web/user_web/middlewares"
	"shop_web/user_web/models"
	"shop_web/user_web/proto"
	"strconv"
	"strings"
	"time"
)


func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields{
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	//处理表单验证出现的错误
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error" : removeTopStruct(errs.Translate(global.Trans)),
	})
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":"该用户已存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("目前访问用户：%d", currentUser.ID)

	pn, _ := strconv.Atoi(ctx.DefaultQuery("pn", "0"))
	psize, _ := strconv.Atoi(ctx.DefaultQuery("psize", "0"))

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pn),
		PSize: uint32(psize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.Nickname,
			//Birthday: time.Time(time.Unix(int64(value.Birthday), 0)).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		//data["id"] = value.Id
		//data["name"] = value.Nickname
		//data["birthday"] = value.Birthday
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func PasswordLogin(c *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil{
		HandleValidatorError(c, err)
		return
	}
	//验证码验证
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false){
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha":"验证码错误",
		})
		return
	}

	//登录的逻辑
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err !=nil{
		if e, ok := status.FromError(err); ok{
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
				"mobile":"用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile":"登陆失败",
				})
			}
			return
		}
	}else {
		//用户存在并进一步检查用户密码以及生成签名
		if passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:passwordLoginForm.PassWord,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password":"登录失败",
			})
		}else {
			//检查密码
			if passRsp.Success{
				//若密码正确则生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID: uint(rsp.Id),
					NickName: rsp.Nickname,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),//签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*7, //7天过期(秒/分/时/天）
						Issuer: "hugo",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil{
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":"生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":rsp.Id,
					"nick_name":rsp.Nickname,
					"token": token,
					"expired_at": (time.Now().Unix() + 60*60*24*7)*1000,
				})
			}else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg":"密码错误",
				})
			}
		}
	}
}

func SmsLogin(c *gin.Context) {
	//表单验证
	SmsLoginForm := forms.SmsLoginForm{}
	if err := c.ShouldBind(&SmsLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	//检查用户是否存在
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: SmsLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登陆失败",
				})
			}
			return
		}
	} else {
		//若用户存在则对短信验证码进行校验
		rdb := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		})
		if value, err := rdb.Get(context.Background(), SmsLoginForm.Mobile).Result(); err == redis.Nil {
			//验证码不存在
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		} else {
			//验证码存在但不正确
			if value != SmsLoginForm.Code {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": "验证码错误",
				})
				return
			}
		}

		//短信验证码校验成功之后生成token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(rsp.Id),
			NickName:    rsp.Nickname,
			AuthorityId: uint(rsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),               //签名的生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24*7, //7天过期(秒/分/时/天）
				Issuer:    "hugo",
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":         rsp.Id,
			"nick_name":  rsp.Nickname,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*7) * 1000,
		})
	}
}

func Register(c *gin.Context)  {
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil{
		HandleValidatorError(c, err)
		return
	}
	//验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	if value, err := rdb.Get(context.Background(), registerForm.Mobile).Result();err == redis.Nil{
		//验证码不存在
		c.JSON(http.StatusBadRequest, gin.H{
			"code":"验证码错误",
		})
		return
	}else {
		//验证码存在但不正确
		if value != registerForm.Code{
			c.JSON(http.StatusBadRequest, gin.H{
				"code":"验证码错误",
			})
			return
		}
	}

	var nickN string
	if registerForm.NickName != ""{
		nickN = registerForm.NickName
	}else {
		nickN = registerForm.Mobile
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: nickN ,
		Password: registerForm.PassWord,
		Mobile: registerForm.Mobile,
	})
	if err != nil {
		//zap.S().Errorf("[Register] 执行【新建用户】失败:%s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//注册后登录并生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: uint(user.Id),
		NickName: user.Nickname,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),//签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, //7天过期
			Issuer: "hugo",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":"生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":user.Id,
		"nick_name":user.Nickname,
		"token": token,
		"expired_at": (time.Now().Unix() + 60*60*24*7)*1000,
	})
}

func GetUserDetail(c *gin.Context)  {
	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("目前访问用户：%d", currentUser.ID)

	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil{
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":rsp.Nickname,
		"birthday":time.Unix(int64(rsp.Birthday), 0).Format("2006-01-01"),
		"gender":rsp.Gender,
		"mobile":rsp.Mobile,
	})
}

func UpdateUser(c *gin.Context)  {
	updateUserForm := forms.UpdateUserForm{}
	if err := c.ShouldBind(&updateUserForm); err != nil{
		HandleValidatorError(c, err)
		return
	}

	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("目前访问用户：%d", currentUser.ID)

	//将前端传递来的日期格式转换成int
	loc, _ := time.LoadLocation("Local")
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		NickName: updateUserForm.Name,
		Birthday: uint64(birthDay.Unix()),
		Gender:   updateUserForm.Gender,
	})
	if err != nil{
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":"修改成功",
	})
}