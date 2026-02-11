import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { getProfile } from '../../../lib/api/users'

export const Route = createFileRoute('/dashboard/settings/profile')({
  component: ProfilePage,
})

function ProfilePage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['profile'],
    queryFn: getProfile,
  })

  if (isLoading) return <p className="text-slate-400">Loading profile...</p>
  if (error) return <p className="text-red-400">Failed to load profile.</p>

  const user = data?.user
  if (!user) return null

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-white">Profile</h1>
      <div className="rounded-lg border border-slate-700 bg-slate-800/50 p-4 max-w-md space-y-2">
        <p className="text-slate-400 text-sm">Email</p>
        <p className="text-white">{user.email}</p>
        {user.full_name && (
          <>
            <p className="text-slate-400 text-sm mt-3">Name</p>
            <p className="text-white">{user.full_name}</p>
          </>
        )}
      </div>
    </div>
  )
}
