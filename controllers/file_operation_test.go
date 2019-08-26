package controllers

import (
	"bytes"
	"io/ioutil"
	"encoding/json"
	_ "errors"
	"reflect"
	"mednicklab/models"
	"testing"
	"net/http"
	"net/http/httptest"
	"net/url"
	"math/rand"
	"mime/multipart"
	"os"
    "time"
	"fmt"
	"log"
)

const charset = "abcdef" +"ABCDEF0123456789"

const IdLength = 6

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

func randID(length int) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

// modified from code samples on stackoverflow 
func SetField(obj interface{}, name string, value interface{}) error {
    structValue := reflect.ValueOf(obj).Elem()
    structFieldValue := structValue.FieldByName(name)

    if !structFieldValue.IsValid() {
        return fmt.Errorf("No such field: %s in obj", name)
    }

    if !structFieldValue.CanSet() {
        return fmt.Errorf("Cannot set %s field value", name)
    }

    structFieldType := structFieldValue.Type()
    val := reflect.ValueOf(value)
    if structFieldType != val.Type() {

        return fmt.Errorf("Provided value %v type %v didn't match obj field type %v",val,val.Type(),structFieldType)
    }

    structFieldValue.Set(val)
    return nil
}


func GenFileInfo(m map[string]interface{}) (models.FileInfo,error) {
    var f models.FileInfo
	for k, v := range m {
		_n :=getFieldName(k,f)
        err := SetField(&f, _n, v)
        if err != nil {
            log.Fatal(err)
        }
    }
    return f,nil
}

func getFieldName(tag string, s interface{}) (fieldname string) {
    rt := reflect.TypeOf(s)
    if rt.Kind() != reflect.Struct {
        panic("bad type")
    }
    for i := 0; i < rt.NumField(); i++ {
        f := rt.Field(i)
        v := f.Tag.Get("json") // use split to ignore tag "options" like omitempty, etc.
        if v == tag {
            return f.Name
        }
    }
    return ""
}


type mockDB struct {
	Files map[string]models.FileInfo
}

func (db *mockDB) Insert(i interface{}) (interface{}, error){
	t,_:= i.(map[string]interface{})
	fmt.Println("start genFile")
	f,err:=GenFileInfo(t)
	if err!=nil {
		return nil,err
	}
	//f.Fid=randID(IdLength)
	fmt.Println("db insertion complete!")
	db.Files[f.Fid]=f
	return f.Fid,nil
}

func (db *mockDB) Find(i interface{}) ([]interface{} , error){
	log.Printf("Find parameter %#v \n",i)
	m,prs:= i.(url.Values)
	if !prs {
		fmt.Println("cannot convert to map string interface")
	}
	log.Println("find query",i,m)
	
	id:= map[string][]string(m)["_id"][0]
	fmt.Println(id)
	/*
	id,_:=ids.(string)
	_,p:=db.Files[id]
	if !p {
		return nil,nil
	}
	*/
	return []interface{}{db.Files[id]},nil
}

func (db *mockDB) Update(i interface{},t interface{}) (interface{} , error) {
	/*
	m,_:=t.(map[string]interface{})
	id,_:=i.(string)
	for k,v := range m {
		switch v.(type){
			case bool:
				// use reflect to retrieve field name through tag name
		}
	}*/
	return nil,nil
}

func (db *mockDB) Delete(i interface{}) (interface{},error) {
	id,_ := i.(string)
	_,p:=db.Files[id]
	if !p {
		return 0,nil
	}
	delete(db.Files,id)
	return 1,nil
}

var db mockDB
func init() {
	db.Files=make(map[string]models.FileInfo)
}


/*
type FileInfo struct {
	Fid             string `json:"_id"`
	FileDir         string `json:"filedir"`
	FilePath        string `json:"filepath"`
	FileBase        string `json:"filebase"`
	FileName        string `json:"filename"`
	FileNameVersion string `json:"filename_version"`
	StudyID         string `json:"studyid"`
	VersionID       int    `json:"versionid"`
	FileType        string `json:"filetype"`
	FileFormat      string `json:"fileformat"`
	SubjectID       int    `json:"subjectid"`
	VisitID         int    `json:"visitid"`
	Deleted         bool   `json:"deleted"`
	Parsed          bool   `json:"parsed"`
	Active          bool   `json:"active"`
	Expired         bool   `json:"expired"`
	DateModified    int64  `json:"datemodified"`
	DateExpired     int64  `json:"dateexpired"`
}
*/
// It is more convenient to post data through request body.
// However, since the existing python client uses url-encoded data
// I will try to follow that convention here

func TestUploadFiles(t *testing.T) {

    	f1:=models.FileInfo{"f11f22f33f44f55f66f77f88","fdir","fpath.json","fbase","fname","fnv","studyid",1,"ftype","fformat",1,1,false,true,false,true,10000000,10000000}
		js,_:=json.Marshal(f1)
		
	//req:=httptest.NewRequest("POST","/files/?data="+url.PathEscape(string(js)),nil)
	req:=newfileUploadRequest("/files/?data="+url.PathEscape(string(js)),"fileobj")

	res:=httptest.NewRecorder()
	CreateFiles(&db,res,req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusOK)
	}
	
	
	var ti map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	err = json.Unmarshal(body, &ti)
	if err != nil {
		panic(err)
	}
	if ti["_id"]!="f11f22f33f44f55f66f77f88" {
		t.Errorf("handler returned wrong fid: got %v want %v",ti["_id"], "fid")
	}
	
	dat,err := ioutil.ReadFile("fpath.json")
	if err!=nil {
		//log.Fatal(err)
	}
	exp, err := ioutil.ReadFile("fpath_compare.json")
	if !bytes.Equal(dat,exp) {
		t.Errorf("file uploaded is corrupted. Expecting %s, but get %s ",string(exp),string(dat))
	}
	
}

func TestRetrieveFiles_DownloadingFile(t *testing.T){
	req:=httptest.NewRequest("GET","/files/f11f22f33f44f55f66f77f88",nil)
	res:=httptest.NewRecorder()
	RetrieveFiles(&db,res,req)
	
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	
	//var ti models.FileInfo
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	
	/*
	err = json.Unmarshal(body, &ti)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v \n",ti)
	*/
	//f1:=models.FileInfo{"f11f22f33f44f55f66f77f88","fdir","fpath.json","fbase","fname","fnv","studyid",1,"ftype","fformat",1,1,false,true,false,true,10000000,10000000}
	
	dat, err := ioutil.ReadFile("fpath.json")
	if !bytes.Equal(dat,body) {
		t.Errorf("handler returned wrong file content")
	}
	
	/*
	if ti!=f1 {
			t.Errorf("handler returned wrong fid: got %v want %v",ti, f1)

	}
	*/
}


// upload formfile, copied from this post 
// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84

func newfileUploadRequest(uri string,  paramName string) *http.Request {
	
	file, err := os.Open("fpath_compare.json")
	if err != nil {
		log.Fatal(err)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	
	/*fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	*/
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, "fpath.json")
	if err != nil {
		log.Fatal(err)
	}
	part.Write(fileContents)

	/*
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	*/
	
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	req:=httptest.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	log.Printf("created http request:\n %#v \n",req)
	return req
}