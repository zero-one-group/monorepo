import * as Lucide from 'lucide-react'
import { Accordion as AccordionPrimitive } from 'radix-ui'
import * as React from 'react'
import { accordionStyles } from './accordion.css'

const Accordion = AccordionPrimitive.Root

const AccordionItem = React.forwardRef<
  React.ComponentRef<typeof AccordionPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof AccordionPrimitive.Item>
>(({ className, ...props }, ref) => {
  const styles = accordionStyles()
  return <AccordionPrimitive.Item ref={ref} className={styles.item({ className })} {...props} />
})

const AccordionTrigger = React.forwardRef<
  React.ComponentRef<typeof AccordionPrimitive.Trigger>,
  React.ComponentPropsWithoutRef<typeof AccordionPrimitive.Trigger>
>(({ className, children, ...props }, ref) => {
  const styles = accordionStyles()
  return (
    <AccordionPrimitive.Header className={styles.headerWrapper()}>
      <AccordionPrimitive.Trigger
        ref={ref}
        className={styles.headerTrigger({ className })}
        {...props}
      >
        {children}
        <Lucide.ChevronDown className={styles.headerIcon()} strokeWidth={2} />
      </AccordionPrimitive.Trigger>
    </AccordionPrimitive.Header>
  )
})

const AccordionContent = React.forwardRef<
  React.ComponentRef<typeof AccordionPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof AccordionPrimitive.Content>
>(({ className, children, ...props }, ref) => {
  const styles = accordionStyles()
  return (
    <AccordionPrimitive.Content ref={ref} className={styles.contentWrapper()} {...props}>
      <div className={styles.contentChildren({ className })}>{children}</div>
    </AccordionPrimitive.Content>
  )
})

AccordionItem.displayName = 'AccordionItem'
AccordionTrigger.displayName = AccordionPrimitive.Trigger.displayName
AccordionContent.displayName = AccordionPrimitive.Content.displayName

export { Accordion, AccordionItem, AccordionTrigger, AccordionContent }
