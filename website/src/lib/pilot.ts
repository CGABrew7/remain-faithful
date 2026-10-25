// Turns a church pilot sign-up from /partners into the subject and message
// the contact endpoint emails to the support inbox.

const MAX_FIELD = 200

type PilotFields = {
  name?: unknown
  email?: unknown
  church?: unknown
  role?: unknown
  numGroups?: unknown
  referral?: unknown
}

function clean(value: unknown) {
  if (typeof value !== 'string') return ''
  return value.replace(/[\u0000-\u001f\u007f]/g, ' ').replace(/\s+/g, ' ').trim().slice(0, MAX_FIELD)
}

export function pilotContact(
  body: PilotFields,
): { ok: true; contact: { name: string; email: string; subject: string; message: string } } | { ok: false; error: string } {
  const name = clean(body.name)
  const email = clean(body.email)
  const church = clean(body.church)
  const role = clean(body.role)
  const numGroups = clean(body.numGroups)
  const referral = clean(body.referral)

  if (!name || !email || !church || !role || !numGroups) {
    return { ok: false, error: 'Please fill in your name, church, email, role, and number of groups.' }
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return { ok: false, error: 'Please enter a valid email address.' }
  }

  const message = [
    'New church pilot sign-up from the Partners page.',
    '',
    `Name: ${name}`,
    `Church or organization: ${church}`,
    `Email: ${email}`,
    `Role: ${role}`,
    `Estimated groups: ${numGroups}`,
    `Heard about us: ${referral || 'Not given'}`,
  ].join('\n')

  return { ok: true, contact: { name, email, subject: `Church pilot: ${church}`, message } }
}
