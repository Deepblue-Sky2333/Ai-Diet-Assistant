'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [connectionError, setConnectionError] = useState<string | null>(null);
  const router = useRouter();
  const { toast } = useToast();

  const isDemoMode = process.env.NEXT_PUBLIC_DEMO_MODE === 'true';

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setConnectionError(null);

    try {
      const result = await apiClient.login(username, password);
      
      if (result.code === 0) {
        toast({
          title: 'Login successful',
          description: isDemoMode ? 'Welcome to AI Diet Assistant (Demo Mode)' : 'Welcome to AI Diet Assistant',
        });
        router.push('/dashboard');
      } else {
        toast({
          title: 'Login failed',
          description: result.message || 'Invalid credentials',
          variant: 'destructive',
        });
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unable to connect to server';
      setConnectionError(errorMessage);
      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleTestLogin = async () => {
    setLoading(true);
    setConnectionError(null);
    
    try {
      const result = await apiClient.loginWithTestAccount();
      
      if (result.code === 0) {
        toast({
          title: 'Test login successful',
          description: isDemoMode ? 'Logged in with test account (Demo Mode)' : 'Logged in with test account',
        });
        router.push('/dashboard');
      } else {
        toast({
          title: 'Test login failed',
          description: result.message || 'Unable to login with test account',
          variant: 'destructive',
        });
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unable to connect to server';
      setConnectionError(errorMessage);
      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-emerald-50 via-white to-teal-50 dark:from-gray-950 dark:via-gray-900 dark:to-emerald-950 p-4">
      <Card className="w-full max-w-md border-emerald-200 dark:border-emerald-900">
        <CardHeader className="space-y-1 text-center">
          <div className="flex justify-center mb-4">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center">
              <svg className="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
          </div>
          <CardTitle className="text-3xl font-bold tracking-tight bg-gradient-to-r from-emerald-600 to-teal-600 bg-clip-text text-transparent">
            AI Diet Assistant
          </CardTitle>
          <CardDescription className="text-base">
            Your personalized nutrition companion
          </CardDescription>
          {isDemoMode && (
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800">
              <svg className="w-4 h-4 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span className="text-xs font-medium text-blue-600 dark:text-blue-400">Demo Mode</span>
            </div>
          )}
        </CardHeader>
        <CardContent>
          {connectionError && !isDemoMode && (
            <Alert variant="destructive" className="mb-4">
              <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <AlertTitle>Connection Error</AlertTitle>
              <AlertDescription className="text-sm">
                {connectionError}
                <div className="mt-2 text-xs">
                  Make sure the backend API is running and NEXT_PUBLIC_API_URL is configured correctly in your environment variables.
                </div>
              </AlertDescription>
            </Alert>
          )}

          {isDemoMode && (
            <Alert className="mb-4 border-blue-200 bg-blue-50 dark:border-blue-800 dark:bg-blue-950">
              <svg className="h-4 w-4 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <AlertTitle className="text-blue-900 dark:text-blue-100">Demo Mode Active</AlertTitle>
              <AlertDescription className="text-sm text-blue-800 dark:text-blue-200">
                You can explore the UI with sample data. No backend connection required.
              </AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleLogin} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                type="text"
                placeholder="Enter your username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
                disabled={loading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="Enter your password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={loading}
              />
            </div>
            <Button 
              type="submit" 
              className="w-full bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700" 
              disabled={loading}
            >
              {loading ? 'Signing in...' : 'Sign in'}
            </Button>
          </form>

          <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-800">
            <p className="text-xs text-gray-500 text-center mb-3">
              Development use only
            </p>
            <Button
              type="button"
              variant="outline"
              className="w-full border-yellow-200 bg-yellow-50 hover:bg-yellow-100 text-yellow-800 dark:border-yellow-900 dark:bg-yellow-950 dark:hover:bg-yellow-900 dark:text-yellow-200"
              onClick={handleTestLogin}
              disabled={loading}
            >
              <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              Login with Test Account
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
