package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"around/model"
	"around/service"

	jwt "github.com/form3tech-oss/jwt-go" // 前面是import重命名
)

var mySigningKey = []byte("secret")

// 传*更快一些，
// Response为啥不用* -> ResponseWriter定义
func signinHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")

    //  Get User information from client (from requst body)
    decoder := json.NewDecoder(r.Body) // 解析json
    var user model.User
    if err := decoder.Decode(&user); err != nil { // 为何传& -> 后面还要用
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

	// 判断是否出错
    exists, err := service.CheckUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

	// 判断是否存在该user
    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }

	// claims == payload 传乱码 -> https://jwt.io/
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // 这个不是加密token,是go struct
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(), //24小时过期
    })

    tokenString, err := token.SignedString(mySigningKey) // 这里才是加密
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }

    w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")

    decoder := json.NewDecoder(r.Body)
    var user model.User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

	// sanity check Username不为空 Password不为空 Username必须是小写字母或者数字
    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    success, err := service.AddUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }
    fmt.Printf("User added successfully: %s.\n", user.Username)
}