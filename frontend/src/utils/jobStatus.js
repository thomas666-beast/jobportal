export function statusLabel(status) {
  switch (status) {
    case 'draft':          return 'Draft'
    case 'pending_review': return 'Pending review'
    case 'published':      return 'Published'
    case 'closed':         return 'Closed'
    case 'archived':       return 'Archived'
    default:               return status
  }
}

export function statusVariant(status) {
  switch (status) {
    case 'draft':          return 'neutral'
    case 'pending_review': return 'warning'
    case 'published':      return 'success'
    case 'closed':         return 'neutral'
    case 'archived':       return 'dark'
    default:               return 'neutral'
  }
}

export function applicationStatusLabel(status) {
  switch (status) {
    case 'submitted':   return 'Submitted'
    case 'viewed':      return 'Viewed'
    case 'shortlisted': return 'Shortlisted'
    case 'rejected':    return 'Rejected'
    case 'hired':       return 'Hired'
    case 'withdrawn':   return 'Withdrawn'
    default:            return status
  }
}

export function applicationStatusVariant(status) {
  switch (status) {
    case 'submitted':   return 'neutral'
    case 'viewed':      return 'info'
    case 'shortlisted': return 'primary'
    case 'rejected':    return 'danger'
    case 'hired':       return 'success'
    case 'withdrawn':   return 'neutral'
    default:            return 'neutral'
  }
}

// Kept for backwards compat with any old imports.
export function statusClasses(status) {
  const map = {
    neutral: 'bg-slate-100 text-slate-700',
    warning: 'bg-amber-100 text-amber-700',
    success: 'bg-emerald-100 text-emerald-700',
    dark:    'bg-slate-800 text-white',
    info:    'bg-blue-100 text-blue-700',
    primary: 'bg-primary-100 text-primary-700',
    danger:  'bg-red-100 text-red-700'
  }
  return map[statusVariant(status)] || map.neutral
}

export function applicationStatusClasses(status) {
  const map = {
    neutral: 'bg-slate-100 text-slate-700',
    info:    'bg-blue-100 text-blue-700',
    primary: 'bg-primary-100 text-primary-700',
    danger:  'bg-red-100 text-red-700',
    success: 'bg-emerald-100 text-emerald-700'
  }
  return map[applicationStatusVariant(status)] || map.neutral
}

// Human-friendly relative date: "just now", "2h ago", "3d ago", "Mar 5"
export function timeAgo(input) {
  if (!input) return ''
  const d = new Date(input)
  const now = Date.now()
  const diff = Math.floor((now - d.getTime()) / 1000)

  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  if (diff < 86400 * 7) return `${Math.floor(diff / 86400)}d ago`
  if (diff < 86400 * 30) return `${Math.floor(diff / 86400 / 7)}w ago`

  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

// Employment type label
export function employmentLabel(et) {
  if (et === 'full_time') return 'Full-time'
  if (et === 'part_time') return 'Part-time'
  if (et === 'contract') return 'Contract'
  if (et === 'internship') return 'Internship'
  return et
}

// Remote type label
export function remoteLabel(rt) {
  if (rt === 'remote') return 'Remote'
  if (rt === 'hybrid') return 'Hybrid'
  if (rt === 'on_site') return 'On-site'
  return rt
}

// Money format — returns null if both endpoints are null
export function formatSalary(min, max, currency) {
  if (min == null && max == null) return null
  const fmt = (n) => (n == null ? '?' : n.toLocaleString())
  const cur = currency ? ` ${currency}` : ''
  return `${fmt(min)} – ${fmt(max)}${cur}`
}

// Returns true if the published date is within the last 24 hours.
export function isNewJob(publishedAt) {
  if (!publishedAt) return false
  const d = new Date(publishedAt)
  const diff = Date.now() - d.getTime()
  return diff >= 0 && diff < 24 * 60 * 60 * 1000
}
