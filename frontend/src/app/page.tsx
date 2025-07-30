"use client"

import * as React from "react"
import { useRouter } from "next/navigation"
import { Mail, Loader2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { ThemeToggle } from "@/components/theme-toggle"
import { PlanCard } from "@/components/plan-card"
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp"
import { apiClient, Plan, Currency } from "@/lib/api"

export default function HomePage() {
  const router = useRouter()
  const [email, setEmail] = React.useState("")
  const [plans, setPlans] = React.useState<Plan[]>([])
  const [currencies, setCurrencies] = React.useState<Currency[]>([])
  const [selectedPlan, setSelectedPlan] = React.useState<Plan | null>(null)
  const [isLoading, setIsLoading] = React.useState(false)
  const [showVerification, setShowVerification] = React.useState(false)
  const [isVerifying, setIsVerifying] = React.useState(false)
  const [isResending, setIsResending] = React.useState(false)

  const [otpCode, setOtpCode] = React.useState("")

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const isEmailValid = email && emailRegex.test(email)
  const canContinue = isEmailValid && selectedPlan

  React.useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [plansResponse, currenciesResponse] = await Promise.all([
        apiClient.getPlans(),
        apiClient.getCurrencies(),
      ])

      if (plansResponse.status === "ok" && plansResponse.data) {
        setPlans(plansResponse.data)
      } else {
        toast.error("Failed to load plans")
      }

      if (currenciesResponse.status === "ok" && currenciesResponse.data) {
        setCurrencies(currenciesResponse.data)
      }
    } catch (error) {
      console.error("Error loading data:", error)
      toast.error("Failed to load data")
    }
  }

  const handleSendVerification = async () => {
    if (!canContinue) return

    setIsLoading(true)
    try {
      const response = await apiClient.sendVerificationCode(email)
      
      if (response.status === "ok") {
        setShowVerification(true)
        toast.success("Verification code sent to your email")
      } else {
        toast.error(response.message || "Failed to send verification code")
      }
    } catch (error) {
      console.error("Error sending verification code:", error)
      toast.error("Failed to send verification code")
    } finally {
      setIsLoading(false)
    }
  }

  const handleVerifyCode = async (code: string) => {
    if (code.length !== 6) return

    setOtpCode(code)
    setIsVerifying(true)

    try {
      const response = await apiClient.verifyCode(email, code)
      
      if (response.status === "ok" && response.data?.is_valid) {
        if (response.data.verification_token) {
          toast.success("Email verified successfully!")
          await createPayment(response.data.verification_token)
        }
      } else {
        toast.error(response.message || "Invalid verification code")
        setOtpCode("")
      }
    } catch (error) {
      console.error("Error verifying code:", error)
      toast.error("Failed to verify code")
      setOtpCode("")
    } finally {
      setIsVerifying(false)
    }
  }

  const createPayment = async (token: string) => {
    if (!selectedPlan) return

    const trxCurrency = currencies.find((c) => c.code === "TRX")
    if (!trxCurrency) {
      toast.error("TRX currency not available")
      return
    }

    try {
      const response = await apiClient.createPayment(
        email,
        selectedPlan.id,
        trxCurrency.code,
        token
      )

      if (response.status === "ok" && response.data) {
        router.push(`/pay?id=${response.data.payment_id}`)
      } else {
        toast.error(response.message || "Failed to create payment")
      }
    } catch (error) {
      console.error("Error creating payment:", error)
      toast.error("Failed to create payment")
    }
  }

  const handleResendCode = async () => {
    setIsResending(true)
    try {
      const response = await apiClient.sendVerificationCode(email)
      
      if (response.status === "ok") {
        toast.success("New verification code sent")
        setOtpCode("")
      } else {
        toast.error(response.message || "Failed to resend code")
      }
    } catch (error) {
      console.error("Error resending code:", error)
      toast.error("Failed to resend code")
    } finally {
      setIsResending(false)
    }
  }

  const handleBackToEmail = () => {
    setShowVerification(false)
    setOtpCode("")
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b bg-card">
        <div className="max-w-3xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="text-lg font-semibold">Byte Payments</div>
          <ThemeToggle />
        </div>
      </div>

      <div className="max-w-3xl mx-auto px-4 py-8">
        {!showVerification ? (
          <div className="max-w-xl mx-auto space-y-8">
            {/* Header */}
            <div className="text-center space-y-2">
              <h1 className="text-2xl font-bold">Select Your Plan</h1>
              <p className="text-muted-foreground">
                Choose a plan that fits your needs
              </p>
            </div>

            {/* Email Input */}
            <Card>
              <CardContent className="pt-6">
                <div className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="email">Email Address</Label>
                    <Input
                      id="email"
                      type="email"
                      placeholder="Enter your email address"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === "Enter" && canContinue) {
                          handleSendVerification()
                        }
                      }}
                    />
                    {email && !isEmailValid && (
                      <p className="text-sm text-destructive">
                        Please enter a valid email address
                      </p>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Plans */}
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Available Plans</h3>
              <div className="space-y-3">
                {plans.length === 0 ? (
                  <div className="flex justify-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin" />
                  </div>
                ) : (
                  plans.map((plan) => (
                    <div
                      key={plan.id}
                      className={`border rounded-lg p-4 cursor-pointer transition-all hover:border-blue-500 ${
                        selectedPlan?.id === plan.id 
                          ? "border-blue-500 bg-blue-50 dark:bg-blue-950/20" 
                          : "border-border"
                      }`}
                      onClick={() => setSelectedPlan(plan)}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-3">
                            <div 
                              className={`w-4 h-4 rounded-full border-2 flex items-center justify-center ${
                                selectedPlan?.id === plan.id
                                  ? "border-blue-500 bg-blue-500"
                                  : "border-muted-foreground"
                              }`}
                            >
                              {selectedPlan?.id === plan.id && (
                                <div className="w-2 h-2 bg-white rounded-full" />
                              )}
                            </div>
                            <div>
                              <h4 className="font-medium">{plan.name}</h4>
                              <p className="text-sm text-muted-foreground">{plan.description}</p>
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-lg font-bold text-blue-600 dark:text-blue-400">
                            ${plan.price_usd}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {plan.duration_days} days
                          </div>
                        </div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            {/* Continue Button */}
            <Button
              onClick={handleSendVerification}
              disabled={!canContinue || isLoading}
              className="w-full"
              size="lg"
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {selectedPlan && isEmailValid
                ? `Continue with ${selectedPlan.name}`
                : !isEmailValid
                ? "Enter your email address"
                : "Select a plan to continue"}
            </Button>

            {/* Footer */}
            <div className="text-center pt-4">
              <div className="flex items-center justify-center gap-6 text-xs text-muted-foreground">
                <span>🔒 Secure payments</span>
                <span>⚡ Instant activation</span>
                <span>🌍 Global access</span>
              </div>
            </div>
          </div>
        ) : (
          /* Verification Step */
          <div className="max-w-md mx-auto">
            <Card>
              <CardHeader className="text-center">
                <div className="w-12 h-12 mx-auto mb-4 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center">
                  <Mail className="w-6 h-6 text-blue-600 dark:text-blue-400" />
                </div>
                <CardTitle>Verify your email</CardTitle>
                <p className="text-sm text-muted-foreground">
                  Enter the 6-digit code sent to<br />
                  <span className="font-medium text-foreground">{email}</span>
                </p>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="flex justify-center">
                  <InputOTP
                    maxLength={6}
                    value={otpCode}
                    onChange={setOtpCode}
                    onComplete={handleVerifyCode}
                  >
                    <InputOTPGroup>
                      <InputOTPSlot index={0} />
                      <InputOTPSlot index={1} />
                      <InputOTPSlot index={2} />
                      <InputOTPSlot index={3} />
                      <InputOTPSlot index={4} />
                      <InputOTPSlot index={5} />
                    </InputOTPGroup>
                  </InputOTP>
                </div>

                <Button
                  onClick={() => handleVerifyCode(otpCode)}
                  disabled={otpCode.length !== 6 || isVerifying}
                  className="w-full"
                >
                  {isVerifying && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Verify Code
                </Button>

                <div className="flex items-center justify-center gap-4 text-sm">
                  <Button
                    variant="link"
                    onClick={handleResendCode}
                    disabled={isResending}
                    className="p-0 h-auto text-xs"
                  >
                    {isResending ? "Sending..." : "Resend code"}
                  </Button>
                  <span className="text-muted-foreground">•</span>
                  <Button
                    variant="link"
                    onClick={handleBackToEmail}
                    className="p-0 h-auto text-xs"
                  >
                    Change email
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}
