"use client"

import * as React from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Image from "next/image"
import { 
  Copy, 
  Clock, 
  AlertCircle, 
  CheckCircle, 
  XCircle, 
  Loader2,
  ArrowLeft
} from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Separator } from "@/components/ui/separator"
import { Alert, AlertDescription } from "@/components/ui/alert"

import { ThemeToggle } from "@/components/theme-toggle"
import { apiClient, PaymentData } from "@/lib/api"

function PaymentPageContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const paymentId = searchParams.get("id")
  
  const [paymentData, setPaymentData] = React.useState<PaymentData | null>(null)
  const [isLoading, setIsLoading] = React.useState(true)
  const [timeLeft, setTimeLeft] = React.useState(15 * 60) // 15 minutes
  const [isCancelling, setIsCancelling] = React.useState(false)
  const previousStatus = React.useRef<string | null>(null)

  const intervalRef = React.useRef<NodeJS.Timeout | null>(null)
  const statusCheckRef = React.useRef<NodeJS.Timeout | null>(null)

  React.useEffect(() => {
    if (!paymentId) {
      showError()
      return
    }

    loadPaymentData()
    
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current)
      if (statusCheckRef.current) clearInterval(statusCheckRef.current)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [paymentId])

  const loadPaymentData = async () => {
    if (!paymentId) return

    try {
      const response = await apiClient.getPaymentStatus(paymentId)
      
      if (response.status === "ok" && response.data) {
        setPaymentData(response.data)
        previousStatus.current = response.data.status.toLowerCase()
        calculateTimeLeft(response.data.created_at)
        startStatusChecking()
      } else {
        showError()
      }
    } catch (error) {
      console.error("Error loading payment data:", error)
      showError()
    } finally {
      setIsLoading(false)
    }
  }

  const calculateTimeLeft = (createdAt: string) => {
    const createdTime = new Date(createdAt).getTime()
    const currentTime = new Date().getTime()
    const timeDifference = currentTime - createdTime
    const remaining = Math.max(0, Math.floor((15 * 60 * 1000 - timeDifference) / 1000))
    
    setTimeLeft(remaining)
    startTimer(remaining)
  }

  const startTimer = (initialTime: number) => {
    if (intervalRef.current) clearInterval(intervalRef.current)
    
    let currentTime = initialTime
    
    intervalRef.current = setInterval(() => {
      currentTime -= 1
      setTimeLeft(currentTime)
      
      if (currentTime <= 0) {
        if (intervalRef.current) clearInterval(intervalRef.current)
        checkIfExpired()
      }
    }, 1000)
  }

  const startStatusChecking = () => {
    if (statusCheckRef.current) clearInterval(statusCheckRef.current)
    
    statusCheckRef.current = setInterval(async () => {
      if (!paymentId) return
      
      try {
        const response = await apiClient.getPaymentStatus(paymentId)
        
        if (response.status === "ok" && response.data) {
          const newStatus = response.data.status.toLowerCase()
          
          previousStatus.current = newStatus
          setPaymentData(response.data)
          
          if (["completed", "cancelled", "expired"].includes(newStatus)) {
            if (statusCheckRef.current) clearInterval(statusCheckRef.current)
          }
        }
      } catch (error) {
        console.error("Error checking payment status:", error)
      }
    }, 10000) // Check every 10 seconds
  }

  const checkIfExpired = async () => {
    if (!paymentId) return
    
    try {
      const response = await apiClient.getPaymentStatus(paymentId)
      
      if (response.status === "ok" && response.data?.status.toLowerCase() === "pending") {
        toast.error("Payment time has expired!")
      }
    } catch (error) {
      console.error("Error validating expiration:", error)
    }
  }

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    return `${minutes.toString().padStart(2, "0")}:${remainingSeconds.toString().padStart(2, "0")}`
  }

  const copyToClipboard = async (text: string, label: string) => {
    try {
      await navigator.clipboard.writeText(text)
      toast.success(`${label} copied to clipboard!`)
    } catch {
      toast.error(`Failed to copy ${label.toLowerCase()}`)
    }
  }

  const handleCancelPayment = async () => {
    if (!paymentId) return
    
    const confirmed = window.confirm("Are you sure you want to cancel this payment?")
    if (!confirmed) return

    setIsCancelling(true)
    try {
      const response = await apiClient.cancelPayment(paymentId)
      
      if (response.status === "ok") {
        toast.success("Payment cancelled successfully")
        setPaymentData(prev => prev ? { ...prev, status: "cancelled" } : null)
      } else {
        toast.error(response.message || "Failed to cancel payment")
      }
    } catch (error) {
      console.error("Error cancelling payment:", error)
      toast.error("Failed to cancel payment")
    } finally {
      setIsCancelling(false)
    }
  }

  const getStatusBadge = (status: string) => {
    const statusLower = status.toLowerCase()
    
    switch (statusLower) {
      case "pending":
        return (
          <Badge variant="secondary" className="bg-amber-100 text-amber-800 dark:bg-amber-900/20 dark:text-amber-400">
            <Clock className="w-3 h-3 mr-1" />
            Awaiting Payment
          </Badge>
        )
      case "completed":
        return (
          <Badge variant="secondary" className="bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400">
            <CheckCircle className="w-3 h-3 mr-1" />
            Payment Completed
          </Badge>
        )
      case "cancelled":
        return (
          <Badge variant="secondary" className="bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400">
            <XCircle className="w-3 h-3 mr-1" />
            Payment Cancelled
          </Badge>
        )
      case "expired":
        return (
          <Badge variant="secondary" className="bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400">
            <AlertCircle className="w-3 h-3 mr-1" />
            Payment Expired
          </Badge>
        )
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  const showError = () => {
    setIsLoading(false)
    // We'll show an error state in the render
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    )
  }

  if (!paymentData) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center p-4">
        <Card className="max-w-md">
          <CardContent className="p-8 text-center">
            <div className="w-16 h-16 bg-red-100 dark:bg-red-900/20 rounded-full flex items-center justify-center mx-auto mb-4">
              <AlertCircle className="w-8 h-8 text-red-500 dark:text-red-400" />
            </div>
            <h2 className="text-xl font-bold mb-2">Payment Not Found</h2>
            <p className="text-muted-foreground mb-6">
              The payment you&apos;re looking for doesn&apos;t exist or has expired. Please create a new payment.
            </p>
            <div className="space-y-3">
              <Button onClick={() => router.push("/")} className="w-full">
                Select New Plan
              </Button>
              <Button variant="outline" onClick={() => router.back()} className="w-full">
                Go Back
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  const amount = parseFloat(paymentData.trx_amount).toFixed(2)
  const isExpired = ["completed", "cancelled", "expired"].includes(paymentData.status.toLowerCase())
  
  // Calculate progress based on time left
  const totalTime = 15 * 60 // 15 minutes in seconds
  const progressPercentage = Math.max(0, (timeLeft / totalTime) * 100)

  return (
    <div className="min-h-screen bg-background">
      {/* Minimal Header */}
      <div className="border-b bg-card">
        <div className="max-w-2xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => router.push("/")}
            >
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back
            </Button>
            <div className="text-lg font-semibold">Byte Payments</div>
          </div>
          <ThemeToggle />
        </div>
      </div>

      {/* Main Payment Card */}
      <div className="max-w-2xl mx-auto p-4 py-8">
        <Card>
          <CardHeader className="text-center pb-6">
            <div className="flex items-center justify-center gap-2 mb-2">
              <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">
                {amount} TRX
              </div>
              {getStatusBadge(paymentData.status)}
            </div>
            <p className="text-muted-foreground text-sm">
              Send exactly this amount to complete your payment
            </p>
          </CardHeader>

          <CardContent className="space-y-6">
            {/* Timer */}
            {!isExpired && (
              <div className="text-center">
                <div className="flex items-center justify-center gap-2 mb-2">
                  <Clock className="w-4 h-4 text-muted-foreground" />
                  <span className="text-sm text-muted-foreground">Time remaining</span>
                </div>
                <div className={`text-2xl font-mono font-bold mb-3 ${
                  timeLeft < 60 ? "text-red-500" : timeLeft < 300 ? "text-amber-500" : "text-green-600"
                }`}>
                  {formatTime(timeLeft)}
                </div>
                <Progress 
                  value={progressPercentage} 
                  className={`h-1 ${
                    timeLeft < 60 ? "[&>div]:bg-red-500" : 
                    timeLeft < 300 ? "[&>div]:bg-amber-500" : 
                    "[&>div]:bg-green-600"
                  }`}
                />
              </div>
            )}

            {/* QR Code */}
            <div className="text-center">
              <div className="inline-block bg-white dark:bg-slate-50 rounded-lg p-4 shadow-sm border">
                <Image
                  src={paymentData.qr_image}
                  alt="Payment QR Code"
                  width={180}
                  height={180}
                  className="cursor-pointer transition-transform hover:scale-105 rounded"
                  onClick={() => copyToClipboard(paymentData.trx_wallet_address, "Address")}
                />
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                Scan with your TRON wallet
              </p>
            </div>

            {/* Payment Address */}
            <div className="space-y-3">
              <div className="flex items-center justify-between bg-muted/30 rounded-lg p-3">
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-1">TRON Address</p>
                  <p className="font-mono text-sm break-all pr-2">
                    {paymentData.trx_wallet_address}
                  </p>
                </div>
                <Button
                  onClick={() => copyToClipboard(paymentData.trx_wallet_address, "Address")}
                  size="sm"
                  variant="outline"
                >
                  <Copy className="w-4 h-4" />
                </Button>
              </div>
            </div>

            {/* Important Message */}
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="text-sm">
                <strong>Important:</strong> Send exactly {amount} TRX. Incorrect amounts may cause delays or loss of funds.
              </AlertDescription>
            </Alert>

            {/* Payment Details */}
            <div className="bg-muted/20 rounded-lg p-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Payment ID</span>
                <span 
                  className="font-mono cursor-pointer hover:text-blue-500 transition-colors"
                  onClick={() => copyToClipboard(paymentData.payment_id, "Payment ID")}
                  title="Click to copy"
                >
                  {paymentData.payment_id.substring(0, 16)}...
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Network</span>
                <span>TRON</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Status Updates</span>
                <span>{isExpired ? "Completed" : "Every 10s"}</span>
              </div>
            </div>

            {/* Cancel Button */}
            {!isExpired && (
              <Button
                variant="outline"
                onClick={handleCancelPayment}
                disabled={isCancelling}
                className="w-full"
              >
                {isCancelling && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Cancel Payment
              </Button>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default function PaymentPage() {
  return (
    <React.Suspense fallback={
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    }>
      <PaymentPageContent />
    </React.Suspense>
  )
}