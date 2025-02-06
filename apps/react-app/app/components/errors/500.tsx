import { useNavigate } from 'react-router'
import { Link } from '#/components/link'
import { errorStyles } from './error.css'

interface InternalErrorProps {
  message: string
  details: string
  stack?: string
}

export default function InternalError({ message, details, stack }: InternalErrorProps) {
  const navigate = useNavigate()
  const styles = errorStyles()

  const handleBack = () => {
    if (window.history.length > 1) {
      navigate(-1)
    } else {
      navigate('/')
    }
  }

  return (
    <div className={styles.wrapper()}>
      <div className={styles.decorativeGradient()}>
        <div className={styles.gradientInner()}>
          <div className={styles.gradientBg()} />
        </div>
      </div>
      <div className={styles.decorativeCode()}>
        <h2 className={styles.decorativeText()}>500</h2>
      </div>
      <div className={styles.content()}>
        <div className={styles.container()}>
          <p className={styles.errorCode()}>{message}</p>
          <h1 className={styles.title()}>Internal Server Error</h1>
          <p className={styles.description()}>{details}</p>
          {stack && (
            <pre className={styles.pre()}>
              <code className={styles.code()}>{stack}</code>
            </pre>
          )}
          <div className={styles.actions()}>
            <button type="button" onClick={handleBack} className={styles.primaryButton()}>
              Go back
            </button>
            <Link href="#" className={styles.secondaryButton()}>
              Troubleshooting Guide
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
