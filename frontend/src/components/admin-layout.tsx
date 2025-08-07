'use client';

import { useState, ReactNode } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { useAdminAuth } from '@/contexts/admin-auth-context';
import { useTheme } from 'next-themes';
import { 
  LayoutDashboard, 
  CreditCard, 
  Wallet, 
  Coins, 
  Settings, 
  LogOut,
  Menu,
  ChevronLeft,
  ChevronRight,
  Sun,
  Moon
} from 'lucide-react';

interface AdminLayoutProps {
  children: ReactNode;
  activeTab: string;
  onTabChange: (tab: string) => void;
  stats: {
    plans: number;
    payments: number;
    wallets: number;
    currencies: number;
  };
}

export function AdminLayout({ children, activeTab, onTabChange, stats }: AdminLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const { admin, logout } = useAdminAuth();
  const { theme, setTheme } = useTheme();

  const navigation = [
    { key: 'plans', label: 'Plans', icon: LayoutDashboard, count: stats.plans },
    { key: 'payments', label: 'Payments', icon: CreditCard, count: stats.payments },
    { key: 'wallets', label: 'Wallets', icon: Wallet, count: stats.wallets },
    { key: 'currencies', label: 'Currencies', icon: Coins, count: stats.currencies },
    { key: 'settings', label: 'Settings', icon: Settings, count: null },
  ];

  const handleLogout = () => {
    logout();
    window.location.href = '/admin';
  };

  return (
    <TooltipProvider>
      <div className="min-h-screen bg-background">
        {/* Mobile sidebar backdrop */}
        {sidebarOpen && (
          <div 
            className="fixed inset-0 z-40 lg:hidden bg-black/50"
            onClick={() => setSidebarOpen(false)}
          />
        )}

        {/* Sidebar */}
        <div className={`fixed inset-y-0 left-0 z-50 ${sidebarCollapsed ? 'w-16' : 'w-64'} bg-card/95 backdrop-blur-sm border-r border-border/50 transition-all duration-300 ease-in-out transform lg:translate-x-0 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}>
          {/* Sidebar header */}
          <div className="flex items-center justify-between p-4 border-b">
            {!sidebarCollapsed && (
              <h1 className="text-lg font-semibold">BytePayments</h1>
            )}
            <div className="flex items-center gap-2">
              <Button 
                variant="ghost" 
                size="sm" 
                className="hidden lg:flex"
                onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              >
                {sidebarCollapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
              </Button>
              <Button 
                variant="ghost" 
                size="sm" 
                className="lg:hidden"
                onClick={() => setSidebarOpen(false)}
              >
                <Menu className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Navigation */}
          <nav className="mt-6 px-2">
            <ul className="space-y-1">
              {navigation.map((item) => {
                const Icon = item.icon;
                const isActive = activeTab === item.key;
                
                const NavButton = (
                  <Button
                    variant={isActive ? "secondary" : "ghost"}
                    size="sm"
                    className={`w-full justify-start h-10 ${sidebarCollapsed ? 'px-2' : 'px-3'}`}
                    onClick={() => {
                      onTabChange(item.key);
                      setSidebarOpen(false);
                    }}
                  >
                    <Icon className="w-4 h-4 flex-shrink-0" />
                    {!sidebarCollapsed && (
                      <>
                        <span className="ml-3 truncate">{item.label}</span>
                        {item.count !== null && (
                          <Badge variant="outline" className="ml-auto text-xs">
                            {item.count}
                          </Badge>
                        )}
                      </>
                    )}
                  </Button>
                );
                
                return (
                  <li key={item.key}>
                    {sidebarCollapsed ? (
                      <Tooltip>
                        <TooltipTrigger asChild>
                          {NavButton}
                        </TooltipTrigger>
                        <TooltipContent side="right">
                          <p>{item.label}</p>
                        </TooltipContent>
                      </Tooltip>
                    ) : (
                      NavButton
                    )}
                  </li>
                );
              })}
            </ul>
          </nav>

          {/* Bottom spacing */}
          <div className="h-4"></div>
        </div>

        {/* Main content */}
        <div className={`transition-all duration-300 ${sidebarCollapsed ? 'lg:ml-16' : 'lg:ml-64'}`}>
          {/* Top bar */}
          <div className="bg-card/95 backdrop-blur-sm border-b border-border/50 sticky top-0 z-10">
            <div className="flex items-center justify-between px-6 py-4">
              <div className="flex items-center">
                <Button
                  variant="ghost"
                  size="sm"
                  className="lg:hidden mr-4"
                  onClick={() => setSidebarOpen(true)}
                >
                  <Menu className="w-4 h-4" />
                </Button>
                <h2 className="text-lg font-semibold">
                  {navigation.find(item => item.key === activeTab)?.label}
                </h2>
              </div>
              
              {/* Top right controls */}
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2">
                  <Avatar className="w-8 h-8">
                    <AvatarFallback className="text-sm">
                      {admin?.username?.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="hidden sm:block">
                    <p className="text-sm font-medium">{admin?.username}</p>
                    <p className="text-xs text-muted-foreground">{admin?.email}</p>
                  </div>
                </div>
                
                <div className="flex items-center gap-1">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                        className="h-9 w-9 p-0"
                      >
                        {theme === 'dark' ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Toggle theme</p>
                    </TooltipContent>
                  </Tooltip>
                  
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleLogout}
                        className="h-9 w-9 p-0 text-destructive hover:text-destructive"
                      >
                        <LogOut className="w-4 h-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Logout</p>
                    </TooltipContent>
                  </Tooltip>
                </div>
              </div>
            </div>
          </div>

          {/* Page content */}
          <div className="p-6">
            {children}
          </div>
        </div>
      </div>
    </TooltipProvider>
  );
}