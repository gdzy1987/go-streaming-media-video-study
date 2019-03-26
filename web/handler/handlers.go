package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/common"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"go-streaming-media-video-study/web/client"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func HomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err != nil || err2 != nil {
		p := &HomePage{
			Name: "张三",
		}
		t, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			logger.Info("parsing template home.html err:\t", err)
			return
		}
		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
	return
}

func UserHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fname := r.FormValue("username")

	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{
			Name: cname.Value,
		}
	} else if len(fname) != 0 {
		p = &UserPage{
			Name: fname,
		}
	}

	t, err := template.ParseFiles("./templates/userhome.html")
	if err != nil {
		logger.Info("parsing userhome.html error:\t", err)
	}
	t.Execute(w, p)
}

func ApiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(common.ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &common.ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		re, _ := json.Marshal(common.ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	client.Request(apibody, w, r)
	defer r.Body.Close()
}

func ProxyVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse("http://" + config.DefaultConfig.Address + ":9000/")
	logger.Info("请求VideoHandler:\t", u)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func ProxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse("http://" + config.DefaultConfig.Address + ":9000/")
	logger.Info("请求UploadHandler:\t", u)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
