import { type SubmitHandler, focus, reset, useForm, valiForm } from '@modular-forms/react'
import { Suspense } from 'react'
import { useNavigate } from 'react-router'
import { toast } from 'sonner'
import * as v from 'valibot'
import { Button, Input, Link } from '#/components/base'
import PageLoader from '#/components/page-loader'
import { cn } from '#/utils/helper'
import logger from '#/utils/logger'

const ResetPasswordSchema = v.pipe(
  v.object({
    password: v.pipe(
      v.string(),
      v.nonEmpty('Please enter your password.'),
      v.minLength(8, 'You password must have 8 characters or more.')
    ),
    confirmPassword: v.pipe(v.string(), v.nonEmpty('Please confirm your password.')),
    agreeTerms: v.boolean(),
  }),
  v.forward(
    v.partialCheck(
      [['password'], ['confirmPassword']],
      (input) => input.password === input.confirmPassword,
      'The passwords entered do not match.'
    ),
    ['confirmPassword']
  )
)

type ResetPasswordForm = v.InferInput<typeof ResetPasswordSchema>

export default function SignUpPage() {
  const navigate = useNavigate()

  const [resetPasswordForm, { Form, Field }] = useForm<ResetPasswordForm>({
    validate: valiForm(ResetPasswordSchema),
    validateOn: 'change',
  })

  const handleSubmit: SubmitHandler<ResetPasswordForm> = async (values, event) => {
    event.preventDefault()
    try {
      logger.info(values.password)
      toast.success('Your password has been reset successfully, redirecting to login page...', {
        duration: 2000,
      })
      setTimeout(() => navigate('/auth/login', { replace: true }), 2000)
    } catch (error) {
      logger.error('[ERROR] handleSubmit', error)
      if (error instanceof Error) {
        toast.error(error.message)
      }
      focus(resetPasswordForm, 'password')
    } finally {
      reset(resetPasswordForm)
    }
  }

  return (
    <Suspense fallback={<PageLoader />}>
      <div className="mx-auto w-full max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
        <div className="mx-auto w-full max-w-sm md:max-w-lg lg:max-w-5xl">
          <div className="text-center">
            <h1 className="font-semibold text-3xl text-neutral-800 sm:text-4xl">
              Create your account
            </h1>
            <p className="mt-4 leading-7">
              Already have an account? <Link href="/auth/login">Sign in</Link>
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
                <div className="grid gap-4">
                  <Field name="password">
                    {(field, props) => (
                      <div>
                        <label htmlFor={field.name} className="sr-only">
                          Password
                        </label>
                        <Input
                          {...props}
                          id={field.name}
                          value={field.value.value || ''}
                          disabled={resetPasswordForm.submitting.value}
                          placeholder="New Password"
                          type="password"
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

                  <Field name="confirmPassword">
                    {(field, props) => (
                      <div>
                        <label htmlFor={field.name} className="sr-only">
                          Password
                        </label>
                        <Input
                          {...props}
                          id={field.name}
                          value={field.value.value || ''}
                          disabled={resetPasswordForm.submitting.value}
                          placeholder="Confirm Password"
                          type="password"
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

                  <Field name="agreeTerms" type="boolean">
                    {(field, props) => (
                      <div className="flex w-full flex-col">
                        <label className="flex items-center">
                          <div className="flex">
                            <input
                              {...props}
                              type="checkbox"
                              className="mt-0.5 shrink-0 rounded border-neutral-200 text-primary-600 focus:ring-primary-500 dark:border-neutral-700 dark:bg-neutral-800 dark:focus:ring-offset-neutral-800 dark:checked:border-primary-500 dark:checked:bg-primary-500"
                              disabled={resetPasswordForm.submitting.value}
                              checked={!!field.value.value}
                            />
                          </div>
                          <div className="ms-2 text-sm dark:text-white">
                            <span className="font-medium decoration-2 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-neutral-600">
                              Terminate other sessions
                            </span>
                          </div>
                        </label>
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
                </div>

                <div className="mt-4 grid">
                  <Button type="submit" disabled={resetPasswordForm.submitting.value}>
                    Continue
                  </Button>
                </div>
              </Form>
            </div>
          </div>
        </div>
      </div>
    </Suspense>
  )
}
