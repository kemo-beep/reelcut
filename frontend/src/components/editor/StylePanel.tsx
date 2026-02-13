import { useQuery } from '@tanstack/react-query'
import { CaptionEditor, DEFAULT_CAPTION_FONTS } from './CaptionEditor'
import { getCaptionFonts } from '../../lib/api/config'
import { listTranscriptionsByVideo } from '../../lib/api/transcriptions'
import { Label } from '../ui/label'
import { Slider } from '../ui/slider'
import { cn } from '../../lib/utils'
import type { ClipStyle } from '../../types'

const LANGUAGE_NAMES: Record<string, string> = {
  en: 'English',
  es: 'Spanish',
  fr: 'French',
  de: 'German',
  it: 'Italian',
  pt: 'Portuguese',
  ja: 'Japanese',
  zh: 'Chinese',
  ko: 'Korean',
  ar: 'Arabic',
  hi: 'Hindi',
}

function languageLabel(code: string): string {
  return LANGUAGE_NAMES[code] ?? code
}

export interface StylePanelProps {
  style: Partial<ClipStyle>
  onChange: (updates: Partial<ClipStyle>) => void
  /** When set, fetches transcriptions for caption language selection. */
  videoId?: string
  className?: string
}

export function StylePanel({ style, onChange, videoId, className }: StylePanelProps) {
  const { data: fontsData } = useQuery({
    queryKey: ['config', 'caption-fonts'],
    queryFn: getCaptionFonts,
    staleTime: 5 * 60 * 1000,
  })
  const fontOptions = fontsData?.fonts?.length ? fontsData.fonts : DEFAULT_CAPTION_FONTS

  const { data: listData } = useQuery({
    queryKey: ['transcriptions-list', videoId],
    queryFn: () => listTranscriptionsByVideo(videoId!),
    enabled: !!videoId,
    staleTime: 60 * 1000,
  })
  const transcriptions = listData?.transcriptions ?? []
  const hasMultipleLanguages = transcriptions.length > 1

  return (
    <div className={cn('space-y-6', className)}>
      <h3 className="font-semibold text-[var(--app-fg)]">Style</h3>
      <CaptionEditor
        enabled={style.caption_enabled ?? true}
        fontSize={style.caption_size}
        color={style.caption_color}
        position={style.caption_position as 'top' | 'center' | 'bottom'}
        font={style.caption_font ?? 'Inter'}
        bgColor={style.caption_bg_color ?? ''}
        animation={style.caption_animation ?? ''}
        fontOptions={fontOptions}
        onEnabledChange={(v) => onChange({ caption_enabled: v })}
        onFontSizeChange={(v) => onChange({ caption_size: v })}
        onColorChange={(v) => onChange({ caption_color: v })}
        onPositionChange={(v) => onChange({ caption_position: v })}
        onFontChange={(v) => onChange({ caption_font: v })}
        onBgColorChange={(v) => onChange({ caption_bg_color: v || undefined })}
        onAnimationChange={(v) => onChange({ caption_animation: v || undefined })}
      />
      {hasMultipleLanguages && (
        <div className="space-y-2">
          <Label>Caption language</Label>
          <select
            className="w-full rounded-lg border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-2 text-sm text-[var(--app-fg)]"
            value={style.caption_language ?? ''}
            onChange={(e) => onChange({ caption_language: e.target.value || undefined })}
          >
            <option value="">Default (source)</option>
            {transcriptions.map((t) => (
              <option key={t.id} value={t.language}>
                {languageLabel(t.language)}
              </option>
            ))}
          </select>
        </div>
      )}
      <div className="space-y-2">
        <Label>Background music volume</Label>
        <Slider
          value={[style.background_music_volume ?? 0.3]}
          onValueChange={(v) => onChange({ background_music_volume: (Array.isArray(v) ? v[0] : v) ?? 0.3 })}
          min={0}
          max={1}
          step={0.1}
        />
      </div>
    </div>
  )
}
