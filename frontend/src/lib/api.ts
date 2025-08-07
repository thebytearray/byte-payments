// Types based on your Go backend DTOs
export interface Plan {
  id: string;
  name: string;
  description: string;
  price_usd: number;
  duration_days: number;
}

export interface Currency {
  code: string;
  name: string;
  network?: string;
  is_token?: boolean;
  contract_addr?: string;
  enabled: boolean;
  // For compatibility with existing components
  id: string;
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

export interface AdminInfo {
  id: number;
  username: string;
  email: string;
}

export interface LoginResponse {
  token: string;
  admin: AdminInfo;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface CreatePlanRequest {
  name: string;
  description: string;
  price_usd: number;  // Backend expects PriceUSD but we'll handle this in the request
  duration_days: number;  // Backend expects DurationDays but we'll handle this in the request
}

export interface UpdatePlanRequest {
  name: string;
  description: string;
  price_usd: number;
  duration_days: number;
}

export interface Wallet {
  ID: string;
  Email: string;
  WalletAddress: string;
  WalletSecret: string;
  // For compatibility with existing components
  id: number;
  email: string;
  tron_address: string;
  private_key: string;
  created_at: string;
  updated_at: string;
}

export interface Payment {
  ID: string;
  PlanID: string;
  Plan?: Plan;
  WalletID: string;
  Wallet?: Wallet;
  CurrencyCode: string;
  Currency?: Currency;
  AmountUSD: number;
  AmountTRX: number;
  UserEmail: string;
  Status: string;
  PaidAmountTRX: number;
  CreatedAt: string;
  UpdatedAt: string;
  // For compatibility with existing components
  id: string;
  user_email: string;
  plan_id: string;
  plan?: Plan;
  wallet?: Wallet;
  status: string;
  currency_code: string;
  required_amount_trx: number;
  paid_amount_trx: number;
  created_at: string;
  updated_at: string;
  expires_at: string;
}

export interface CreateCurrencyRequest {
  code: string;
  name: string;
  network: string;
  is_token: boolean;
  contract_addr: string;
  enabled: boolean;
}

export interface UpdateCurrencyRequest {
  code: string;
  name: string;
  network: string;
  is_token: boolean;
  contract_addr: string;
  enabled: boolean;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}



const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<ApiResponse<T>> {
    const url = `${API_BASE_URL}${endpoint}`;

    // Ensure proper headers without duplication
    const defaultHeaders: Record<string, string> = {
      "Content-Type": "application/json",
    };
    
    const headers = {
      ...defaultHeaders,
      ...options.headers,
    };

    console.log('🚀 API Request:', {
      url,
      method: options.method || 'GET',
      headers,
      body: options.body ? JSON.parse(options.body as string) : null
    });

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const responseData = await response.json();
    console.log('📥 API Response:', { 
      status: response.status, 
      ok: response.ok,
      statusText: response.statusText,
      data: responseData 
    });

    // Don't throw error, let components handle the response
    return responseData;
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

  // Admin endpoints
  async adminLogin(username: string, password: string): Promise<ApiResponse<LoginResponse>> {
    return this.request<LoginResponse>("/api/v1/admin/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    });
  }

  async createPlan(planData: CreatePlanRequest, token: string): Promise<ApiResponse<Plan>> {
    // Ensure exact Go DTO format
    const requestBody = {
      name: String(planData.name).trim(),
      description: String(planData.description).trim(),
      price_usd: Number(planData.price_usd),
      duration_days: Number(planData.duration_days),
    };
    
    // Validate before sending
    if (!requestBody.name || !requestBody.description) {
      throw new Error('Name and description are required');
    }
    if (isNaN(requestBody.price_usd) || requestBody.price_usd <= 0) {
      throw new Error('Price must be a positive number');
    }
    if (isNaN(requestBody.duration_days) || requestBody.duration_days <= 0) {
      throw new Error('Duration must be a positive number');
    }
    
    console.log('🚀 Creating plan API call:', requestBody);
    
    return this.request<Plan>("/api/v1/admin/plans", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestBody),
    });
  }

  async updatePlan(planId: string, planData: UpdatePlanRequest, token: string): Promise<ApiResponse<Plan>> {
    // Ensure exact Go DTO format
    const requestBody = {
      name: String(planData.name).trim(),
      description: String(planData.description).trim(),
      price_usd: Number(planData.price_usd),
      duration_days: Number(planData.duration_days),
    };
    
    // Validate before sending
    if (!requestBody.name || !requestBody.description) {
      throw new Error('Name and description are required');
    }
    if (isNaN(requestBody.price_usd) || requestBody.price_usd <= 0) {
      throw new Error('Price must be a positive number');
    }
    if (isNaN(requestBody.duration_days) || requestBody.duration_days <= 0) {
      throw new Error('Duration must be a positive number');
    }
    
    console.log('📝 Updating plan API call:', { planId, requestBody });
    
    return this.request<Plan>(`/api/v1/admin/plans/${planId}`, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestBody),
    });
  }

  async deletePlan(planId: string, token: string): Promise<ApiResponse<null>> {
    return this.request<null>(`/api/v1/admin/plans/${planId}`, {
      method: "DELETE", 
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }

  // Wallet management
  async getAllWallets(token: string): Promise<ApiResponse<Wallet[]>> {
    return this.request<Wallet[]>("/api/v1/admin/wallets", {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }

  async deleteWallet(walletId: string, token: string): Promise<ApiResponse<null>> {
    return this.request<null>(`/api/v1/admin/wallets/${walletId}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }

  // Payment management
  async getAllPayments(token: string): Promise<ApiResponse<Payment[]>> {
    return this.request<Payment[]>("/api/v1/admin/payments", {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }

  async deletePayment(paymentId: string, token: string): Promise<ApiResponse<null>> {
    return this.request<null>(`/api/v1/admin/payments/${paymentId}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }

  // Currency management
  async createCurrency(currencyData: CreateCurrencyRequest, token: string): Promise<ApiResponse<Currency>> {
    // Convert to backend expected format
    const backendData = {
      code: currencyData.code,
      name: currencyData.name,
      network: currencyData.network,
      is_token: currencyData.is_token,
      contract_addr: currencyData.contract_addr,
      enabled: currencyData.enabled,
    };
    
    return this.request<Currency>("/api/v1/admin/currencies", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(backendData),
    });
  }

  async updateCurrency(currencyCode: string, currencyData: UpdateCurrencyRequest, token: string): Promise<ApiResponse<Currency>> {
    // Convert to backend expected format
    const backendData = {
      code: currencyData.code,
      name: currencyData.name,
      network: currencyData.network,
      is_token: currencyData.is_token,
      contract_addr: currencyData.contract_addr,
      enabled: currencyData.enabled,
    };
    
    return this.request<Currency>(`/api/v1/admin/currencies/${currencyCode}`, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(backendData),
    });
  }

  async deleteCurrency(currencyCode: string, token: string): Promise<ApiResponse<null>> {
    return this.request<null>(`/api/v1/admin/currencies/${currencyCode}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }

  async changePassword(passwordData: ChangePasswordRequest, token: string): Promise<ApiResponse<null>> {
    return this.request<null>("/api/v1/admin/change-password", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(passwordData),
    });
  }
}

export const apiClient = new ApiClient();
