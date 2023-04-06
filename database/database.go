package database

import (
	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	var err error

	DB, err = gorm.Open(mysql.Open(config.DATABASE_DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(err.Error())
	}

	err = DB.SetupJoinTable(&model.User{}, "Followers", &model.Follower{})

	if err != nil {
		panic(err.Error())
	}

	err = DB.SetupJoinTable(&model.User{}, "Following", &model.Follower{})

	if err != nil {
		panic(err.Error())
	}

	err = DB.SetupJoinTable(&model.User{}, "PostLikes", &model.PostLike{})

	if err != nil {
		panic(err.Error())
	}

	err = DB.SetupJoinTable(&model.Post{}, "Likes", &model.PostLike{})

	if err != nil {
		panic(err.Error())
	}

	err = DB.SetupJoinTable(&model.PostComment{}, "SubComments", &model.PostSubComment{})

	if err != nil {
		panic(err.Error())
	}

	DB.AutoMigrate(
		&model.User{},
		&model.UserCredential{},
		&model.Session{},
		&model.Post{},
		&model.PostImage{},
		&model.PostComment{},
	)
}
