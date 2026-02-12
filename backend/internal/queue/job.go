package queue

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeVideoMetadata   = "video:metadata"
	TypeVideoThumbnail  = "video:thumbnail"
	TypeVideoFetchURL   = "video:fetch_url"
	TypeTranscription   = "transcription"
	TypeAnalysis        = "analysis"
	TypeRender          = "render"
)

type VideoMetadataPayload struct {
	VideoID string `json:"video_id"`
}

type VideoThumbnailPayload struct {
	VideoID string `json:"video_id"`
}

type VideoFetchURLPayload struct {
	VideoID string `json:"video_id"`
	URL     string `json:"url"`
}

type TranscriptionPayload struct {
	VideoID string `json:"video_id"`
	TranscriptionID string `json:"transcription_id"`
}

type AnalysisPayload struct {
	VideoID string `json:"video_id"`
}

type RenderPayload struct {
	ClipID string `json:"clip_id"`
	JobID  string `json:"job_id"`
}

func NewVideoMetadataTask(videoID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(VideoMetadataPayload{VideoID: videoID.String()})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeVideoMetadata, payload), nil
}

func NewVideoThumbnailTask(videoID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(VideoThumbnailPayload{VideoID: videoID.String()})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeVideoThumbnail, payload), nil
}

func NewVideoFetchURLTask(videoID uuid.UUID, url string) (*asynq.Task, error) {
	payload, err := json.Marshal(VideoFetchURLPayload{VideoID: videoID.String(), URL: url})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeVideoFetchURL, payload), nil
}

func NewTranscriptionTask(videoID, transcriptionID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(TranscriptionPayload{VideoID: videoID.String(), TranscriptionID: transcriptionID.String()})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeTranscription, payload), nil
}

func NewAnalysisTask(videoID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(AnalysisPayload{VideoID: videoID.String()})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAnalysis, payload), nil
}

func NewRenderTask(clipID, jobID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(RenderPayload{ClipID: clipID.String(), JobID: jobID.String()})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRender, payload), nil
}

func ParseVideoMetadataPayload(b []byte) (VideoMetadataPayload, error) {
	var p VideoMetadataPayload
	err := json.Unmarshal(b, &p)
	return p, err
}

func ParseVideoThumbnailPayload(b []byte) (VideoThumbnailPayload, error) {
	var p VideoThumbnailPayload
	err := json.Unmarshal(b, &p)
	return p, err
}

func ParseVideoFetchURLPayload(b []byte) (VideoFetchURLPayload, error) {
	var p VideoFetchURLPayload
	err := json.Unmarshal(b, &p)
	return p, err
}

func ParseTranscriptionPayload(b []byte) (TranscriptionPayload, error) {
	var p TranscriptionPayload
	err := json.Unmarshal(b, &p)
	return p, err
}

func ParseAnalysisPayload(b []byte) (AnalysisPayload, error) {
	var p AnalysisPayload
	err := json.Unmarshal(b, &p)
	return p, err
}

func ParseRenderPayload(b []byte) (RenderPayload, error) {
	var p RenderPayload
	err := json.Unmarshal(b, &p)
	return p, err
}

type QueueClient struct {
	client *asynq.Client
}

func NewQueueClient(redisURL string) (*QueueClient, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	client := asynq.NewClient(opt)
	return &QueueClient{client: client}, nil
}

func (q *QueueClient) EnqueueVideoMetadata(videoID uuid.UUID) error {
	task, err := NewVideoMetadataTask(videoID)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *QueueClient) EnqueueVideoThumbnail(videoID uuid.UUID) error {
	task, err := NewVideoThumbnailTask(videoID)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *QueueClient) EnqueueVideoFetchURL(videoID uuid.UUID, url string) error {
	task, err := NewVideoFetchURLTask(videoID, url)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *QueueClient) EnqueueTranscription(videoID, transcriptionID uuid.UUID) error {
	task, err := NewTranscriptionTask(videoID, transcriptionID)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *QueueClient) EnqueueAnalysis(videoID uuid.UUID) error {
	task, err := NewAnalysisTask(videoID)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *QueueClient) EnqueueRender(clipID, jobID uuid.UUID) error {
	task, err := NewRenderTask(clipID, jobID)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *QueueClient) Close() error {
	return q.client.Close()
}
