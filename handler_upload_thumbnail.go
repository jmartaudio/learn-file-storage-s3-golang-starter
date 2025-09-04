package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse form", err)
		return
	}
	f, h, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could get data", err)
		return
	}
	mediaType := h.Header.Get("Content-Type")

	data, err := io.ReadAll(f)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not read file", err)
		return
	}
	meta, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find file in db", err)
		return
	}
	if userID != meta.UserID {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}
	thumb := thumbnail{
		data:      data,
		mediaType: mediaType,
	}

	videoThumbnails[videoID] = thumb
	url := "http://localhost:8091/api/thumbnails/" + videoIDString
	fmt.Println(url)

	err = cfg.db.UpdateVideo(database.Video{
		ID:           videoID,
		ThumbnailURL: &url,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find file in db", err)
		return
	}

	respondWithJSON(w, http.StatusOK, database.Video{
		ID:           meta.ID, `json:"id"`
		ThumbnailURL: &url, `json:"thumbnail_url"`
	})
}
