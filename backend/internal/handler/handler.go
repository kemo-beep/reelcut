package handler

type Handler struct {
	Auth          *AuthHandler
	User          *UserHandler
	Project       *ProjectHandler
	Video         *VideoHandler
	Transcription *TranscriptionHandler
	Analysis      *AnalysisHandler
	Clip          *ClipHandler
	Template      *TemplateHandler
	Job           *JobHandler
	Subscription  *SubscriptionHandler
	Webhook       *WebhookHandler
	WebSocket     *WebSocketHandler
}
