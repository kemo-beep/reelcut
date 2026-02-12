Document Control
VersionDateAuthorChanges1.02026-02-11Engineering TeamInitial design document

Table of Contents

Executive Summary
Product Overview
System Architecture
Data Models & Database Design
API Design
Frontend Architecture
Core Features & Workflows
AI/ML Components
Video Processing Pipeline
Security & Authentication
Performance & Scalability
Third-party Integrations
Deployment Architecture
Monitoring & Observability
Error Handling & Recovery
Testing Strategy
Cost Analysis
Risk Assessment
Future Enhancements
Appendix


1. Executive Summary
1.1 Purpose
This document outlines the technical design for an AI-powered video clip generation platform that automatically transforms long-form video content into engaging short clips optimized for social media platforms.
1.2 Scope
The platform will enable users to:

Upload long-form videos (podcasts, webinars, presentations)
Automatically transcribe and analyze content using AI
Generate AI-suggested clips with virality scoring
Customize clips with captions, branding, and styles
Export clips in multiple formats and aspect ratios
Manage projects and collaborate with team members

1.3 Goals

Time Efficiency: Reduce clip creation time from hours to minutes
Quality: Generate professional-quality clips with minimal manual editing
Scalability: Handle thousands of concurrent users and video processing jobs
Reliability: 99.9% uptime with automated failover and recovery
User Experience: Intuitive interface requiring no video editing expertise

1.4 Success Metrics

Average clip creation time < 5 minutes
User satisfaction score > 4.5/5
Video processing completion rate > 95%
AI clip suggestion acceptance rate > 60%
Monthly active users growth rate > 20%


2. Product Overview
2.1 Target Audience
Primary Users

Content Creators: YouTubers, podcasters, streamers
Marketing Teams: Social media managers, digital marketers
Agencies: Video production agencies, marketing agencies
Enterprises: Corporate communications, HR, training departments

User Personas
Persona 1: Solo Content Creator

Creates 2-3 long-form videos per week
Needs to create 10-15 short clips from each video
Limited video editing skills
Budget-conscious
Primary platforms: YouTube, TikTok, Instagram

Persona 2: Marketing Manager

Manages content for multiple brands
Works with team of 3-5 people
Needs consistent branding across clips
Values collaboration features
Primary platforms: LinkedIn, Twitter, Instagram

Persona 3: Enterprise User

Large volume of video content (webinars, training)
Needs approval workflows
Requires brand compliance
Values API access and integrations
Primary platforms: Internal platforms, YouTube

2.2 Key Features
Core Features (MVP)

Video Upload & Management

Drag-and-drop upload
URL import (YouTube, Vimeo)
Project organization
Video library with search/filter


AI Transcription

Automatic speech-to-text
Word-level timestamps
Multi-language support
Speaker diarization


AI Clip Detection

Automatic clip suggestions
Virality scoring
Hook detection
Natural break point identification


Clip Editor

Timeline-based editing
Trim/cut functionality
Caption customization
Brand overlay


Export & Rendering

Multiple aspect ratios (9:16, 1:1, 16:9)
Quality settings
Format selection
Batch export



Advanced Features (Post-MVP)

Auto-reframing for vertical video
Multi-language caption generation
Template marketplace
Team collaboration
Analytics dashboard
API access
White-label solution

2.3 User Journey
┌─────────────┐
│  Sign Up    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Create Project│
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Upload Video │──────┐
└──────┬──────┘      │
       │             │
       ▼             │
┌─────────────┐      │
│Transcription│      │
│  (Auto)     │      │
└──────┬──────┘      │
       │             │
       ▼             │
┌─────────────┐      │
│AI Analysis  │      │
│  (Auto)     │      │
└──────┬──────┘      │
       │             │
       ▼             │
┌─────────────┐      │
│Review AI    │      │
│Suggestions  │      │
└──────┬──────┘      │
       │             │
       ├─Accept──┐   │
       │         │   │
       ▼         ▼   │
    Create    Edit   │
    Clip      Clip   │
       │         │   │
       └────┬────┘   │
            │        │
            ▼        │
       ┌─────────┐   │
       │Customize│   │
       │ Style   │   │
       └────┬────┘   │
            │        │
            ▼        │
       ┌─────────┐   │
       │ Render  │   │
       └────┬────┘   │
            │        │
            ▼        │
       ┌─────────┐   │
       │Download │   │
       │ & Share │   │
       └─────────┘   │
            │        │
            └────────┘
         (Repeat for
         more clips)

3. System Architecture
3.1 High-Level Architecture
┌─────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Web Browser  │  │ Mobile App   │  │  API Client  │      │
│  │ (React SPA)  │  │  (Future)    │  │   (Future)   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTPS
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                      CDN / CLOUDFRONT                         │
│  • Static Assets (JS, CSS, Images)                           │
│  • Video Thumbnails & Previews                               │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                    LOAD BALANCER (ALB)                        │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                   APPLICATION LAYER                           │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌────────────────────────────────────────────────────┐     │
│  │         API Gateway (Go Gin)                       │     │
│  │  • Authentication & Authorization                  │     │
│  │  • Rate Limiting                                   │     │
│  │  • Request Validation                              │     │
│  │  • Response Formatting                             │     │
│  └────────────┬───────────────────────┬────────────┬──┘     │
│               │                       │            │         │
│               ▼                       ▼            ▼         │
│  ┌──────────────────┐   ┌──────────────────┐  ┌──────────┐ │
│  │ Auth Service     │   │ Video Service    │  │User Svc  │ │
│  ├──────────────────┤   ├──────────────────┤  └──────────┘ │
│  │ Clip Service     │   │ Transcription    │               │
│  ├──────────────────┤   │    Service       │               │
│  │ Template Service │   ├──────────────────┤               │
│  └──────────────────┘   │ Analysis Service │               │
│                         ├──────────────────┤               │
│                         │ Rendering Svc    │               │
│                         └──────────────────┘               │
│                                                               │
└───────────────┬───────────────────────┬─────────────────────┘
                │                       │
                ▼                       ▼
┌─────────────────────────┐  ┌──────────────────────────────┐
│   DATA LAYER            │  │   WORKER LAYER               │
├─────────────────────────┤  ├──────────────────────────────┤
│                         │  │                              │
│ ┌─────────────────┐    │  │ ┌──────────────────────────┐ │
│ │   PostgreSQL    │    │  │ │   Worker Pool Manager    │ │
│ │   (Primary DB)  │    │  │ └───────────┬──────────────┘ │
│ └─────────────────┘    │  │             │                │
│                         │  │   ┌─────────▼────────┐      │
│ ┌─────────────────┐    │  │   │ Transcription    │      │
│ │  Redis Cache    │    │  │   │    Worker        │      │
│ │  • Session      │    │  │   ├──────────────────┤      │
│ │  • Job Queue    │    │  │   │  Analysis        │      │
│ │  • Rate Limit   │    │  │   │    Worker        │      │
│ └─────────────────┘    │  │   ├──────────────────┤      │
│                         │  │   │  Rendering       │      │
│ ┌─────────────────┐    │  │   │    Worker        │      │
│ │   Amazon S3     │    │  │   └──────────────────┘      │
│ │  • Raw Videos   │    │  │                              │
│ │  • Rendered     │    │  │   Uses: FFmpeg, Whisper,    │
│ │  • Thumbnails   │    │  │   OpenCV, PyTorch           │
│ └─────────────────┘    │  │                              │
│                         │  └──────────────────────────────┘
└─────────────────────────┘
                │
                ▼
┌──────────────────────────────────────────────────┐
│         EXTERNAL SERVICES                         │
├──────────────────────────────────────────────────┤
│                                                   │
│  • OpenAI Whisper API (Transcription)            │
│  • Stripe (Payments)                             │
│  • SendGrid (Email)                              │
│  • Sentry (Error Tracking)                       │
│  • DataDog (Monitoring)                          │
│                                                   │
└──────────────────────────────────────────────────┘
3.2 Component Descriptions
3.2.1 Frontend (React + TanStack Start)
Responsibilities:

User interface rendering
Client-side routing
State management (Zustand)
API communication
Real-time updates (WebSocket)
File upload handling

Technology Stack:

TanStack Start (React framework)
TypeScript
Zustand (State management)
TanStack Query (Data fetching)
Tailwind CSS (Styling)
Video.js (Video player)

3.2.2 API Gateway (Go Gin)
Responsibilities:

HTTP request routing
Authentication/authorization
Rate limiting
Request validation
Response formatting
CORS handling
WebSocket connections

Technology Stack:

Go 1.21+
Gin web framework
JWT for authentication
Redis for rate limiting

3.2.3 Services Layer
Auth Service

User registration/login
Token generation/validation
Password reset
Email verification
Session management

Video Service

Video upload orchestration
Metadata extraction
Storage management
Video listing/filtering
Thumbnail generation

Transcription Service

Whisper API integration
Transcript management
Word-level timestamp processing
Speaker diarization
Multi-language support

Analysis Service

Scene detection
Face/speaker tracking
Sentiment analysis
Topic extraction
Engagement scoring

Clip Service

Clip creation/management
Style management
Template application
Virality scoring

Rendering Service

Video cutting/trimming
Aspect ratio conversion
Caption burning
Logo overlay
Audio mixing
Export format handling

3.2.4 Data Layer
PostgreSQL

Primary relational database
User data
Video metadata
Transcriptions
Clips and styles
Templates
Subscriptions

Redis

Session storage
Job queue (using Asynq)
Rate limiting counters
Cache layer
Real-time pub/sub

Amazon S3

Original video files
Rendered clips
Thumbnails
User avatars
Brand assets

3.2.5 Worker Layer
Worker Pool Manager

Job distribution
Worker health monitoring
Load balancing
Queue management

Specialized Workers

Transcription Worker: Processes audio → text
Analysis Worker: AI analysis of content
Rendering Worker: Video encoding/rendering

3.3 Data Flow Diagrams
3.3.1 Video Upload Flow
User Browser                API Gateway              Storage Service           Database
     │                           │                           │                    │
     │  1. Request Upload URL    │                           │                    │
     ├──────────────────────────>│                           │                    │
     │                           │                           │                    │
     │                           │  2. Generate Presigned    │                    │
     │                           │     URL                   │                    │
     │                           ├──────────────────────────>│                    │
     │                           │                           │                    │
     │                           │  3. Return Presigned URL  │                    │
     │                           │<──────────────────────────┤                    │
     │  4. Upload URL            │                           │                    │
     │<──────────────────────────┤                           │                    │
     │                           │                           │                    │
     │  5. Direct Upload to S3   │                           │                    │
     ├───────────────────────────────────────────────────────>│                    │
     │                           │                           │                    │
     │  6. Upload Complete       │                           │                    │
     │<───────────────────────────────────────────────────────┤                    │
     │                           │                           │                    │
     │  7. Notify Upload Done    │                           │                    │
     ├──────────────────────────>│                           │                    │
     │                           │                           │                    │
     │                           │  8. Create Video Record   │                    │
     │                           ├────────────────────────────────────────────────>│
     │                           │                           │                    │
     │                           │  9. Enqueue Processing    │                    │
     │                           │     Jobs (Metadata,       │                    │
     │                           │     Thumbnail, etc.)      │                    │
     │                           ├──────────────────────────>│                    │
     │                           │                           │                    │
     │  10. Success Response     │                           │                    │
     │<──────────────────────────┤                           │                    │
     │                           │                           │                    │
3.3.2 AI Clip Suggestion Flow
User                  API                Worker               AI Services         Database
 │                     │                   │                      │                  │
 │ 1. Request Clips    │                   │                      │                  │
 ├────────────────────>│                   │                      │                  │
 │                     │                   │                      │                  │
 │                     │ 2. Enqueue Job    │                      │                  │
 │                     ├──────────────────>│                      │                  │
 │                     │                   │                      │                  │
 │ 3. Job Queued       │                   │                      │                  │
 │<────────────────────┤                   │                      │                  │
 │                     │                   │                      │                  │
 │                     │                   │ 4. Get Transcription │                  │
 │                     │                   ├──────────────────────────────────────────>│
 │                     │                   │                      │                  │
 │                     │                   │ 5. Transcript Data   │                  │
 │                     │                   │<──────────────────────────────────────────┤
 │                     │                   │                      │                  │
 │                     │                   │ 6. Analyze Sentiment │                  │
 │                     │                   ├─────────────────────>│                  │
 │                     │                   │                      │                  │
 │                     │                   │ 7. Detect Scenes     │                  │
 │                     │                   ├─────────────────────>│                  │
 │                     │                   │                      │                  │
 │                     │                   │ 8. Score Segments    │                  │
 │                     │                   │  (Internal AI Logic) │                  │
 │                     │                   │                      │                  │
 │                     │                   │ 9. Save Suggestions  │                  │
 │                     │                   ├──────────────────────────────────────────>│
 │                     │                   │                      │                  │
 │                     │ 10. Job Complete  │                      │                  │
 │                     │    (WebSocket)    │                      │                  │
 │<─ ─ ─ ─ ─ ─ ─ ─ ─ ─┤                   │                      │                  │
 │                     │                   │                      │                  │
 │ 11. Get Suggestions │                   │                      │                  │
 ├────────────────────>│                   │                      │                  │
 │                     │                   │                      │                  │
 │                     │ 12. Query DB      │                      │                  │
 │                     ├──────────────────────────────────────────────────────────────>│
 │                     │                   │                      │                  │
 │                     │ 13. Suggestions   │                      │                  │
 │                     │<──────────────────────────────────────────────────────────────┤
 │ 14. Display Clips   │                   │                      │                  │
 │<────────────────────┤                   │                      │                  │
3.3.3 Clip Rendering Flow
User               API           Queue          Rendering Worker      FFmpeg        Storage       DB
 │                  │              │                    │                │             │           │
 │ 1. Render Clip   │              │                    │                │             │           │
 ├─────────────────>│              │                    │                │             │           │
 │                  │              │                    │                │             │           │
 │                  │ 2. Validate  │                    │                │             │           │
 │                  │    Settings  │                    │                │             │           │
 │                  │              │                    │                │             │           │
 │                  │ 3. Enqueue   │                    │                │             │           │
 │                  ├─────────────>│                    │                │             │           │
 │                  │              │                    │                │             │           │
 │ 4. Render Started│              │                    │                │             │           │
 │<─────────────────┤              │                    │                │             │           │
 │                  │              │                    │                │             │           │
 │                  │              │ 5. Dequeue Job     │                │             │           │
 │                  │              ├───────────────────>│                │             │           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 6. Get Source  │             │           │
 │                  │              │                    ├────────────────────────────>│           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 7. Download    │             │           │
 │                  │              │                    │<────────────────────────────┤           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 8. Build FFmpeg Command      │           │
 │                  │              │                    │   (trim, scale, captions,    │           │
 │                  │              │                    │    overlays, etc.)           │           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 9. Execute     │             │           │
 │                  │              │                    ├───────────────>│             │           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 10. Progress   │             │           │
 │                  │              │                    │     Updates    │             │           │
 │<─ ─ ─ ─ ─ ─ ─ ─ ─│─ ─ ─ ─ ─ ─ ─│─ ─ ─ ─ ─ ─ ─ ─ ─ ─│ (WebSocket)    │             │           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 11. Rendered   │             │           │
 │                  │              │                    │<───────────────┤             │           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 12. Upload Result            │           │
 │                  │              │                    ├──────────────────────────────>│           │
 │                  │              │                    │                │             │           │
 │                  │              │                    │ 13. Update DB  │             │           │
 │                  │              │                    ├────────────────────────────────────────>│
 │                  │              │                    │                │             │           │
 │                  │              │ 14. Job Complete   │                │             │           │
 │                  │              │<───────────────────┤                │             │           │
 │                  │              │                    │                │             │           │
 │ 15. Render Done  │              │                    │                │             │           │
 │<─────────────────┤              │                    │                │             │           │
 │  (WebSocket)     │              │                    │                │             │           │
3.4 Technology Decisions
Backend: Go (Golang)
Rationale:

High performance and low memory footprint
Excellent concurrency support (goroutines)
Strong typing reduces runtime errors
Fast compilation and execution
Rich ecosystem for web development
Native support for building microservices

Alternatives Considered:

Node.js: Less performant for CPU-intensive tasks
Python: Slower execution, GIL limitations
Java: More verbose, higher memory usage

Frontend: TanStack Start + React
Rationale:

Modern React framework with SSR/SSG support
Built-in routing and data fetching
TypeScript support out of the box
Great developer experience
Growing ecosystem

Alternatives Considered:

Next.js: More opinionated, larger bundle
Remix: Newer, smaller ecosystem
Vue/Nuxt: Different paradigm, team expertise

Database: PostgreSQL
Rationale:

ACID compliance
Excellent JSON/JSONB support
Full-text search capabilities
Mature and stable
Rich indexing options
Strong community support

Alternatives Considered:

MySQL: Less feature-rich
MongoDB: Not ideal for relational data
CockroachDB: Overkill for initial scale

Queue: Redis + Asynq
Rationale:

In-memory speed
Native pub/sub support
Atomic operations
Asynq provides Go-native queue with retry logic
Used for multiple purposes (cache, queue, sessions)

Alternatives Considered:

RabbitMQ: More complex setup
SQS: AWS vendor lock-in
Kafka: Overkill for use case

Storage: Amazon S3
Rationale:

Industry standard for object storage
99.999999999% durability
Scalable and cost-effective
CDN integration (CloudFront)
Lifecycle policies for cost optimization

Alternatives Considered:

Google Cloud Storage: Similar but less familiar
Azure Blob: Similar but less familiar
Self-hosted MinIO: More operational overhead


4. Data Models & Database Design
4.1 Entity Relationship Diagram
┌─────────────────┐
│     users       │
├─────────────────┤
│ id (PK)         │
│ email           │
│ password_hash   │
│ full_name       │
│ avatar_url      │
│ subscription_tier│
│ created_at      │
└────────┬────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐         ┌─────────────────┐
│   projects      │         │ subscriptions   │
├─────────────────┤         ├─────────────────┤
│ id (PK)         │         │ id (PK)         │
│ user_id (FK)    │         │ user_id (FK)    │
│ name            │         │ tier            │
│ description     │         │ status          │
│ created_at      │         │ stripe_id       │
└────────┬────────┘         └─────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐
│     videos      │
├─────────────────┤
│ id (PK)         │
│ project_id (FK) │
│ user_id (FK)    │
│ storage_path    │
│ duration        │
│ width, height   │
│ status          │
│ created_at      │
└────────┬────────┘
         │
         │ 1:1
         ▼
┌─────────────────┐
│transcriptions   │
├─────────────────┤
│ id (PK)         │
│ video_id (FK)   │
│ language        │
│ status          │
└────────┬────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐         ┌─────────────────┐
│transcript_      │         │transcript_      │
│  segments       │    1:N  │   words         │
├─────────────────┤────────>├─────────────────┤
│ id (PK)         │         │ id (PK)         │
│ transcription_id│         │ segment_id (FK) │
│ start_time      │         │ word            │
│ end_time        │         │ start_time      │
│ text            │         │ end_time        │
│ speaker_id      │         │ confidence      │
└─────────────────┘         └─────────────────┘

┌─────────────────┐
│ video_analysis  │
├─────────────────┤
│ id (PK)         │
│ video_id (FK)   │───┐
│ scenes (JSONB)  │   │
│ faces (JSONB)   │   │ 1:1
│ topics (JSONB)  │   │
│ sentiment(JSONB)│   │
└─────────────────┘   │
                      │
┌─────────────────┐   │
│     clips       │<──┘
├─────────────────┤ N:1
│ id (PK)         │
│ video_id (FK)   │
│ user_id (FK)    │
│ name            │
│ start_time      │
│ end_time        │
│ aspect_ratio    │
│ virality_score  │
│ status          │
│ storage_path    │
└────────┬────────┘
         │
         │ 1:1
         ▼
┌─────────────────┐         ┌─────────────────┐
│  clip_styles    │         │   templates     │
├─────────────────┤         ├─────────────────┤
│ id (PK)         │         │ id (PK)         │
│ clip_id (FK)    │         │ user_id (FK)    │
│ caption_enabled │         │ name            │
│ caption_font    │         │ category        │
│ caption_size    │         │ is_public       │
│ caption_color   │         │ style_config    │
│ brand_logo_url  │         │   (JSONB)       │
│ overlay_template│         │ usage_count     │
└─────────────────┘         └─────────────────┘

┌─────────────────┐
│processing_jobs  │
├─────────────────┤
│ id (PK)         │
│ user_id (FK)    │
│ job_type        │
│ entity_type     │
│ entity_id       │
│ status          │
│ progress        │
│ error_message   │
│ metadata (JSONB)│
│ started_at      │
│ completed_at    │
└─────────────────┘
4.2 Key Database Tables
4.2.1 users
sqlCREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url TEXT,
    subscription_tier VARCHAR(50) DEFAULT 'free',
    credits_remaining INT DEFAULT 0,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_subscription_tier ON users(subscription_tier);
Design Decisions:

UUID for primary key (better for distributed systems, no sequential guessing)
Soft delete with deleted_at (GDPR compliance, data recovery)
Email verification flag for security
Credits system for flexible pricing models

4.2.2 videos
sqlCREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    storage_path TEXT NOT NULL,
    thumbnail_url TEXT,
    duration_seconds DECIMAL(10, 2),
    width INT,
    height INT,
    fps DECIMAL(5, 2),
    file_size_bytes BIGINT,
    codec VARCHAR(50),
    bitrate INT,
    status VARCHAR(50) DEFAULT 'uploading',
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_project_id ON videos(project_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
Design Decisions:

Cascade delete on project deletion (cleanup)
Store technical metadata for processing decisions
Status field for tracking upload/processing state
JSONB for flexible metadata storage

4.2.3 transcriptions & segments
sqlCREATE TABLE transcriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    status VARCHAR(50) DEFAULT 'pending',
    error_message TEXT,
    word_count INT,
    duration_seconds DECIMAL(10, 2),
    confidence_avg DECIMAL(5, 4),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transcript_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transcription_id UUID REFERENCES transcriptions(id) ON DELETE CASCADE,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    text TEXT NOT NULL,
    confidence DECIMAL(5, 4),
    speaker_id INT,
    sequence_order INT NOT NULL,
    CONSTRAINT check_time_order CHECK (end_time > start_time)
);

CREATE TABLE transcript_words (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    segment_id UUID REFERENCES transcript_segments(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    confidence DECIMAL(5, 4),
    sequence_order INT NOT NULL
);

CREATE INDEX idx_transcriptions_video_id ON transcriptions(video_id);
CREATE INDEX idx_segments_transcription_id ON transcript_segments(transcription_id);
CREATE INDEX idx_segments_time ON transcript_segments(start_time, end_time);
CREATE INDEX idx_words_segment_id ON transcript_words(segment_id);
CREATE INDEX idx_words_time ON transcript_words(start_time);
Design Decisions:

Three-level hierarchy: transcription → segments → words
Word-level granularity for caption generation
Speaker diarization support
Time-based indexing for fast seeking
Sequence order for maintaining text order

4.2.4 clips & clip_styles
sqlCREATE TABLE clips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    duration_seconds DECIMAL(10, 2),
    aspect_ratio VARCHAR(10) DEFAULT '9:16',
    virality_score DECIMAL(5, 2),
    status VARCHAR(50) DEFAULT 'draft',
    storage_path TEXT,
    thumbnail_url TEXT,
    is_ai_suggested BOOLEAN DEFAULT false,
    suggestion_reason TEXT,
    view_count INT DEFAULT 0,
    download_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT check_clip_time CHECK (end_time > start_time),
    CONSTRAINT check_virality_score CHECK (virality_score >= 0 AND virality_score <= 100)
);

CREATE TABLE clip_styles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clip_id UUID UNIQUE REFERENCES clips(id) ON DELETE CASCADE,
    
    -- Caption settings
    caption_enabled BOOLEAN DEFAULT true,
    caption_font VARCHAR(100) DEFAULT 'Inter',
    caption_size INT DEFAULT 48,
    caption_color VARCHAR(7) DEFAULT '#FFFFFF',
    caption_bg_color VARCHAR(9),
    caption_position VARCHAR(50) DEFAULT 'bottom',
    caption_animation VARCHAR(50),
    caption_max_words INT DEFAULT 3,
    
    -- Brand settings
    brand_logo_url TEXT,
    brand_logo_position VARCHAR(50),
    brand_logo_scale DECIMAL(3, 2) DEFAULT 1.0,
    brand_watermark_opacity DECIMAL(3, 2) DEFAULT 0.8,
    
    -- Overlay & effects
    overlay_template VARCHAR(100),
    transition_effect VARCHAR(50),
    
    -- Audio settings
    background_music_url TEXT,
    background_music_volume DECIMAL(3, 2) DEFAULT 0.3,
    original_audio_volume DECIMAL(3, 2) DEFAULT 1.0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_clips_video_id ON clips(video_id);
CREATE INDEX idx_clips_user_id ON clips(user_id);
CREATE INDEX idx_clips_status ON clips(status);
CREATE INDEX idx_clips_is_ai_suggested ON clips(is_ai_suggested);
CREATE INDEX idx_clips_virality_score ON clips(virality_score DESC);
Design Decisions:

Separate style table (1:1) to keep clip data lean
Comprehensive style options for customization
Virality score for ranking
Track AI suggestions separately
Usage metrics (views, downloads) for analytics

4.2.5 processing_jobs
sqlCREATE TABLE processing_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL, -- transcription, analysis, rendering
    entity_type VARCHAR(50) NOT NULL, -- video, clip
    entity_id UUID NOT NULL,
    priority INT DEFAULT 5, -- 1-10, higher = more important
    status VARCHAR(50) DEFAULT 'pending',
    progress INT DEFAULT 0,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    metadata JSONB,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_jobs_status_priority ON processing_jobs(status, priority DESC, created_at);
CREATE INDEX idx_jobs_entity ON processing_jobs(entity_type, entity_id);
CREATE INDEX idx_jobs_user_id ON processing_jobs(user_id);
CREATE INDEX idx_jobs_created_at ON processing_jobs(created_at DESC);
Design Decisions:

Generic job table for all processing types
Priority queue support
Retry logic built-in
Progress tracking for user feedback
JSONB metadata for job-specific data

4.3 Data Migration Strategy
4.3.1 Migration Tool: golang-migrate
bash# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir db/migrations -seq create_users_table

# Run migrations
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" down 1
4.3.2 Migration Best Practices
Rules:

Never modify existing migrations after deployment
Always provide down migrations
Test migrations on copy of production data
Use transactions where possible
Avoid data migrations in schema migrations
Monitor migration duration
Have rollback plan ready

Example Migration:
sql-- 001_create_users_table.up.sql
BEGIN;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

COMMIT;

-- 001_create_users_table.down.sql
BEGIN;

DROP TABLE IF EXISTS users CASCADE;

COMMIT;
4.4 Database Optimization
4.4.1 Indexing Strategy
Rules:

Index foreign keys
Index columns used in WHERE clauses
Index columns used in ORDER BY
Composite indexes for common query patterns
Avoid over-indexing (impacts write performance)

Example:
sql-- Good: Composite index for common query
CREATE INDEX idx_clips_user_status ON clips(user_id, status);

-- Query that benefits
SELECT * FROM clips WHERE user_id = ? AND status = 'ready' ORDER BY created_at DESC;
4.4.2 Query Optimization
Techniques:

Use EXPLAIN ANALYZE to understand query plans
Avoid SELECT *
Use LIMIT for pagination
Leverage materialized views for complex aggregations
Use connection pooling
Implement read replicas for heavy read workloads

Example:
sql-- Bad: Full table scan
SELECT * FROM videos WHERE LOWER(original_filename) LIKE '%presentation%';

-- Good: Use full-text search with GIN index
ALTER TABLE videos ADD COLUMN filename_tsvector tsvector;
CREATE INDEX idx_videos_filename_search ON videos USING GIN(filename_tsvector);

UPDATE videos SET filename_tsvector = to_tsvector('english', original_filename);

SELECT * FROM videos WHERE filename_tsvector @@ to_tsquery('presentation');
4.4.3 Connection Pooling
go// Database configuration
const (
    MaxOpenConns    = 25
    MaxIdleConns    = 5
    ConnMaxLifetime = 5 * time.Minute
    ConnMaxIdleTime = 10 * time.Minute
)

db.SetMaxOpenConns(MaxOpenConns)
db.SetMaxIdleConns(MaxIdleConns)
db.SetConnMaxLifetime(ConnMaxLifetime)
db.SetConnMaxIdleTime(ConnMaxIdleTime)
4.4.4 Partitioning Strategy
For large tables (videos, clips, usage_logs), implement time-based partitioning:
sql-- Partition videos by month
CREATE TABLE videos (
    id UUID,
    created_at TIMESTAMP NOT NULL,
    -- ... other columns
) PARTITION BY RANGE (created_at);

CREATE TABLE videos_2026_02 PARTITION OF videos
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE videos_2026_03 PARTITION OF videos
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
```

---

## 5. API Design

### 5.1 API Principles

**RESTful Conventions:**
- Use standard HTTP methods (GET, POST, PUT, DELETE, PATCH)
- Resource-based URLs
- Proper status codes
- Consistent error format
- Versioned endpoints (/api/v1)

**Design Standards:**
- JSON request/response bodies
- ISO 8601 timestamps
- Pagination for list endpoints
- Filtering and sorting support
- Rate limiting headers

### 5.2 API Endpoints

#### 5.2.1 Authentication Endpoints
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
POST   /api/v1/auth/verify-email
Example: POST /api/v1/auth/register
Request:
json{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
Response (201 Created):
json{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "full_name": "John Doe",
    "avatar_url": null,
    "subscription_tier": "free",
    "created_at": "2026-02-11T10:30:00Z"
  },
  "token": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJl...",
    "expires_in": 3600
  }
}
Error Response (400 Bad Request):
json{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Email already exists"
      }
    ]
  }
}
```

#### 5.2.2 Video Endpoints
```
POST   /api/v1/videos/upload             - Upload video
POST   /api/v1/videos/upload/url         - Upload from URL
GET    /api/v1/videos                    - List videos
GET    /api/v1/videos/:id                - Get video details
PUT    /api/v1/videos/:id                - Update video
DELETE /api/v1/videos/:id                - Delete video
GET    /api/v1/videos/:id/metadata       - Get technical metadata
GET    /api/v1/videos/:id/thumbnail      - Get thumbnail
POST   /api/v1/videos/:id/regenerate-thumbnail
```

**Example: POST /api/v1/videos/upload**

Request (Multipart Form Data):
```
Content-Type: multipart/form-data

project_id: 123e4567-e89b-12d3-a456-426614174000
file: [binary video file]
Response (202 Accepted):
json{
  "video": {
    "id": "789e4567-e89b-12d3-a456-426614174000",
    "project_id": "123e4567-e89b-12d3-a456-426614174000",
    "original_filename": "podcast_episode_1.mp4",
    "status": "uploading",
    "created_at": "2026-02-11T10:35:00Z"
  },
  "upload": {
    "progress_url": "/api/v1/videos/789e4567-e89b-12d3-a456-426614174000/upload-progress"
  }
}
```

**Example: GET /api/v1/videos**

Query Parameters:
```
?page=1
&per_page=20
&project_id=123e4567-e89b-12d3-a456-426614174000
&status=ready
&sort_by=created_at
&sort_order=desc
&search=podcast
Response (200 OK):
json{
  "videos": [
    {
      "id": "789e4567-e89b-12d3-a456-426614174000",
      "project_id": "123e4567-e89b-12d3-a456-426614174000",
      "original_filename": "podcast_episode_1.mp4",
      "thumbnail_url": "https://cdn.example.com/thumbnails/789e4567.jpg",
      "duration_seconds": 3600.5,
      "width": 1920,
      "height": 1080,
      "file_size_bytes": 524288000,
      "status": "ready",
      "created_at": "2026-02-11T10:35:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_count": 87
  }
}
```

#### 5.2.3 Transcription Endpoints
```
POST   /api/v1/transcriptions/videos/:videoId    - Start transcription
GET    /api/v1/transcriptions/:id                - Get transcription
GET    /api/v1/transcriptions/videos/:videoId    - Get by video ID
PUT    /api/v1/transcriptions/:id/segments/:segmentId  - Edit segment
POST   /api/v1/transcriptions/:id/export         - Export (SRT/VTT/TXT)
Example: POST /api/v1/transcriptions/videos/:videoId
Request:
json{
  "language": "en",
  "enable_diarization": true
}
Response (202 Accepted):
json{
  "transcription": {
    "id": "abc12345-e89b-12d3-a456-426614174000",
    "video_id": "789e4567-e89b-12d3-a456-426614174000",
    "language": "en",
    "status": "pending",
    "created_at": "2026-02-11T10:40:00Z"
  },
  "job": {
    "id": "def67890-e89b-12d3-a456-426614174000",
    "status_url": "/api/v1/jobs/def67890-e89b-12d3-a456-426614174000"
  }
}
Example: GET /api/v1/transcriptions/:id
Response (200 OK):
json{
  "transcription": {
    "id": "abc12345-e89b-12d3-a456-426614174000",
    "video_id": "789e4567-e89b-12d3-a456-426614174000",
    "language": "en",
    "status": "completed",
    "word_count": 5420,
    "duration_seconds": 3600.5,
    "confidence_avg": 0.94,
    "created_at": "2026-02-11T10:40:00Z",
    "updated_at": "2026-02-11T10:55:00Z",
    "segments": [
      {
        "id": "seg-001",
        "start_time": 0.0,
        "end_time": 5.2,
        "text": "Welcome to episode one of our podcast.",
        "confidence": 0.96,
        "speaker_id": 1,
        "words": [
          {
            "word": "Welcome",
            "start_time": 0.0,
            "end_time": 0.5,
            "confidence": 0.98
          },
          // ... more words
        ]
      }
      // ... more segments
    ]
  }
}
```

#### 5.2.4 Analysis & Clip Suggestion Endpoints
```
POST   /api/v1/analysis/videos/:videoId            - Start AI analysis
GET    /api/v1/analysis/videos/:videoId            - Get analysis results
POST   /api/v1/analysis/videos/:videoId/suggest-clips  - Generate clip suggestions
Example: POST /api/v1/analysis/videos/:videoId/suggest-clips
Request:
json{
  "min_duration": 30,
  "max_duration": 90,
  "min_virality_score": 60,
  "max_suggestions": 10,
  "preferences": {
    "include_hooks": true,
    "avoid_long_pauses": true,
    "prefer_emotional_peaks": true
  }
}
Response (200 OK):
json{
  "suggestions": [
    {
      "id": "sug-001",
      "start_time": 120.5,
      "end_time": 185.2,
      "duration": 64.7,
      "virality_score": 85,
      "confidence": 0.89,
      "reason": "Strong hook with emotional peak and complete thought",
      "highlights": [
        "Contains question hook at start",
        "High emotional sentiment (joy)",
        "Natural conclusion",
        "Mentions trending topic: AI"
      ],
      "preview_url": "https://api.example.com/previews/sug-001.mp4",
      "transcript": "But here's the crazy part about AI that nobody talks about..."
    },
    {
      "id": "sug-002",
      "start_time": 450.0,
      "end_time": 523.5,
      "duration": 73.5,
      "virality_score": 78,
      "confidence": 0.82,
      "reason": "Controversial statement with strong reaction potential",
      "highlights": [
        "Controversial opinion",
        "Contains statistics",
        "Speaker change adds dynamic",
        "Strong closing statement"
      ],
      "preview_url": "https://api.example.com/previews/sug-002.mp4",
      "transcript": "Most people get this completely wrong. The data shows..."
    }
  ],
  "metadata": {
    "total_suggestions": 8,
    "analysis_duration": 45.2,
    "confidence_avg": 0.84
  }
}
```

#### 5.2.5 Clip Endpoints
```
POST   /api/v1/clips                      - Create clip
GET    /api/v1/clips                      - List clips
GET    /api/v1/clips/:id                  - Get clip details
PUT    /api/v1/clips/:id                  - Update clip
DELETE /api/v1/clips/:id                  - Delete clip
POST   /api/v1/clips/:id/duplicate        - Duplicate clip
POST   /api/v1/clips/:id/render           - Start rendering
GET    /api/v1/clips/:id/render/status    - Get render status
GET    /api/v1/clips/:id/download         - Download rendered clip
POST   /api/v1/clips/:id/cancel-render    - Cancel rendering
Example: POST /api/v1/clips
Request:
json{
  "video_id": "789e4567-e89b-12d3-a456-426614174000",
  "name": "Best moment - AI discussion",
  "start_time": 120.5,
  "end_time": 185.2,
  "aspect_ratio": "9:16",
  "from_suggestion": "sug-001"
}
Response (201 Created):
json{
  "clip": {
    "id": "clip-001",
    "video_id": "789e4567-e89b-12d3-a456-426614174000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Best moment - AI discussion",
    "start_time": 120.5,
    "end_time": 185.2,
    "duration_seconds": 64.7,
    "aspect_ratio": "9:16",
    "virality_score": 85,
    "status": "draft",
    "is_ai_suggested": true,
    "created_at": "2026-02-11T11:00:00Z"
  }
}
```

#### 5.2.6 Clip Style Endpoints
```
GET    /api/v1/clips/:clipId/style                      - Get style
PUT    /api/v1/clips/:clipId/style                      - Update style
POST   /api/v1/clips/:clipId/style/apply-template/:templateId  - Apply template
Example: PUT /api/v1/clips/:clipId/style
Request:
json{
  "caption_enabled": true,
  "caption_font": "Montserrat",
  "caption_size": 52,
  "caption_color": "#FFFFFF",
  "caption_bg_color": "#000000CC",
  "caption_position": "bottom",
  "caption_animation": "slide-up",
  "caption_max_words": 3,
  "brand_logo_url": "https://cdn.example.com/logos/brand.png",
  "brand_logo_position": "top-right",
  "brand_logo_scale": 0.8,
  "background_music_url": "https://cdn.example.com/music/upbeat.mp3",
  "background_music_volume": 0.25,
  "original_audio_volume": 1.0
}
Response (200 OK):
json{
  "style": {
    "id": "style-001",
    "clip_id": "clip-001",
    "caption_enabled": true,
    "caption_font": "Montserrat",
    "caption_size": 52,
    "caption_color": "#FFFFFF",
    "caption_bg_color": "#000000CC",
    "caption_position": "bottom",
    "caption_animation": "slide-up",
    "caption_max_words": 3,
    "brand_logo_url": "https://cdn.example.com/logos/brand.png",
    "brand_logo_position": "top-right",
    "brand_logo_scale": 0.8,
    "background_music_url": "https://cdn.example.com/music/upbeat.mp3",
    "background_music_volume": 0.25,
    "original_audio_volume": 1.0,
    "updated_at": "2026-02-11T11:05:00Z"
  }
}
5.2.7 Rendering Endpoint
Example: POST /api/v1/clips/:id/render
Request:
json{
  "quality": "high",
  "format": "mp4",
  "codec": "h264",
  "bitrate": 5000,
  "resolution": "1080x1920",
  "fps": 30,
  "audio_bitrate": 192
}
Response (202 Accepted):
json{
  "render_job": {
    "id": "render-001",
    "clip_id": "clip-001",
    "status": "queued",
    "priority": 5,
    "estimated_duration": 120,
    "created_at": "2026-02-11T11:10:00Z"
  },
  "websocket_url": "wss://api.example.com/ws/render-001",
  "status_url": "/api/v1/clips/clip-001/render/status"
}
Example: GET /api/v1/clips/:id/render/status
Response (200 OK):
json{
  "render_job": {
    "id": "render-001",
    "clip_id": "clip-001",
    "status": "rendering",
    "progress": 45,
    "current_step": "burning captions",
    "estimated_completion": "2026-02-11T11:12:00Z",
    "started_at": "2026-02-11T11:10:30Z"
  }
}
When complete:
json{
  "render_job": {
    "id": "render-001",
    "clip_id": "clip-001",
    "status": "completed",
    "progress": 100,
    "output_url": "https://cdn.example.com/clips/clip-001_rendered.mp4",
    "file_size_bytes": 15728640,
    "duration_seconds": 64.7,
    "started_at": "2026-02-11T11:10:30Z",
    "completed_at": "2026-02-11T11:12:15Z",
    "render_duration": 105
  }
}
```

#### 5.2.8 Template Endpoints
```
GET    /api/v1/templates                  - List templates
GET    /api/v1/templates/public           - List public templates
POST   /api/v1/templates                  - Create template
GET    /api/v1/templates/:id              - Get template
PUT    /api/v1/templates/:id              - Update template
DELETE /api/v1/templates/:id              - Delete template
POST   /api/v1/templates/:id/duplicate    - Duplicate template
5.3 Error Response Format
Standard Error Structure:
json{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": [
      {
        "field": "field_name",
        "message": "Field-specific error"
      }
    ],
    "request_id": "req-123456",
    "timestamp": "2026-02-11T11:15:00Z"
  }
}
```

**Error Codes:**
- `VALIDATION_ERROR` - Invalid input data (400)
- `AUTHENTICATION_ERROR` - Authentication failed (401)
- `AUTHORIZATION_ERROR` - Insufficient permissions (403)
- `NOT_FOUND` - Resource not found (404)
- `CONFLICT` - Resource conflict (409)
- `RATE_LIMIT_EXCEEDED` - Too many requests (429)
- `INTERNAL_ERROR` - Server error (500)
- `SERVICE_UNAVAILABLE` - Service temporarily unavailable (503)

### 5.4 Rate Limiting

**Rate Limits by Tier:**
- Free: 100 requests/hour
- Pro: 1000 requests/hour
- Enterprise: 10000 requests/hour

**Headers:**
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 847
X-RateLimit-Reset: 1644583200
```

**Implementation:**
- Redis-based token bucket algorithm
- Per-user and per-IP tracking
- Graceful degradation for burst traffic
- WebSocket connections don't count toward limit

### 5.5 Pagination

**Query Parameters:**
```
?page=1              # Page number (1-indexed)
&per_page=20         # Items per page (max 100)
&sort_by=created_at  # Sort field
&sort_order=desc     # Sort direction (asc/desc)
Response Metadata:
json{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total_pages": 10,
    "total_count": 200,
    "has_next": true,
    "has_prev": false
  }
}
```

### 5.6 WebSocket Protocol

**Connection:**
```
wss://api.example.com/ws?token=<jwt_token>
Message Format:
json{
  "type": "JOB_UPDATE",
  "payload": {
    "job_id": "render-001",
    "status": "rendering",
    "progress": 45,
    "message": "Burning captions..."
  },
  "timestamp": "2026-02-11T11:15:30Z"
}
```

**Event Types:**
- `JOB_UPDATE` - Processing job status update
- `TRANSCRIPTION_COMPLETE` - Transcription finished
- `RENDER_COMPLETE` - Rendering finished
- `ERROR` - Error occurred
- `NOTIFICATION` - General notification

---

## 6. Frontend Architecture

### 6.1 Component Architecture
```
┌─────────────────────────────────────────────┐
│           TanStack Start App                │
└─────────────────┬───────────────────────────┘
                  │
    ┌─────────────┴────────────┐
    │                          │
    ▼                          ▼
┌─────────┐              ┌──────────┐
│ Routes  │              │  Layout  │
└────┬────┘              └─────┬────┘
     │                         │
     ├─ Auth Routes            ├─ Header
     ├─ Dashboard Routes       ├─ Sidebar
     ├─ Editor Routes          └─ Footer
     └─ Public Routes
                  │
    ┌─────────────┴────────────┐
    │                          │
    ▼                          ▼
┌─────────────┐        ┌──────────────┐
│  Features   │        │    Shared    │
├─────────────┤        │  Components  │
│ • Video     │        ├──────────────┤
│ • Clips     │        │ • Button     │
│ • Editor    │        │ • Input      │
│ • Templates │        │ • Modal      │
│ • Analytics │        │ • Card       │
└─────────────┘        └──────────────┘
       │
       ▼
┌──────────────────────┐
│   Business Logic     │
├──────────────────────┤
│ • API Clients        │
│ • Hooks              │
│ • Utils              │
│ • Stores (Zustand)   │
└──────────────────────┘
6.2 State Management Strategy
6.2.1 Global State (Zustand)
Auth Store:
typescriptinterface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
}
Video Store:
typescriptinterface VideoState {
  videos: Video[];
  currentVideo: Video | null;
  isLoading: boolean;
  error: Error | null;
  fetchVideos: (filters?: VideoFilters) => Promise<void>;
  uploadVideo: (file: File, projectId: string) => Promise<void>;
  selectVideo: (id: string) => void;
}
Editor Store:
typescriptinterface EditorState {
  currentClip: Clip | null;
  currentTime: number;
  isPlaying: boolean;
  zoom: number;
  selectedLayer: string | null;
  history: EditorAction[];
  historyIndex: number;
  
  // Actions
  setClip: (clip: Clip) => void;
  updateStyle: (updates: Partial<ClipStyle>) => void;
  undo: () => void;
  redo: () => void;
  play: () => void;
  pause: () => void;
}
UI Store:
typescriptinterface UIState {
  modals: {
    [key: string]: boolean;
  };
  toasts: Toast[];
  sidebarCollapsed: boolean;
  theme: 'light' | 'dark';
  
  openModal: (modalId: string) => void;
  closeModal: (modalId: string) => void;
  showToast: (toast: Toast) => void;
  toggleSidebar: () => void;
}
6.2.2 Server State (TanStack Query)
typescript// Video queries
const useVideos = (filters?: VideoFilters) => {
  return useQuery({
    queryKey: ['videos', filters],
    queryFn: () => videoApi.getVideos(filters),
    staleTime: 30000,
  });
};

const useVideo = (id: string) => {
  return useQuery({
    queryKey: ['video', id],
    queryFn: () => videoApi.getVideo(id),
  });
};

// Mutations
const useUploadVideo = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: UploadVideoData) => videoApi.upload(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['videos'] });
    },
  });
};

const useRenderClip = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (clipId: string) => clipApi.render(clipId),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['clip', data.clip.id] });
    },
  });
};
6.3 Routing Structure
typescript// app/routes/__root.tsx
export const Route = createRootRoute({
  component: RootLayout,
});

// app/routes/index.tsx
export const Route = createRoute({
  getParentRoute: () => RootRoute,
  path: '/',
  component: LandingPage,
});

// app/routes/auth/login.tsx
export const Route = createRoute({
  getParentRoute: () => RootRoute,
  path: '/auth/login',
  component: LoginPage,
});

// app/routes/dashboard/_layout.tsx
export const Route = createRoute({
  getParentRoute: () => RootRoute,
  path: '/dashboard',
  component: DashboardLayout,
  beforeLoad: ({ context }) => {
    if (!context.auth.isAuthenticated) {
      throw redirect({ to: '/auth/login' });
    }
  },
});

// app/routes/dashboard/videos/$videoId.tsx
export const Route = createRoute({
  getParentRoute: () => DashboardRoute,
  path: '/videos/$videoId',
  component: VideoDetailPage,
  loader: async ({ params }) => {
    const video = await videoApi.getVideo(params.videoId);
    return { video };
  },
});
6.4 Key Components
6.4.1 Video Player Component
typescriptinterface VideoPlayerProps {
  src: string;
  poster?: string;
  currentTime?: number;
  onTimeUpdate?: (time: number) => void;
  onPlay?: () => void;
  onPause?: () => void;
  markers?: TimelineMarker[];
  showControls?: boolean;
  aspectRatio?: '16:9' | '9:16' | '1:1';
}

export const VideoPlayer: React.FC<VideoPlayerProps> = ({
  src,
  poster,
  currentTime = 0,
  onTimeUpdate,
  markers = [],
  showControls = true,
  aspectRatio = '16:9',
}) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [volume, setVolume] = useState(1);
  const [playbackRate, setPlaybackRate] = useState(1);
  
  // Implement video player logic
  
  return (
    <div className="video-player-container">
      <video
        ref={videoRef}
        src={src}
        poster={poster}
        className={`aspect-${aspectRatio}`}
      />
      {showControls && (
        <VideoControls
          isPlaying={isPlaying}
          volume={volume}
          playbackRate={playbackRate}
          currentTime={currentTime}
          duration={duration}
          onPlayPause={togglePlay}
          onSeek={handleSeek}
          onVolumeChange={setVolume}
          onPlaybackRateChange={setPlaybackRate}
        />
      )}
      <TimelineMarkers markers={markers} />
    </div>
  );
};
6.4.2 Timeline Editor Component
typescriptinterface TimelineEditorProps {
  clip: Clip;
  transcript: Transcription;
  onClipUpdate: (updates: Partial<Clip>) => void;
}

export const TimelineEditor: React.FC<TimelineEditorProps> = ({
  clip,
  transcript,
  onClipUpdate,
}) => {
  const [zoom, setZoom] = useState(1);
  const [selectedSegment, setSelectedSegment] = useState<string | null>(null);
  
  const handleDragStart = (position: 'start' | 'end') => {
    // Handle drag start
  };
  
  const handleDrag = (newTime: number) => {
    // Update clip boundary
  };
  
  return (
    <div className="timeline-editor">
      <TimelineRuler
        duration={clip.duration_seconds}
        zoom={zoom}
      />
      <TranscriptTimeline
        segments={transcript.segments}
        selectedSegment={selectedSegment}
        onSegmentSelect={setSelectedSegment}
      />
      <ClipBoundaries
        startTime={clip.start_time}
        endTime={clip.end_time}
        onStartDrag={handleDragStart}
        onEndDrag={handleDragStart}
      />
      <Waveform
        audioUrl={audioUrl}
        zoom={zoom}
      />
    </div>
  );
};
6.4.3 Caption Editor Component
typescriptinterface CaptionEditorProps {
  clip: Clip;
  style: ClipStyle;
  onStyleUpdate: (updates: Partial<ClipStyle>) => void;
}

export const CaptionEditor: React.FC<CaptionEditorProps> = ({
  clip,
  style,
  onStyleUpdate,
}) => {
  return (
    <div className="caption-editor">
      <div className="preview-panel">
        <CaptionPreview
          clip={clip}
          style={style}
        />
      </div>
      
      <div className="controls-panel">
        <Toggle
          label="Enable Captions"
          checked={style.caption_enabled}
          onChange={(enabled) => onStyleUpdate({ caption_enabled: enabled })}
        />
        
        <Select
          label="Font"
          value={style.caption_font}
          options={FONT_OPTIONS}
          onChange={(font) => onStyleUpdate({ caption_font: font })}
        />
        
        <Slider
          label="Size"
          value={style.caption_size}
          min={24}
          max={72}
          onChange={(size) => onStyleUpdate({ caption_size: size })}
        />
        
        <ColorPicker
          label="Text Color"
          value={style.caption_color}
          onChange={(color) => onStyleUpdate({ caption_color: color })}
        />
        
        <ColorPicker
          label="Background Color"
          value={style.caption_bg_color}
          onChange={(color) => onStyleUpdate({ caption_bg_color: color })}
          includeOpacity
        />
        
        <Select
          label="Position"
          value={style.caption_position}
          options={['top', 'center', 'bottom']}
          onChange={(position) => onStyleUpdate({ caption_position: position })}
        />
        
        <Select
          label="Animation"
          value={style.caption_animation}
          options={ANIMATION_OPTIONS}
          onChange={(animation) => onStyleUpdate({ caption_animation: animation })}
        />
      </div>
    </div>
  );
};
6.5 Performance Optimization
6.5.1 Code Splitting
typescript// Lazy load heavy components
const ClipEditor = lazy(() => import('./components/editor/ClipEditor'));
const TemplateLibrary = lazy(() => import('./components/templates/TemplateLibrary'));

// Route-based code splitting automatically handled by TanStack Start
6.5.2 Memoization
typescript// Expensive calculations
const sortedClips = useMemo(() => {
  return clips.sort((a, b) => b.virality_score - a.virality_score);
}, [clips]);

// Callback stability
const handleClipUpdate = useCallback((clipId: string, updates: Partial<Clip>) => {
  updateClip({ id: clipId, ...updates });
}, [updateClip]);
6.5.3 Virtual Scrolling
typescriptimport { useVirtualizer } from '@tanstack/react-virtual';

export const VideoList: React.FC<{ videos: Video[] }> = ({ videos }) => {
  const parentRef = useRef<HTMLDivElement>(null);
  
  const virtualizer = useVirtualizer({
    count: videos.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 100,
    overscan: 5,
  });
  
  return (
    <div ref={parentRef} className="video-list">
      <div style={{ height: virtualizer.getTotalSize() }}>
        {virtualizer.getVirtualItems().map((virtualRow) => (
          <VideoCard
            key={virtualRow.index}
            video={videos[virtualRow.index]}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              transform: `translateY(${virtualRow.start}px)`,
            }}
          />
        ))}
      </div>
    </div>
  );
};
6.5.4 Image Optimization
typescript// Responsive images with Next/Image-like optimization
export const OptimizedImage: React.FC<ImageProps> = ({
  src,
  alt,
  width,
  height,
}) => {
  const [loaded, setLoaded] = useState(false);
  
  return (
    <div className="relative">
      {!loaded && <Skeleton width={width} height={height} />}
      <img
        src={src}
        alt={alt}
        loading="lazy"
        onLoad={() => setLoaded(true)}
        className={loaded ? 'opacity-100' : 'opacity-0'}
      />
    </div>
  );
};
```

---

## 7. Core Features & Workflows

### 7.1 Video Upload Workflow
```
┌──────────────────────────────────────────────────────┐
│ 1. User selects file or provides URL                 │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 2. Client-side validation                            │
│    • File type (mp4, mov, avi, etc.)                 │
│    • File size (max 10GB)                            │
│    • Duration (if extractable)                       │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 3. Request presigned upload URL from API             │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 4. Direct upload to S3 with progress tracking        │
│    • Multipart upload for large files                │
│    • Resume capability                               │
│    • Real-time progress updates                      │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 5. Notify API of upload completion                   │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 6. API creates video record in database              │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 7. Enqueue processing jobs                           │
│    • Extract metadata (FFmpeg)                       │
│    • Generate thumbnail                              │
│    • Extract audio for transcription                 │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 8. Workers process jobs asynchronously               │
│    • Metadata worker extracts video properties       │
│    • Thumbnail worker generates preview images       │
│    • Results saved to database                       │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 9. Video marked as "ready" in database               │
│    • WebSocket notification sent to user             │
│    • User can now start transcription/analysis       │
└──────────────────────────────────────────────────────┘
Implementation Details:
typescript// Frontend upload implementation
export const useVideoUpload = () => {
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadStatus, setUploadStatus] = useState<'idle' | 'uploading' | 'processing' | 'complete' | 'error'>('idle');
  
  const uploadVideo = async (file: File, projectId: string) => {
    try {
      setUploadStatus('uploading');
      
      // 1. Get presigned URL
      const { uploadUrl, videoId } = await videoApi.getUploadUrl({
        filename: file.name,
        fileSize: file.size,
        projectId,
      });
      
      // 2. Upload to S3 with progress
      await uploadToS3(uploadUrl, file, (progress) => {
        setUploadProgress(progress);
      });
      
      // 3. Notify backend
      await videoApi.notifyUploadComplete(videoId);
      
      setUploadStatus('processing');
      
      // 4. Wait for processing (via WebSocket or polling)
      await waitForProcessing(videoId);
      
      setUploadStatus('complete');
      
    } catch (error) {
      setUploadStatus('error');
      throw error;
    }
  };
  
  return { uploadVideo, uploadProgress, uploadStatus };
};
```

### 7.2 Transcription Workflow
```
┌──────────────────────────────────────────────────────┐
│ 1. User initiates transcription                      │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 2. API creates transcription record                  │
│    • Status: "pending"                               │
│    • Creates processing job                          │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 3. Transcription worker picks up job                 │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 4. Extract audio from video (FFmpeg)                 │
│    ffmpeg -i video.mp4 -vn -acodec pcm_s16le audio.wav│
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 5. Split audio into chunks (if > 25MB)               │
│    • Whisper API has 25MB limit                      │
│    • Split at silence points                         │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 6. Send chunks to Whisper API                        │
│    • Process in parallel                             │
│    • Retry on failure                                │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 7. Parse Whisper response                            │
│    • Extract word-level timestamps                   │
│    • Identify speakers (if diarization enabled)      │
│    • Calculate confidence scores                     │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 8. Save to database                                  │
│    • Create transcript segments                      │
│    • Create transcript words                         │
│    • Update transcription status to "completed"      │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 9. Notify user via WebSocket                         │
│    • User can now view/edit transcript               │
│    • Can proceed to AI analysis                      │
└──────────────────────────────────────────────────────┘
Whisper API Integration:
go// Go implementation
type WhisperClient struct {
    apiKey string
    client *http.Client
}

func (w *WhisperClient) Transcribe(audioPath string, language string) (*TranscriptionResult, error) {
    // Read audio file
    audioData, err := os.ReadFile(audioPath)
    if err != nil {
        return nil, err
    }
    
    // Check file size
    if len(audioData) > 25*1024*1024 {
        return nil, errors.New("audio file exceeds 25MB limit")
    }
    
    // Prepare multipart request
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    
    // Add file
    part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
    if err != nil {
        return nil, err
    }
    part.Write(audioData)
    
    // Add parameters
    writer.WriteField("model", "whisper-1")
    writer.WriteField("language", language)
    writer.WriteField("response_format", "verbose_json")
    writer.WriteField("timestamp_granularities[]", "word")
    
    writer.Close()
    
    // Make request
    req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", body)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+w.apiKey)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    resp, err := w.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Parse response
    var result TranscriptionResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

### 7.3 AI Clip Detection Workflow
```
┌──────────────────────────────────────────────────────┐
│ 1. User requests clip suggestions                    │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 2. Load transcription and video analysis             │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 3. Hook Detection                                    │
│    • Questions ("What if...", "Did you know...")     │
│    • Bold statements ("The truth is...")             │
│    • Numbers/statistics ("95% of people...")         │
│    • Controversy ("Most people get this wrong...")   │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 4. Natural Break Detection                           │
│    • Sentence boundaries                             │
│    • Topic changes                                   │
│    • Speaker changes                                 │
│    • Long pauses (> 1 second)                        │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 5. Sentiment Analysis                                │
│    • Identify emotional peaks                        │
│    • Detect sentiment shifts                         │
│    • Find inspirational moments                      │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 6. Topic Coherence Check                             │
│    • Ensure complete thoughts                        │
│    • Avoid mid-sentence cuts                         │
│    • Maintain context                                │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 7. Duration Filtering                                │
│    • Apply min/max duration constraints              │
│    • Prefer 30-90 second clips                       │
│    • Adjust boundaries if needed                     │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 8. Virality Scoring                                  │
│    • Hook strength: 30%                              │
│    • Emotional impact: 25%                           │
│    • Topic relevance: 20%                            │
│    • Complete thought: 15%                           │
│    • Visual quality: 10%                             │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 9. Rank and return top N suggestions                 │
│    • Sort by virality score                          │
│    • Remove duplicates/overlaps                      │
│    • Include reasoning for each suggestion           │
└──────────────────────────────────────────────────────┘
Scoring Algorithm:
gotype ClipSuggestion struct {
    StartTime      float64
    EndTime        float64
    ViralityScore  float64
    HookScore      float64
    EmotionalScore float64
    TopicScore     float64
    CoherenceScore float64
    Reasoning      []string
}

func (s *AnalysisService) ScoreClipCandidate(
    segments []TranscriptSegment,
    analysis VideoAnalysis,
    startIdx, endIdx int,
) (*ClipSuggestion, error) {
    
    clip := &ClipSuggestion{
        StartTime: segments[startIdx].StartTime,
        EndTime:   segments[endIdx].EndTime,
        Reasoning: []string{},
    }
    
    // 1. Hook Detection (30% of score)
    hookScore := s.detectHooks(segments[startIdx:startIdx+3])
    clip.HookScore = hookScore
    if hookScore > 0.7 {
        clip.Reasoning = append(clip.Reasoning, "Strong hook detected")
    }
    
    // 2. Emotional Analysis (25% of score)
    emotionalScore := s.analyzeEmotion(segments[startIdx:endIdx+1], analysis)
    clip.EmotionalScore = emotionalScore
    if emotionalScore > 0.8 {
        clip.Reasoning = append(clip.Reasoning, "High emotional impact")
    }
    
    // 3. Topic Relevance (20% of score)
    topicScore := s.scoreTopicRelevance(segments[startIdx:endIdx+1], analysis.Topics)
    clip.TopicScore = topicScore
    if topicScore > 0.7 {
        clip.Reasoning = append(clip.Reasoning, "Discusses trending topic")
    }
    
    // 4. Coherence (15% of score)
    coherenceScore := s.checkCoherence(segments[startIdx:endIdx+1])
    clip.CoherenceScore = coherenceScore
    if coherenceScore > 0.9 {
        clip.Reasoning = append(clip.Reasoning, "Complete, coherent thought")
    }
    
    // 5. Visual Quality (10% of score)
    visualScore := s.scoreVisuals(analysis, clip.StartTime, clip.EndTime)
    
    // Calculate weighted score
    clip.ViralityScore = (
        hookScore * 0.30 +
        emotionalScore * 0.25 +
        topicScore * 0.20 +
        coherenceScore * 0.15 +
        visualScore * 0.10
    ) * 100
    
    return clip, nil
}

func (s *AnalysisService) detectHooks(segments []TranscriptSegment) float64 {
    text := strings.ToLower(strings.Join(extractText(segments), " "))
    
    score := 0.0
    
    // Question hooks
    questionPatterns := []string{"what if", "did you know", "have you ever", "imagine if"}
    for _, pattern := range questionPatterns {
        if strings.Contains(text, pattern) {
            score += 0.3
        }
    }
    
    // Bold statements
    boldPatterns := []string{"the truth is", "here's the thing", "let me tell you"}
    for _, pattern := range boldPatterns {
        if strings.Contains(text, pattern) {
            score += 0.25
        }
    }
    
    // Statistics/numbers
    if regexp.MustCompile(`\d+%|\d+ percent`).MatchString(text) {
        score += 0.2
    }
    
    // Controversy
    controversialWords := []string{"wrong", "mistake", "secret", "nobody tells you"}
    for _, word := range controversialWords {
        if strings.Contains(text, word) {
            score += 0.15
        }
    }
    
    return math.Min(score, 1.0)
}
```

### 7.4 Rendering Workflow
```
┌──────────────────────────────────────────────────────┐
│ 1. User clicks "Render" with final clip settings     │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 2. API validates render settings                     │
│    • Check aspect ratio                              │
│    • Validate quality settings                       │
│    • Ensure user has sufficient credits              │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 3. Create render job and enqueue                     │
│    • Assign priority based on user tier              │
│    • Estimate render time                            │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 4. Rendering worker picks up job                     │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 5. Download source video from S3                     │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 6. Build FFmpeg command                              │
│    a. Extract clip segment (trim)                    │
│    b. Scale to target resolution                     │
│    c. Crop to aspect ratio                           │
│    d. Generate caption file (SRT/ASS)                │
│    e. Burn in captions with styling                  │
│    f. Add logo overlay                               │
│    g. Mix audio (background music + original)        │
│    h. Encode to target format/codec                  │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 7. Execute FFmpeg with progress tracking             │
│    • Parse FFmpeg output for progress                │
│    • Update job progress in database                 │
│    • Send WebSocket updates to user                  │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 8. Upload rendered clip to S3                        │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 9. Update database                                   │
│    • Set clip status to "ready"                      │
│    • Save storage path                               │
│    • Mark job as completed                           │
│    • Deduct user credits                             │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│ 10. Notify user via WebSocket                        │
│     • Clip ready for download                        │
│     • Generate download link                         │
└──────────────────────────────────────────────────────┘
FFmpeg Command Builder:
gotype RenderConfig struct {
    SourcePath      string
    OutputPath      string
    StartTime       float64
    EndTime         float64
    Width           int
    Height          int
    AspectRatio     string
    Captions        string // Path to subtitle file
    CaptionStyle    CaptionStyleConfig
    LogoPath        string
    LogoPosition    string
    BackgroundMusic string
    MusicVolume     float64
}

func (r *RenderingService) BuildFFmpegCommand(config RenderConfig) []string {
    args := []string{
        "-i", config.SourcePath,
    }
    
    // Add background music if provided
    if config.BackgroundMusic != "" {
        args = append(args, "-i", config.BackgroundMusic)
    }
    
    // Add logo if provided
    if config.LogoPath != "" {
        args = append(args, "-i", config.LogoPath)
    }
    
    // Complex filter for video processing
    filterComplex := []string{}
    
    // 1. Trim video
    filterComplex = append(filterComplex, 
        fmt.Sprintf("[0:v]trim=start=%f:end=%f,setpts=PTS-STARTPTS[trimmed]",
            config.StartTime, config.EndTime))
    
    // 2. Scale and crop for aspect ratio
    scaleFilter := r.buildScaleFilter(config.AspectRatio, config.Width, config.Height)
    filterComplex = append(filterComplex, 
        fmt.Sprintf("[trimmed]%s[scaled]", scaleFilter))
    
    // 3. Add subtitles with styling
    if config.Captions != "" {
        subtitleFilter := r.buildSubtitleFilter(config.Captions, config.CaptionStyle)
        filterComplex = append(filterComplex, 
            fmt.Sprintf("[scaled]%s[captioned]", subtitleFilter))
    } else {
        filterComplex = append(filterComplex, "[scaled]null[captioned]")
    }
    
    // 4. Add logo overlay
    if config.LogoPath != "" {
        logoFilter := r.buildLogoFilter(config.LogoPosition)
        filterComplex = append(filterComplex, 
            fmt.Sprintf("[captioned][2:v]%s[final]", logoFilter))
    } else {
        filterComplex = append(filterComplex, "[captioned]null[final]")
    }
    
    // Audio processing
    audioFilters := []string{}
    
    // 5. Trim audio
    audioFilters = append(audioFilters,
        fmt.Sprintf("[0:a]atrim=start=%f:end=%f,asetpts=PTS-STARTPTS[original_audio]",
            config.StartTime, config.EndTime))
    
    // 6. Mix with background music if provided
    if config.BackgroundMusic != "" {
        audioFilters = append(audioFilters,
            fmt.Sprintf("[1:a]volume=%f[bg_music]", config.MusicVolume),
            "[original_audio][bg_music]amix=inputs=2:duration=first[audio]")
    } else {
        audioFilters = append(audioFilters, "[original_audio]anull[audio]")
    }
    
    // Combine all filters
    allFilters := append(filterComplex, audioFilters...)
    args = append(args, "-filter_complex", strings.Join(allFilters, ";"))
    
    // Map outputs
    args = append(args,
        "-map", "[final]",
        "-map", "[audio]",
    )
    
    // Encoding settings
    args = append(args,
        "-c:v", "libx264",
        "-preset", "medium",
        "-crf", "23",
        "-c:a", "aac",
        "-b:a", "192k",
        "-movflags", "+faststart",
        config.OutputPath,
    )
    
    return args
}
