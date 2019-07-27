package routes

import (
	"fmt"
	"mednicklab/controllers"
	"mednicklab/models"
	"net/http"
)

func FileHandler(db models.DatabaseAccess) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL, r.Method)
		switch r.Method {
		case "POST":
			controllers.CreateFiles(db, w, r)
		case "GET":
			controllers.RetrieveFiles(db, w, r)
		case "PUT":
			controllers.UpdateFiles(db, w, r)
		case "DELETE":
			controllers.DeleteFiles(db, w, r)
		default:
			http.Error(w, http.StatusText(405), 405)
			return
		}
	})
}

/*
func DataHandler( db models.DatabaseAccess ) http.Handler {
	return http.HandleFunc(
		func(w http.ResponseWriter, r *http.Request){
			switch r.Method {
			case "POST":
				controllers.CreateData(db,w,r)
			case "GET":
				controllers.RetrieveData(db,w,r)
			case "PUT":
				controllers.UpdateData(db,w,r)
			case "DELETE":
				controllers.DeleteData(db,w,r)
			default:
				http.Error(w, http.StatusText(405), 405)
				return
			}
		}
	)
}
*/
