import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import { AdminAuthProvider } from "@/contexts/admin-auth-context";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Byte Payments",
  description: "Fast, secure, and global payment gateway by The Byte Array",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <AdminAuthProvider>
            {children}
            <Toaster richColors position="top-right" />
          </AdminAuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
