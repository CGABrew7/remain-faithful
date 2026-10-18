import type { Metadata } from 'next'
import AdminDashboard from './admin-dashboard'

export const metadata: Metadata = {
  title: 'Owner dashboard',
  robots: {
    index: false,
    follow: false,
    nocache: true,
    googleBot: { index: false, follow: false },
  },
}

export default function AdminPage() {
  return <AdminDashboard />
}
