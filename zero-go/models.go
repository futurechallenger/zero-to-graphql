package main

import "time"

/*
CREATE TABLE IF NOT EXISTS "person" (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"password" varchar(128) NOT NULL,
	"last_login" datetime NULL,
	"is_superuser" bool NOT NULL,
	"username" varchar(30) NOT NULL UNIQUE,
	"first_name" varchar(30) NOT NULL,
	"last_name" varchar(30) NOT NULL,
	"email" varchar(254) NOT NULL,
	"is_staff" bool NOT NULL,
	"is_staff" bool NOT NULL,
	"date_joined" datetime NOT NULL);
*/
type Person struct {
	ID          int64      `json:"id"`
	Password    string     `json:"password"`
	LastLogin   *time.Time `json:"last_login"`
	IsSuperuser bool       `json:"is_superuser"`
	Username    string     `json:"username"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	IsStaff     bool       `json:"is_staff"`
	IsActive    bool       `json:"is_active"`
	DateJoined  *time.Time `json:"date_joined"`
}
