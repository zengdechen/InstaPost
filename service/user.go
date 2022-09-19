package service

import (
    "fmt"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"

)


func CheckUser(username, password string) (bool, error) {
	// 去数据库读username,password
	// select * from post where user = xxx && password = xxx
    query := elastic.NewBoolQuery() 
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))

    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    var utype model.User
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
        u := item.(model.User)
        if u.Password == password {
            fmt.Printf("Login as %s\n", username)
            return true, nil
        }
    }
    return false, nil
}

func AddUser(user *model.User) (bool, error) {
    query := elastic.NewTermQuery("username", user.Username)
	// 找找这个username是否已经存在
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

	// 找到了也不能创建(已经有这个username了)
    if searchResult.TotalHits() > 0 {
        return false, nil
    }

	// 加到ES
    err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }

	// 加入成功
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}

