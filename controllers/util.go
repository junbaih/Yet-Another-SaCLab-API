package controllers

import (
	"regexp"
	_ "errors"
)

func fileValidate(f map[string]interface{}) bool {
	key := []string{"studyid", "filename", "fileformat", "filetype", "versionid"}
	for _, k := range key {
		if _, prs := f[k]; !prs {
			return false
		}
	}
	return true
}

func expirePreviousVersion(f map[string]interface{}) {
	return
}


func urlContainsID(s string) (bool,error){
	 return regexp.MatchString(`(?i)^[a-f\d]{24}$`, s)
}


