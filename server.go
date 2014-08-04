package main

import (
    "log"
    "net/http"
    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/render"
    "gowebserver/sessionauth"
)

func main() {
    m := martini.Classic()

    store := sessions.NewCookieStore([]byte("secret123"))
    store.Options(sessions.Options{
        MaxAge: 0,
    })
    m.Use(sessions.Sessions("my_session", store))
    m.Use(render.Renderer())
    m.Use(sessionauth.Sessionsauth())

    m.Get("/", func() string {
        return "Hello world!"
    })

    m.Get("/users/:userId", sessionauth.LoginRequired(), func(params martini.Params,
            sess sessions.Session, u sessionauth.User, res http.ResponseWriter,
            r render.Render) {
        us := u.Get("Username")
        var username string = ""
        if us.(string) != "" {
            username = us.(string)
        }
        data := map[string]interface{}{
          "msg": "Hello, " + username + "!",
        }
        // Set Header
        res.Header().Add("X-Powered-By", "goweb")
        r.JSON(200, data)
    })

    m.Get("/login", func(r render.Render, u sessionauth.User) {
        if (u.IsAuthenticated()){
            r.Redirect("/")
            return
        }
        page, _ := r.Template().Parse(`{{define "login"}}
                            <h1>Login Page</h1>
                                <form action="/login" method="POST">
                                    <input id="user" name="user" type="text">
                                    <input id="password" name="password" type="password">
                                    <input type="submit" id="submit" value="login">
                                </form>
                            {{end}}`)
        r.HTML(200, "login", page)
    })

    m.Post("/login", binding.Form(sessionauth.PostedUser{}),
        func(sess sessions.Session, posted sessionauth.PostedUser,
              r render.Render, res http.ResponseWriter, req *http.Request) {
        log.Println("user: ", posted.Username)
        log.Println("password: ", posted.Password)
        if posted.Username != "nick" || posted.Password != `1234` {
            r.Redirect("/login")
        } else {
            u := &sessionauth.PostedUser{
                Username: posted.Username,
                Password: posted.Password,
                Authenticated: true,
            }
            log.Println(u)
            sessionauth.Login(sess, u)

            log.Println("LOGIN Success! ", u)
            r.Redirect("/users/" + posted.Username)
            return
        }
    })

    m.Get("/logout", func(sess sessions.Session) {
        sess.Clear()
    })

    http.ListenAndServe(":3333", m)
}
