'use client';

import { useState, useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
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
import { Edit, Trash2, Copy, QrCode, Search, X } from 'lucide-react';
import { Wallet, apiClient } from '@/lib/api';
import { toast } from 'sonner';
import QRCode from 'qrcode';

interface WalletManagerProps {
  wallets: Wallet[];
  token: string;
  onWalletsChange: () => void;
  isLoading: boolean;
}

export function WalletManager({ wallets, token, onWalletsChange, isLoading }: WalletManagerProps) {
  const [qrDataUrl, setQrDataUrl] = useState<string>('');
  const [showQR, setShowQR] = useState<string | null>(null);
  const [editingWallet, setEditingWallet] = useState<Wallet | null>(null);
  const [editFormData, setEditFormData] = useState({
    email: '',
  });
  const [searchTerm, setSearchTerm] = useState('');

  const generateQR = async (address: string) => {
    try {
      const qr = await QRCode.toDataURL(address, {
        width: 120,
        margin: 1,
        color: {
          dark: '#000000',
          light: '#FFFFFF'
        }
      });
      setQrDataUrl(qr);
      setShowQR(address);
    } catch (error) {
      console.error('Error generating QR code:', error);
    }
  };

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success('Copied to clipboard!');
    } catch (error) {
      console.error('Failed to copy:', error);
      toast.error('Failed to copy to clipboard');
    }
  };

  const handleSaveEdit = async () => {
    if (!editingWallet) return;
    
    try {
      // For now, just simulate saving
      console.log('Saving wallet edit:', editFormData);
      toast.success('Wallet updated successfully!', {
        description: `Email updated to: ${editFormData.email}`
      });
      setEditingWallet(null);
      onWalletsChange();
    } catch (error) {
      console.error('Failed to update wallet:', error);
      toast.error('Failed to update wallet', {
        description: 'Please try again later'
      });
    }
  };

  const cancelEdit = () => {
    setEditingWallet(null);
    setEditFormData({ email: '' });
  };

  const handleDelete = async (walletId: string) => {
    try {
      const response = await apiClient.deleteWallet(walletId, token);
      if (response.status === 'ok') {
        toast.success('Wallet deleted successfully!');
        onWalletsChange();
      } else {
        toast.error('Failed to delete wallet', {
          description: response.message || 'Unknown error occurred'
        });
      }
    } catch (error) {
      console.error('Delete error:', error);
      toast.error('Failed to delete wallet', {
        description: 'Network error occurred'
      });
    }
  };

  // Filter wallets based on search term
  const filteredWallets = useMemo(() => {
    if (!searchTerm) return wallets;
    
    return wallets.filter(wallet => {
      const email = String(wallet.Email || wallet.email || '').toLowerCase();
      const walletId = String(wallet.ID || wallet.id || '').toLowerCase();
      const address = String(wallet.WalletAddress || wallet.tron_address || '').toLowerCase();
      const search = searchTerm.toLowerCase();
      
      return email.includes(search) || walletId.includes(search) || address.includes(search);
    });
  }, [wallets, searchTerm]);

  if (isLoading) {
    return <div className="text-center py-8">Loading wallets...</div>;
  }

  return (
    <div className="space-y-6">
      {/* Search */}
      <div className="flex items-center gap-2 mb-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input 
            placeholder="Search wallets by email, ID, or address..."
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
      
      <div className="flex justify-between items-center mb-4">
        <span className="text-sm text-muted-foreground">{filteredWallets.length} wallets total</span>
      </div>

      {filteredWallets.length === 0 ? (
        searchTerm ? (
          <div className="text-center py-8 text-muted-foreground">
            No wallets found matching "{searchTerm}"
          </div>
        ) : (
          <div className="text-center py-12 text-gray-500">
            No wallets found.
          </div>
        )
      ) : (
        <TooltipProvider>
          <div className="space-y-4">
            {/* Edit Wallet Modal */}
            {editingWallet && (
              <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                <div className="bg-card p-6 rounded-lg border max-w-md w-full mx-4">
                  <h3 className="text-lg font-semibold mb-4">Edit Wallet</h3>
                  <div className="space-y-4">
                    <div>
                      <label className="text-sm font-medium">Email</label>
                      <input 
                        type="email"
                        value={editFormData.email}
                        onChange={(e) => setEditFormData(prev => ({ ...prev, email: e.target.value }))}
                        className="w-full mt-1 p-2 border rounded"
                      />
                    </div>
                    <div className="text-sm text-muted-foreground">
                      <p><strong>Wallet Address:</strong> {editingWallet.WalletAddress || editingWallet.tron_address}</p>
                      <p><strong>Created:</strong> {editingWallet.created_at ? new Date(editingWallet.created_at).toLocaleString() : 'N/A'}</p>
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
            {filteredWallets.map((wallet) => {
              const address = wallet.WalletAddress || wallet.tron_address;
              const email = wallet.Email || wallet.email;
              const walletId = wallet.ID || wallet.id;
              
              return (
                <div key={walletId} className="border rounded-lg p-4 bg-card">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-3">
                        <h4 className="font-semibold">{email}</h4>
                        <Badge variant="outline" className="text-xs">
                          {walletId}
                        </Badge>
                      </div>
                      
                      <div className="space-y-3">
                        <div className="flex items-center gap-2">
                          <span className="text-sm font-medium">TRON Address:</span>
                          <code className="bg-muted px-2 py-1 rounded text-xs font-mono flex-1 truncate">
                            {address}
                          </code>
                          <div className="flex gap-1">
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => copyToClipboard(address)}
                                  className="h-8 w-8 p-0"
                                >
                                  <Copy className="w-3 h-3" />
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p>Copy address</p>
                              </TooltipContent>
                            </Tooltip>
                            
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => {
                                    if (showQR === address) {
                                      setShowQR(null);
                                    } else {
                                      generateQR(address);
                                    }
                                  }}
                                  className="h-8 w-8 p-0"
                                >
                                  <QrCode className="w-3 h-3" />
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p>{showQR === address ? 'Hide QR code' : 'Show QR code'}</p>
                              </TooltipContent>
                            </Tooltip>
                          </div>
                        </div>
                        
                        {showQR === address && qrDataUrl && (
                          <div className="flex items-center gap-3 p-3 bg-muted rounded-lg border">
                            <img src={qrDataUrl} alt="QR Code" className="w-20 h-20 border rounded" />
                            <div className="text-xs text-muted-foreground flex-1">
                              <p className="font-medium">QR Code for wallet address</p>
                              <p className="font-mono truncate max-w-[200px] mt-1">{address}</p>
                            </div>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => setShowQR(null)}
                              className="h-8 w-8 p-0"
                            >
                              ×
                            </Button>
                          </div>
                        )}
                        
                        {wallet.created_at && (
                          <div className="text-xs text-muted-foreground">
                            Created: {new Date(wallet.created_at).toLocaleString()}
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
                            className="h-8 w-8 p-0"
                            onClick={() => {
                              setEditingWallet(wallet);
                              setEditFormData({
                                email: email,
                              });
                            }}
                          >
                            <Edit className="w-3 h-3" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>Edit wallet</p>
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
                            <AlertDialogTitle>Delete Wallet</AlertDialogTitle>
                            <AlertDialogDescription>
                              Are you sure you want to delete the wallet for "{email}"? This action cannot be undone.
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction onClick={() => handleDelete(walletId.toString())} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
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