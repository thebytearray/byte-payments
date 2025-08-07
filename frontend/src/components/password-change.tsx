'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ChangePasswordRequest, apiClient } from '@/lib/api';
import { toast } from 'sonner';

interface PasswordChangeProps {
  token: string;
}

export function PasswordChange({ token }: PasswordChangeProps) {
  const [formData, setFormData] = useState({
    old_password: '',
    new_password: '',
    confirm_password: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (!formData.old_password || !formData.new_password) {
      setError('All fields are required');
      return;
    }

    if (formData.new_password.length < 6) {
      setError('New password must be at least 6 characters');
      return;
    }

    if (formData.new_password !== formData.confirm_password) {
      setError('New passwords do not match');
      return;
    }

    setIsLoading(true);

    try {
      const passwordData: ChangePasswordRequest = {
        old_password: formData.old_password.trim(),
        new_password: formData.new_password.trim(),
      };

      console.log('Submitting password change:', { ...passwordData, old_password: '[HIDDEN]', new_password: '[HIDDEN]' }); // Debug log

      const response = await apiClient.changePassword(passwordData, token);
      
      console.log('Password change response:', response); // Debug log
      
      if (response.status === 'ok') {
        toast.success('Password changed successfully!', {
          description: 'Your password has been updated securely'
        });
        setSuccess('Password changed successfully');
        setFormData({ old_password: '', new_password: '', confirm_password: '' });
      } else {
        toast.error('Failed to change password');
        setError(response.message || 'Failed to change password');
      }
    } catch (error) {
      console.error('Password change error:', error);
      setError('Failed to change password. Please check your old password.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>Change Password</CardTitle>
        <CardDescription>Update your admin password</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="old_password">Current Password</Label>
            <Input
              id="old_password"
              type="password"
              value={formData.old_password}
              onChange={(e) => setFormData(prev => ({ ...prev, old_password: e.target.value }))}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="new_password">New Password</Label>
            <Input
              id="new_password"
              type="password"
              value={formData.new_password}
              onChange={(e) => setFormData(prev => ({ ...prev, new_password: e.target.value }))}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="confirm_password">Confirm New Password</Label>
            <Input
              id="confirm_password"
              type="password"
              value={formData.confirm_password}
              onChange={(e) => setFormData(prev => ({ ...prev, confirm_password: e.target.value }))}
              required
            />
          </div>

          {error && (
            <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
              {error}
            </div>
          )}

          {success && (
            <div className="text-sm text-green-600 bg-green-50 dark:bg-green-900/20 dark:text-green-400 p-3 rounded-md">
              {success}
            </div>
          )}

          <Button type="submit" disabled={isLoading} className="w-full">
            {isLoading ? 'Changing...' : 'Change Password'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}