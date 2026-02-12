package ai

import "encoding/json"

// FaceDetection is a placeholder for face/speaker detection (MediaPipe/OpenCV later).
type FaceDetection struct {
	Timestamp  float64   `json:"timestamp"`
	Bbox       []float64 `json:"bbox"` // x, y, w, h
	Confidence float64   `json:"confidence"`
}

// DetectFaces stub returns empty list. Replace with MediaPipe/OpenCV when needed.
func DetectFaces(_ string) ([]FaceDetection, error) {
	return nil, nil
}

// FacesToJSON returns JSON for video_analysis.faces_detected.
func FacesToJSON(faces []FaceDetection) ([]byte, error) {
	return json.Marshal(faces)
}
