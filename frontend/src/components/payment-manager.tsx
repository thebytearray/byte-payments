'use client';

import { useState, useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
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
import { Edit, Trash2, Copy, Search, X } from 'lucide-react';
import { Payment, apiClient } from '@/lib/api';
import { toast } from 'sonner';

interface PaymentManagerProps {
  payments: Payment[];
  token: string;
  onPaymentsChange: () => void;
  isLoading: boolean;
}

function getStatusColor(status: string) {
  switch (status.toLowerCase()) {
    case 'completed':
      return 'bg-green-100 text-green-800';
    case 'pending':
      return 'bg-yellow-100 text-yellow-800';
    case 'expired':
      return 'bg-red-100 text-red-800';
    case 'cancelled':
      return 'bg-gray-100 text-gray-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
}

export function PaymentManager({ payments, token, onPaymentsChange, isLoading }: PaymentManagerProps) {
  const [editingPayment, setEditingPayment] = useState<Payment | null>(null);
  const [editFormData, setEditFormData] = useState({
    status: '',
    paid_amount: '',
  });
  const [searchTerm, setSearchTerm] = useState('');

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success('Copied to clipboard!');
    } catch (error) {
      console.error('Failed to copy:', error);
      toast.error('Failed to copy');
    }
  };

  // Filter payments based on search term
  const filteredPayments = useMemo(() => {
    if (!searchTerm) return payments;
    
    return payments.filter(payment => {
      const email = (payment.UserEmail || payment.user_email || '').toLowerCase();
      const paymentId = (payment.ID || payment.id || '').toLowerCase();
      const status = (payment.Status || payment.status || '').toLowerCase();
      const search = searchTerm.toLowerCase();
      
      return email.includes(search) || paymentId.includes(search) || status.includes(search);
    });
  }, [payments, searchTerm]);

  const handleEditClick = (payment: Payment) => {
    const status = payment.Status || payment.status;
    const paidAmount = payment.PaidAmountTRX || payment.paid_amount_trx;
    
    setEditingPayment(payment);
    setEditFormData({
      status: status,
      paid_amount: paidAmount?.toString() || '0',
    });
  };

  const handleSaveEdit = async () => {
    if (!editingPayment) return;
    
    try {
      // For now, just simulate saving
      console.log('Saving payment edit:', editFormData);
      alert('Payment updated successfully!');
      setEditingPayment(null);
      onPaymentsChange();
    } catch (error) {
      console.error('Failed to update payment:', error);
      alert('Failed to update payment');
    }
  };

  const cancelEdit = () => {
    setEditingPayment(null);
    setEditFormData({ status: '', paid_amount: '' });
  };

  const handleDelete = async (paymentId: string) => {
    if (!confirm('Are you sure you want to delete this payment?')) {
      return;
    }

    try {
      await apiClient.deletePayment(paymentId, token);
      onPaymentsChange();
    } catch (error) {
      alert('Failed to delete payment. Please try again.');
    }
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading payments...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center mb-4">
        <span className="text-sm text-gray-600">{payments.length} payments</span>
      </div>

      {payments.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          No payments found.
        </div>
            ) : (
        <TooltipProvider>
          <div className="space-y-4">
            {/* Edit Payment Modal */}
            {editingPayment && (
              <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                <div className="bg-card p-6 rounded-lg border max-w-md w-full mx-4">
                  <h3 className="text-lg font-semibold mb-4">Edit Payment</h3>
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="status">Status</Label>
                      <Select 
                        value={editFormData.status} 
                        onValueChange={(value) => setEditFormData(prev => ({ ...prev, status: value }))}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select status" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="pending">Pending</SelectItem>
                          <SelectItem value="completed">Completed</SelectItem>
                          <SelectItem value="cancelled">Cancelled</SelectItem>
                          <SelectItem value="expired">Expired</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="paid_amount">Paid Amount (TRX)</Label>
                      <Input 
                        type="number"
                        step="0.000001"
                        value={editFormData.paid_amount}
                        onChange={(e) => setEditFormData(prev => ({ ...prev, paid_amount: e.target.value }))}
                        placeholder="Enter amount"
                      />
                    </div>
                  </div>
                  <div className="flex gap-2 mt-6">
                    <Button onClick={handleSaveEdit} className="flex-1">
                      Save Changes
                    </Button>
                    <Button variant="outline" onClick={cancelEdit} className="flex-1">
                      Cancel
                    </Button>
                  </div>
                </div>
              </div>
            )}
            {payments.map((payment) => {
              const paymentId = payment.ID || payment.id;
              const email = payment.UserEmail || payment.user_email;
              const status = payment.Status || payment.status;
              const currency = payment.CurrencyCode || payment.currency_code;
              const requiredAmount = payment.AmountTRX || payment.required_amount_trx;
              const paidAmount = payment.PaidAmountTRX || payment.paid_amount_trx;
              const walletAddress = payment.Wallet?.WalletAddress || payment.wallet?.tron_address;
              
              return (
                <div key={paymentId} className="border rounded-lg p-4 bg-card">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-3">
                        <h4 className="font-semibold">{email}</h4>
                        <Badge className={getStatusColor(status)}>
                          {status.toUpperCase()}
                        </Badge>
                        <Badge variant="outline">{currency}</Badge>
                      </div>
                      
                      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 text-sm">
                        <div className="space-y-2">
                          <div className="flex items-center gap-2">
                            <span className="font-medium">Payment ID:</span>
                            <code className="bg-muted px-2 py-1 rounded text-xs font-mono flex-1 truncate">
                              {paymentId}
                            </code>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => copyToClipboard(paymentId)}
                                  className="h-6 w-6 p-0"
                                >
                                  <Copy className="w-3 h-3" />
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p>Copy payment ID</p>
                              </TooltipContent>
                            </Tooltip>
                          </div>
                          
                          <div>
                            <span className="font-medium">Plan:</span> 
                            <span className="ml-2">{payment.Plan?.name || payment.plan?.name || 'N/A'}</span>
                          </div>
                          
                          <div className="flex justify-between">
                            <span><span className="font-medium">Required:</span> {requiredAmount} TRX</span>
                            <span><span className="font-medium">Paid:</span> {paidAmount} TRX</span>
                          </div>
                        </div>
                        
                        <div className="space-y-2">
                          <div>
                            <span className="font-medium">Created:</span>
                            <span className="ml-2 text-muted-foreground">
                              {new Date(payment.CreatedAt || payment.created_at).toLocaleString()}
                            </span>
                          </div>
                          
                          {walletAddress && (
                            <div className="flex items-center gap-2">
                              <span className="font-medium">Wallet:</span>
                              <code className="bg-muted px-2 py-1 rounded text-xs font-mono flex-1 truncate">
                                {walletAddress}
                              </code>
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => copyToClipboard(walletAddress)}
                                    className="h-6 w-6 p-0"
                                  >
                                    <Copy className="w-3 h-3" />
                                  </Button>
                                </TooltipTrigger>
                                <TooltipContent>
                                  <p>Copy wallet address</p>
                                </TooltipContent>
                              </Tooltip>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex gap-2 ml-4">
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button 
                            variant="outline" 
                            size="sm" 
                            className="h-8 w-8 p-0"
                            onClick={() => handleEditClick(payment)}
                          >
                            <Edit className="w-3 h-3" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>Edit payment</p>
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
                            <AlertDialogTitle>Delete Payment</AlertDialogTitle>
                            <AlertDialogDescription>
                              Are you sure you want to delete payment for "{email}"? This action cannot be undone.
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction onClick={() => handleDelete(paymentId)} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </TooltipProvider>
      )}
    </div>
  );
}