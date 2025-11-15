'use client';

import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { Loader2, TrendingUp, CalendarIcon } from 'lucide-react';
import { BarChart, Bar, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { Calendar } from '@/components/ui/calendar';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { format } from 'date-fns';

interface NutritionValue {
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
}

interface DailyNutritionData {
  actual: NutritionValue;
  goal: NutritionValue;
}

interface MonthlyNutritionData {
  date: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
}

export default function NutritionPage() {
  const [loading, setLoading] = useState(true);
  const [view, setView] = useState<'daily' | 'monthly'>('daily');
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const [dailyData, setDailyData] = useState<DailyNutritionData | null>(null);
  const [monthlyData, setMonthlyData] = useState<MonthlyNutritionData[]>([]);
  const { toast } = useToast();

  const loadDailyNutrition = React.useCallback(async () => {
    setLoading(true);
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const result = await apiClient.getDailyNutrition(dateStr);
      if (result.code === 0 && result.data) {
        setDailyData(result.data as DailyNutritionData);
      }
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to load nutrition data',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  }, [selectedDate, toast]);

  const loadMonthlyNutrition = React.useCallback(async () => {
    setLoading(true);
    try {
      const year = selectedDate.getFullYear();
      const month = selectedDate.getMonth() + 1;
      const result = await apiClient.getMonthlyNutrition(year, month);
      if (result.code === 0 && result.data && Array.isArray(result.data)) {
        setMonthlyData(result.data as MonthlyNutritionData[]);
      }
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to load nutrition data',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  }, [selectedDate, toast]);

  useEffect(() => {
    if (view === 'daily') {
      loadDailyNutrition();
    } else {
      loadMonthlyNutrition();
    }
  }, [view, loadDailyNutrition, loadMonthlyNutrition]);



  const chartData = view === 'daily' && dailyData ? [
    { name: 'Protein', actual: dailyData.actual?.protein || 0, goal: dailyData.goal?.protein || 0 },
    { name: 'Carbs', actual: dailyData.actual?.carbs || 0, goal: dailyData.goal?.carbs || 0 },
    { name: 'Fat', actual: dailyData.actual?.fat || 0, goal: dailyData.goal?.fat || 0 },
    { name: 'Fiber', actual: dailyData.actual?.fiber || 0, goal: dailyData.goal?.fiber || 0 },
  ] : [];

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-white">
            Nutrition Analysis
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            Track and analyze your nutrition intake over time
          </p>
        </div>
        <div className="flex gap-3">
          <Select value={view} onValueChange={(v) => setView(v as 'daily' | 'monthly')}>
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="daily">Daily View</SelectItem>
              <SelectItem value="monthly">Monthly View</SelectItem>
            </SelectContent>
          </Select>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" className="w-64 justify-start text-left">
                <CalendarIcon className="mr-2 h-4 w-4" />
                {format(selectedDate, view === 'daily' ? 'PPP' : 'MMMM yyyy')}
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
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
        </div>
      ) : (
        <>
          {view === 'daily' && dailyData && (
            <>
              <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
                      Total Calories
                    </CardTitle>
                    <TrendingUp className="h-4 w-4 text-emerald-600" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-gray-900 dark:text-white">
                      {dailyData.actual?.calories || 0}
                    </div>
                    <p className="text-xs text-gray-600 dark:text-gray-400">
                      Goal: {dailyData.goal?.calories || 0} kcal
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
                      Protein
                    </CardTitle>
                    <TrendingUp className="h-4 w-4 text-blue-600" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-gray-900 dark:text-white">
                      {dailyData.actual?.protein || 0}g
                    </div>
                    <p className="text-xs text-gray-600 dark:text-gray-400">
                      Goal: {dailyData.goal?.protein || 0}g
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
                      Carbohydrates
                    </CardTitle>
                    <TrendingUp className="h-4 w-4 text-amber-600" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-gray-900 dark:text-white">
                      {dailyData.actual?.carbs || 0}g
                    </div>
                    <p className="text-xs text-gray-600 dark:text-gray-400">
                      Goal: {dailyData.goal?.carbs || 0}g
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
                      Fat
                    </CardTitle>
                    <TrendingUp className="h-4 w-4 text-rose-600" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-gray-900 dark:text-white">
                      {dailyData.actual?.fat || 0}g
                    </div>
                    <p className="text-xs text-gray-600 dark:text-gray-400">
                      Goal: {dailyData.goal?.fat || 0}g
                    </p>
                  </CardContent>
                </Card>
              </div>

              <Card>
                <CardHeader>
                  <CardTitle>Actual vs Goal Comparison</CardTitle>
                </CardHeader>
                <CardContent>
                  <ResponsiveContainer width="100%" height={300}>
                    <BarChart data={chartData}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="name" />
                      <YAxis />
                      <Tooltip />
                      <Legend />
                      <Bar dataKey="actual" fill="#10b981" name="Actual" />
                      <Bar dataKey="goal" fill="#94a3b8" name="Goal" />
                    </BarChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>
            </>
          )}

          {view === 'monthly' && monthlyData.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Monthly Nutrition Trends</CardTitle>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={400}>
                  <LineChart data={monthlyData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="date" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Line type="monotone" dataKey="calories" stroke="#10b981" name="Calories" />
                    <Line type="monotone" dataKey="protein" stroke="#3b82f6" name="Protein" />
                    <Line type="monotone" dataKey="carbs" stroke="#f59e0b" name="Carbs" />
                    <Line type="monotone" dataKey="fat" stroke="#ef4444" name="Fat" />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          )}

          {view === 'monthly' && monthlyData.length === 0 && (
            <Card>
              <CardContent className="text-center py-12">
                <TrendingUp className="h-16 w-16 mx-auto mb-4 text-gray-400" />
                <p className="text-gray-600 dark:text-gray-400">
                  No nutrition data for this month
                </p>
              </CardContent>
            </Card>
          )}
        </>
      )}
    </div>
  );
}
