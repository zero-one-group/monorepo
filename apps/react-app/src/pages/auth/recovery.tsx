import { type SubmitHandler, focus, reset, useForm, valiForm } from '@modular-forms/react'
import { Button, Input, Label, cn } from '@myorg/shared-ui'
import { toast } from 'sonner'
import * as v from 'valibot'
import { Link } from '#/components/link'
import logger from '#/utils/logger'

const ForgotPasswordSchema = v.object({
  email: v.pipe(
    v.string(),
    v.nonEmpty('Please enter your email.'),
    v.email('The email address is badly formatted.')
  ),
})

type ForgotPasswordForm = v.InferInput<typeof ForgotPasswordSchema>

export default function ForgotPasswordPage() {
  const [forgotPasswordForm, { Form, Field }] = useForm<ForgotPasswordForm>({
    validate: valiForm(ForgotPasswordSchema),
    validateOn: 'change',
  })

  const handleSubmit: SubmitHandler<ForgotPasswordForm> = async (values, event) => {
    event.preventDefault()

    try {
      logger.info(values.email)
      toast.success("Check your inbox! We've sent you an email with a link to reset your password.")
    } catch (error) {
      logger.error('[ERROR] handleSubmit', error)
      if (error instanceof Error) {
        toast.error(error.message)
      }
      focus(forgotPasswordForm, 'email')
    } finally {
      reset(forgotPasswordForm)
    }
  }

  return (
    <div className="mx-auto w-full max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
      <div className="mx-auto w-full max-w-sm md:max-w-lg lg:max-w-5xl">
        <div className="text-center">
          <h1 className="font-semibold text-3xl text-neutral-800 sm:text-4xl">Forgot password?</h1>
          <p className="mt-4 leading-7">
            Remember your password? <Link href="/auth/login">Sign in</Link>
          </p>
        </div>

        <div className="mx-auto max-w-md items-center">
          <div
            className={cn(
              'mt-8 flex flex-col p-4 sm:p-6 md:mt-12 lg:p-8',
              'lg:rounded-xl lg:bg-white lg:shadow-md lg:ring-1 lg:ring-primary-950/5 dark:lg:bg-primary-900 dark:lg:ring-white/10'
            )}
          >
            <Form onSubmit={handleSubmit}>
              <Field name="email">
                {(field, props) => (
                  <div>
                    <Label htmlFor={field.name}>Email address</Label>
                    <Input
                      {...props}
                      id={field.name}
                      value={field.value.value || ''}
                      disabled={forgotPasswordForm.submitting.value}
                      placeholder="somebody@example.com"
                      type="email"
                      required
                    />
                    <span
                      className={cn(
                        field.error?.value ? 'block px-1 py-2 text-red-500 text-xs' : 'sr-only'
                      )}
                    >
                      {field.error.toString()}
                    </span>
                  </div>
                )}
              </Field>

              <div className="mt-4 grid">
                <Button type="submit" disabled={forgotPasswordForm.submitting.value}>
                  Continue
                </Button>
              </div>
            </Form>
          </div>
        </div>
      </div>
    </div>
  )
}
