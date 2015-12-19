package main

import (
	"fmt"
	"github.com/godwhoa/login/crypt"
	"github.com/godwhoa/login/store"
	"github.com/godwhoa/login/upload"
	"github.com/gorilla/sessions"
	"github.com/imdario/mergo"
	"github.com/kennygrant/sanitize"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var db store.SqliteDB
var cstore = sessions.NewCookieStore([]byte("s1!!"))

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := cstore.Get(r, "login")

	if r.Method == "GET" {

		if session.Values["user"] != nil {
			http.Redirect(w, r, "/profile", http.StatusFound)
		} else {
			http.ServeFile(w, r, "./views/index.html")
		}
	}
	if r.Method == "POST" {
		user := sanitize.HTML(r.FormValue("user"))
		pass := sanitize.HTML(r.FormValue("pass"))
		Pass := store.QueryPass(db, user)

		//db has hashed pass and we have raw pass
		//store lib returns "1" if user doesn't exist
		if Pass != "1" && crypt.Check(Pass, pass) {
			session.Values["user"] = user
			session.Save(r, w)
			log.Printf("%s logged in\n", user)
			http.Redirect(w, r, "/profile", http.StatusFound)
		} else {
			io.WriteString(w, "neg")
		}
	}
}

func LoginOut(w http.ResponseWriter, r *http.Request) {
	session, _ := cstore.Get(r, "login")
	if session.Values["user"] != nil {
		log.Printf("%s logged out\n", session.Values["user"].(string))
		session.Values["user"] = nil
		session.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusFound)

	} else {
		fmt.Fprintf(w, "Not logged in")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {

	session, _ := cstore.Get(r, "login")
	user := session.Values["user"]
	ustring := user.(string)
	if user != nil {
		u := store.QueryProfile(db, ustring)

		t := template.Must(template.ParseFiles("./views/pro_index.html"))
		err := t.Execute(w, u)
		if err != nil {
			log.Printf("template execution: %s\n", err)
		}
		log.Printf("%s visted profile\n", ustring)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {

	session, _ := cstore.Get(r, "login")
	cuser := session.Values["user"].(string)
	u := store.QueryProfile(db, cuser)

	if r.Method == "GET" {
		if session.Values["user"] == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			t := template.Must(template.ParseFiles("./views/edit.html"))
			err := t.Execute(w, u)
			if err != nil {
				log.Printf("template execution: %s\n", err)
			}
		}

	}
	if r.Method == "POST" {

		user := sanitize.HTML(r.FormValue("user"))
		pass := sanitize.HTML(r.FormValue("pass"))
		about := sanitize.HTML(r.FormValue("about"))
		log.Printf("Updated user: %s pass: %s about: %s pic: %s\n", user, pass, about)
		if u.User != "" { //Make sure user exists
			//Handle file verification
			r := upload.Upload(r)
			//Delete old file
			//Ignore if same as old one
			if r.Filename != u.Pic && r.Filename != "" {
				_ = os.Remove(u.Pic)
			}
			if r.Res == "pos" {
				//Hash password
				pass = crypt.Hash(pass)

				//Save it all to db now
				upuser := store.Users{
					User:  user,
					Pass:  pass,
					About: about,
					Pic:   r.Filename,
				}
				if r.Filename == "" {
					upuser.Pic = u.Pic
				}
				mergo.Merge(&upuser, u) //just so we only update whats need to updated.
				store.UpdateUser(db, upuser)
				log.Printf("Updated: %+v\n", upuser)
				//Give client the ok sign
				io.WriteString(w, "pos")
			} else {
				io.WriteString(w, "neg")
			}
		} else { //If user don't exist
			io.WriteString(w, "neg")
		}

	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	//muh if soup
	if r.Method == "GET" {
		http.ServeFile(w, r, "./views/reg.html")
	}
	if r.Method == "POST" {
		session, _ := cstore.Get(r, "login")
		user := sanitize.HTML(r.FormValue("user"))
		pass := sanitize.HTML(r.FormValue("pass"))
		about := sanitize.HTML(r.FormValue("about"))

		u := store.QueryProfile(db, user)
		if u.User == "" {
			//Handle file verification
			//re = Result
			re := upload.Upload(r)
			if re.Res == "pos" {
				//Hash password
				pass = crypt.Hash(pass)
				//Save it all to db now
				newuser := store.Users{User: user, Pass: pass, About: about, Pic: re.Filename}
				store.AddUser(db, newuser)
				log.Printf("Added user: %s pass: %s about: %s\n", user, pass, about)
				//Set session
				session.Values["user"] = user
				session.Save(r, w)
				//Give client the ok sign
				io.WriteString(w, "pos")
			} else {
				io.WriteString(w, "neg")
			}
		} else { //Non empty struct == user already exists
			io.WriteString(w, "neg")
		}

	}
}

func main() {
	db.Init("store.db")
	fs := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	}
	http.HandleFunc("/public/", fs)

	http.HandleFunc("/", Login)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", LoginOut)

	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/register", Register)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
