import { Label } from '../ui/label'
import { Input } from '../ui/input'
import { Switch } from '../ui/switch'
import { cn } from '../../lib/utils'

/** Default allowed caption fonts (fallback if API unavailable); must match backend AllowedCaptionFonts. */
export const DEFAULT_CAPTION_FONTS = ['Arial', 'Inter', 'Montserrat', 'Open Sans', 'Roboto', 'Helvetica']

export const CAPTION_ANIMATION_PRESETS = [
  { value: '', label: 'None' },
  { value: 'fade', label: 'Fade' },
  { value: 'slide-up', label: 'Slide up' },
  { value: 'slide-down', label: 'Slide down' },
] as const

export interface CaptionEditorProps {
  enabled: boolean
  text?: string
  fontSize?: number
  color?: string
  position?: 'top' | 'center' | 'bottom'
  font?: string
  bgColor?: string
  animation?: string
  fontOptions?: string[]
  onEnabledChange?: (enabled: boolean) => void
  onTextChange?: (text: string) => void
  onFontSizeChange?: (size: number) => void
  onColorChange?: (color: string) => void
  onPositionChange?: (position: 'top' | 'center' | 'bottom') => void
  onFontChange?: (font: string) => void
  onBgColorChange?: (color: string) => void
  onAnimationChange?: (animation: string) => void
  className?: string
}

export function CaptionEditor({
  enabled,
  text = '',
  fontSize = 48,
  color = '#FFFFFF',
  position = 'bottom',
  font = 'Inter',
  bgColor = '',
  animation = '',
  fontOptions = DEFAULT_CAPTION_FONTS,
  onEnabledChange,
  onTextChange,
  onFontSizeChange,
  onColorChange,
  onPositionChange,
  onFontChange,
  onBgColorChange,
  onAnimationChange,
  className,
}: CaptionEditorProps) {
  return (
    <div className={cn('space-y-4', className)}>
      <div className="flex items-center justify-between">
        <Label>Captions</Label>
        <Switch
          checked={enabled}
          onCheckedChange={onEnabledChange}
          aria-label="Toggle captions"
        />
      </div>
      {enabled && (
        <>
          {onFontChange && (
            <div className="space-y-2">
              <Label>Font</Label>
              <select
                value={font}
                onChange={(e) => onFontChange(e.target.value)}
                className="h-9 w-full rounded-md border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-1 text-sm text-[var(--app-fg)]"
              >
                {fontOptions.map((f) => (
                  <option key={f} value={f}>
                    {f}
                  </option>
                ))}
              </select>
            </div>
          )}
          {onTextChange && (
            <div className="space-y-2">
              <Label>Preview text</Label>
              <Input
                value={text}
                onChange={(e) => onTextChange(e.target.value)}
                className="bg-[var(--app-bg)] text-[var(--app-fg)]"
              />
            </div>
          )}
          {onFontSizeChange && (
            <div className="space-y-2">
              <Label>Font size</Label>
              <Input
                type="number"
                min={12}
                max={96}
                value={fontSize}
                onChange={(e) => onFontSizeChange(Number(e.target.value))}
                className="bg-[var(--app-bg)] text-[var(--app-fg)]"
              />
            </div>
          )}
          {onColorChange && (
            <div className="space-y-2">
              <Label>Color</Label>
              <div className="flex gap-2">
                <input
                  type="color"
                  value={color}
                  onChange={(e) => onColorChange(e.target.value)}
                  className="h-10 w-14 cursor-pointer rounded border border-[var(--app-border)]"
                />
                <Input
                  value={color}
                  onChange={(e) => onColorChange(e.target.value)}
                  className="flex-1 bg-[var(--app-bg)] text-[var(--app-fg)] font-mono"
                />
              </div>
            </div>
          )}
          {onBgColorChange && (
            <div className="space-y-2">
              <Label>Background color</Label>
              <div className="flex gap-2">
                <input
                  type="color"
                  value={bgColor && bgColor.length >= 7 ? bgColor.slice(0, 7) : '#000000'}
                  onChange={(e) => {
                    const hex = e.target.value
                    onBgColorChange(bgColor.length === 9 ? hex + bgColor.slice(7) : hex)
                  }}
                  className="h-10 w-14 cursor-pointer rounded border border-[var(--app-border)]"
                />
                <Input
                  value={bgColor}
                  onChange={(e) => onBgColorChange(e.target.value)}
                  placeholder="#000000 or #000000CC"
                  className="flex-1 bg-[var(--app-bg)] text-[var(--app-fg)] font-mono"
                />
              </div>
            </div>
          )}
          {onPositionChange && (
            <div className="space-y-2">
              <Label>Position</Label>
              <select
                value={position}
                onChange={(e) => onPositionChange(e.target.value as 'top' | 'center' | 'bottom')}
                className="h-9 w-full rounded-md border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-1 text-sm text-[var(--app-fg)]"
              >
                <option value="top">Top</option>
                <option value="center">Center</option>
                <option value="bottom">Bottom</option>
              </select>
            </div>
          )}
          {onAnimationChange && (
            <div className="space-y-2">
              <Label>Animation</Label>
              <select
                value={animation}
                onChange={(e) => onAnimationChange(e.target.value)}
                className="h-9 w-full rounded-md border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-1 text-sm text-[var(--app-fg)]"
              >
                {CAPTION_ANIMATION_PRESETS.map(({ value, label }) => (
                  <option key={value || 'none'} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
          )}
        </>
      )}
    </div>
  )
}
