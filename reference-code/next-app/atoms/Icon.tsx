'use client'

import { icons, type LucideProps } from 'lucide-react'
import { cn } from '@/lib/utils'

export type IconName = keyof typeof icons

export interface IconProps extends Omit<LucideProps, 'ref'> {
  name: IconName
  className?: string
}

export function Icon({ name, className, size = 24, ...props }: IconProps) {
  const LucideIcon = icons[name]

  if (!LucideIcon) {
    console.warn(`Icon "${name}" not found in Lucide icons`)
    return null
  }

  return (
    <LucideIcon
      size={size}
      className={cn('inline-block', className)}
      {...props}
    />
  )
}
