import { z } from 'zod'

const envSchema = z.object({
  // Server-side only (never sent to client)
  NODE_ENV: z
    .enum(['development', 'test', 'production'])
    .default('development'),

  // Client-side (prefixed with NEXT_PUBLIC_)
  NEXT_PUBLIC_APP_URL: z.string().url().default('http://localhost:3001'),
  NEXT_PUBLIC_API_URL: z.string().url().default('http://localhost:8080'),

  // Add more as needed in future milestones
})

const envParsed = envSchema.safeParse(process.env)

if (!envParsed.success) {
  console.error(
    '❌ Invalid environment variables:',
    envParsed.error.flatten().fieldErrors
  )
  throw new Error('Invalid environment variables')
}

export const env = envParsed.data

// Usage example:
// import { env } from '@/lib/env'
// console.log(env.NEXT_PUBLIC_APP_URL)
