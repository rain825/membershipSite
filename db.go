package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type user struct {
	ID       int64  `db:"id"`
	UserID   string `db:"userID"`
	UserName string `db:"userName"`
	Password string `db:"password"`
}

func (db DB) FetchUsers() ([]*user, error) {
	var users []*user
	q := "SELECT `id`, `userID`,`userName`,`password` FROM `user`"
	if err := db.Select(&users, q); err != nil {
		fmt.Errorf("FetchUsers Error: %v\n", err)
		return nil, err
	}
	//for debug
	for _, u := range users {
		fmt.Printf("%s : %s\n", u.UserID, u.Password)
	}
	return users, nil
}

func (db DB) FetchUserByID(id string) (user, error) {
	var user user
	q := "SELECT `id`, `userID`,`userName`,`password` FROM `user` WHERE `userID` = ?"
	if err := db.Get(&user, q, id); err != nil {
		fmt.Errorf("FetchUserByID Error: %v\n", err)
		return user, err
	}
	return user, nil
}

func (db DB) UpdateUserNameByID(id string, name string) {
	q := "UPDATE `user` SET `userName` = ? WHERE `id` = ?"
	res, err := db.Exec(q, name, id)
	if err != nil {
		log.Fatalf("UpdateUserNameByID Error:%v\n", err)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		log.Fatalf("UpdateUserNameByID Error: %v\n", err)
	} else if rowsAffected > 0 {
		fmt.Printf("Successfully Updated. ID:%s\n", id)
	} else {
		fmt.Printf("%s is not exist\n", id)
	}
	return
}

func (db DB) InsertUser(UserID, userName, password string) (int64, error) {
	q := "INSERT INTO `user` (`UserID`,`userName`,`password`) VALUES (?,?,?) "
	res, err := db.Exec(q, UserID, userName, password)
	if err != nil {
		log.Fatalf("Insert UserByID Error:%v\n", err)
		return 0, err
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Insert userByID Error: %v\n", err)
	}
	fmt.Printf("Successfully Inserted. ID is %d\n", lastInsertedID)
	return int64(lastInsertedID), nil
}
