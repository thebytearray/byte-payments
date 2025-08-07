'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAdminAuth } from '@/contexts/admin-auth-context';
import { Plan, Currency, Payment, Wallet, apiClient } from '@/lib/api';
import { AdminLayout } from '@/components/admin-layout';
import { PlanManager } from '@/components/plan-manager';
import { PaymentManager } from '@/components/payment-manager';
import { WalletManager } from '@/components/wallet-manager';
import { CurrencyManager } from '@/components/currency-manager';
import { PasswordChange } from '@/components/password-change';

export default function AdminDashboardPage() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [currencies, setCurrencies] = useState<Currency[]>([]);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [wallets, setWallets] = useState<Wallet[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<'plans' | 'payments' | 'wallets' | 'currencies' | 'settings'>('plans');
  
  const { token, isAuthenticated } = useAdminAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/admin');
      return;
    }
    
    loadAllData();
  }, [isAuthenticated, router]);

  const loadPlans = async () => {
    try {
      const response = await apiClient.getPlans();
      if (response.status === 'ok' && response.data) {
        setPlans(response.data);
      }
    } catch (err) {
      console.error('Failed to load plans:', err);
    }
  };

  const loadCurrencies = async () => {
    try {
      const response = await apiClient.getCurrencies();
      if (response.status === 'ok' && response.data) {
        setCurrencies(response.data);
      }
    } catch (err) {
      console.error('Failed to load currencies:', err);
    }
  };

  const loadPayments = async () => {
    try {
      const response = await apiClient.getAllPayments(token!);
      if (response.status === 'ok' && response.data) {
        setPayments(response.data);
      }
    } catch (err) {
      console.error('Failed to load payments:', err);
    }
  };

  const loadWallets = async () => {
    try {
      const response = await apiClient.getAllWallets(token!);
      if (response.status === 'ok' && response.data) {
        setWallets(response.data);
      }
    } catch (err) {
      console.error('Failed to load wallets:', err);
    }
  };

  const loadAllData = async () => {
    setIsLoading(true);
    try {
      await Promise.all([
        loadPlans(),
        loadCurrencies(),
        loadPayments(),
        loadWallets(),
      ]);
    } catch (err) {
      setError('Failed to load data');
    } finally {
      setIsLoading(false);
    }
  };



  if (!isAuthenticated) {
    return null;
  }

  const stats = {
    plans: plans.length,
    payments: payments.length,
    wallets: wallets.length,
    currencies: currencies.length,
  };

  return (
    <AdminLayout 
      activeTab={activeTab} 
      onTabChange={setActiveTab} 
      stats={stats}
    >
      <div className="bg-card rounded-lg border shadow-sm">
        <div className="p-6">
            {error && (
              <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md mb-4 border border-destructive/20">
                {error}
              </div>
            )}
            
            {activeTab === 'plans' && (
              <PlanManager 
                plans={plans}
                token={token!}
                onPlansChange={loadPlans}
                isLoading={isLoading}
              />
            )}
            
            {activeTab === 'payments' && (
              <PaymentManager 
                payments={payments}
                token={token!}
                onPaymentsChange={loadPayments}
                isLoading={isLoading}
              />
            )}
            
            {activeTab === 'wallets' && (
              <WalletManager 
                wallets={wallets}
                token={token!}
                onWalletsChange={loadWallets}
                isLoading={isLoading}
              />
            )}
            
            {activeTab === 'currencies' && (
              <CurrencyManager 
                currencies={currencies}
                token={token!}
                onCurrenciesChange={loadCurrencies}
                isLoading={isLoading}
              />
            )}
            
          {activeTab === 'settings' && (
            <div className="max-w-md">
              <PasswordChange token={token!} />
            </div>
          )}
        </div>
      </div>
    </AdminLayout>
  );
}