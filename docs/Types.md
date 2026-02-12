// types/video.ts
export interface Video {
  id: string;
  projectId: string;
  userId: string;
  originalFilename: string;
  storagePath: string;
  thumbnailUrl?: string;
  durationSeconds: number;
  width: number;
  height: number;
  fileSizeBytes: number;
  status: 'uploading' | 'processing' | 'ready' | 'failed';
  createdAt: string;
  updatedAt: string;
}

// types/clip.ts
export interface Clip {
  id: string;
  videoId: string;
  userId: string;
  name: string;
  startTime: number;
  endTime: number;
  durationSeconds: number;
  aspectRatio: '9:16' | '1:1' | '16:9';
  viralityScore?: number;
  status: 'draft' | 'rendering' | 'ready' | 'failed';
  storagePath?: string;
  thumbnailUrl?: string;
  isAiSuggested: boolean;
  style?: ClipStyle;
  createdAt: string;
  updatedAt: string;
}

export interface ClipStyle {
  captionEnabled: boolean;
  captionFont: string;
  captionSize: number;
  captionColor: string;
  captionBgColor?: string;
  captionPosition: 'top' | 'center' | 'bottom';
  captionAnimation?: string;
  brandLogoUrl?: string;
  brandLogoPosition?: string;
  overlayTemplate?: string;
  backgroundMusicUrl?: string;
  backgroundMusicVolume: number;
}

// types/transcription.ts
export interface Transcription {
  id: string;
  videoId: string;
  language: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  segments: TranscriptSegment[];
  createdAt: string;
}

export interface TranscriptSegment {
  id: string;
  startTime: number;
  endTime: number;
  text: string;
  confidence: number;
  speakerId?: number;
  words: TranscriptWord[];
}

export interface TranscriptWord {
  id: string;
  word: string;
  startTime: number;
  endTime: number;
  confidence: number;
}