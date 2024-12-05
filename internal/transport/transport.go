package transport

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/testJWT/internal/repository"
	uc "github.com/SerzhLimon/testJWT/internal/usecase"
)

type Server struct {
	Usecase uc.UseCase
}

func NewServer(database *sql.DB) *Server {
	pgClient := repository.NewPGRepository(database)
	uc := uc.NewUsecase(pgClient)

	return &Server{
		Usecase: uc,
	}
}

func (s *Server) CreatePairTokens(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		err := errors.New("parametr 'id' is empty")
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'id' is empty"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %s", userID)

	userIP := c.ClientIP()
	result, err := s.Usecase.CreatePairTokens(userID, userIP)
	if err == nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create tokens"})
		return
	}

	c.JSON(http.StatusOK, result)
}
