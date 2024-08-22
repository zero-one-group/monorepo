import { EyeClosedIcon, EyeOpenIcon } from '@radix-ui/react-icons'
import * as React from 'react'
import { cn } from '#/utils'

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    const [isVisible, setIsVisible] = React.useState(false)
    if (type === 'password') {
      return (
        <div className="relative rounded-md shadow-sm">
          <input
            id="price"
            name="price"
            type={isVisible ? 'text' : 'password'}
            placeholder="0.00"
            aria-describedby="price-currency"
            className={cn(
              'flex h-9 w-full rounded-md border border-input bg-transparent py-1 pr-12 pl-3 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:font-medium file:text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
              className
            )}
            ref={ref}
            {...props}
          />
          <div className="absolute inset-y-0 right-0 flex items-center">
            <button
              type="button"
              className={cn(
                'inline-flex h-9 w-9 items-center justify-center whitespace-nowrap rounded-md',
                'font-medium text-sm transition-colors hover:bg-transparent hover:text-accent-foreground',
                'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring',
                'disabled:pointer-events-none disabled:opacity-50'
              )}
              onClick={() => setIsVisible(!isVisible)}
            >
              {isVisible ? (
                <EyeOpenIcon className="h-3.5 w-3.5 fill-primary" />
              ) : (
                <EyeClosedIcon className="h-3.5 w-3.5 fill-primary" />
              )}
            </button>
          </div>
        </div>
      )
    }

    return (
      <input
        type={type}
        className={cn(
          'flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:font-medium file:text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)

export { Input }
