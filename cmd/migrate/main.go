package main

import (
	"fmt"
	"os"
	"path/filepath"

	"tamago/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func main() {
	goPath := os.Getenv("GOPATH")
	viper.SetConfigName("development_config")
	viper.AddConfigPath(filepath.Join(goPath, "src/force/configs/"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	config := viper.GetViper()

	dbConfig := config.GetStringMapString("postgres")
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig["host"], dbConfig["port"], dbConfig["user"], dbConfig["password"], dbConfig["db"], dbConfig["ssl_mode"]))

	if err != nil {
		fmt.Println(err)
		panic("Failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Session{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.UserProduct{}).AddForeignKey("host_id", "users(id)", "CASCADE", "NO ACTION").AddForeignKey(
		"product_id", "products(id)", "CASCADE", "NO ACTION")
	db.AutoMigrate(&models.UserProductAccount{}).AddForeignKey(
		"user_id", "users(id)", "CASCADE", "NO ACTION").AddForeignKey(
		"product_id", "products(id)", "CASCADE", "NO ACTION")
	db.AutoMigrate(&models.UserProductMember{}).AddForeignKey(
		"user_id", "users(id)", "CASCADE", "NO ACTION").AddForeignKey(
		"user_product_id", "user_products(id)", "CASCADE", "NO ACTION")
	db.AutoMigrate(&models.UserThirdPartyPayment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "NO ACTION")
	db.AutoMigrate(&models.Notification{}).AddForeignKey("user_id", "users(id)", "CASCADE", "NO ACTION")
	db.AutoMigrate(&models.BillingRequest{}).AddForeignKey(
		"user_product_id", "user_products(id)", "CASCADE", "NO ACTION").AddForeignKey(
		"user_product_member_id", "user_product_members(id)", "CASCADE", "NO ACTION")
}
