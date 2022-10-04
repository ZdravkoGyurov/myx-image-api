package handlers

import (
	"net/http"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/api/response"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/controller"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/errors"
)

const maxUploadSize = 10 * 1024 * 1024 // 10MB

type Image struct {
	Controller controller.Controller
}

func (h *Image) Upload(writer http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(writer, request.Body, maxUploadSize)
	if err := request.ParseMultipartForm(maxUploadSize); err != nil {
		err = errors.Newf("failed to upload file, file is too big: %w", err)
		response.SendError(writer, request, err)
		return
	}

	file, fileHandler, err := request.FormFile("file")
	if err != nil {
		err = errors.Newf("failed to retrieve file: %w", err)
		response.SendError(writer, request, err)
		return
	}
	defer file.Close()

	if fileHandler.Header.Get("Content-Type") != "image/jpeg" {
		err := errors.Newf("failed to upload file, only jpeg images are allowed: %w", errors.ErrInvalidEntity)
		response.SendError(writer, request, err)
		return

	}

	image, err := h.Controller.UploadImage(request.Context(), fileHandler.Filename, file)
	if err != nil {
		err = errors.Newf("failed to upload file: %w", err)
		response.SendError(writer, request, err)
		return
	}

	response.SendData(writer, request, http.StatusOK, image)
}

func (h *Image) GetAll(writer http.ResponseWriter, request *http.Request) {
	bbox := request.URL.Query().Get("bbox")
	if bbox == "" {
		err := errors.Newf("failed to get images, bbox query parameter is required: %w", errors.ErrInvalidEntity)
		response.SendError(writer, request, err)
		return
	}

	images, err := h.Controller.GetImages(request.Context(), bbox)
	if err != nil {
		response.SendError(writer, request, err)
		return
	}

	response.SendData(writer, request, http.StatusOK, images)
}

func (h *Image) Delete(writer http.ResponseWriter, request *http.Request) {
	imageName := pathParam(request, "name")
	if err := h.Controller.DeleteImage(request.Context(), imageName); err != nil {
		err = errors.Newf("failed to delete file: %w", err)
		response.SendError(writer, request, err)
		return
	}

	response.SendData(writer, request, http.StatusNoContent, struct{}{})
}
