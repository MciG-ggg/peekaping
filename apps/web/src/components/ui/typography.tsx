import clsx from "clsx"
import type { PropsWithChildren } from "react"

export function TypographyH1({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <h1 className={clsx("scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl", className)}>
      {children}
    </h1>
  )
}

export function TypographyH2({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <h2 className={clsx("scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0", className)}>
      {children}
    </h2>
  )
}

export function TypographyH3({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <h3 className={clsx("scroll-m-20 text-2xl font-semibold tracking-tight", className)}>
     {children}
    </h3>
  )
}

export function TypographyH4({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <h4 className={clsx("scroll-m-20 text-xl font-semibold tracking-tight", className)}>
      {children}
    </h4>
  )
}

export function TypographyH5({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <h5 className={clsx("scroll-m-20 tracking-tight", className)}>
      {children}
    </h5>
  )
}

export function TypographyP({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <p className={clsx("leading-7 [&:not(:first-child)]:mt-6", className)}>
      {children}
    </p>
  )
}

