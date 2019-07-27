package controllers

func fileValidate(f *map[string]interface{}) bool {
	key := []string{"studyid", "filename", "fileformat", "filetype", "versionid"}
	for _, k := range key {
		if _, prs := (*f)[k]; !prs {
			return false
		}
	}
	return true
}

func expirePreviousVersion(f *map[string]interface{}) {
	return
}
