import "@/styles/globals.css"
import { fontSans } from "@/lib/fonts"
import { cn } from "@/lib/utils"


interface RootLayoutProps {
  children: React.ReactNode
}

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <>
      <html suppressHydrationWarning={true} lang="en">
        <head />
        <body
          suppressHydrationWarning={true}
          className={cn(
            "min-h-screen bg-background font-sans antialiased",
            fontSans.variable
          )}
        >
          <div className="flex-1">{children}</div>
        </body>
      </html>
    </>
  )
}
