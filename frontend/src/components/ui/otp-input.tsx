"use client"

import * as React from "react"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

interface OTPInputProps {
  length?: number
  onComplete?: (code: string) => void
  onChange?: (code: string) => void
  className?: string
}

export function OTPInput({ 
  length = 6, 
  onComplete, 
  onChange, 
  className 
}: OTPInputProps) {
  const [values, setValues] = React.useState<string[]>(Array(length).fill(""))
  const inputRefs = React.useRef<(HTMLInputElement | null)[]>([])

  const handleChange = React.useCallback((index: number, value: string) => {
    const newValue = value.replace(/\D/g, "").slice(0, 1)
    const newValues = [...values]
    newValues[index] = newValue
    setValues(newValues)

    const code = newValues.join("")
    onChange?.(code)

    if (newValue && index < length - 1) {
      inputRefs.current[index + 1]?.focus()
    }

    if (code.length === length) {
      onComplete?.(code)
    }
  }, [values, length, onChange, onComplete])

  const handleKeyDown = React.useCallback((index: number, e: React.KeyboardEvent) => {
    if (e.key === "Backspace" && !values[index] && index > 0) {
      inputRefs.current[index - 1]?.focus()
    }
    if (e.key === "ArrowLeft" && index > 0) {
      inputRefs.current[index - 1]?.focus()
    }
    if (e.key === "ArrowRight" && index < length - 1) {
      inputRefs.current[index + 1]?.focus()
    }
  }, [values, length])

  const handlePaste = React.useCallback((e: React.ClipboardEvent) => {
    e.preventDefault()
    const pastedData = e.clipboardData
      .getData("text")
      .replace(/\D/g, "")
      .slice(0, length)

    const newValues = Array(length).fill("")
    pastedData.split("").forEach((char, i) => {
      if (i < length) {
        newValues[i] = char
      }
    })

    setValues(newValues)
    
    const code = newValues.join("")
    onChange?.(code)
    
    if (pastedData.length === length) {
      onComplete?.(code)
      inputRefs.current[length - 1]?.focus()
    } else if (pastedData.length > 0) {
      inputRefs.current[Math.min(pastedData.length, length - 1)]?.focus()
    }
  }, [length, onChange, onComplete])



  return (
    <div className={cn("flex justify-center gap-2", className)}>
      {Array.from({ length }, (_, index) => (
        <Input
          key={index}
          ref={(el) => {
            inputRefs.current[index] = el
          }}
          type="text"
          inputMode="numeric"
          value={values[index]}
          onChange={(e) => handleChange(index, e.target.value)}
          onKeyDown={(e) => handleKeyDown(index, e)}
          onPaste={handlePaste}
          className="w-12 h-12 text-center text-xl font-semibold"
          maxLength={1}
        />
      ))}
    </div>
  )
}