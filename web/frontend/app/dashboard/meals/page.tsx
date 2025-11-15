'use client';

import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from '@/components/ui/dialog';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { Loader2, Plus, Trash2, CalendarIcon, Coffee, Sun, Moon, Cookie } from 'lucide-react';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { Calendar } from '@/components/ui/calendar';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { format } from 'date-fns';

const MEAL_TYPES = [
  { value: 'breakfast', label: 'Breakfast', icon: Coffee, color: 'amber' },
  { value: 'lunch', label: 'Lunch', icon: Sun, color: 'blue' },
  { value: 'dinner', label: 'Dinner', icon: Moon, color: 'indigo' },
  { value: 'snack', label: 'Snack', icon: Cookie, color: 'rose' },
];

interface FoodItem {
  id: number;
  name: string;
  food_id?: number;
  food_name?: string;
  portion?: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
}

interface Meal {
  id: number;
  date: string;
  meal_type: string;
  foods: FoodItem[];
  total_calories: number;
  total_protein: number;
  total_carbs: number;
  total_fat: number;
}

export default function MealsPage() {
  const [meals, setMeals] = useState<Meal[]>([]);
  const [foods, setFoods] = useState<FoodItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const { toast } = useToast();

  const [formData, setFormData] = useState({
    meal_type: 'breakfast',
    food_items: [{ food_id: '', portion: '' }],
  });

  const loadMeals = React.useCallback(async () => {
    setLoading(true);
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const result = await apiClient.getMeals(dateStr);
      if (result.code === 0 && result.data && Array.isArray(result.data)) {
        setMeals(result.data as Meal[]);
      }
    } catch {
      // Error handled silently
    } finally {
      setLoading(false);
    }
  }, [selectedDate]);

  const loadFoods = React.useCallback(async () => {
    try {
      const result = await apiClient.getFoods(1, 100);
      if (result.code === 0 && result.data && Array.isArray(result.data)) {
        setFoods(result.data as FoodItem[]);
      }
    } catch {
      // Error handled silently
    }
  }, []);

  useEffect(() => {
    loadMeals();
    loadFoods();
  }, [loadMeals, loadFoods]);



  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const mealData = {
      date: format(selectedDate, 'yyyy-MM-dd'),
      meal_type: formData.meal_type,
      foods: formData.food_items
        .filter(item => item.food_id && item.portion)
        .map(item => ({
          food_id: parseInt(item.food_id),
          portion: parseFloat(item.portion),
        })),
    };

    if (mealData.foods.length === 0) {
      toast({
        title: 'Error',
        description: 'Please add at least one food item',
        variant: 'destructive',
      });
      return;
    }

    try {
      const result = await apiClient.createMeal(mealData);
      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'Meal logged successfully',
        });
        setIsAddOpen(false);
        resetForm();
        loadMeals();
      } else {
        toast({
          title: 'Error',
          description: result.message,
          variant: 'destructive',
        });
      }
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to log meal',
        variant: 'destructive',
      });
    }
  };

  const handleDelete = async () => {
    if (!deleteId) return;

    try {
      const result = await apiClient.deleteMeal(deleteId);
      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'Meal deleted successfully',
        });
        setDeleteId(null);
        loadMeals();
      } else {
        toast({
          title: 'Error',
          description: result.message,
          variant: 'destructive',
        });
      }
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to delete meal',
        variant: 'destructive',
      });
    }
  };

  const resetForm = () => {
    setFormData({
      meal_type: 'breakfast',
      food_items: [{ food_id: '', portion: '' }],
    });
  };

  const addFoodItem = () => {
    setFormData({
      ...formData,
      food_items: [...formData.food_items, { food_id: '', portion: '' }],
    });
  };

  const removeFoodItem = (index: number) => {
    setFormData({
      ...formData,
      food_items: formData.food_items.filter((_, i) => i !== index),
    });
  };

  const updateFoodItem = (index: number, field: string, value: string) => {
    const newItems = [...formData.food_items];
    newItems[index] = { ...newItems[index], [field]: value };
    setFormData({ ...formData, food_items: newItems });
  };

  const groupedMeals = MEAL_TYPES.reduce((acc, type) => {
    acc[type.value] = meals.filter(m => m.meal_type === type.value);
    return acc;
  }, {} as Record<string, Meal[]>);

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-white">
            Meal Logs
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            Track your daily meals and nutrition intake
          </p>
        </div>
        <div className="flex gap-3">
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" className="w-64 justify-start text-left">
                <CalendarIcon className="mr-2 h-4 w-4" />
                {format(selectedDate, 'PPP')}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="end">
              <Calendar
                mode="single"
                selected={selectedDate}
                onSelect={(date) => date && setSelectedDate(date)}
                initialFocus
              />
            </PopoverContent>
          </Popover>
          <Dialog open={isAddOpen} onOpenChange={(open) => {
            setIsAddOpen(open);
            if (!open) resetForm();
          }}>
            <DialogTrigger asChild>
              <Button className="bg-gradient-to-r from-emerald-600 to-teal-600">
                <Plus className="h-4 w-4 mr-2" />
                Log Meal
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle>Log a Meal</DialogTitle>
                <DialogDescription>
                  Record what you ate for {format(selectedDate, 'PPP')}
                </DialogDescription>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <Label htmlFor="meal_type">Meal Type *</Label>
                  <Select value={formData.meal_type} onValueChange={(value) => setFormData({ ...formData, meal_type: value })}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {MEAL_TYPES.map((type) => (
                        <SelectItem key={type.value} value={type.value}>
                          {type.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-3">
                  <Label>Food Items *</Label>
                  {formData.food_items.map((item, index) => (
                    <div key={index} className="flex gap-3">
                      <Select
                        value={item.food_id}
                        onValueChange={(value) => updateFoodItem(index, 'food_id', value)}
                      >
                        <SelectTrigger className="flex-1">
                          <SelectValue placeholder="Select food" />
                        </SelectTrigger>
                        <SelectContent>
                          {foods.map((food) => (
                            <SelectItem key={food.id} value={String(food.id)}>
                              {food.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Input
                        type="number"
                        step="0.1"
                        placeholder="Portion (100g)"
                        value={item.portion}
                        onChange={(e) => updateFoodItem(index, 'portion', e.target.value)}
                        className="w-32"
                      />
                      {formData.food_items.length > 1 && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => removeFoodItem(index)}
                        >
                          <Trash2 className="h-4 w-4 text-red-600" />
                        </Button>
                      )}
                    </div>
                  ))}
                  <Button type="button" variant="outline" size="sm" onClick={addFoodItem}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Food
                  </Button>
                </div>

                <DialogFooter>
                  <Button type="button" variant="outline" onClick={() => setIsAddOpen(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" className="bg-gradient-to-r from-emerald-600 to-teal-600">
                    Log Meal
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
        </div>
      ) : (
        <div className="space-y-6">
          {MEAL_TYPES.map((type) => {
            const Icon = type.icon;
            const typeMeals = groupedMeals[type.value] || [];
            
            return (
              <Card key={type.value}>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Icon className={`h-5 w-5 text-${type.color}-600`} />
                    {type.label}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {typeMeals.length > 0 ? (
                    <div className="space-y-4">
                      {typeMeals.map((meal) => (
                        <div key={meal.id} className="border border-gray-200 dark:border-gray-800 rounded-lg p-4">
                          <div className="flex justify-between items-start mb-3">
                            <div className="space-y-1">
                              {meal.foods.map((food, idx) => (
                                <p key={idx} className="text-sm text-gray-900 dark:text-white">
                                  {food.food_name} - {food.portion}g
                                </p>
                              ))}
                            </div>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => setDeleteId(meal.id)}
                            >
                              <Trash2 className="h-4 w-4 text-red-600" />
                            </Button>
                          </div>
                          <div className="grid grid-cols-4 gap-4 pt-3 border-t border-gray-200 dark:border-gray-800">
                            <div>
                              <p className="text-xs text-gray-600 dark:text-gray-400">Calories</p>
                              <p className="text-sm font-semibold text-gray-900 dark:text-white">{meal.total_calories}</p>
                            </div>
                            <div>
                              <p className="text-xs text-gray-600 dark:text-gray-400">Protein</p>
                              <p className="text-sm font-semibold text-gray-900 dark:text-white">{meal.total_protein}g</p>
                            </div>
                            <div>
                              <p className="text-xs text-gray-600 dark:text-gray-400">Carbs</p>
                              <p className="text-sm font-semibold text-gray-900 dark:text-white">{meal.total_carbs}g</p>
                            </div>
                            <div>
                              <p className="text-xs text-gray-600 dark:text-gray-400">Fat</p>
                              <p className="text-sm font-semibold text-gray-900 dark:text-white">{meal.total_fat}g</p>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-center py-8 text-gray-500 dark:text-gray-400 text-sm">
                      No meals logged for {type.label.toLowerCase()} yet
                    </p>
                  )}
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}

      <AlertDialog open={deleteId !== null} onOpenChange={(open) => !open && setDeleteId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete meal?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently remove this meal from your logs.
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
