import { forwardRef } from 'react'
import { Icon, type IconName } from '@/components/atoms/Icon'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

export interface InputWithIconProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  icon?: IconName
  iconPosition?: 'left' | 'right'
  error?: boolean
}

function IconWrapper({
  icon,
  position,
  error,
}: {
  icon: IconName
  position: 'left' | 'right'
  error?: boolean
}) {
  const positionClasses = position === 'left' ? 'left-0 pl-3' : 'right-0 pr-3'
  const iconColor = error ? 'text-destructive' : 'text-muted-foreground'

  return (
    <div
      className={cn(
        'pointer-events-none absolute inset-y-0 flex items-center',
        positionClasses
      )}
    >
      <Icon name={icon} className={cn('h-5 w-5', iconColor)} />
    </div>
  )
}

const InputWithIcon = forwardRef<HTMLInputElement, InputWithIconProps>(
  ({ className, icon, iconPosition = 'left', error, ...props }, ref) => {
    const inputPaddingClass =
      icon && (iconPosition === 'left' ? 'pl-10' : 'pr-10')
    const errorClass =
      error && 'border-destructive focus-visible:ring-destructive'

    return (
      <div className="relative">
        {icon && (
          <IconWrapper icon={icon} position={iconPosition} error={error} />
        )}
        <Input
          ref={ref}
          className={cn(inputPaddingClass, errorClass, className)}
          aria-invalid={error}
          {...props}
        />
      </div>
    )
  }
)

InputWithIcon.displayName = 'InputWithIcon'

export { InputWithIcon }
