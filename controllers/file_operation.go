package controllers

import (
	"os"
	"io"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
	"net/http"
	"mednicklab/models"
) 

var filedir string = ""

func CreateFiles(db models.DatabaseAccess, w http.ResponseWriter,r *http.Request) {
	if r.URL.Path!="/files/"{
		http.Error(w, http.StatusText(400), 400)
		return
	}
	r.ParseMultipartForm(32<<20)
	fmt.Println("create Files:")
	fmt.Println(r.Header)
	fmt.Println(r.Form)
	fmt.Println(r.MultipartForm)
	
	var t map[string]interface{}
    
	fileinfo := r.FormValue("data")
	if fileinfo=="" {
		panic("file info is empty")
	}
	log.Println(fileinfo)

	err := json.Unmarshal([]byte(fileinfo), &t)
	if err != nil {
			panic(err)
	}
	
	if !fileValidate(&t) {
		http.Error(w, "File info not complete, unable to process", 400)
		return
	}
	
	file, fhandler, err := r.FormFile("fileobj")
    if err != nil {
		panic(err)
    }
    defer file.Close()
	fmt.Println(fhandler.Filename,fhandler.Header)
	
	fname:=filedir+fhandler.Filename
	// fname=generateTimeStampName(fname)

	
	// store the file
    f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    io.Copy(f, file)
	
	// err = prepareDocument(&t)
	
	/*
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	*/
	
	// err = expirePreviousVersion(&t)
	
	fid,err := db.Insert(t)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	
	t["_id"] = fid

	
    fmt.Printf("%T,%+v",t,t)
	
	res, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(res)
	

}

func RetrieveFiles(db models.DatabaseAccess, w http.ResponseWriter,r *http.Request) {
	r.ParseForm()
	fmt.Println("retrieve Files:")
	fmt.Println(r.Form)
	//w.Write([]byte("file read"))
	http.ServeFile(w, r, "test.json")
}

func UpdateFiles(db models.DatabaseAccess, w http.ResponseWriter,r *http.Request) {
	r.ParseMultipartForm(32<<20)
	fmt.Println("Update Files:")
	fmt.Println(r.Form)
	w.Write([]byte("file updated"))
		var t interface{}
    body, err := ioutil.ReadAll(r.Body)
	if err!=nil {
		panic(err)
	}
	log.Println(string(body))
    err = json.Unmarshal(body, &t)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%T,%+v",t,t)
    if err != nil {
        panic(err)
    }
}

func DeleteFiles (db models.DatabaseAccess, w http.ResponseWriter,r *http.Request) {
	r.ParseForm()
	fmt.Println("delete Files:")
	fmt.Println(r.Form)
	w.Write([]byte("file deleted"))
}