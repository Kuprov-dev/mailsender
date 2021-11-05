package mailsender

import (
	"fmt"
	"log"
	"net/http"

	"simplemailsender/pkg/db"
)

func SenderHandler(dao *db.MongoDBUserDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err, templates := dao.PrintTemplates()
		fmt.Println(err)
		if err != nil {
			log.Println(err)

			return
		}

		fmt.Println(templates)
	}

}
