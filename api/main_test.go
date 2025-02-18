package api

import (
	"context"
	"log"
	"os"
	"testing"

	db "github.com/Tech-With-Tim/cdn/db/sqlc"
	"github.com/Tech-With-Tim/cdn/server"
	"github.com/Tech-With-Tim/cdn/utils"
	"github.com/go-chi/chi/v5"
)

var config *utils.Config
var s *server.Server

func TestMain(m *testing.M) {
	conf, err := utils.LoadConfig("../", "test")
	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}
	config = &conf
	s = server.NewServer(conf)
	CdnRouter := chi.NewRouter()
	//Add Routes to Routers Here
	MainRouter(CdnRouter, s.Store, conf)
	//Mount Routers here
	s.Router.Mount("/", CdnRouter)
	err = createTestUser()
	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	os.Exit(m.Run())

}

func createTestUser() error {
	user := db.CreateUserParams{
		ID:            328604827967815690,
		Username:      utils.RandomString(4),
		Discriminator: "3212",
	}
	err := s.Store.CreateUser(context.Background(), user)
	return err
}
