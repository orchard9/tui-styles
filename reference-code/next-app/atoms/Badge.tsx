import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const badgeVariants = cva(
  'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-all focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
  {
    variants: {
      variant: {
        default:
          'border-transparent bg-primary text-primary-foreground shadow-sm hover:glow-primary-sm',
        secondary:
          'border-transparent bg-secondary text-secondary-foreground shadow-sm hover:glow-primary-sm',
        destructive:
          'border-transparent bg-destructive text-destructive-foreground shadow-sm hover:glow-error-sm',
        outline: 'text-foreground border-border',
        success:
          'border-transparent bg-[hsl(var(--semantic-success))] text-[hsl(var(--success-foreground))] shadow-sm hover:glow-success-sm',
        warning:
          'border-transparent bg-[hsl(var(--semantic-warning))] text-[hsl(var(--warning-foreground))] shadow-sm',
        info: 'border-transparent bg-[hsl(var(--semantic-info))] text-[hsl(var(--info-foreground))] shadow-sm',
        error:
          'border-transparent bg-[hsl(var(--semantic-error))] text-white shadow-sm hover:glow-error-sm',
      },
      size: {
        sm: 'px-2 py-0.5 text-xs',
        default: 'px-2.5 py-0.5 text-xs',
        lg: 'px-3 py-1 text-sm',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, size, ...props }: BadgeProps) {
  return (
    <div
      className={cn(badgeVariants({ variant, size }), className)}
      {...props}
    />
  )
}
