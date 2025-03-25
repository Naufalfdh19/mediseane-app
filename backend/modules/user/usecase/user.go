package usecase

import (
	"context"
	"fmt"
	"montelukast/modules/user/entity"
	"montelukast/modules/user/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	jwttoken "montelukast/pkg/jwt_token"
	"montelukast/pkg/logger"
	"montelukast/pkg/transaction"
	"os"
	"strconv"

	"github.com/resendlabs/resend-go"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(c context.Context, user entity.User) error
	Login(c context.Context, user entity.User) (string, error)
	ForgetPassword(c context.Context, user entity.User) error
	ResetPassword(c context.Context, user entity.ResetPassword) error
	VerifyEmail(c context.Context, token string) error
	CheckResetPassToken(c context.Context, token string) error
	GetProfile(c context.Context, userID int) (*entity.User, error)
	UpdateName(c context.Context, user entity.User) error
	ResendVerificationEmail(c context.Context, email string) error
}

type userUsecaseImpl struct {
	r  repository.UserRepo
	tr transaction.TransactorRepoImpl
	rc *resend.Client
}

func NewUserUsecase(r repository.UserRepo, tr transaction.TransactorRepoImpl, rc *resend.Client) userUsecaseImpl {
	return userUsecaseImpl{
		r:  r,
		tr: tr,
		rc: rc,
	}
}

func (u userUsecaseImpl) Register(c context.Context, user entity.User) error {
	err := apperror.ValidateUsernameAndPassword(user.Name, user.Password)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrRegister, err, apperror.ErrInvalidEmailOrPassword)
	}

	isUserExists, err := u.r.IsUserExistsByEmail(c, user.Email)
	if err != nil {
		return err
	}
	if isUserExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrRegister, apperror.ErrEmailAlreadyExists, apperror.ErrEmailAlreadyExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	user.Password = string(hash)
	user.Role = appconstant.ROLE_USER

	err = u.r.AddUser(c, user)
	if err != nil {
		return err
	}

	err = u.SendEmailForVerificationAccount(c, user.Email)
	if err != nil {
		logger.Log.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrRegister, apperror.ErrSendEmail, err))
	}
	return nil
}

func (u userUsecaseImpl) SendEmailForVerificationAccount(c context.Context, email string) error {
	jwt_token := jwttoken.NewJWT()

	userData, err := u.r.GetUserByEmail(c, email)
	if err != nil {
		return err
	}
	userIDStr := strconv.Itoa(userData.ID)

	token, err := jwt_token.GenerateJwtTokenForAuth(appconstant.JwtTokenVerifyEmailType, userIDStr, "")
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	err = u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		err = u.r.UpdateTokenForVerifyEmail(txCtx, userData.ID, token)
		if err != nil {
			return err
		}
		err = u.sendVerificationEmail(email, token)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u userUsecaseImpl) sendVerificationEmail(email string, token string) error {
	params := &resend.SendEmailRequest{
		From:    "mediSEAne <no-reply@mediseane.store>",
		To:      []string{email},
		Subject: "[mediSEAne] Account Verification",
		Html: fmt.Sprintf(
			`<p style="font-size:3rem;font-weigth:bold;margin:0px">
				medi<span style="color:#008081">SEA</span>ne
			</p>
			<p style="font-weight:bold">
				All Your <span style="color:#008081">Healthcare</span> Needs at Your Fingertips
			</p>
			<hr>
			<p style="font-weight:bold">Hello, welcome to mediSEAne!</p>
			<p>Verify your email with the button bellow!<p>
			<a href="%s%s" style="text-decoration:none !important;cursor:pointer !important;color:#FAFAFA !important">
				<div style="display:inline;background-color:#008081;padding:10px 30px;border:none;border-radius:30px;font-weight:600">
					Verify email
				</div>
			</a>`,
			os.Getenv("WEB_URL"),
			appconstant.VERIFY_EMAIL_PATH+token,
		),
	}

	_, err := u.rc.Emails.Send(params)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrRegister, apperror.ErrSendEmail, err)
	}
	return nil
}

func (u userUsecaseImpl) Login(c context.Context, user entity.User) (string, error) {
	jwt_token := jwttoken.NewJWT()
	isUserExists, err := u.r.IsUserExistsByEmail(c, user.Email)
	if err != nil {
		return "", err
	}
	if !isUserExists {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}

	hashPass, err := u.r.GetUserPasswordByEmail(c, user.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(user.Password))
	if err != nil {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}

	userData, err := u.r.GetUserByEmail(c, user.Email)
	if err != nil {
		return "", err
	}
	if userData.Role != appconstant.ROLE_USER {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrCredentialWrong, nil)
	}

	userIDStr := strconv.Itoa(userData.ID)
	token, err := jwt_token.GenerateJwtTokenForAuth(appconstant.JwtTokenAuthType, userIDStr, userData.Role)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return token, nil
}

func (u userUsecaseImpl) ForgetPassword(c context.Context, user entity.User) error {
	isUserExists, err := u.r.IsUserExistsByEmail(c, user.Email)
	if err != nil {
		return err
	}
	if !isUserExists {
		return nil
	}

	userData, err := u.r.GetUserByEmail(c, user.Email)
	if err != nil && err != apperror.ErrUserNotExists {
		return err
	}

	jwtToken := jwttoken.NewJWT()
	token, err := jwtToken.GenerateJwtTokenForAuth(appconstant.JwtTokenResetPassType, "", "")
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrForgetPassword, apperror.ErrTokenFailedToGenerated, err)
	}

	err = u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		err = u.r.AddResetPassToken(txCtx, userData.ID, token)
		if err != nil {
			return err
		}

		err = u.sendResetPasswordEmail(user.Email, token)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u userUsecaseImpl) sendResetPasswordEmail(email string, token string) error {
	params := &resend.SendEmailRequest{
		From:    "mediSEAne <no-reply@mediseane.store>",
		To:      []string{email},
		Subject: "[mediSEAne] Reset Password",
		Html: fmt.Sprintf(
			`<p style="font-size:3rem;font-weigth:bold;margin:0px">
				medi<span style="color:#008081">SEA</span>ne
			</p>
			<p style="font-weight:bold">
				All Your <span style="color:#008081">Healthcare</span> Needs at Your Fingertips
			</p>
			<hr>
			<p style="font-weight:bold">Hello, a reset password request has just been requested for this email account!</p>
			<p>Click the button bellow to continue the reset password process!</p>
			<a href="%s%s" style="text-decoration:none !important;cursor:pointer !important;color:#FAFAFA !important;">
				<div style="display:inline;background-color:#008081;padding:10px 30px;border:none;border-radius:30px;color:#FAFAFA !important;font-weight:600">
					Reset password
				</div>
			</a>
			<p>If you didn't make this request, you can ignore this email.</p>`,
			os.Getenv("WEB_URL"),
			appconstant.RESET_PASSWORD_PATH+token,
		),
	}

	_, err := u.rc.Emails.Send(params)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrForgetPassword, apperror.ErrSendEmail, err)
	}
	return nil
}

func (u userUsecaseImpl) CheckResetPassToken(c context.Context, token string) error {
	jwt_token := jwttoken.JwtTokenImpl{}
	jwtTokenClaims, err := jwt_token.ParseJwtTokenForAuth(c, token)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrCheckResetPassToken, apperror.ErrParseJwtToken, err)
	}
	if jwtTokenClaims.UserID == "0" || jwtTokenClaims.Type != appconstant.JwtTokenResetPassType {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrIsResetPassTokenExists, apperror.ErrUserUnauthorized, apperror.ErrUserUnauthorized)
	}

	isTokenExists, err := u.r.IsResetPasswordTokenExists(c, token)
	if err != nil {
		return err
	}
	if !isTokenExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrIsResetPassTokenExists, apperror.ErrTokenNotExists, err)
	}

	return nil
}

func (u userUsecaseImpl) ResetPassword(c context.Context, user entity.ResetPassword) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		isPasswordValid := apperror.IsPasswordValid(user.NewPassword)
		if !isPasswordValid {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrResetPassword, apperror.ErrPasswordNotValid, apperror.ErrPasswordNotValid)
		}

		isTokenExists, err := u.r.IsResetPasswordTokenExists(txCtx, user.Token)
		if err != nil {
			return err
		}
		if !isTokenExists {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrResetPassword, apperror.ErrTokenNotExists, err)
		}

		userID, err := u.r.GetUserIDBYResetPassToken(c, user.Token)
		if err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		user.NewPassword = string(hash)

		err = u.r.UpdateUserPassword(txCtx, user.NewPassword, userID)
		if err != nil {
			return err
		}

		err = u.r.DeleteResetPasswordToken(txCtx, user.Token)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u userUsecaseImpl) VerifyEmail(c context.Context, token string) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		jwt_token := jwttoken.JwtTokenImpl{}
		jwtTokenClaims, err := jwt_token.ParseJwtTokenForAuth(c, token)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrVerifyEmail, apperror.ErrParseJwtToken, err)
		}
		if jwtTokenClaims.UserID == "0" || jwtTokenClaims.Type != appconstant.JwtTokenVerifyEmailType {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrVerifyEmail, apperror.ErrUserUnauthorized, apperror.ErrUserUnauthorized)
		}

		isTokenExists, err := u.r.IsVerifyEmailTokenExists(txCtx, token)
		if err != nil {
			return err
		}
		if !isTokenExists {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrVerifyEmail, apperror.ErrTokenNotExists, apperror.ErrTokenNotExists)
		}

		userID, err := strconv.Atoi(jwtTokenClaims.UserID)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrVerifyEmail, apperror.ErrConvertVariableType, err)
		}

		err = u.r.UpdateEmailVerifyStatus(txCtx, userID)
		if err != nil {
			return err
		}

		err = u.r.DeleteVerifyEmailToken(txCtx, token)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u userUsecaseImpl) GetProfile(c context.Context, userID int) (*entity.User, error) {
	isUserExists, err := u.r.IsUserExistsByID(c, userID)
	if err != nil {
		return nil, err
	}
	if !isUserExists {
		return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUser, apperror.ErrUserNotExists, err)
	}

	user, err := u.r.GetUserByID(c, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userUsecaseImpl) UpdateName(c context.Context, user entity.User) error {
	isExists, err := u.r.IsUserExistsByID(c, user.ID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdateUser, apperror.ErrUserNotExists, apperror.ErrUserNotExists)
	}

	err = u.r.UpdateNameByID(c, user)
	if err != nil {
		return err
	}

	return nil
}

func (u userUsecaseImpl) ResendVerificationEmail(c context.Context, email string) error {
	isUserExists, err := u.r.IsUserExistsByEmail(c, email)
	if err != nil {
		return err
	}
	if !isUserExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrEmail, apperror.ErrUserNotExists, apperror.ErrUserNotExists)
	}
	err = u.SendEmailForVerificationAccount(c, email)
	if err != nil {
		return err
	}

	return nil
}
