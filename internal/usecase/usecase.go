package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/SerzhLimon/testJWT/internal/models"
	"github.com/SerzhLimon/testJWT/internal/repository"
)

var jwtSecretKey = []byte("secretkeyforaccesstoken")

type Usecase struct {
	pgPepo repository.Repository
}

type UseCase interface {
	CreatePairTokens(ID, IP string) (models.CreatePairTokensResponse, error)
}

func NewUsecase(pgPepo repository.Repository) UseCase {
	return &Usecase{pgPepo: pgPepo}
}

func (u *Usecase) CreatePairTokens(ID, IP string) (models.CreatePairTokensResponse, error) {

	userID, err := u.parsedUUID(ID)
	if err != nil {
		err = errors.Errorf("usecase.CreatePairTokens: incorrect userID  %v", err)
		return models.CreatePairTokensResponse{}, err
	}

	accessToken, err := u.createAccessToken(userID, IP)
	if err != nil {
		err = errors.Errorf("usecase.CreatePairTokens: fail to create access token %v", err)
		return models.CreatePairTokensResponse{}, err
	}

	refreshToken, err := u.createRefreshToken()
	if err != nil {
		err = errors.Errorf("usecase.CreatePairTokens: fail to create refresh token %v", err)
		return models.CreatePairTokensResponse{}, err
	}

	hashedRefreshToken, err := u.createHashedRefreshToken(refreshToken)
	if err != nil {
		err = errors.Errorf("usecase.CreatePairTokens: fail to create hashed refresh token %v", err)
		return models.CreatePairTokensResponse{}, err
	}
	
	err = u.pgPepo.SetUserInfo(userID, IP, hashedRefreshToken)
	if err != nil {
		return models.CreatePairTokensResponse{}, err
	}

	return models.CreatePairTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *Usecase) parsedUUID(data string) (uuid.UUID, error) {
	id, err := uuid.Parse(data)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (u *Usecase) createAccessToken(userID uuid.UUID, userIP string) (string, error) {
	tokenClaims := models.TokenClaims{
		UserID: userID,
		IP:     userIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaims)
	return token.SignedString(jwtSecretKey)
}

func (u *Usecase) createRefreshToken() (string, error) {

	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", nil
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (u *Usecase) createHashedRefreshToken(refreshToken string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
}
