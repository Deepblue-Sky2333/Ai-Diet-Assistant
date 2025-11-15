'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from '@/components/ui/dialog';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { Loader2, Plus, Check, Trash2, Sparkles, Calendar, Utensils } from 'lucide-react';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { format } from 'date-fns';

interface Plan {
  id: number;
  date: string;
  meal_type: string;
  foods: Array<{
    food_id: number;
    food_name: string;
    portion: number;
  }>;
  total_calories: number;
  total_protein: number;
  total_carbs: number;
  total_fat: number;
  reason: string;
}

export default function PlansPage() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [isGenerateOpen, setIsGenerateOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const { toast } = useToast();

  const [formData, setFormData] = useState({
    days: '2',
    preferences: '',
  });

  useEffect(() => {
    loadPlans();
  }, []);

  const loadPlans = async () => {
    setLoading(true);
    try {
      const result = await apiClient.getPlans();
      if (result.code === 0 && result.data) {
        setPlans(result.data);
      }
    } catch (error) {
      console.error('[v0] Load plans error:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async (e: React.FormEvent) => {
    e.preventDefault();
    setGenerating(true);

    try {
      const result = await apiClient.generatePlan(
        parseInt(formData.days),
        formData.preferences
      );

      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'AI diet plan generated successfully',
        });
        setIsGenerateOpen(false);
        setFormData({ days: '2', preferences: '' });
        loadPlans();
      } else {
        toast({
          title: 'Error',
          description: result.message || 'Failed to generate plan',
          variant: 'destructive',
        });
      }
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to generate plan',
        variant: 'destructive',
      });
    } finally {
      setGenerating(false);
    }
  };

  const handleComplete = async (id: number) => {
    try {
      const result = await apiClient.completePlan(id);
      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'Plan completed and converted to meal log',
        });
        loadPlans();
      } else {
        toast({
          title: 'Error',
          description: result.message,
          variant: 'destructive',
        });
      }
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to complete plan',
        variant: 'destructive',
      });
    }
  };

  const handleDelete = async () => {
    if (!deleteId) return;

    try {
      const result = await apiClient.deletePlan(deleteId);
      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'Plan deleted successfully',
        });
        setDeleteId(null);
        loadPlans();
      } else {
        toast({
          title: 'Error',
          description: result.message,
          variant: 'destructive',
        });
      }
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete plan',
        variant: 'destructive',
      });
    }
  };

  const groupedPlans = plans.reduce((acc, plan) => {
    if (!acc[plan.date]) acc[plan.date] = [];
    acc[plan.date].push(plan);
    return acc;
  }, {} as Record<string, Plan[]>);

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-white">
            Diet Plans
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            AI-generated personalized meal plans based on your preferences
          </p>
        </div>
        <Dialog open={isGenerateOpen} onOpenChange={setIsGenerateOpen}>
          <DialogTrigger asChild>
            <Button className="bg-gradient-to-r from-purple-600 to-pink-600">
              <Sparkles className="h-4 w-4 mr-2" />
              Generate AI Plan
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Generate AI Diet Plan</DialogTitle>
              <DialogDescription>
                Let AI create a personalized meal plan based on your preferences
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleGenerate} className="space-y-4">
              <div>
                <Label htmlFor="days">Number of Days *</Label>
                <Input
                  id="days"
                  type="number"
                  min="1"
                  max="7"
                  value={formData.days}
                  onChange={(e) => setFormData({ ...formData, days: e.target.value })}
                  required
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Generate plans for 1-7 days
                </p>
              </div>
              <div>
                <Label htmlFor="preferences">Preferences & Requirements</Label>
                <Textarea
                  id="preferences"
                  placeholder="e.g., I prefer high protein meals, vegetarian options, avoid dairy..."
                  value={formData.preferences}
                  onChange={(e) => setFormData({ ...formData, preferences: e.target.value })}
                  rows={4}
                />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setIsGenerateOpen(false)}>
                  Cancel
                </Button>
                <Button 
                  type="submit" 
                  className="bg-gradient-to-r from-purple-600 to-pink-600"
                  disabled={generating}
                >
                  {generating ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Generating...
                    </>
                  ) : (
                    <>
                      <Sparkles className="h-4 w-4 mr-2" />
                      Generate
                    </>
                  )}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
        </div>
      ) : plans.length > 0 ? (
        <div className="space-y-6">
          {Object.entries(groupedPlans).map(([date, datePlans]) => (
            <Card key={date}>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Calendar className="h-5 w-5 text-emerald-600" />
                  {format(new Date(date), 'EEEE, MMMM d, yyyy')}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {datePlans.map((plan) => (
                  <div key={plan.id} className="border border-gray-200 dark:border-gray-800 rounded-lg p-4 space-y-3">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <Utensils className="h-4 w-4 text-emerald-600" />
                          <span className="font-semibold text-gray-900 dark:text-white capitalize">
                            {plan.meal_type}
                          </span>
                        </div>
                        <div className="space-y-1 mb-3">
                          {plan.foods.map((food, idx) => (
                            <p key={idx} className="text-sm text-gray-700 dark:text-gray-300">
                              • {food.food_name} ({food.portion}g)
                            </p>
                          ))}
                        </div>
                        {plan.reason && (
                          <div className="bg-purple-50 dark:bg-purple-950 rounded-lg p-3 mb-3">
                            <p className="text-xs font-medium text-purple-900 dark:text-purple-100 mb-1">
                              AI Recommendation
                            </p>
                            <p className="text-sm text-purple-700 dark:text-purple-300">
                              {plan.reason}
                            </p>
                          </div>
                        )}
                        <div className="grid grid-cols-4 gap-3 pt-3 border-t border-gray-200 dark:border-gray-800">
                          <div>
                            <p className="text-xs text-gray-600 dark:text-gray-400">Calories</p>
                            <p className="text-sm font-semibold text-gray-900 dark:text-white">{plan.total_calories}</p>
                          </div>
                          <div>
                            <p className="text-xs text-gray-600 dark:text-gray-400">Protein</p>
                            <p className="text-sm font-semibold text-gray-900 dark:text-white">{plan.total_protein}g</p>
                          </div>
                          <div>
                            <p className="text-xs text-gray-600 dark:text-gray-400">Carbs</p>
                            <p className="text-sm font-semibold text-gray-900 dark:text-white">{plan.total_carbs}g</p>
                          </div>
                          <div>
                            <p className="text-xs text-gray-600 dark:text-gray-400">Fat</p>
                            <p className="text-sm font-semibold text-gray-900 dark:text-white">{plan.total_fat}g</p>
                          </div>
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          className="border-emerald-600 text-emerald-600 hover:bg-emerald-50 dark:hover:bg-emerald-950"
                          onClick={() => handleComplete(plan.id)}
                        >
                          <Check className="h-4 w-4 mr-1" />
                          Complete
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => setDeleteId(plan.id)}
                        >
                          <Trash2 className="h-4 w-4 text-red-600" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <Card>
          <CardContent className="text-center py-12">
            <Sparkles className="h-16 w-16 mx-auto mb-4 text-gray-400" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              No diet plans yet
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-6">
              Generate your first AI-powered meal plan to get started
            </p>
            <Button 
              onClick={() => setIsGenerateOpen(true)}
              className="bg-gradient-to-r from-purple-600 to-pink-600"
            >
              <Sparkles className="h-4 w-4 mr-2" />
              Generate AI Plan
            </Button>
          </CardContent>
        </Card>
      )}

      <AlertDialog open={deleteId !== null} onOpenChange={(open) => !open && setDeleteId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete plan?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently remove this meal plan.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete} className="bg-red-600 hover:bg-red-700">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
