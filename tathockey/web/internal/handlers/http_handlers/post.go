package http_handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository"
	"time"
)

type PostManager struct {
	log         *slog.Logger
	postRepo    *repository.Post
	videoRepo   *repository.Video
	sessionRepo *repository.Session
}

func NewPostManager(log *slog.Logger, postRepo *repository.Post, vidRepo *repository.Video, sesRepo *repository.Session) *PostManager {
	return &PostManager{
		log:         log,
		postRepo:    postRepo,
		videoRepo:   vidRepo,
		sessionRepo: sesRepo,
	}
}

const fileNameLength = 20 // Length of the random filename

// generateRandomString generates a random string of given length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

type VideoRequest struct {
	VideoName string `json:"video_name"`
}

// VideoResponse represents the outgoing response format
type VideoResponse struct {
	Seconds   []int    `json:"seconds"`
	EventType []string `json:"event_type"`
}

func FetchAnalyticsData(videoName string, wg *sync.WaitGroup, resultChan chan<- VideoResponse) {
	defer wg.Done()

	// Создание запроса
	requestPayload := map[string]string{"video_name": videoName}
	requestBody, _ := json.Marshal(requestPayload)
	response, err := http.Post("0.0.0.0:5000/infer", "application/json", bytes.NewReader(requestBody))
	if err != nil {
		log.Printf("Failed to fetch analytics data: %v", err)
		resultChan <- VideoResponse{}
		return
	}
	defer response.Body.Close()

	var analyticsResponse VideoResponse
	if err := json.NewDecoder(response.Body).Decode(&analyticsResponse); err != nil {
		log.Printf("Failed to decode analytics response: %v", err)
		resultChan <- VideoResponse{}
		return
	}

	resultChan <- analyticsResponse
}

func (p *PostManager) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form to retrieve file and other data
	err := r.ParseMultipartForm(10 << 31) // Limit the size of the form to 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve title from the form
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Process photo (if provided)
	photo, ph, err := r.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving photo", http.StatusBadRequest)
		return
	}
	var photoFileName string
	if photo != nil {
		defer photo.Close()

		// Generate a unique filename for the photo
		s, _ := generateRandomString(fileNameLength)

		photoFileName = s + filepath.Ext(ph.Filename)
		if err != nil {
			http.Error(w, "Error generating filename", http.StatusInternalServerError)
			return
		}

		// Save the photo
		photoPath := filepath.Join("static", "preview", photoFileName)
		out, err := os.Create(photoPath)
		if err != nil {
			http.Error(w, "Error saving photo", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, photo)
		if err != nil {
			http.Error(w, "Error saving photo", http.StatusInternalServerError)
			return
		}
	}

	// Process video
	video, vh, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving video", http.StatusBadRequest)
		return
	}
	defer video.Close()

	// Generate a unique filename for the video
	s, _ := generateRandomString(fileNameLength)
	videoFileName := s + filepath.Ext(vh.Filename)

	// Save the video
	videoPath := filepath.Join("uploads", videoFileName)
	out, err := os.Create(videoPath)
	if err != nil {
		http.Error(w, "Error saving video", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, video)
	if err != nil {
		http.Error(w, "Error saving video", http.StatusInternalServerError)
		return
	}

	// Асинхронный запрос к аналитическому сервису
	var wg sync.WaitGroup
	wg.Add(1)
	resultChan := make(chan VideoResponse)
	go FetchAnalyticsData(videoFileName, &wg, resultChan)

	// Отправка ответа клиенту сразу же
	response := VideoResponse{
		Seconds:   []int{}, // Можно указать пустые или предварительные данные
		EventType: []string{},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	vId, err := p.videoRepo.Save(&models.Video{
		Name: videoFileName,
	})
	uid, _ := p.sessionRepo.GetUserIDbySessionID(r.Context())
	postID, err := p.postRepo.CreatePost(r.Context(), vId, title, photoFileName, uid)
	if err != nil {
		return
	}

	// Ожидание завершения асинхронного запроса
	go func() {
		defer wg.Done()
		analyticsResponse := <-resultChan
		// Логирование или дальнейшая обработка результатов аналитики
		for i, sec := range analyticsResponse.Seconds {
			path := foo(videoPath, sec, postID)
			p.postRepo.SaveGIF(r.Context(), path, analyticsResponse.EventType[i], postID)
		}
	}()
}

func foo(videoPath string, centralSecond int, postID int) string {

	startTime := centralSecond - 10
	duration := 20

	// Проверка, чтобы начальное время было не меньше нуля
	if startTime < 0 {
		startTime = 0
	}

	// Выходной файл
	outputGif := fmt.Sprintf("/static/gifs/%dgif%d.gif", postID, centralSecond)

	// Команда ffmpeg для вырезки и конвертации видео в gif
	cmd := exec.Command("ffmpeg", "-ss", strconv.Itoa(startTime), "-t", strconv.Itoa(duration), "-i", videoPath, "-vf", "fps=10,scale=1200:-1:flags=lanczos", "-y", outputGif)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	// Запуск команды
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Ошибка при создании GIF: %v\n", err)
		fmt.Printf("stdout: %s\n", outbuf.String())
		fmt.Printf("stderr: %s\n", errbuf.String())
		os.Exit(1)
	}

	fmt.Printf("GIF успешно создан: %s\n", outputGif)
	return outputGif
}

// Example data
var examplePost = models.Post{
	ID:        1,
	Title:     "News Title",
	CreatedAt: time.Now(),
	GIFs: []models.Gif{
		{Name: "1gif180.gif", EventClass: "123"},
		{Name: "1gif60.gif", EventClass: "чета другое"},
	},
}

// Parse the template
var tmpl = template.Must(template.ParseFiles("/app/tmp/post.html"))

// PostHandler handles the "/post" route
func (p *PostManager) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from query params
	idStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil || postID < 1 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Dummy post data, normally fetched from a database
	post, _ := p.postRepo.GetPostByID(postID)

	// Execute the template
	err = tmpl.Execute(w, post)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}

// GifHandler serves the GIFs
func (p *PostManager) GifHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
	gifIDStr := r.URL.Query().Get("gifID")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 1 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	gifID, err := strconv.Atoi(gifIDStr)
	if err != nil || gifID < 1 {
		http.Error(w, "Invalid GIF ID", http.StatusBadRequest)
		return
	}
	// Construct the file path based on postID and gifID
	filePath := fmt.Sprintf("/app/gifs/%dgif%d.gif", postID, gifID)

	// Serve the GIF file
	http.ServeFile(w, r, filePath)
}
