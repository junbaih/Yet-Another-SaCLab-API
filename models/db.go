package models

import (
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type FileInfo struct {
	Fid          	string  `json:"_id"`
	FileDir      	string  `json:"filedir"`
	FilePath     	string  `json:"filepath"`
	FileBase	 	string	`json:"filebase"`
	FileName		string	`json:"filename"`
	FileNameVersion	string	`json:"filename_version"`
	StudyID			string	`json:"studyid"`
	VersionID		int		`json:"versionid"`
	FileType		string	`json:"filetype"`
	FileFormat   	string	`json:"fileformat"`
	SubjectID		int		`json:"subjectid"`
	VisitID			int		`json:"visitid"`
	Deleted			bool	`json:"deleted"`
	Parsed			bool	`json:"parsed"`
	Active			bool	`json:"active"`
	Expired			bool	`json:"expired"`
	DateModified	int64	`json:"datemodified"`
	DateExpired		int64	`json:"dateexpired"`
} 

type DatabaseAccess interface{
	Insert(interface{}) (interface{}, error)
	//Find(interface{}) (interface{} , error)
	//Update(interface{}) (interface{} , error)
	//Delete(interface{}) error
}

type DatabaseClient struct{
	*mongo.Client
}

func (db *DatabaseClient) Insert(i interface{}) (interface{},error) {
	res,err:=db.Database("mednick").Collection("fileUploads").InsertOne(context.TODO(),i)
	if err!=nil {
		return nil,err
	}
	return res.InsertedID,nil
}


func NewDB( dbAddress string ) (*DatabaseClient, error) {
	clientOptions := options.Client().ApplyURI(dbAddress)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
		
	if err != nil {
		return nil,err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	
	if err != nil {
		return nil,err
	}

	fmt.Println("Connected to MongoDB!")
	return &DatabaseClient{client},nil
}