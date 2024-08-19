package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Task struct to hold frame extraction parameters
type Task struct {
	VideoPath string
	OutputDir string
	StartTime string
	Duration  string
	FrameRate int
}

// Worker function to handle frame extraction tasks from the channel
func Worker(taskChan <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		// Ensure the output directory ends with a "/"
		if task.OutputDir[len(task.OutputDir)-1:] != "/" {
			task.OutputDir += "/"
		}

		// Create the output directory if it doesn't exist
		err := exec.Command("mkdir", "-p", task.OutputDir).Run()
		if err != nil {
			fmt.Printf("failed to create output directory: %v\n", err)
			continue
		}

		// Generate the ffmpeg command with segment-specific start time and duration
		outputPattern := filepath.Join(task.OutputDir, "frame_%04d.png")
		//cmd := exec.Command("ffmpeg", "-ss", task.StartTime, "-t", task.Duration, "-i", task.VideoPath, "-vf", fmt.Sprintf("fps=%d", task.FrameRate), outputPattern)
		cmd := exec.Command("ffmpeg", "-threads", "8", "-ss", task.StartTime, "-t", task.Duration, "-i", task.VideoPath, "-vf", fmt.Sprintf("fps=%d", task.FrameRate), outputPattern)

		// Run the command and capture any error
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("failed to extract frames for segment starting at %s: %v, output: %s\n", task.StartTime, err, string(output))
		} else {
			fmt.Printf("Frames extracted successfully for segment starting at %s\n", task.StartTime)
		}
	}
}

// Function to get the duration of the video
func getVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get video duration: %v", err)
	}

	// Convert the output to a float and trim any whitespace/newline characters
	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse video duration: %v", err)
	}
	return duration, nil
}

func main() {
	check := time.Now()

	videoPath := "test.mp4" // Path to the video file
	outputDir := "out"      // Base directory to save frames
	frameRate := 5          // Number of frames per second
	segmentDuration := 60   // Duration of each segment in seconds

	// Get the total duration of the video
	videoDuration, err := getVideoDuration(videoPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Create a channel to pass tasks to workers
	taskChan := make(chan Task, 64)
	var wg sync.WaitGroup

	// Number of worker goroutines (adjust this based on system resources)
	numWorkers := 32

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go Worker(taskChan, &wg)
	}

	// Determine the number of segments
	numSegments := int(videoDuration) / segmentDuration

	// Loop through each segment and create a task to extract frames
	for i := 0; i <= numSegments; i++ {
		// Calculate the start time for the current segment
		startTime := fmt.Sprintf("%02d:%02d:%02d", i*segmentDuration/3600, (i*segmentDuration/60)%60, (i*segmentDuration)%60)

		// Directory for the current segment
		segmentOutputDir := fmt.Sprintf("%s/segment_%d", outputDir, i+1)

		// Create a task and send it to the channel
		taskChan <- Task{
			VideoPath: videoPath,
			OutputDir: segmentOutputDir,
			StartTime: startTime,
			Duration:  fmt.Sprintf("%ds", segmentDuration),
			FrameRate: frameRate,
		}
	}

	// Close the task channel after sending all tasks
	close(taskChan)

	// Wait for all workers to finish
	wg.Wait()

	fmt.Printf("Total time elapsed: %s\n", time.Since(check))
}
