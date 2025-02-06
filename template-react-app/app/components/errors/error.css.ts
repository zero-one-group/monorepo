import { tv } from 'tailwind-variants'

export const errorStyles = tv({
  slots: {
    wrapper: 'relative min-h-screen overflow-hidden bg-black',
    decorativeGradient: 'absolute inset-0 overflow-hidden',
    gradientInner: '-inset-[10px] absolute opacity-50',
    gradientBg: [
      'absolute top-0 h-[40rem] w-full',
      'bg-gradient-to-b from-gray-500/20 via-transparent to-transparent dark:from-gray-900/30 dark:via-transparent dark:to-transparent',
      'before:absolute before:inset-0 before:bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] before:from-gray-400/10 before:via-transparent before:to-transparent dark:before:from-gray-500/10',
      'after:absolute after:inset-0 after:bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] after:from-indigo-400/10 after:via-transparent after:to-transparent dark:after:from-indigo-500/10',
    ],
    content:
      'relative flex min-h-screen flex-col items-center justify-center px-4 py-16 sm:px-6 lg:px-8',
    container: 'relative z-20 text-center',
    errorCode: 'font-bold text-2xl text-red-500 dark:text-red-400',
    title: 'mt-4 font-bold text-3xl text-gray-900 tracking-tight sm:text-5xl dark:text-gray-100',
    description: 'mt-6 text-base text-gray-600 leading-7 dark:text-gray-400',
    pre: 'w-full overflow-x-auto p-4',
    code: 'text-gray-600 dark:text-gray-400',
    actions: 'mt-10 flex items-center justify-center gap-x-4',
    primaryButton: [
      'min-w-[140px] cursor-pointer rounded-lg bg-gray-900 px-4 py-2.5 font-semibold text-sm text-white',
      'transition-all duration-200 hover:bg-gray-800 hover:shadow-gray-500/20 hover:shadow-lg',
      'focus:outline-none focus:ring-2 focus:ring-gray-400/50 focus:ring-offset-2',
      'focus:ring-offset-gray-50 dark:focus:ring-offset-gray-950',
      'border border-gray-200 dark:border-gray-800',
    ],
    secondaryButton: [
      'min-w-[140px] rounded-lg border border-gray-200 bg-white/80 px-4 py-2.5 dark:border-gray-800 dark:bg-gray-900/80',
      'cursor-pointer font-semibold text-gray-700 text-sm dark:text-gray-200',
      'transition-all duration-200 hover:border-gray-500/30 hover:bg-gray-50 dark:hover:bg-gray-800',
      'hover:text-gray-500 hover:shadow-gray-500/10 hover:shadow-lg dark:hover:text-gray-400',
      'focus:outline-none focus:ring-2 focus:ring-gray-400/50 focus:ring-offset-2',
      'focus:ring-offset-gray-50 dark:focus:ring-offset-gray-950',
    ],
    decorativeCode:
      'pointer-events-none fixed inset-0 z-10 flex select-none items-center justify-center',
    decorativeText:
      'font-black text-[12rem] text-red-100/30 mix-blend-overlay sm:text-[16rem] md:text-[20rem] dark:text-red-900/20',
  },
})
