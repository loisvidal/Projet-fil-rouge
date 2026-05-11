package controllers

import (
    "encoding/json"
    "os"
    "log"
    "RedProject/models" 
)

const filePath = "Data/userList.json"
var CheckUser models.User

func getUsers() []models.User {
    file, err := os.ReadFile(filePath)
    if err != nil {
        return []models.User{}
    }
    var users []models.User
    json.Unmarshal(file, &users)
    return users
}

func CheckUserConnect(user string, pwd string) bool {
    users := getUsers()
    for _, u := range users {
        if u.NameUser == user && u.Password == pwd {
			CheckUser = u
			CheckUser.IsConnect = true
            return true
        }
    }
    return false
}

func WriteUserConnect(newUser models.User) {
    users := getUsers()
    newUser.IdUser = len(users) + 1
    users = append(users, newUser)

    data, err := json.MarshalIndent(users, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    os.WriteFile(filePath, data, 0644)
}