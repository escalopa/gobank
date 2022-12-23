package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/escalopa/gobank/api/handlers/response"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

type loginUserReq struct {
	Username string `json:"username" binding:"required,min=6,max=16,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type loginUserRes struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_expires_at"`
	User                  userResponse `json:"user"`
}

// Login godoc
//
//	@Summary		Login user and return session
//	@Description	Login user and return session
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginUserReq	true	"Login user"
//	@Success		200		{object}	response.JSON{data=loginUserRes}
//	@Failure		400,500	{object}	response.JSON{}
//	@Router			/users/login [post]
func (s *GinServer) loginUser(ctx *gin.Context) {
	var req loginUserReq
	if err := parseBody(ctx, &req); err != nil {
		return
	}

	// Get User from DB by Username
	user, found := s.getUserIfExists(ctx, req.Username)
	if !found {
		return
	}

	// Check User's Password
	err := util.CheckHashedPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Err(err))
		return
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := s.tm.CreateToken(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	// Generate New Refresh Token for User
	refreshToken, refreshPayload, err := s.tm.CreateRefreshToken(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	// Create new session for User
	session, err := s.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpireAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	res := loginUserRes{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpireAt,
		User:                  mapUserToResponse(user),
	}
	ctx.JSON(http.StatusAccepted, res)
}

type createUserReq struct {
	Username        string `json:"username"  binding:"required,min=6,max=16,alphanum"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=16"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

// Register godoc
//
//	@Summary		Register user
//	@Description	Register user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createUserReq	true	"Create user"
//	@Success		200		{object}	response.JSON{data=userResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Router			/users/register [post]
func (s *GinServer) register(ctx *gin.Context) {
	var req createUserReq
	if err := parseBody(ctx, &req); err != nil {
		return
	}

	hashPassword, err := util.GenerateHashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("couldn't hash password, err: %v", err))
		return
	}

	user, err := s.db.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, response.Err(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(mapUserToResponse(&user)))
}

// Get User godoc
//
//	@Summary		Get current user info
//	@Description	Get current user info
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createUserReq	true	"Create user"
//	@Success		200		{object}	response.JSON{data=userResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/users [get]
func (s *GinServer) getUser(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, found := s.getUserIfExists(ctx, payload.Username)
	if !found {
		return
	}

	ctx.JSON(http.StatusOK, response.Success(mapUserToResponse(user)))
}

func (s *GinServer) getUserIfExists(ctx *gin.Context, username string) (*db.User, bool) {
	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, response.Err(err))
			return nil, false
		}

		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return nil, false
	}
	return &user, true
}

type updateUserReq struct {
	FullName    string `json:"full_name" binding:"alpha,required"`
	Email       string `json:"email" binding:"email"`
	OldPassword string `json:"old_password" binding:"min=6,max=16,required_with=NewPassword"`
	NewPassword string `json:"new_password" binding:"min=6,max=16,require"`
}

// UpdateUser godoc
//
//	@Summary		Update current user info
//	@Description	Update current user info
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		updateUserReq	true	"Update user"
//	@Success		200		{object}	response.JSON{data=userResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/users [patch]
func (s *GinServer) updateUser(ctx *gin.Context) {
	var req updateUserReq
	if err := parseBody(ctx, &req); err != nil {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, found := s.getUserIfExists(ctx, payload.Username)
	if !found {
		return
	}

	var params db.UpdateUserParams
	if req.Email == user.Email {
		ctx.JSON(http.StatusBadRequest, response.Err(ErrEmailSameAsOld))
		return
	}

	if err := util.CheckHashedPassword(req.OldPassword, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(ErrPasswordWrong))
		return
	}

	hashedPassword, err := util.GenerateHashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}
	params.HashedPassword = sql.NullString{String: hashedPassword, Valid: true}

	dbUser, err := s.db.UpdateUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(mapUserToResponse(&dbUser)))
}
