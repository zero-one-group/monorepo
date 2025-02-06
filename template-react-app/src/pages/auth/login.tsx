import { type SubmitHandler, focus, reset, useForm, valiForm } from '@modular-forms/react'
import { useNavigate, useSearchParams } from 'react-router'
import { toast } from 'sonner'
import * as v from 'valibot'
import { Button, Input, Label, Link } from '#/components/base'
import { useAuth } from '#/context/hooks/use-auth'
import { cn } from '#/utils/helper'
import logger from '#/utils/logger'

const LoginSchema = v.object({
  identity: v.pipe(
    v.string(),
    v.nonEmpty('Please enter your email.'),
    v.email('The email address is badly formatted.')
  ),
  password: v.pipe(v.string(), v.nonEmpty('Please enter your password.')),
})

type LoginForm = v.InferInput<typeof LoginSchema>

export default function SignInPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const returnTo = searchParams.get('returnTo')
  const auth = useAuth()

  const [loginForm, { Form, Field }] = useForm<LoginForm>({
    validate: valiForm(LoginSchema),
    validateOn: 'change',
  })

  const handleSubmit: SubmitHandler<LoginForm> = async (values, event) => {
    event.preventDefault()

    // TODO: enhance security by count of failed login attempts (submitCount)
    logger.info('LOGIN', 'Submit Count:', loginForm.submitCount.value)

    await auth
      .login(values.identity, values.password)
      .then((user) => {
        toast.success(`Welcome back ${user?.first_name}`)
        navigate(returnTo || '/')
      })
      .catch((error) => {
        logger.error('[ERROR] handleSubmit', error)
        if (error instanceof Error) {
          toast.error(error.message)
        }
        focus(loginForm, 'identity')
      })
      .finally(() => reset(loginForm))
  }

  return (
    <div className="mx-auto w-full max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
      <div className="mx-auto w-full max-w-sm md:max-w-lg lg:max-w-5xl">
        <div className="text-center">
          <h1 className="font-semibold text-3xl text-neutral-800 sm:text-4xl">Welcome back</h1>
          <p className="mt-4 leading-7">
            Don&apos;t have an account? <Link href="/auth/register">Create account</Link>
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
              <Field name="identity">
                {(field, props) => (
                  <div>
                    <Label htmlFor={field.name}>Email address</Label>
                    <Input
                      {...props}
                      id={field.name}
                      value={field.value.value || ''}
                      disabled={loginForm.submitting.value}
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
                  <div className="mt-4">
                    <label htmlFor={field.name}>Password</label>
                    <Input
                      {...props}
                      id={field.name}
                      value={field.value.value || ''}
                      disabled={loginForm.submitting.value}
                      placeholder="************"
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

              <div className="mt-4 grid">
                <Button type="submit" disabled={loginForm.submitting.value}>
                  Continue
                </Button>
              </div>

              <div className="mt-6 text-center">
                <p className="text-neutral-500 text-sm">
                  <Link href="/auth/recovery">Forgot your password?</Link>
                </p>
              </div>
            </Form>
          </div>
        </div>
      </div>
    </div>
  )
}
