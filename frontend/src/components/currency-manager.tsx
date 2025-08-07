'use client';

import { useState, useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Edit, Trash2, Plus, Coins, Network, Search, X } from 'lucide-react';
import { Currency, CreateCurrencyRequest, UpdateCurrencyRequest, apiClient } from '@/lib/api';
import { toast } from 'sonner';

interface CurrencyManagerProps {
  currencies: Currency[];
  token: string;
  onCurrenciesChange: () => void;
  isLoading: boolean;
}

interface CurrencyFormData {
  code: string;
  name: string;
  network: string;
  is_token: boolean;
  contract_addr: string;
  enabled: boolean;
}

export function CurrencyManager({ currencies, token, onCurrenciesChange, isLoading }: CurrencyManagerProps) {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingCurrency, setEditingCurrency] = useState<Currency | null>(null);
  const [formData, setFormData] = useState<CurrencyFormData>({
    code: '',
    name: '',
    network: '',
    is_token: false,
    contract_addr: '',
    enabled: true,
  });
  const [formError, setFormError] = useState('');
  const [formLoading, setFormLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const resetForm = () => {
    setFormData({
      code: '',
      name: '',
      network: '',
      is_token: false,
      contract_addr: '',
      enabled: true,
    });
    setFormError('');
    setShowCreateForm(false);
    setEditingCurrency(null);
  };

  const handleCreateClick = () => {
    resetForm();
    setShowCreateForm(true);
  };

  const handleEditClick = (currency: Currency) => {
    setFormData({
      code: currency.code,
      name: currency.name,
      network: currency.network || '',
      is_token: currency.is_token || false,
      contract_addr: currency.contract_addr || '',
      enabled: currency.enabled || currency.is_active,
    });
    setEditingCurrency(currency);
    setShowCreateForm(true);
  };

  const validateForm = (): boolean => {
    if (!formData.code.trim()) {
      setFormError('Currency code is required');
      return false;
    }
    if (!formData.name.trim()) {
      setFormError('Currency name is required');
      return false;
    }
    if (!formData.network.trim()) {
      setFormError('Network is required');
      return false;
    }
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');

    if (!validateForm()) {
      return;
    }

    setFormLoading(true);

    try {
      const currencyData = {
        code: formData.code.trim().toUpperCase(),
        name: formData.name.trim(),
        network: formData.network.trim(),
        is_token: formData.is_token,
        contract_addr: formData.contract_addr.trim(),
        enabled: formData.enabled,
      };

      console.log('Submitting currency data:', currencyData); // Debug log

      let response;
      if (editingCurrency) {
        // Update existing currency
        response = await apiClient.updateCurrency(editingCurrency.code, currencyData as UpdateCurrencyRequest, token);
      } else {
        // Create new currency
        response = await apiClient.createCurrency(currencyData as CreateCurrencyRequest, token);
      }

      console.log('Currency API Response:', response); // Debug log

      if (response.status === 'ok') {
        toast.success(
          editingCurrency ? 'Currency updated successfully!' : 'Currency created successfully!',
          { description: `${currencyData.code} - ${currencyData.name}` }
        );
        onCurrenciesChange();
        resetForm();
      } else {
        toast.error('Failed to save currency');
        setFormError(response.message || 'Failed to save currency');
      }
    } catch (error) {
      setFormError('Failed to save currency. Please try again.');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (currencyCode: string) => {
    try {
      const response = await apiClient.deleteCurrency(currencyCode, token);
      if (response.status === 'ok') {
        toast.success('Currency deleted successfully!');
        onCurrenciesChange();
      } else {
        toast.error('Failed to delete currency', {
          description: response.message || 'Unknown error occurred'
        });
      }
    } catch (error) {
      console.error('Currency delete error:', error);
      toast.error('Failed to delete currency', {
        description: 'Network error occurred'
      });
    }
  };

  // Filter currencies based on search term
  const filteredCurrencies = useMemo(() => {
    if (!searchTerm) return currencies;
    
    return currencies.filter(currency => {
      const code = (currency.code || '').toLowerCase();
      const name = (currency.name || '').toLowerCase();
      const network = (currency.network || '').toLowerCase();
      const search = searchTerm.toLowerCase();
      
      return code.includes(search) || name.includes(search) || network.includes(search);
    });
  }, [currencies, searchTerm]);

  if (isLoading) {
    return <div className="text-center py-8">Loading currencies...</div>;
  }

  return (
    <div className="space-y-6">
      {/* Search */}
      <div className="flex items-center gap-2 mb-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input 
            placeholder="Search currencies by code, name, or network..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10 pr-10"
          />
          {searchTerm && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setSearchTerm('')}
              className="absolute right-1 top-1/2 transform -translate-y-1/2 h-6 w-6 p-0"
            >
              <X className="w-3 h-3" />
            </Button>
          )}
        </div>
      </div>
      {/* Create/Edit Form */}
      {showCreateForm && (
        <div className="border rounded-lg p-4 mb-4">
          <h3 className="font-semibold mb-4">
            {editingCurrency ? 'Edit Currency' : 'Create Currency'}
          </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="code">Currency Code</Label>
                  <Input
                    id="code"
                    value={formData.code}
                    onChange={(e) => setFormData(prev => ({ ...prev, code: e.target.value }))}
                    placeholder="e.g. BTC, ETH, TRX"
                    required
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="name">Currency Name</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                    placeholder="e.g. Bitcoin, Ethereum, Tron"
                    required
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="network">Network</Label>
                <Input
                  id="network"
                  value={formData.network}
                  onChange={(e) => setFormData(prev => ({ ...prev, network: e.target.value }))}
                  placeholder="e.g. TRC20, ERC20, BEP20"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="contract_addr">Contract Address (Optional)</Label>
                <Input
                  id="contract_addr"
                  value={formData.contract_addr}
                  onChange={(e) => setFormData(prev => ({ ...prev, contract_addr: e.target.value }))}
                  placeholder="Smart contract address if token"
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Switch
                    id="is_token"
                    checked={formData.is_token}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, is_token: checked }))}
                  />
                  <Label htmlFor="is_token">Is Token</Label>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    id="enabled"
                    checked={formData.enabled}
                    onCheckedChange={(checked) => setFormData(prev => ({ ...prev, enabled: checked }))}
                  />
                  <Label htmlFor="enabled">Enabled</Label>
                </div>
              </div>

              {formError && (
                <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md">
                  {formError}
                </div>
              )}

              <div className="flex gap-2">
                <Button type="submit" disabled={formLoading}>
                  {formLoading ? 'Saving...' : (editingCurrency ? 'Update Currency' : 'Create Currency')}
                </Button>
                <Button type="button" variant="outline" onClick={resetForm}>
                  Cancel
                </Button>
              </div>
            </form>
        </div>
      )}

      {/* Currencies List */}
      <div className="flex justify-between items-center mb-4">
        <span className="text-sm text-muted-foreground">{filteredCurrencies.length} currencies total</span>
        {!showCreateForm && (
          <Button size="sm" onClick={handleCreateClick}>
            <Plus className="w-4 h-4 mr-2" />
            Add Currency
          </Button>
        )}
      </div>

      {filteredCurrencies.length === 0 ? (
        searchTerm ? (
          <div className="text-center py-8 text-muted-foreground">
            No currencies found matching "{searchTerm}"
          </div>
        ) : (
          <div className="text-center py-12 text-gray-500">
            No currencies found. Create your first currency to get started.
          </div>
        )
      ) : (
                <TooltipProvider>
          <div className="space-y-4">
            {filteredCurrencies.map((currency) => (
              <div key={currency.code} className="border rounded-lg p-4 bg-card">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-3">
                      <h4 className="font-semibold flex items-center gap-2">
                        <Coins className="w-4 h-4" />
                        {currency.code}
                      </h4>
                      <Badge variant="secondary">
                        {currency.name}
                      </Badge>
                      <Badge variant={currency.enabled ? "default" : "secondary"}>
                        {currency.enabled ? 'Active' : 'Inactive'}
                      </Badge>
                      {currency.is_token && (
                        <Badge variant="outline">
                          Token
                        </Badge>
                      )}
                    </div>
                    
                    <div className="space-y-2 text-sm">
                      <div className="flex items-center gap-2">
                        <Network className="w-3 h-3" />
                        <span className="font-medium">Network:</span>
                        <span className="text-muted-foreground">{currency.network || 'N/A'}</span>
                      </div>
                      
                      {currency.contract_addr && (
                        <div className="flex items-start gap-2">
                          <span className="font-medium mt-0.5">Contract:</span>
                          <code className="bg-muted px-2 py-1 rounded text-xs font-mono break-all">
                            {currency.contract_addr}
                          </code>
                        </div>
                      )}
                    </div>
                  </div>
                  
                  <div className="flex gap-2 ml-4">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleEditClick(currency)}
                          className="h-8 w-8 p-0"
                        >
                          <Edit className="w-3 h-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>Edit currency</p>
                      </TooltipContent>
                    </Tooltip>
                    
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button
                          variant="outline"
                          size="sm"
                          className="h-8 w-8 p-0 text-destructive hover:text-destructive"
                        >
                          <Trash2 className="w-3 h-3" />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Delete Currency</AlertDialogTitle>
                          <AlertDialogDescription>
                            Are you sure you want to delete the currency "{currency.code}"? This action cannot be undone and may affect existing payments.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleDelete(currency.code)} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                            Delete
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </TooltipProvider>
      )}
    </div>
  );
}