import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Download, Loader2 } from 'lucide-react'
import { Button } from '../ui/button'
import { Label } from '../ui/label'
import { ProgressBar } from '../processing/ProgressBar'
import { getExportPresets } from '../../lib/api/config'
import { cn } from '../../lib/utils'

export interface ExportPanelProps {
  aspectRatios?: string[]
  onRender?: (options: { preset?: string; aspectRatio?: string }) => void
  renderStatus?: 'idle' | 'rendering' | 'ready' | 'failed'
  renderProgress?: number
  downloadUrl?: string | null
  className?: string
}

const DEFAULT_ASPECT_RATIOS = ['9:16', '1:1', '16:9']
const CUSTOM_VALUE = ''

export function ExportPanel({
  aspectRatios = DEFAULT_ASPECT_RATIOS,
  onRender,
  renderStatus = 'idle',
  renderProgress = 0,
  downloadUrl,
  className,
}: ExportPanelProps) {
  const [presetOrCustom, setPresetOrCustom] = useState<string>(CUSTOM_VALUE)
  const [aspectRatio, setAspectRatio] = useState(aspectRatios[0])

  const { data: presetsData } = useQuery({
    queryKey: ['config', 'export-presets'],
    queryFn: getExportPresets,
    staleTime: 5 * 60 * 1000,
  })
  const presets = presetsData?.presets ?? []
  const isCustom = presetOrCustom === CUSTOM_VALUE

  const handleRender = () => {
    if (presetOrCustom && presetOrCustom !== CUSTOM_VALUE) {
      onRender?.({ preset: presetOrCustom })
    } else {
      onRender?.({ aspectRatio })
    }
  }

  return (
    <div className={cn('space-y-4', className)}>
      <h3 className="font-semibold text-[var(--app-fg)]">Export</h3>
      <div className="space-y-2">
        <Label>Platform / preset</Label>
        <select
          value={presetOrCustom}
          onChange={(e) => setPresetOrCustom(e.target.value)}
          className="h-9 w-full rounded-md border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-1 text-sm text-[var(--app-fg)]"
        >
          <option value={CUSTOM_VALUE}>Custom</option>
          {presets.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name} ({p.aspect_ratio})
            </option>
          ))}
        </select>
      </div>
      {isCustom && (
        <div className="space-y-2">
          <Label>Aspect ratio</Label>
          <select
            value={aspectRatio}
            onChange={(e) => setAspectRatio(e.target.value)}
            className="h-9 w-full rounded-md border border-[var(--app-border)] bg-[var(--app-bg)] px-3 py-1 text-sm text-[var(--app-fg)]"
          >
            {aspectRatios.map((ar) => (
              <option key={ar} value={ar}>
                {ar}
              </option>
            ))}
          </select>
        </div>
      )}
      {renderStatus === 'idle' && (
        <Button
          className="w-full bg-[var(--app-accent)] text-[#0a0a0b] hover:bg-[var(--app-accent-hover)]"
          onClick={handleRender}
        >
          Render clip
        </Button>
      )}
      {renderStatus === 'rendering' && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-[var(--app-fg-muted)]">
            <Loader2 size={16} className="animate-spin" />
            Rendering…
          </div>
          <ProgressBar value={renderProgress} showLabel />
        </div>
      )}
      {(renderStatus === 'ready' || downloadUrl) && (
        <Button asChild className="w-full" variant="outline">
          <a href={downloadUrl ?? '#'} download target="_blank" rel="noreferrer">
            <Download size={18} className="mr-2" />
            Download
          </a>
        </Button>
      )}
      {renderStatus === 'failed' && (
        <p className="text-sm text-[var(--app-destructive)]">Render failed. Try again.</p>
      )}
    </div>
  )
}
