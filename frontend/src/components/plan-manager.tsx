'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
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
import { Edit, Trash2, Plus, DollarSign, Calendar } from 'lucide-react';
import { Plan, CreatePlanRequest, UpdatePlanRequest, apiClient } from '@/lib/api';
import { toast } from 'sonner';

interface PlanManagerProps {
  plans: Plan[];
  token: string;
  onPlansChange: () => void;
  isLoading: boolean;
}

interface PlanFormData {
  name: string;
  description: string;
  price_usd: string;
  duration_days: string;
}

export function PlanManager({ plans, token, onPlansChange, isLoading }: PlanManagerProps) {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingPlan, setEditingPlan] = useState<Plan | null>(null);
  const [formData, setFormData] = useState<PlanFormData>({
    name: '',
    description: '',
    price_usd: '',
    duration_days: '',
  });
  const [formError, setFormError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      price_usd: '',
      duration_days: '',
    });
    setFormError('');
    setShowCreateForm(false);
    setEditingPlan(null);
  };

  const handleCreateClick = () => {
    resetForm();
    setShowCreateForm(true);
  };

  const handleEditClick = (plan: Plan) => {
    console.log('✅ Edit button clicked for plan:', plan);
    setFormData({
      name: plan.name,
      description: plan.description,
      price_usd: plan.price_usd.toString(),
      duration_days: plan.duration_days.toString(),
    });
    setEditingPlan(plan);
    setShowCreateForm(true);
    setFormError(''); // Clear any previous errors
  };

  const validateForm = (): boolean => {
    if (!formData.name.trim()) {
      setFormError('Plan name is required');
      return false;
    }
    if (!formData.description.trim()) {
      setFormError('Description is required');
      return false;
    }
    const price = parseFloat(formData.price_usd);
    if (isNaN(price) || price <= 0) {
      setFormError('Price must be a positive number');
      return false;
    }
    const duration = parseInt(formData.duration_days);
    if (isNaN(duration) || duration <= 0) {
      setFormError('Duration must be a positive number');
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
      // Create data matching Go DTO expectations exactly
      const priceValue = parseFloat(formData.price_usd);
      const durationValue = parseInt(formData.duration_days, 10);

      // Validate parsed values
      if (isNaN(priceValue) || priceValue <= 0) {
        setFormError('Price must be a valid positive number');
        return;
      }
      if (isNaN(durationValue) || durationValue <= 0) {
        setFormError('Duration must be a valid positive number');
        return;
      }

      const planData = {
        name: formData.name.trim(),
        description: formData.description.trim(),
        price_usd: priceValue,     // Go expects float64
        duration_days: durationValue,  // Go expects int64
      };

      console.log('🚀 Submitting plan data:', planData);
      console.log('📊 Data types:', {
        price_usd_type: typeof planData.price_usd,
        duration_days_type: typeof planData.duration_days,
        name_length: planData.name.length,
        description_length: planData.description.length
      });

      let response;
      if (editingPlan) {
        console.log('📝 Updating existing plan:', editingPlan.id);
        response = await apiClient.updatePlan(editingPlan.id, planData as UpdatePlanRequest, token);
      } else {
        console.log('➕ Creating new plan');
        response = await apiClient.createPlan(planData as CreatePlanRequest, token);
      }

      console.log('✅ API Response:', response);

      if (response.status === 'ok') {
        console.log('✅ Plan saved successfully!');
        toast.success(
          editingPlan ? 'Plan updated successfully!' : 'Plan created successfully!',
          { description: `${planData.name} - $${planData.price_usd}` }
        );
        onPlansChange();
        resetForm();
      } else {
        console.error('❌ Plan save failed:', response);
        toast.error('Failed to save plan');
        setFormError(response.message || 'Failed to save plan');
      }
    } catch (error) {
      console.error('Submit error:', error);
      setFormError('Failed to save plan. Please try again.');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (planId: string) => {
    try {
      const response = await apiClient.deletePlan(planId, token);
      if (response.status === 'ok') {
        toast.success('Plan deleted successfully!');
        onPlansChange();
      } else {
        toast.error('Failed to delete plan');
      }
    } catch (error) {
      console.error('Delete error:', error);
      toast.error('Failed to delete plan');
    }
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading plans...</div>;
  }

  return (
    <div className="space-y-6">
      {/* Create/Edit Form */}
      {showCreateForm && (
        <div className="border rounded-lg p-4 mb-4">
          <h3 className="font-semibold mb-4">
            {editingPlan ? 'Edit Plan' : 'Create Plan'}
          </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Plan Name</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                    placeholder="e.g. Premium Plan"
                    required
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="price">Price (USD)</Label>
                  <Input
                    id="price"
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.price_usd}
                    onChange={(e) => setFormData(prev => ({ ...prev, price_usd: e.target.value }))}
                    placeholder="9.99"
                    required
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Input
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="Describe the plan features and benefits"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="duration">Duration (Days)</Label>
                <Input
                  id="duration"
                  type="number"
                  min="1"
                  value={formData.duration_days}
                  onChange={(e) => setFormData(prev => ({ ...prev, duration_days: e.target.value }))}
                  placeholder="30"
                  required
                />
              </div>

              {formError && (
                <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md">
                  {formError}
                </div>
              )}

              <div className="flex gap-2">
                <Button type="submit" disabled={formLoading}>
                  {formLoading ? 'Saving...' : (editingPlan ? 'Update Plan' : 'Create Plan')}
                </Button>
                <Button type="button" variant="outline" onClick={resetForm}>
                  Cancel
                </Button>
              </div>
            </form>
        </div>
      )}

      {/* Plans List */}
      <div className="flex justify-between items-center mb-4">
        <span className="text-sm text-muted-foreground">{plans.length} plans total</span>
        {!showCreateForm && (
          <Button size="sm" onClick={handleCreateClick}>
            <Plus className="w-4 h-4 mr-2" />
            Add Plan
          </Button>
        )}
      </div>

      {plans.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          No plans found. Create your first plan to get started.
        </div>
      ) : (
                <TooltipProvider>
          <div className="space-y-4">
            {plans.map((plan) => (
              <div key={plan.id} className="border rounded-lg p-4 bg-card">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-3">
                      <h4 className="font-semibold">{plan.name}</h4>
                      <Badge variant="secondary" className="flex items-center gap-1">
                        <DollarSign className="w-3 h-3" />
                        {plan.price_usd}
                      </Badge>
                      <Badge variant="outline" className="flex items-center gap-1">
                        <Calendar className="w-3 h-3" />
                        {plan.duration_days} days
                      </Badge>
                    </div>
                    <p className="text-muted-foreground mb-2">{plan.description}</p>
                    <div className="text-xs text-muted-foreground">
                      Plan ID: {plan.id}
                    </div>
                  </div>
                  
                  <div className="flex gap-2 ml-4">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleEditClick(plan)}
                          className="h-8 w-8 p-0"
                        >
                          <Edit className="w-3 h-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>Edit plan</p>
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
                          <AlertDialogTitle>Delete Plan</AlertDialogTitle>
                          <AlertDialogDescription>
                            Are you sure you want to delete the plan "{plan.name}"? This action cannot be undone.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleDelete(plan.id)} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
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