package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/ltheinrich/gorum/internal/pkg/config"
	"github.com/ltheinrich/gorum/internal/pkg/db"
)

// UploadAvatar http handler function
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	var err error

	// security headers
	SecurityHeaders(w, r)

	// get username and password
	rawUsername, _ := base64.StdEncoding.DecodeString(r.FormValue("username"))
	username, password := string(rawUsername), r.FormValue("password")

	// check if provided and login is correct
	if username != "" && password != "" && login(username, password) {
		// get file
		var file multipart.File
		var header *multipart.FileHeader
		file, header, err = r.FormFile("avatar")

		// check if file provided
		if err != nil || file == nil || header == nil {
			// write header
			w.Header().Add("content-type", "text/html")
			w.WriteHeader(200)

			// write content
			w.Write([]byte(`<form method="post" enctype="multipart/form-data"><input name="avatar" type="file" size="50" accept="image/png"><button type="submit">Avatar Upload</button></form>`))
		} else {
			// check avatar size limit
			var avatarSizeLimit int
			avatarSizeLimit, _ = strconv.Atoi(config.Get("limit", "avatar"))
			if int(header.Size) > avatarSizeLimit {
				// write header
				w.Header().Add("content-type", "text/html")
				w.WriteHeader(200)

				// write content and return
				w.Write([]byte(`<h3>Avatar size limit exceeded</h3><form method="post" enctype="multipart/form-data"><input name="avatar" type="file" size="50" accept="image/png"><button type="submit">Avatar Upload</button></form>`))
				return
			}

			// read file
			fileData := make([]byte, header.Size)
			io.ReadAtLeast(file, fileData, int(header.Size))
			defer file.Close()

			// get user id
			var userID int
			err = db.DB.QueryRow("SELECT id from users WHERE username = $1;", username).Scan(&userID)
			if err != nil {
				// write header
				w.Header().Add("content-type", "text/html")
				w.WriteHeader(200)

				// write content
				w.Write([]byte(err.Error()))
				return
			}

			// open avatar file
			var avatar *os.File
			avatarName := fmt.Sprintf("%s/%v.png", config.Get("data", "avatar"), userID)
			avatar, err = os.OpenFile(avatarName, os.O_RDWR|os.O_CREATE, os.ModePerm)

			// create directories
			if os.IsNotExist(err) {
				err = os.MkdirAll(config.Get("data", "avatar"), os.ModePerm)
				if err != nil {
					// write header
					w.Header().Add("content-type", "text/html")
					w.WriteHeader(200)

					// write content
					w.Write([]byte(err.Error()))
					return
				}

				// open avatar file again
				avatar, err = os.OpenFile(avatarName, os.O_RDWR|os.O_CREATE, os.ModePerm)
				if err != nil {
					// write header
					w.Header().Add("content-type", "text/html")
					w.WriteHeader(200)

					// write content
					w.Write([]byte(err.Error()))
					return
				}
			} else if err != nil {
				// write header
				w.Header().Add("content-type", "text/html")
				w.WriteHeader(200)

				// write content
				w.Write([]byte(err.Error()))
				return
			}

			// write avatar file
			_, err := avatar.Write(fileData)
			if err != nil {
				// write header
				w.Header().Add("content-type", "text/html")
				w.WriteHeader(200)

				// write content
				w.Write([]byte(err.Error()))
				return
			}

			// redirect
			http.Redirect(w, r, "/edit-profile", http.StatusSeeOther)
		}
	}
}
