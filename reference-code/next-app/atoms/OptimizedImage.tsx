import Image, { ImageProps } from 'next/image'

/**
 * OptimizedImage component wraps Next.js Image with performance best practices
 * - Lazy loading by default
 * - Blur placeholder during load
 * - Automatic AVIF/WebP format selection
 */
export function OptimizedImage({ alt, ...props }: ImageProps) {
  return (
    <Image
      {...props}
      alt={alt}
      loading="lazy"
      placeholder="blur"
      blurDataURL="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2VlZSIvPjwvc3ZnPg=="
    />
  )
}
