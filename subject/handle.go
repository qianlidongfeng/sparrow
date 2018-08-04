package subject

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/qianlidongfeng/log"
	"encoding/json"
	"io/ioutil"
	"time"
	"regexp"
)


func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err :=template.New("index").Parse(tmpIndex)
	if err != nil{
		log.Warn(err)
		http.Redirect(w, r, "/unknownerror", http.StatusFound)
		return
	}
	t.Execute(w, mSubjects)
}

func serverSubjects(w http.ResponseWriter, r *http.Request) {
	data,err:=json.Marshal(mSubjects)
	if err != nil{
		log.Warn(err)
		fmt.Fprint(w,err)
		return
	}
	fmt.Fprint(w,string(data))
}

func clientSubjects(w http.ResponseWriter, r *http.Request) {
	data,err:=json.Marshal(rmSubjects)
	if err != nil{
		log.Warn(err)
		fmt.Fprint(w,err)
		return
	}
	fmt.Fprint(w,string(data))
}

func login(w http.ResponseWriter, r *http.Request) {
	t, err :=template.New("login").Parse(tmpLogin)
	if err != nil{
		log.Warn(err)
		http.Redirect(w, r, "/unknownerror", http.StatusFound)
		return
	}
	t.Execute(w, struct{}{})
}

func loginservice(w http.ResponseWriter, r *http.Request){
	username := r.PostFormValue("name")
	password := r.PostFormValue("password")
	data, err := ioutil.ReadFile(dir+"/subject_conf/user.json")
	if err != nil {
		log.Warn(err)
		http.Redirect(w, r, "/failed?err="+err.Error(), http.StatusFound)
		return
	}
	users := make(map[string]interface{})
	err = json.Unmarshal(data, &users)
	if err != nil{
		log.Warn(err)
		http.Redirect(w, r, "/failed?err="+err.Error(), http.StatusFound)
		return
	}
	if pwd,ok:= users[username];ok{
		if pwd.(string) == password{
			cookie:=&http.Cookie{
				Name:   "username",
				Value:    username,
				Path:     "/",
				HttpOnly: false,
				MaxAge:   int(time.Hour * 2 / time.Second),
			}
			http.SetCookie(w, cookie)
			cookie=&http.Cookie{
				Name:   "password",
				Value:   password,
				Path:     "/",
				HttpOnly: false,
				MaxAge:   int(time.Hour * 2 / time.Second),
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/failed?err=username or password is invalid", http.StatusFound)
}

func failed(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	if len(r.Form["err"]) > 0 {
		fmt.Fprint(w, r.Form["err"][0])
	}else{
		fmt.Fprint(w, "unkown err,please check the log")
	}
}


func add(w http.ResponseWriter, r *http.Request){
	id := r.PostFormValue("id")
	subject := r.PostFormValue("subject")
	if _,ok:=mSubjects[id];ok{
		fmt.Fprint(w, "failed,the key is existing")
		return
	}
	if _,ok:=rmSubjects[subject];ok{
		fmt.Fprint(w, "failed,the subject is existing")
		return
	}
	reg, err := regexp.Compile("[1-9][0-9]*")
	if err!=nil{
		log.Warn(err)
		fmt.Fprint(w, err)
		return
	}
	s:=reg.FindAllString(id, -1)
	temp:=""
	for i:=0;i<len(s);i++{
		temp+=s[i]
	}
	if temp != id{
		fmt.Fprint(w, "id must be a number")
		return
	}
	mSubjects[id]=subject
	rmSubjects[subject]=id
	err=saveSubjects()
	if err != nil{
		log.Warn(err)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, "success")
}

func update(w http.ResponseWriter, r *http.Request){
	id := r.PostFormValue("id")
	reg, err := regexp.Compile("[1-9][0-9]*")
	if err!=nil{
		log.Warn(err)
		fmt.Fprint(w, err)
		return
	}
	s:=reg.FindAllString(id, -1)
	temp:=""
	for i:=0;i<len(s);i++{
		temp+=s[i]
	}
	if temp != id{
		fmt.Fprint(w, "id must be a number")
		return
	}
	subject := r.PostFormValue("subject")
	oldId:=r.PostFormValue("oldId")
	if _,ok:= mSubjects[oldId];!ok{
		fmt.Fprint(w, "failed,unkown the key")
		return
	}
	oldSubject:=mSubjects[oldId].(string)
	if _,ok:= mSubjects[id];ok &&id != oldId{
		fmt.Fprint(w, "failed,the key is existing")
		return
	}
	if _,ok:= rmSubjects[subject];ok&&subject!=oldSubject{
		fmt.Fprint(w, "failed,the subject is existing")
		return
	}
	delete(mSubjects,oldId)
	mSubjects[id]=subject
	delete(rmSubjects,oldSubject)
	rmSubjects[subject]=id
	err=saveSubjects()
	if err != nil{
		log.Warn(err)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, "success")
}

func del(w http.ResponseWriter, r *http.Request){
	id := r.PostFormValue("id")
	if _,ok:=mSubjects[id];!ok{
		fmt.Fprint(w,"failed,the key is not existing")
		return
	}
	subject:=mSubjects[id].(string)
	delete(mSubjects,id)
	delete(rmSubjects,subject)
	err:=saveSubjects()
	if err != nil{
		log.Warn(err)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, "success")
}