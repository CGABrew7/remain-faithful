import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import AdminDashboard from '../admin-dashboard'
import { previewHealth, previewOverview } from '../preview-data'

export const metadata: Metadata = {
  title: 'Owner dashboard sample',
  robots: { index: false, follow: false, nocache: true },
}

export default function AdminPreviewPage() {
  if (process.env.NODE_ENV === 'production') notFound()
  const ranges = {
    '7d': previewOverview('7d'),
    '30d': previewOverview('30d'),
    '90d': previewOverview('90d'),
    all: previewOverview('all'),
  }
  return <AdminDashboard preview ranges={ranges} healthPreview={previewHealth} />
}
