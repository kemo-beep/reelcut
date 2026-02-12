import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getProfile, updateProfile, changePassword, uploadAvatar } from '../../../lib/api/users'
import { useAuthStore } from '../../../stores/authStore'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import { toast } from 'sonner'
import { ApiError } from '../../../types'

export const Route = createFileRoute('/dashboard/settings/profile')({
  component: ProfilePage,
})

function ProfilePage() {
  const queryClient = useQueryClient()
  const setAuth = useAuthStore((s) => s.setAuth)
  const user = useAuthStore((s) => s.user)
  const [fullName, setFullName] = useState('')
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [avatarFile, setAvatarFile] = useState<File | null>(null)

  const { data, isLoading, error } = useQuery({
    queryKey: ['profile'],
    queryFn: getProfile,
  })
  const profileUser = data?.user
  useEffect(() => {
    if (profileUser?.full_name != null) setFullName(profileUser.full_name)
  }, [profileUser?.full_name])

  const updateMutation = useMutation({
    mutationFn: (body: { full_name?: string | null }) => updateProfile(body),
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      const token = useAuthStore.getState().getAccessToken()
      const refresh = useAuthStore.getState().refreshToken
      if (res.user && token) setAuth(res.user, token, refresh ?? '')
      toast.success('Profile updated')
    },
    onError: (e) => toast.error(e instanceof ApiError ? e.message : 'Update failed'),
  })
  const passwordMutation = useMutation({
    mutationFn: () => changePassword(currentPassword, newPassword),
    onSuccess: () => {
      setCurrentPassword('')
      setNewPassword('')
      toast.success('Password changed')
    },
    onError: (e) => toast.error(e instanceof ApiError ? e.message : 'Password change failed'),
  })
  const avatarMutation = useMutation({
    mutationFn: (file: File) => uploadAvatar(file),
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      const token = useAuthStore.getState().getAccessToken()
      const refresh = useAuthStore.getState().refreshToken
      if (res.user && token) setAuth(res.user, token, refresh ?? '')
      setAvatarFile(null)
      toast.success('Avatar updated')
    },
    onError: () => toast.error('Avatar upload failed'),
  })

  if (isLoading) return <p className="text-slate-400">Loading profile...</p>
  if (error) return <p className="text-red-400">Failed to load profile.</p>

  const u = profileUser ?? user
  if (!u) return null

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold text-[var(--app-fg)]">Profile</h1>

      <section className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6 max-w-md space-y-4">
        <h2 className="font-semibold text-[var(--app-fg)]">Account</h2>
        <p className="text-caption text-[var(--app-fg-muted)]">{u.email}</p>
        <div className="space-y-2">
          <Label htmlFor="full_name">Display name</Label>
          <Input
            id="full_name"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            className="bg-[var(--app-bg)] text-[var(--app-fg)]"
          />
          <Button
            size="sm"
            onClick={() => updateMutation.mutate({ full_name: fullName || null })}
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? 'Saving…' : 'Save name'}
          </Button>
        </div>
      </section>

      <section className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6 max-w-md space-y-4">
        <h2 className="font-semibold text-[var(--app-fg)]">Avatar</h2>
        <input
          type="file"
          accept="image/jpeg,image/png,image/webp"
          onChange={(e) => setAvatarFile(e.target.files?.[0] ?? null)}
          className="text-sm text-[var(--app-fg)]"
        />
        {avatarFile && (
          <Button
            size="sm"
            onClick={() => avatarMutation.mutate(avatarFile)}
            disabled={avatarMutation.isPending}
          >
            {avatarMutation.isPending ? 'Uploading…' : 'Upload avatar'}
          </Button>
        )}
      </section>

      <section className="rounded-xl border border-[var(--app-border)] bg-[var(--app-bg-raised)] p-6 max-w-md space-y-4">
        <h2 className="font-semibold text-[var(--app-fg)]">Change password</h2>
        <div className="space-y-2">
          <Label htmlFor="current_password">Current password</Label>
          <Input
            id="current_password"
            type="password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            className="bg-[var(--app-bg)] text-[var(--app-fg)]"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="new_password">New password</Label>
          <Input
            id="new_password"
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="bg-[var(--app-bg)] text-[var(--app-fg)]"
          />
        </div>
        <Button
          size="sm"
          onClick={() => passwordMutation.mutate()}
          disabled={passwordMutation.isPending || !currentPassword || !newPassword}
        >
          {passwordMutation.isPending ? 'Changing…' : 'Change password'}
        </Button>
      </section>
    </div>
  )
}
