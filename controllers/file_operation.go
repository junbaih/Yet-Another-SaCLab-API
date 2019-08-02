package controllers

import (
	//"controllers/parser"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mednicklab/models"
	"net/http"
	"net/url"
	"os"
)

var filedir string = ""

/*
 To be consistent with existing server&client apps implementation, 
 request body uses a nested data form.  ie => request.data = {"data":json_serialized_string<actual_data>}
 This allows server to get query by directly parsing url without reading Request.Body 
*/
func CreateFiles(db models.DatabaseAccess, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/files/" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	r.ParseMultipartForm(32 << 20)
	fmt.Println("create Files:")
	fmt.Println(r.Header)
	fmt.Println(r.Form)
	fmt.Println(r.MultipartForm)

	var t map[string]interface{}

	fileinfo := r.FormValue("data")
	if fileinfo == "" {
		panic("file info is empty")
	}
	log.Println(fileinfo)

	err := json.Unmarshal([]byte(fileinfo), &t)
	if err != nil {
		panic(err)
	}

	if !fileValidate(t) {
		http.Error(w, "File info not complete, unable to process", 400)
		return
	}

	file, fhandler, err := r.FormFile("fileobj")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Println(fhandler.Filename, fhandler.Header)

	fname := filedir + fhandler.Filename
	// fname=generateTimeStampName(fname)

	// store the file
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(f, file)

	// err = prepareDocument(t)

	/*
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
	*/

	// err = expirePreviousVersion(t)

	fid, err := db.Insert(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	t["_id"] = fid

	fmt.Printf("%T,%+v", t, t)

	res, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func RetrieveFiles(db models.DatabaseAccess, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("retrieve Files:")
	fmt.Println(r.Form)
	fmt.Println(r.URL.RequestURI())
	
	// get file with fid ==> GET /files/<fid>
	if match,_ := urlContainsID(r.URL.RequestURI()[len("/files/"):]); match {
		query:=r.URL.RequestURI()[len("/files/"):]
		r.Form["_id"]=[]string{query}
		
	 }
	// get file by query ==> GET /files/?operand1=operator:operand2

	

	res,err:=db.Find(r.Form)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	ret,_:=res.([]models.FileInfo)
	
	fmt.Println(url.PathUnescape(r.URL.RequestURI()[len("/files/?"):]))
	//w.Write([]byte("file read"))
	http.ServeFile(w, r, ret[0].FilePath)
	
}

/*
 A example to process requests with json body 
*/
func UpdateFiles(db models.DatabaseAccess, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	//fmt.Println("Update Files:")
	fmt.Println(r.Form)
	var _id string
	if match,_ := urlContainsID(r.URL.RequestURI()[len("/files/"):]); match {
		_id =r.URL.RequestURI()[len("/files/"):]
		
	} else {
		http.Error(w, http.StatusText(400), 400)
		return
	}
		
	 
	//w.Write([]byte("file updated"))
	var t interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T,%+v", t, t)
	if err != nil {
		panic(err)
	}
	
	res,err:= db.Update(_id,t)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	
	ret, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func DeleteFiles(db models.DatabaseAccess, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("delete Files:")
	fmt.Println(r.Form)
	var _id string
	if match,_ := urlContainsID(r.URL.RequestURI()[len("/files/"):]); match {
		_id =r.URL.RequestURI()[len("/files/"):]
	} else {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	
	_,err:= db.Delete(_id)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.Write([]byte("file deleted"))
}
