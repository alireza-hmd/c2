/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/alireza-hmd/c2/clients"
	cService "github.com/alireza-hmd/c2/clients"
	cRepo "github.com/alireza-hmd/c2/clients/mysql"
	"github.com/alireza-hmd/c2/cmd"
	"github.com/alireza-hmd/c2/listeners"
	lService "github.com/alireza-hmd/c2/listeners"
	lRepo "github.com/alireza-hmd/c2/listeners/mysql"
	"github.com/alireza-hmd/c2/pkg/configs"
	"github.com/alireza-hmd/c2/tasks"
	tService "github.com/alireza-hmd/c2/tasks"
	tRepo "github.com/alireza-hmd/c2/tasks/mysql"
	ms "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := configs.LoadConfig()
	if err != nil {
		log.Fatalln("Error in loading .env file :", err)
	}
	db, err := InitDB()
	if err != nil {
		log.Panic(err)
	}
	tasksRepo := tRepo.NewRepository(db)
	tasksService := tService.NewService(tasksRepo)

	clientRepo := cRepo.NewRepository(db)
	clientService := cService.NewService(clientRepo)
	listenerRepo := lRepo.NewRepository(db)
	listenerService := lService.NewService(listenerRepo, clientRepo, tasksRepo)

	stopChannel := make(map[int](chan listeners.Cancel))
	s := &cmd.Services{
		Listener: listenerService,
		Tasks:    tasksService,
		Clients:  clientService,
		Stop:     stopChannel,
	}
	s.Listener.RunActiveListeners(stopChannel)
	cmd.Init(s)
}

func InitDB() (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", configs.Get("DB_USER"), configs.Get("DB_PASS"), configs.Get("DB_HOST"), configs.Get("DB_PORT"), configs.Get("DB_DATABASE"))
	db, err := gorm.Open(ms.Open(dataSourceName), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Println(err)
		return nil, errors.New("error initializing database connection")
	}

	if err := db.Migrator().AutoMigrate(&listeners.Listener{}, &clients.Client{}, &tasks.Task{}); err != nil {
		log.Println(err)
		return nil, errors.New("error migrating models")
	}
	return db, nil
}
