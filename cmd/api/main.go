package main

import (
	"context"
	"filmoteka/internal/config"
	"filmoteka/internal/infrastructure"
	"filmoteka/internal/repository"
	"filmoteka/internal/service"
	"filmoteka/internal/transport"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("viper OK")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	pgConfig, err := config.GetDBConfig()

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("got bd config")

	db, err := infrastructure.SetUpPostgresDatabase(ctx, pgConfig)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("connected to db")

	//actorRepository := repository.NewActorRepository(db)
	movieRepository := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepository)
	movieHandler := transport.NewMovieHandler(movieService)
	// actorRepository.DeleteActor(ctx, 4)

	// actor1 := core.Actor{Name: "benedict cumberbatch", Sex: 109, Bd: "1976-07-19"}
	// actor2 := core.Actor{Name: "martin freeman", Sex: 109, Bd: "1971-09-08"}
	// actorRepository.CreateActor(ctx, &actor1)
	// actorRepository.CreateActor(ctx, &actor2)

	// movie1 := &core.Movie{Title: "Sherlock Holmes5",
	// 	Descr:   "Dr Watson, a former army doctor, finds himself sharing a flat with Sherlock Holmes, an eccentric individual with a knack for solving crimes. Together, they take on the most unusual cases.",
	// 	Release: "2010-07-25"}

	// actors1 := []string{"benedict cumberbatch", "martin freeman", "andrew scott"}
	// movie1, err = movieRepository.CreateMovie(ctx, movie1, actors1)
	// if err != nil {
	// 	log.Info(err.Error())
	// }

	http.HandleFunc("/movie", movieHandler.CreateMovie)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func SetupViper() error {

	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
