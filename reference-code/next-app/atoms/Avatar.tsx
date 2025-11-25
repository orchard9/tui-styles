'use client'

import Image from 'next/image'
import { User } from 'lucide-react'
import { cn } from '@/lib/utils'

export interface AvatarProps {
  src?: string
  alt?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  className?: string
  fallback?: string
}

const sizeClasses = {
  sm: 'h-8 w-8 text-xs',
  md: 'h-10 w-10 text-sm',
  lg: 'h-14 w-14 text-base',
  xl: 'h-20 w-20 text-lg',
}

const sizePixels = {
  sm: 32,
  md: 40,
  lg: 56,
  xl: 80,
}

export function Avatar({
  src,
  alt = 'Avatar',
  size = 'md',
  className,
  fallback,
}: AvatarProps) {
  if (src) {
    return (
      <div
        className={cn(
          'relative rounded-full overflow-hidden',
          sizeClasses[size],
          className
        )}
      >
        <Image
          src={src}
          alt={alt}
          width={sizePixels[size]}
          height={sizePixels[size]}
          className="object-cover"
        />
      </div>
    )
  }

  // Fallback to initials or icon
  return (
    <div
      className={cn(
        'flex items-center justify-center rounded-full bg-muted text-muted-foreground',
        sizeClasses[size],
        className
      )}
    >
      {fallback ? (
        <span className="font-medium">{fallback}</span>
      ) : (
        <User className="h-1/2 w-1/2" />
      )}
    </div>
  )
}
