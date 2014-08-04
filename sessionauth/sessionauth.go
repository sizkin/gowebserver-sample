package sessionauth

import (
    "encoding/gob"
    "log"
    "reflect"
    "github.com/go-martini/martini"
    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/render"
)

type User interface {
    IsAuthenticated() bool

    Get(field string) interface{}
}

type PostedUser struct {
    Id    int64
    Username    string   `form:"user"`
    Password    string   `form:"password"`
    Authenticated    bool
}

func Sessionsauth() martini.Handler {
    gob.Register(&PostedUser{})
    return func(sess sessions.Session, c martini.Context) {
        sessuser := sess.Get("user")
        u := GenerateAnonymousUser()
        if sessuser != nil {
            u = sessuser.(*PostedUser)
                // log.Println(u.Id)
        }
        c.MapTo(u, (*User)(nil))
    }
}

func (u *PostedUser) Get(field string) interface{} {
    val := reflect.ValueOf(u).Elem()
    return val.FieldByName(field).Interface()
}

func (u *PostedUser) IsAuthenticated() bool {
    return u.Authenticated
}

func GenerateAnonymousUser() User {
    return &PostedUser{}
}

func Login(sess sessions.Session, user *PostedUser) error {
    user.Authenticated = true
    sess.Set("user", &user)
    log.Println("Login: ", &user)
    return nil
}

func LoginRequired() martini.Handler {
    return func(sess sessions.Session, r render.Render, c martini.Context) {
        sessuser := sess.Get("user")
        if sessuser == nil {
            r.Redirect("/login")
            return
        }
        c.Next()
    }
}
