'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { AdminInfo, LoginResponse, apiClient } from '@/lib/api';

interface AdminAuthContextType {
  admin: AdminInfo | null;
  token: string | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  isAuthenticated: boolean;
}

const AdminAuthContext = createContext<AdminAuthContextType | undefined>(undefined);

export const useAdminAuth = () => {
  const context = useContext(AdminAuthContext);
  if (context === undefined) {
    throw new Error('useAdminAuth must be used within an AdminAuthProvider');
  }
  return context;
};

export const AdminAuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [admin, setAdmin] = useState<AdminInfo | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check for existing token in localStorage
    const savedToken = localStorage.getItem('admin_token');
    const savedAdmin = localStorage.getItem('admin_info');
    
    if (savedToken && savedAdmin) {
      setToken(savedToken);
      setAdmin(JSON.parse(savedAdmin));
    }
    setIsLoading(false);
  }, []);

  const login = async (username: string, password: string): Promise<boolean> => {
    try {
      const response = await apiClient.adminLogin(username, password);
      
      if (response.status === 'ok' && response.data) {
        const { token: newToken, admin: adminInfo } = response.data;
        
        setToken(newToken);
        setAdmin(adminInfo);
        
        // Save to localStorage
        localStorage.setItem('admin_token', newToken);
        localStorage.setItem('admin_info', JSON.stringify(adminInfo));
        
        return true;
      }
      return false;
    } catch (error) {
      console.error('Login failed:', error);
      return false;
    }
  };

  const logout = () => {
    setAdmin(null);
    setToken(null);
    localStorage.removeItem('admin_token');
    localStorage.removeItem('admin_info');
  };

  const isAuthenticated = !!(admin && token);

  const value: AdminAuthContextType = {
    admin,
    token, 
    login,
    logout,
    isLoading,
    isAuthenticated,
  };

  return (
    <AdminAuthContext.Provider value={value}>
      {children}
    </AdminAuthContext.Provider>
  );
};