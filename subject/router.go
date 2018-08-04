package subject

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/qianlidongfeng/log"
)

type Router struct{
}

func NewRouter() Router{
	return Router{}
}

func (this Router)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path[len(path)-1] != '/'{
		path += "/"
	}
	if this.isLogin(w,r){
		switch path {
		case "/favicon.ico/":
			return
		case "/login/","/loginservice/":
			http.Redirect(w,r,"/",http.StatusFound)
		case "/":
			indexHandler(w,r)
		case "/failed/":
			failed(w,r)
		case "/update/":
			update(w,r)
		case "/add/":
			add(w,r)
		case "/del/":
			del(w,r)
		case "/subject/server/":
			serverSubjects(w,r)
		case "/subject/client/":
			clientSubjects(w,r)
		default:
			w.Write([]byte("403 forbidden"))
		}
	}else{
		switch path {
		case "/favicon.ico/":
			return
		case "/login/":
			login(w,r)
			return
		case "/loginservice/":
			loginservice(w,r)
			return
		case "/failed/":
			failed(w,r)
			return
		case "/subject/server/":
			serverSubjects(w,r)
		case "/subject/client/":
			clientSubjects(w,r)
		default:
			http.Redirect(w,r,"/login",http.StatusFound)
		}
	}
}

func (this Router) isLogin(w http.ResponseWriter, r *http.Request) bool{
	cookie, err := r.Cookie("username")
	if err != nil{
		return false
	}
	username := cookie.Value
	cookie, err = r.Cookie("password")
	if err != nil{
		return false
	}

	password := cookie.Value
	data, err := ioutil.ReadFile(dir+"/subject_conf/user.json")
	users := make(map[string]interface{})
	err = json.Unmarshal(data, &users)
	if err != nil{
		log.Warn(err)
		http.Redirect(w, r, "/failed?err="+err.Error(), http.StatusFound)
		return false
	}
	if pwd,ok:= users[username];ok{
		if pwd.(string) == password{
			return true
		}
	}
	return false
}
