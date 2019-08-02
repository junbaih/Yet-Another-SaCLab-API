package models

import (
	_ "strings"
	_ "go.mongodb.org/mongo-driver/bson"
)

type Filter struct{
	operator string
	operand string
}

/*
func buildBsonFromStrings(f string, s []string) (bson.E,error) {
	for _,v:= range s {
		temp:= strings.Split(v,":")'
		if len(temp)==1
	}
	
}
*/