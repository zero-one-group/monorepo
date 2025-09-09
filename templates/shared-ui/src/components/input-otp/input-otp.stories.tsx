import type { Meta, StoryObj } from '@storybook/react'
import { InputOTP, InputOTPGroup, InputOTPSeparator, InputOTPSlot } from './input-otp'
import type { InputOtpVariants } from './input-otp.css'

const sizeOptions: NonNullable<InputOtpVariants['size']>[] = ['sm', 'default', 'lg']

const meta: Meta = {
  title: 'Basic Components/InputOTP',
  component: InputOTP,
  argTypes: {
    size: {
      control: 'radio',
      options: sizeOptions,
      table: {
        defaultValue: { summary: 'default' },
        type: { summary: 'InputOtpVariants["size"]' },
      },
    },
  },
} satisfies Meta<typeof InputOTP>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    maxLength: 6,
    children: (
      <InputOTPGroup>
        <InputOTPSlot index={0} />
        <InputOTPSlot index={1} />
        <InputOTPSlot index={2} />
        <InputOTPSeparator />
        <InputOTPSlot index={3} />
        <InputOTPSlot index={4} />
        <InputOTPSlot index={5} />
      </InputOTPGroup>
    ),
  },
  render: ({ size, ...args }) => (
    <InputOTP {...args}>
      <InputOTPGroup>
        <InputOTPSlot index={0} size={size} />
        <InputOTPSlot index={1} size={size} />
        <InputOTPSlot index={2} size={size} />
        <InputOTPSeparator size={size} />
        <InputOTPSlot index={3} size={size} />
        <InputOTPSlot index={4} size={size} />
        <InputOTPSlot index={5} size={size} />
      </InputOTPGroup>
    </InputOTP>
  ),
}

export const Sizes: Story = {
  args: {
    maxLength: 6,
    children: (
      <InputOTPGroup>
        <InputOTPSlot index={0} />
        <InputOTPSlot index={1} />
        <InputOTPSlot index={2} />
        <InputOTPSeparator />
        <InputOTPSlot index={3} />
        <InputOTPSlot index={4} />
        <InputOTPSlot index={5} />
      </InputOTPGroup>
    ),
  },
  parameters: {
    controls: { exclude: ['size'] },
  },
  render: () => (
    <div className="space-y-4">
      <InputOTP maxLength={6} size="sm">
        <InputOTPGroup>
          <InputOTPSlot index={0} size="sm" />
          <InputOTPSlot index={1} size="sm" />
          <InputOTPSlot index={2} size="sm" />
          <InputOTPSeparator size="sm" />
          <InputOTPSlot index={3} size="sm" />
          <InputOTPSlot index={4} size="sm" />
          <InputOTPSlot index={5} size="sm" />
        </InputOTPGroup>
      </InputOTP>

      <InputOTP maxLength={6}>
        <InputOTPGroup>
          <InputOTPSlot index={0} />
          <InputOTPSlot index={1} />
          <InputOTPSlot index={2} />
          <InputOTPSeparator />
          <InputOTPSlot index={3} />
          <InputOTPSlot index={4} />
          <InputOTPSlot index={5} />
        </InputOTPGroup>
      </InputOTP>

      <InputOTP maxLength={6} size="lg">
        <InputOTPGroup>
          <InputOTPSlot index={0} size="lg" />
          <InputOTPSlot index={1} size="lg" />
          <InputOTPSlot index={2} size="lg" />
          <InputOTPSeparator size="lg" />
          <InputOTPSlot index={3} size="lg" />
          <InputOTPSlot index={4} size="lg" />
          <InputOTPSlot index={5} size="lg" />
        </InputOTPGroup>
      </InputOTP>
    </div>
  ),
}
