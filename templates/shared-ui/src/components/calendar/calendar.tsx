import * as Lucide from 'lucide-react'
import { DayPicker } from 'react-day-picker'
import type { DayPickerProps } from 'react-day-picker'
import { clx } from '../../utils'
import { calendarStyles } from './calendar.css'

export type CalendarProps = {
  timezone?: string
} & DayPickerProps

function Calendar({
  className,
  classNames,
  showOutsideDays = true,
  timezone,
  ...props
}: CalendarProps) {
  const styles = calendarStyles()
  return (
    <DayPicker
      timeZone={timezone}
      showOutsideDays={showOutsideDays}
      className={styles.root({ className })}
      classNames={% raw %}{{
        root: styles.root(),
        months: styles.months(),
        month: styles.month(),
        month_caption: styles.month_caption(),
        caption_label: styles.caption_label(),
        nav: styles.nav(),
        button_previous: styles.button_previous(),
        button_next: styles.button_next(),
        month_grid: styles.month_grid(),
        weekdays: styles.weekdays(),
        weekday: styles.weekday(),
        week: styles.week(),
        day: clx(
          styles.cell(),
          props.mode === 'range' ? styles.cell_range() : styles.cell_single()
        ),
        day_button: styles.day_button(),
        range_start: styles.range_start(),
        range_end: styles.range_end(),
        selected: styles.selected(),
        today: styles.today(),
        outside: styles.outside(),
        disabled: styles.disabled(),
        range_middle: styles.range_middle(),
        hidden: styles.hidden(),
        ...classNames,
      }}{% endraw %}
      components={% raw %}{{
        Chevron: ({ orientation, ...props }) => {
          const Icon = orientation === 'left' ? Lucide.ChevronLeft : Lucide.ChevronRight
          return <Icon className={clx(styles.icon(), props.className)} />
        },
      }}{% endraw %}
      {...props}
    />
  )
}

Calendar.displayName = 'Calendar'

export { Calendar }
