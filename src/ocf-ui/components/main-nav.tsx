import Link from "next/link"

import { cn } from "@/lib/utils"

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <nav
      className={cn("flex items-center space-x-4 lg:space-x-6", className)}
      {...props}
    >
      <Link
        href="/app/dashboard"
        className="text-sm font-medium transition-colors hover:text-primary"
      >
        Overview
      </Link>
      <Link
        href="/app/playground"
        className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
      >
        Playground
      </Link>
      <Link
        href="https://ocf.autoai.org"
        className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
      >
        Docs
      </Link>
      <Link
        href="https://ocfstatus.autoai.dev"
        className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
      >
        Status
      </Link>
    </nav>
  )
}
