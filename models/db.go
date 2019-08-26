package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileInfo struct {
	Fid             string `json:"_id"`
	FileDir         string `json:"filedir"`
	FilePath        string `json:"filepath"`
	FileBase        string `json:"filebase"`
	FileName        string `json:"filename"`
	FileNameVersion string `json:"filename_version"`
	StudyID         string `json:"studyid"`
	VersionID       float64    `json:"versionid"`
	FileType        string `json:"filetype"`
	FileFormat      string `json:"fileformat"`
	SubjectID       float64    `json:"subjectid"`
	VisitID         float64    `json:"visitid"`
	Deleted         bool   `json:"deleted"`
	Parsed          bool   `json:"parsed"`
	Active          bool   `json:"active"`
	Expired         bool   `json:"expired"`
	DateModified    float64  `json:"datemodified"`
	DateExpired     float64  `json:"dateexpired"`
}

// kinda like a simple version of a db driver
type DatabaseAccess interface {
	Insert(interface{}) (interface{}, error)
	Find(interface{}) ([]interface{} , error)
	Update(interface{},interface{}) (interface{} , error)
	Delete(interface{}) (interface{},error)
}

type DatabaseClient struct {
	*mongo.Client
}

func (db *DatabaseClient) Insert(i interface{}) (interface{}, error) {
	res, err := db.Database("mednick").Collection("fileUploads").InsertOne(context.TODO(), i)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (db *DatabaseClient) Find(i interface{}) ([]interface{},error) {
	t,ok := i.(map[string][]string)
	if !ok{
		return nil,nil
	}
	if len(t)<1 {
		return nil,nil
	}
	filter:=bson.D{}
	
	// need a customized parser to process queries 
	// such as "id==xx OR subjectid==xx"
	// using url parsing could only solve AND (use & to leverage url parsing) and other binary operators ( for example /?id=lte:100000&topic=nin:[a,b,c] ) 
	// most likely to skip this part as I don't want to code another LL parser like the one I did in the compiler class
	/* 
	for k,v := range t {
		if _,prs:=map[string]bool{"_id":true,"studyid":true,"versionid":true,"subjectid":true,"visitid":true,"filetype":true,"fileformat":true,"filename":true}[k];prs {
			f,err:=buildBsonFromStrings(k,v)
			if err!=nil {
				return nil,err
			}
			filter=append(filter,f)
		}
	}
	*/
	res := []interface{}{}
	cur, err := db.Database("mednick").Collection("fileUploads").Find(context.TODO(), filter)
	if err != nil {
		return nil,err
	}
	// Close the cursor once finished
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		
		var elem FileInfo
		err := cur.Decode(&elem)
		if err != nil {
			return nil,err
		}

		res = append(res, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil,err
	}
	
	return res,nil
}

func (db *DatabaseClient) Update(i interface{}, t interface{}) (interface{},error) {
	id,ok := i.(string)
	if !ok {
		// to simplify the problem, return nil,nil if id is invalid, the returned error is reserved for internal problem
		return nil,nil 
	}
	filter := bson.D{{"_id",id}}
	var update interface{}
	res,err:= db.Database("mednick").Collection("fileUploads").UpdateOne(context.TODO(),filter,update)
	if err!=nil {
		return nil,err
	}
	return res.UpsertedID,nil
} 

func (db *DatabaseClient) Delete(i interface{}) (interface{},error) {
	id,ok := i.(string)
	if !ok {
		// to simplify the problem, return nil,nil if id is invalid, the returned error is reserved for internal problem
		return 0,nil 
	}
	var filter bson.D
	if id=="all" {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"_id",id}}
	}
	res,err:= db.Database("mednick").Collection("fileUploads").DeleteOne(context.TODO(),filter)
	if err!=nil {
		return 0,err
	}
	return res.DeletedCount,nil
}


func NewDB(dbAddress string) (*DatabaseClient, error) {
	clientOptions := options.Client().ApplyURI(dbAddress)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &DatabaseClient{client}, nil
}
