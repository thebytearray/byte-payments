// Types based on your Go backend DTOs
export interface Plan {
  id: string;
  name: string;
  description: string;
  price_usd: number;
  duration_days: number;
}

export interface Currency {
  id: string;
  code: string;
  name: string;
  symbol: string;
  is_active: boolean;
}

export interface PaymentData {
  payment_id: string;
  plan_id: string;
  email: string;
  status: string;
  trx_amount: string;
  trx_wallet_address: string;
  qr_image: string;
  created_at: string;
}

export interface ApiResponse<T> {
  status: string;
  message?: string;
  data?: T;
}

export interface VerificationResponse {
  is_valid: boolean;
  verification_token?: string;
}



const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<ApiResponse<T>> {
    const url = `${API_BASE_URL}${endpoint}`;

    const response = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async getPlans(): Promise<ApiResponse<Plan[]>> {
    return this.request<Plan[]>("/api/v1/plans", {
      method: "GET",
    });
  }

  async getCurrencies(): Promise<ApiResponse<Currency[]>> {
    return this.request<Currency[]>("/api/v1/currencies", {
      method: "GET", 
    });
  }

  async sendVerificationCode(email: string): Promise<ApiResponse<null>> {
    return this.request<null>("/api/v1/verification/send-code", {
      method: "POST",
      body: JSON.stringify({ email }),
    });
  }

  async verifyCode(
    email: string,
    code: string,
  ): Promise<ApiResponse<VerificationResponse>> {
    return this.request<VerificationResponse>(
      "/api/v1/verification/verify-code",
      {
        method: "POST",
        body: JSON.stringify({ email, code }),
      },
    );
  }

  async createPayment(
    email: string,
    planId: string,
    currencyCode: string,
    verificationToken: string,
  ): Promise<ApiResponse<PaymentData>> {
    return this.request<PaymentData>("/api/v1/payments/create", {
      method: "POST",
      body: JSON.stringify({
        email,
        plan_id: planId,
        currency_code: currencyCode,
        verification_token: verificationToken,
      }),
    });
  }

  async getPaymentStatus(paymentId: string): Promise<ApiResponse<PaymentData>> {
    return this.request<PaymentData>(`/api/v1/payments/${paymentId}/status`);
  }

  async cancelPayment(paymentId: string): Promise<ApiResponse<null>> {
    return this.request<null>(`/api/v1/payments/${paymentId}/cancel`, {
      method: "PATCH",
    });
  }
}

export const apiClient = new ApiClient();
