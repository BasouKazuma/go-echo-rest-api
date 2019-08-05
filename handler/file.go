package handler

import (
	"go-echo-rest-api/crypto"
	"go-echo-rest-api/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type CreateFileBody struct {
	UserId	int64	`json:"userId" db:"userId"`
	Name	string	`json:"name" db:"name"`
	Bytes	string	`json:"bytes" db:"bytes"`
}

func (h *Handler) CreateFile(c echo.Context) (err error) {
	// Bind
	fileBody := new(CreateFileBody)
	if err = c.Bind(fileBody); err != nil {
		message := model.ErrorMessage{Message: "Invalid body."}
		return c.JSON(http.StatusBadRequest, message)
	}
	// Params
	fileBody.UserId, err = strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		message := model.ErrorMessage{Message: "Invalid user id."}
		return c.JSON(http.StatusBadRequest, message)
	}

	// Check Connection
	database := h.DB
	err = database.Ping()
	if err != nil {
		message := model.ErrorMessage{Message: "An error occurred."}
		return c.JSON(http.StatusBadRequest, message)
	}

	// Create Hash
	fileHashData := crypto.FileHashData{UserId: fileBody.UserId, Name: fileBody.Name}
	hash := crypto.CreateFileHash(fileHashData)
	
	// Add File
	file := model.File{
		UserId: fileBody.UserId,
		Name: fileBody.Name,
		Hash: hash}
	addFileQuery := `INSERT INTO files ( user_id, name, hash ) 
 		VALUES( $1, $2, $3 ) 
 		RETURNING id`
	err = database.QueryRow(addFileQuery, file.UserId, file.Name, file.Hash).Scan(&file.Id)
	if err != nil {
		message := model.ErrorMessage{Message: err.Error()}
		return c.JSON(http.StatusBadRequest, message)
	}
	return c.JSON(http.StatusCreated, file)
}

func (h *Handler) GetFileList(c echo.Context) (err error) {
	// Params
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return
	}

	// Check Connection
	database := h.DB
	err = database.Ping()
	if err != nil {
		message := model.ErrorMessage{Message: "An error occurred."}
		return c.JSON(http.StatusBadRequest, message)
	}

	// Get list of files
	var fileList []model.File
	// fileList := model.File{}
	findFileListQuery := `SELECT * FROM files WHERE user_id = $1`
	rows, err := database.Query(findFileListQuery, userId)
	if err != nil {
		// handle this error
		panic(err)
	}
	for rows.Next() {
		file := model.File{}
		err = rows.Scan(&file.Id, &file.UserId, &file.Name, &file.Hash)
		if err != nil {
			// handle this error
			panic(err)
		}
		fileList = append(fileList, file)
	}
	if err != nil {
		message := model.ErrorMessage{Message: "File list not found."}
		return c.JSON(http.StatusNotFound, message)
	}
	return c.JSON(http.StatusOK, fileList)
}

// func (h *Handler) GetFile(c echo.Context) (err error) {
// 	//
// }
