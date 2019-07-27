package main

import (
	"log"
	"mednicklab/models"
	"mednicklab/routes"
	"net/http"
)

// Not used for the current implementation
type MyContext struct {
	models.DatabaseAccess
}

func main() {
	dbAddress := ""
	dbClient, err := models.NewDB(dbAddress)
	if err != nil {
		log.Panic(err)
	}
	http.Handle("/login", http.NotFoundHandler()) //TODO
	http.Handle("/", http.NotFoundHandler())      //TODO frontend connectivity

	// GET  	/files/12321fe12d4                      download the file whose fid=12321fe12d4
	// GET 		/files/?id=12321fe12d4&versionid=2 		get the file info whose fid=12321fe12d4 and versionid=2
	// POST 	/files/                                 upload a file with file info enclosed in request body
	// PUT     	/files/12321fe12d4						update the file info whose fid=12321fe12d4 with new info enclosed in request body
	// DELETE	/files/12321fe12d4						delete the file whose fid=12321fe12d4

	http.Handle("/files/", routes.FileHandler(dbClient))

	//http.Handle("/data/",routes.DataHandler(dbClient))
	http.ListenAndServe(":8080", nil)

}
