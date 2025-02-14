import type { Meta, StoryObj } from '@storybook/react'
import { Button } from '../button/button'
import { Toaster, toast } from './toast'

const meta: Meta = {
  title: 'Basic Components/Toast',
  component: Toaster,
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <Button
      variant="outline"
      onClick={() => {
        toast('Event has been created', {
          description: 'Sunday, December 03, 2023 at 9:00 AM',
        })
      }}
    >
      Show Toast
    </Button>
  ),
}

export const ToastShowcase: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <Button
        variant="outline"
        onClick={() => {
          toast.success('Successfully saved!', {
            description: 'Your changes have been saved.',
          })
        }}
      >
        Success Toast
      </Button>

      <Button
        variant="outline"
        onClick={() => {
          toast.error('Error occurred', {
            description: 'There was a problem with your request.',
          })
        }}
      >
        Error Toast
      </Button>

      <Button
        variant="outline"
        onClick={() => {
          toast.promise(new Promise((resolve) => setTimeout(resolve, 2000)), {
            loading: 'Loading...',
            success: 'Successfully loaded',
            error: 'Error loading data',
          })
        }}
      >
        Promise Toast
      </Button>
    </div>
  ),
}
