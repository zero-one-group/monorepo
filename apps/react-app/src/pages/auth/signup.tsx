import { type SubmitHandler, focus, reset, useForm, valiForm } from '@modular-forms/react'
import { Button, Input, cn } from '@myorg/shared-ui'
import { Suspense } from 'react'
import { toast } from 'sonner'
import * as v from 'valibot'
import { Link } from '#/components/link'
import PageLoader from '#/components/page-loader'
import logger from '#/utils/logger'

const SignupSchema = v.pipe(
  v.object({
    firstName: v.pipe(v.string(), v.nonEmpty('Please enter your first name.')),
    lastName: v.pipe(v.string(), v.nonEmpty('Please enter your last name.')),
    email: v.pipe(
      v.string(),
      v.nonEmpty('Please enter your email.'),
      v.email('The email address is badly formatted.')
    ),
    password: v.pipe(
      v.string(),
      v.nonEmpty('Please enter your password.'),
      v.minLength(8, 'You password must have 8 characters or more.')
    ),
    confirmPassword: v.pipe(v.string(), v.nonEmpty('Please confirm your password.')),
    agreeTerms: v.pipe(
      v.boolean(),
      v.check((value) => value === true, 'You must agree to the terms and conditions.')
    ),
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

type SignupForm = v.InferInput<typeof SignupSchema>

export default function SignUpPage() {
  const [signupForm, { Form, Field }] = useForm<SignupForm>({
    validate: valiForm(SignupSchema),
    validateOn: 'change',
  })

  const handleSubmit: SubmitHandler<SignupForm> = async (values, event) => {
    event.preventDefault()
    try {
      logger.info(values.email)
      toast.success("Check your inbox! We've sent you an email with a link to reset your password.")
    } catch (error) {
      logger.error('[ERROR] handleSubmit', error)
      if (error instanceof Error) {
        toast.error(error.message)
      }
      focus(signupForm, 'firstName')
    } finally {
      reset(signupForm)
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
                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <Field name="firstName">
                      {(field, props) => (
                        <div>
                          <label htmlFor={field.name} className="sr-only">
                            First Name
                          </label>
                          <Input
                            {...props}
                            id={field.name}
                            value={field.value.value || ''}
                            disabled={signupForm.submitting.value}
                            placeholder="First name"
                            type="text"
                            required
                          />
                          <span
                            className={cn(
                              field.error?.value
                                ? 'block px-1 py-2 text-red-500 text-xs'
                                : 'sr-only'
                            )}
                          >
                            {field.error.toString()}
                          </span>
                        </div>
                      )}
                    </Field>

                    <Field name="lastName">
                      {(field, props) => (
                        <div>
                          <label htmlFor={field.name} className="sr-only">
                            Last Name
                          </label>
                          <Input
                            {...props}
                            id={field.name}
                            value={field.value.value || ''}
                            disabled={signupForm.submitting.value}
                            placeholder="Last name"
                            type="text"
                            required
                          />
                          <span
                            className={cn(
                              field.error?.value
                                ? 'block px-1 py-2 text-red-500 text-xs'
                                : 'sr-only'
                            )}
                          >
                            {field.error.toString()}
                          </span>
                        </div>
                      )}
                    </Field>
                  </div>

                  <Field name="email">
                    {(field, props) => (
                      <div>
                        <label htmlFor={field.name} className="sr-only">
                          Email address
                        </label>
                        <Input
                          {...props}
                          id={field.name}
                          value={field.value.value || ''}
                          disabled={signupForm.submitting.value}
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
                          disabled={signupForm.submitting.value}
                          placeholder="Password"
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
                          disabled={signupForm.submitting.value}
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
                              disabled={signupForm.submitting.value}
                              className="mt-0.5 shrink-0 rounded border-neutral-200 text-primary-600 focus:ring-primary-500 dark:border-neutral-700 dark:bg-neutral-800 dark:focus:ring-offset-neutral-800 dark:checked:border-primary-500 dark:checked:bg-primary-500"
                              checked={!!field.value.value}
                              required
                            />
                          </div>
                          <div className="ms-2 text-sm dark:text-white">
                            I accept the{' '}
                            <Link href="https://example.com/terms" tabIndex={-1} newTab>
                              Terms and Conditions
                            </Link>
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
                  <Button type="submit" disabled={signupForm.submitting.value}>
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
