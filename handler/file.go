package handler

import (
	"go-echo-rest-api/crypto"
	"go-echo-rest-api/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type CreateFileRequest struct {
	UserId	int64	`json:"userId" db:"userId"`
	Name	string	`json:"name" db:"name"`
	Bytes	[]byte	`json:"bytes" db:"bytes"`
}

type GetFileResponse struct {
	Id		int64	`json:"id" db:"id"`
	UserId	int64	`json:"userId" db:"userId"`
	Name	string	`json:"name" db:"name"`
	Hash	string	`json:"hash" db:"hash"`
	Bytes	[]byte	`json:"bytes" db:"bytes"`
}

func (h *Handler) CreateFile(c echo.Context) (err error) {
	// Bind
	fileBody := new(CreateFileRequest)
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
	fileHashData := crypto.FileHashData{
		UserId: fileBody.UserId,
		Name: fileBody.Name,
		Bytes: fileBody.Bytes}
	hash := crypto.CreateFileHash(fileHashData)

	// Check Hash
	var count int
	checkFileQuery := `SELECT COUNT(*) FROM files WHERE hash = $1`
	err = database.QueryRow(checkFileQuery, hash).Scan(&count)
	if count > 0 {
		message := model.ErrorMessage{Message: "File was already uploaded."}
		return c.JSON(http.StatusConflict, message)
	}

	// Upload File
	err = model.UploadFileToS3(hash, fileBody.Bytes)
	if err != nil {
		message := model.ErrorMessage{Message: "Upload failed."}
		return c.JSON(http.StatusBadRequest, message)
	}
	
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

func (h *Handler) GetFile(c echo.Context) (err error) {
	// Params
	userId := c.Param("user_id")
	fileHash := c.Param("file_hash")

	// Check Connection
	database := h.DB
	err = database.Ping()
	if err != nil {
		message := model.ErrorMessage{Message: "An error occurred."}
		return c.JSON(http.StatusBadRequest, message)
	}

	// Get file
	file := model.File{}
	findFileQuery := `SELECT * FROM files WHERE hash = $1 AND user_id = $2`
	err = database.QueryRow(findFileQuery, fileHash, userId).Scan(&file.Id, &file.UserId, &file.Name, &file.Hash)
	if err != nil {
		message := model.ErrorMessage{Message: "File not found."}
		return c.JSON(http.StatusNotFound, message)
	}
	fileBytes, err := model.GetBytesOfFileFromS3(file.Hash)
	if err != nil {
		message := model.ErrorMessage{Message: "File not downloaded."}
		return c.JSON(http.StatusNotFound, message)
	}

	// Send Response
	fileResponse := GetFileResponse{
		Id: file.Id,
		UserId: file.UserId,
		Name: file.Name,
		Hash: file.Hash,
		Bytes: fileBytes}
	return c.JSON(http.StatusOK, fileResponse)
}
