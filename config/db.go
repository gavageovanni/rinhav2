package config

import (
	"baitadores-rinhav2/model"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var e error

func DatabaseInit() {
	host := "db"
	//host := "localhost"
	user := "rinhav2"
	password := "rinhav2"
	dbName := "rinhav2"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	sqlDB, err := database.DB()
	if err != nil {
		panic("failed to get DB instance")
	}
	// Definindo tamanho máximo do pool de conexões
	sqlDB.SetMaxOpenConns(100)
	// Definindo tamanho máximo do pool de conexões em uso
	sqlDB.SetMaxIdleConns(10)
	// Definindo tempo máximo de vida das conexões ociosas
	sqlDB.SetConnMaxLifetime(time.Hour)

	if e != nil {
		panic(e)
	}
	database.AutoMigrate(&model.Model{}, &model.UserModel{})
	clearData()
}

func DB() *gorm.DB {
	return database
}

func clearData() {
	_ = database.Exec("DELETE FROM models").Error

	_ = database.Exec("DELETE FROM user_models").Error

	if err := database.Exec(`
        INSERT INTO user_models (id, balance, "limit") VALUES (1, 0, -100000);
        INSERT INTO user_models (id, balance, "limit") VALUES (2, 0, -80000);
        INSERT INTO user_models (id, balance, "limit") VALUES (3, 0, -1000000);
        INSERT INTO user_models (id, balance, "limit") VALUES (4, 0, -10000000);
        INSERT INTO user_models (id, balance, "limit") VALUES (5, 0, -500000);
    `).Error; err != nil {
		panic(err)
	}

	if err := database.Exec(`
		SET statement_timeout = 0;
		SET lock_timeout = 0;
		SET idle_in_transaction_session_timeout = 0;
		SET client_encoding = 'UTF8';
		SET standard_conforming_strings = on;
		SET check_function_bodies = false;
		SET xmloption = content;
		SET client_min_messages = warning;
		SET row_security = off;
		
		SET default_tablespace = '';
		
		SET default_table_access_method = heap;
    `).Error; err != nil {
		panic(err)
	}

	if err := database.Exec(`
		CREATE OR REPLACE FUNCTION createtransaction(
			IN idUser integer,
			IN value integer,
			IN description varchar(10)
		) RETURNS RECORD AS $$
		DECLARE
		userfound user_models%rowtype;
			ret RECORD;
		BEGIN
		SELECT * FROM user_models
			INTO userfound
		WHERE id = idUser;
		
		IF not found THEN
				--raise notice'Id user % not found.', idUser;
		select -1 into ret;
		RETURN ret;
		END IF;
		
			--raise notice'transaction by user %.', idUser;
		INSERT INTO models (value, description, created_at, user_id)
		VALUES (value, description, now() at time zone 'utc', idUser);
		UPDATE user_models
		SET balance = balance + value
		WHERE id = idUser AND (value > 0 OR balance + value >= "limit")
			RETURNING balance, "limit"
		INTO ret;
		raise notice'Ret: %', ret;
			IF ret."limit" is NULL THEN
				--raise notice'Id  user % not found.', idUser;
		select -2 into ret;
		END IF;
		RETURN ret;
		END;$$ LANGUAGE plpgsql;

    `).Error; err != nil {
		panic(err)
	}
}
