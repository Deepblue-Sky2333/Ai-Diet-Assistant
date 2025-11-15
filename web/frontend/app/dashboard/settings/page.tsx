'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { Loader2, Save, CheckCircle2, SettingsIcon, User, Lock } from 'lucide-react';

export default function SettingsPage() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const { toast } = useToast();

  const [aiSettings, setAiSettings] = useState({
    provider: 'openai',
    api_key: '',
    api_key_masked: '', // 用于显示的掩码版本
    api_endpoint: '', // 修正：使用 api_endpoint 而不是 base_url
  });

  const [preferences, setPreferences] = useState({
    taste_preferences: '',
    dietary_restrictions: '',
    daily_calories_goal: '2000',
    daily_protein_goal: '150',
    daily_carbs_goal: '250',
    daily_fat_goal: '70',
    daily_fiber_goal: '30',
  });

  const [passwordForm, setPasswordForm] = useState({
    old_password: '',
    new_password: '',
    confirm_password: '',
  });

  useEffect(() => {
    loadSettings();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadSettings = async () => {
    setLoading(true);
    try {
      const result = await apiClient.getSettings();
      if (result.code === 0 && result.data) {
        const data = result.data as {
          ai_config?: {
            provider?: string;
            api_key_masked?: string;
            api_endpoint?: string;
          };
          user_preferences?: typeof preferences;
        };
        
        if (data.ai_config) {
          const apiKey = data.ai_config.api_key_masked || '';
          setAiSettings({
            provider: data.ai_config.provider || 'openai',
            api_key: '',
            api_key_masked: apiKey,
            api_endpoint: data.ai_config.api_endpoint || '',
          });
        }
        if (data.user_preferences) {
          setPreferences(data.user_preferences);
        }
      }
    } catch {
      // Error handled silently
    } finally {
      setLoading(false);
    }
  };

  const handleSaveAI = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // 如果已有密钥且用户没有输入新密钥，不需要更新
    if (!aiSettings.api_key && !aiSettings.api_key_masked) {
      toast({
        title: 'Error',
        description: 'Please enter an API key',
        variant: 'destructive',
      });
      return;
    }

    setSaving(true);

    try {
      const result = await apiClient.updateAISettings(
        aiSettings.provider,
        aiSettings.api_key || undefined, // 只在有新密钥时发送
        aiSettings.api_endpoint || undefined
      );

      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'AI settings saved successfully',
        });
        // 清空输入框，重新加载设置以获取掩码版本
        setAiSettings({ ...aiSettings, api_key: '' });
        loadSettings();
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
        description: 'Failed to save AI settings',
        variant: 'destructive',
      });
    } finally {
      setSaving(false);
    }
  };

  const handleTestAI = async () => {
    setTesting(true);
    try {
      const result = await apiClient.testAIConnection();
      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'AI connection test successful',
        });
      } else {
        toast({
          title: 'Error',
          description: result.message || 'Connection test failed',
          variant: 'destructive',
        });
      }
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to test AI connection',
        variant: 'destructive',
      });
    } finally {
      setTesting(false);
    }
  };

  const handleSavePreferences = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);

    try {
      const result = await apiClient.updateUserPreferences(preferences);
      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'Preferences saved successfully',
        });
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
        description: 'Failed to save preferences',
        variant: 'destructive',
      });
    } finally {
      setSaving(false);
    }
  };

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (passwordForm.new_password !== passwordForm.confirm_password) {
      toast({
        title: 'Error',
        description: 'Passwords do not match',
        variant: 'destructive',
      });
      return;
    }

    setSaving(true);
    try {
      const result = await apiClient.changePassword(
        passwordForm.old_password,
        passwordForm.new_password
      );

      if (result.code === 0) {
        toast({
          title: 'Success',
          description: 'Password changed successfully',
        });
        setPasswordForm({ old_password: '', new_password: '', confirm_password: '' });
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
        description: 'Failed to change password',
        variant: 'destructive',
      });
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    );
  }

  return (
    <div className="p-8 space-y-6">
      <div>
        <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-white">
          Settings
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          Manage your account, preferences, and AI configuration
        </p>
      </div>

      <Tabs defaultValue="ai" className="space-y-6">
        <TabsList>
          <TabsTrigger value="ai" className="gap-2">
            <SettingsIcon className="h-4 w-4" />
            AI Configuration
          </TabsTrigger>
          <TabsTrigger value="preferences" className="gap-2">
            <User className="h-4 w-4" />
            Preferences
          </TabsTrigger>
          <TabsTrigger value="security" className="gap-2">
            <Lock className="h-4 w-4" />
            Security
          </TabsTrigger>
        </TabsList>

        <TabsContent value="ai">
          <Card>
            <CardHeader>
              <CardTitle>AI Provider Configuration</CardTitle>
              <CardDescription>
                Configure your AI service provider for diet planning and chat features
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSaveAI} className="space-y-4">
                <div>
                  <Label htmlFor="provider">AI Provider *</Label>
                  <Select
                    value={aiSettings.provider}
                    onValueChange={(value) => setAiSettings({ ...aiSettings, provider: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="openai">OpenAI</SelectItem>
                      <SelectItem value="deepseek">DeepSeek</SelectItem>
                      <SelectItem value="custom">Custom Provider</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="api_key">API Key *</Label>
                  <div className="space-y-2">
                    {aiSettings.api_key_masked && (
                      <div className="flex items-center gap-2 p-3 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-800">
                        <span className="text-sm text-gray-600 dark:text-gray-400 font-mono">
                          {aiSettings.api_key_masked}
                        </span>
                        <span className="text-xs text-emerald-600 dark:text-emerald-400 ml-auto">
                          Configured
                        </span>
                      </div>
                    )}
                    <Input
                      id="api_key"
                      type="password"
                      placeholder={aiSettings.api_key_masked ? "Enter new API key to update" : "Enter your API key"}
                      value={aiSettings.api_key}
                      onChange={(e) => setAiSettings({ ...aiSettings, api_key: e.target.value })}
                      required={!aiSettings.api_key_masked}
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {aiSettings.api_key_masked 
                        ? "Leave empty to keep current key. Enter a new key to update."
                        : "Your API key will be encrypted and stored securely on the server."}
                    </p>
                  </div>
                </div>
                {aiSettings.provider === 'custom' && (
                  <div>
                    <Label htmlFor="api_endpoint">API Endpoint</Label>
                    <Input
                      id="api_endpoint"
                      placeholder="https://api.example.com/v1"
                      value={aiSettings.api_endpoint}
                      onChange={(e) => setAiSettings({ ...aiSettings, api_endpoint: e.target.value })}
                    />
                  </div>
                )}
                <div className="flex gap-3">
                  <Button type="submit" disabled={saving} className="bg-gradient-to-r from-emerald-600 to-teal-600">
                    {saving ? (
                      <>
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="h-4 w-4 mr-2" />
                        Save Settings
                      </>
                    )}
                  </Button>
                  <Button type="button" variant="outline" onClick={handleTestAI} disabled={testing}>
                    {testing ? (
                      <>
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                        Testing...
                      </>
                    ) : (
                      <>
                        <CheckCircle2 className="h-4 w-4 mr-2" />
                        Test Connection
                      </>
                    )}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="preferences">
          <Card>
            <CardHeader>
              <CardTitle>User Preferences</CardTitle>
              <CardDescription>
                Set your dietary preferences and nutrition goals
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSavePreferences} className="space-y-4">
                <div>
                  <Label htmlFor="taste_preferences">Taste Preferences</Label>
                  <Input
                    id="taste_preferences"
                    placeholder="e.g., spicy, sweet, savory..."
                    value={preferences.taste_preferences}
                    onChange={(e) => setPreferences({ ...preferences, taste_preferences: e.target.value })}
                  />
                </div>
                <div>
                  <Label htmlFor="dietary_restrictions">Dietary Restrictions</Label>
                  <Input
                    id="dietary_restrictions"
                    placeholder="e.g., vegetarian, gluten-free, dairy-free..."
                    value={preferences.dietary_restrictions}
                    onChange={(e) => setPreferences({ ...preferences, dietary_restrictions: e.target.value })}
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="daily_calories_goal">Daily Calories Goal (kcal)</Label>
                    <Input
                      id="daily_calories_goal"
                      type="number"
                      value={preferences.daily_calories_goal}
                      onChange={(e) => setPreferences({ ...preferences, daily_calories_goal: e.target.value })}
                    />
                  </div>
                  <div>
                    <Label htmlFor="daily_protein_goal">Daily Protein Goal (g)</Label>
                    <Input
                      id="daily_protein_goal"
                      type="number"
                      value={preferences.daily_protein_goal}
                      onChange={(e) => setPreferences({ ...preferences, daily_protein_goal: e.target.value })}
                    />
                  </div>
                  <div>
                    <Label htmlFor="daily_carbs_goal">Daily Carbs Goal (g)</Label>
                    <Input
                      id="daily_carbs_goal"
                      type="number"
                      value={preferences.daily_carbs_goal}
                      onChange={(e) => setPreferences({ ...preferences, daily_carbs_goal: e.target.value })}
                    />
                  </div>
                  <div>
                    <Label htmlFor="daily_fat_goal">Daily Fat Goal (g)</Label>
                    <Input
                      id="daily_fat_goal"
                      type="number"
                      value={preferences.daily_fat_goal}
                      onChange={(e) => setPreferences({ ...preferences, daily_fat_goal: e.target.value })}
                    />
                  </div>
                  <div>
                    <Label htmlFor="daily_fiber_goal">Daily Fiber Goal (g)</Label>
                    <Input
                      id="daily_fiber_goal"
                      type="number"
                      value={preferences.daily_fiber_goal}
                      onChange={(e) => setPreferences({ ...preferences, daily_fiber_goal: e.target.value })}
                    />
                  </div>
                </div>
                <Button type="submit" disabled={saving} className="bg-gradient-to-r from-emerald-600 to-teal-600">
                  {saving ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    <>
                      <Save className="h-4 w-4 mr-2" />
                      Save Preferences
                    </>
                  )}
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="security">
          <Card>
            <CardHeader>
              <CardTitle>Change Password</CardTitle>
              <CardDescription>
                Update your account password
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleChangePassword} className="space-y-4">
                <div>
                  <Label htmlFor="old_password">Current Password *</Label>
                  <Input
                    id="old_password"
                    type="password"
                    value={passwordForm.old_password}
                    onChange={(e) => setPasswordForm({ ...passwordForm, old_password: e.target.value })}
                    required
                  />
                </div>
                <div>
                  <Label htmlFor="new_password">New Password *</Label>
                  <Input
                    id="new_password"
                    type="password"
                    value={passwordForm.new_password}
                    onChange={(e) => setPasswordForm({ ...passwordForm, new_password: e.target.value })}
                    required
                  />
                </div>
                <div>
                  <Label htmlFor="confirm_password">Confirm New Password *</Label>
                  <Input
                    id="confirm_password"
                    type="password"
                    value={passwordForm.confirm_password}
                    onChange={(e) => setPasswordForm({ ...passwordForm, confirm_password: e.target.value })}
                    required
                  />
                </div>
                <Button type="submit" disabled={saving} className="bg-gradient-to-r from-emerald-600 to-teal-600">
                  {saving ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Changing...
                    </>
                  ) : (
                    <>
                      <Lock className="h-4 w-4 mr-2" />
                      Change Password
                    </>
                  )}
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
