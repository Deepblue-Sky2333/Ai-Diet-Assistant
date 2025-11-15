'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { apiClient } from '@/lib/api';
import { Loader2, TrendingUp, Calendar, Apple, Utensils, ShoppingCart, MessageSquare } from '@/components/icons';
import Link from 'next/link';

/**
 * Type definitions matching backend response structure
 * Backend: internal/handler/dashboard_handler.go - GetDashboard()
 * 
 * Response structure:
 * {
 *   today_nutrition: { calories, protein, carbs, fat },
 *   nutrition_goal: { calories, protein, carbs, fat },
 *   upcoming_plans: [{ id, date, meal_type, reason }]
 * }
 */
interface NutritionData {
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
}

interface NutritionGoal {
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
}

interface Plan {
  id: number;
  date: string;
  meal_type: string;
  reason?: string;
}

interface DashboardData {
  today_nutrition: NutritionData;
  nutrition_goal: NutritionGoal;
  upcoming_plans: Plan[];
}

export default function DashboardPage() {
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<DashboardData | null>(null);

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      const result = await apiClient.getDashboard();
      if (result.code === 0 && result.data) {
        // Verify the data structure matches our expectations
        setData(result.data as DashboardData);
      }
    } catch (error) {
      console.error('[v0] Dashboard load error:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    );
  }

  // Extract nutrition data with proper defaults
  const todayNutrition: NutritionData = data?.today_nutrition || {
    calories: 0,
    protein: 0,
    carbs: 0,
    fat: 0,
  };
  
  const todayGoal: NutritionGoal = data?.nutrition_goal || {
    calories: 2000,
    protein: 150,
    carbs: 250,
    fat: 70,
  };

  const calculatePercentage = (actual: number, goal: number) => {
    if (!goal) return 0;
    return Math.min(Math.round((actual / goal) * 100), 100);
  };

  return (
    <div className="p-8 space-y-8">
      <div>
        <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-white">
          Dashboard
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          Welcome back! Here's your nutrition overview.
        </p>
      </div>

      {/* Today's Nutrition Stats */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="border-emerald-200 dark:border-emerald-900">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Calories
            </CardTitle>
            <Apple className="h-4 w-4 text-emerald-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900 dark:text-white">
              {todayNutrition.calories || 0}
            </div>
            <div className="mt-2">
              <div className="flex justify-between text-xs text-gray-600 dark:text-gray-400 mb-1">
                <span>Goal: {todayGoal.calories || 2000} kcal</span>
                <span>{calculatePercentage(todayNutrition.calories, todayGoal.calories)}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-800 rounded-full h-2">
                <div
                  className="bg-gradient-to-r from-emerald-500 to-teal-600 h-2 rounded-full transition-all"
                  style={{ width: `${calculatePercentage(todayNutrition.calories, todayGoal.calories)}%` }}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-blue-200 dark:border-blue-900">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Protein
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900 dark:text-white">
              {todayNutrition.protein || 0}g
            </div>
            <div className="mt-2">
              <div className="flex justify-between text-xs text-gray-600 dark:text-gray-400 mb-1">
                <span>Goal: {todayGoal.protein || 150}g</span>
                <span>{calculatePercentage(todayNutrition.protein, todayGoal.protein)}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-800 rounded-full h-2">
                <div
                  className="bg-blue-600 h-2 rounded-full transition-all"
                  style={{ width: `${calculatePercentage(todayNutrition.protein, todayGoal.protein)}%` }}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-amber-200 dark:border-amber-900">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Carbs
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-amber-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900 dark:text-white">
              {todayNutrition.carbs || 0}g
            </div>
            <div className="mt-2">
              <div className="flex justify-between text-xs text-gray-600 dark:text-gray-400 mb-1">
                <span>Goal: {todayGoal.carbs || 250}g</span>
                <span>{calculatePercentage(todayNutrition.carbs, todayGoal.carbs)}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-800 rounded-full h-2">
                <div
                  className="bg-amber-600 h-2 rounded-full transition-all"
                  style={{ width: `${calculatePercentage(todayNutrition.carbs, todayGoal.carbs)}%` }}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-rose-200 dark:border-rose-900">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Fat
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-rose-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900 dark:text-white">
              {todayNutrition.fat || 0}g
            </div>
            <div className="mt-2">
              <div className="flex justify-between text-xs text-gray-600 dark:text-gray-400 mb-1">
                <span>Goal: {todayGoal.fat || 70}g</span>
                <span>{calculatePercentage(todayNutrition.fat, todayGoal.fat)}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-800 rounded-full h-2">
                <div
                  className="bg-rose-600 h-2 rounded-full transition-all"
                  style={{ width: `${calculatePercentage(todayNutrition.fat, todayGoal.fat)}%` }}
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Upcoming Plans */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-xl">Upcoming Diet Plans</CardTitle>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              AI-generated meal plans for the next 2 days
            </p>
          </div>
          <Link href="/dashboard/plans">
            <Button variant="outline" size="sm">
              <Calendar className="h-4 w-4 mr-2" />
              View All
            </Button>
          </Link>
        </CardHeader>
        <CardContent>
          {data?.upcoming_plans && data.upcoming_plans.length > 0 ? (
            <div className="space-y-4">
              {data.upcoming_plans.slice(0, 3).map((plan: Plan) => (
                <div key={plan.id} className="flex items-start gap-4 p-4 rounded-lg border border-gray-200 dark:border-gray-800">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <span className="text-sm font-medium text-gray-900 dark:text-white">
                        {new Date(plan.date).toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })}
                      </span>
                      <span className="px-2 py-1 text-xs rounded-full bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-100">
                        {plan.meal_type}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {plan.reason || 'AI-generated meal plan'}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12 text-gray-500 dark:text-gray-400">
              <Calendar className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>No upcoming plans yet</p>
              <Link href="/dashboard/plans">
                <Button className="mt-4 bg-gradient-to-r from-emerald-600 to-teal-600">
                  Generate AI Plan
                </Button>
              </Link>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <div className="grid gap-6 md:grid-cols-3">
        <Link href="/dashboard/meals">
          <Card className="cursor-pointer hover:shadow-lg transition-shadow">
            <CardContent className="pt-6">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-xl bg-emerald-100 dark:bg-emerald-950 flex items-center justify-center">
                  <Utensils className="h-6 w-6 text-emerald-600" />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900 dark:text-white">Log a Meal</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">Record what you ate</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/supermarket">
          <Card className="cursor-pointer hover:shadow-lg transition-shadow">
            <CardContent className="pt-6">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-xl bg-blue-100 dark:bg-blue-950 flex items-center justify-center">
                  <ShoppingCart className="h-6 w-6 text-blue-600" />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900 dark:text-white">Manage Foods</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">Update your food library</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/chat">
          <Card className="cursor-pointer hover:shadow-lg transition-shadow">
            <CardContent className="pt-6">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-xl bg-purple-100 dark:bg-purple-950 flex items-center justify-center">
                  <MessageSquare className="h-6 w-6 text-purple-600" />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900 dark:text-white">Ask AI</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">Get diet advice</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  );
}
