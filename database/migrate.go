package database

import (
	"snipetty.com/main/repositories"	
)

func init() {
	LoadEnvs()
	InitializeDatabaseLayer()
}

func main() {
    dbInstance.AutoMigrate(&repositories.User{})
}