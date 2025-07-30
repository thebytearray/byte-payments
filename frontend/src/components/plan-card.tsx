"use client"

import * as React from "react"
import { Check } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import { Plan } from "@/lib/api"

interface PlanCardProps {
  plan: Plan
  isSelected?: boolean
  onClick?: () => void
  className?: string
}

export function PlanCard({ plan, isSelected, onClick, className }: PlanCardProps) {
  return (
    <Card
      className={cn(
        "cursor-pointer transition-all duration-200 hover:border-blue-500 hover:shadow-md",
        isSelected && "border-blue-500 bg-blue-50 dark:bg-blue-950/20",
        className
      )}
      onClick={onClick}
    >
      <CardContent className="p-6">
        <div className="text-center">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
            {plan.name}
          </h3>
          <div className="mb-3">
            <span className="text-2xl font-bold text-blue-600 dark:text-blue-400">
              ${plan.price_usd}
            </span>
            <span className="text-gray-500 dark:text-gray-400 text-sm ml-1">
              USD
            </span>
          </div>
          <p className="text-gray-600 dark:text-gray-400 text-sm mb-3">
            {plan.description}
          </p>
          <div className="flex items-center justify-center gap-1 text-sm text-gray-500 dark:text-gray-400">
            <span>Duration</span>
            <span>{plan.duration_days} days</span>
          </div>
        </div>
        <div className="mt-4 flex justify-center">
          <div
            className={cn(
              "w-5 h-5 border-2 rounded-full flex items-center justify-center transition-all",
              isSelected
                ? "bg-blue-500 border-blue-500"
                : "border-gray-300 dark:border-gray-600"
            )}
          >
            {isSelected && <Check className="w-3 h-3 text-white" />}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}