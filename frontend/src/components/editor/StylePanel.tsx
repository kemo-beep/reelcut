import { CaptionEditor } from './CaptionEditor'
import { Label } from '../ui/label'
import { Input } from '../ui/input'
import { Slider } from '../ui/slider'
import { cn } from '../../lib/utils'
import type { ClipStyle } from '../../types'

export interface StylePanelProps {
  style: Partial<ClipStyle>
  onChange: (updates: Partial<ClipStyle>) => void
  className?: string
}

export function StylePanel({ style, onChange, className }: StylePanelProps) {
  return (
    <div className={cn('space-y-6', className)}>
      <h3 className="font-semibold text-[var(--app-fg)]">Style</h3>
      <CaptionEditor
        enabled={style.caption_enabled ?? true}
        fontSize={style.caption_size}
        color={style.caption_color}
        position={style.caption_position as 'top' | 'center' | 'bottom'}
        onEnabledChange={(v) => onChange({ caption_enabled: v })}
        onFontSizeChange={(v) => onChange({ caption_size: v })}
        onColorChange={(v) => onChange({ caption_color: v })}
        onPositionChange={(v) => onChange({ caption_position: v })}
      />
      <div className="space-y-2">
        <Label>Font</Label>
        <Input
          value={style.caption_font ?? 'Inter'}
          onChange={(e) => onChange({ caption_font: e.target.value })}
          className="bg-[var(--app-bg)] text-[var(--app-fg)]"
        />
      </div>
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
