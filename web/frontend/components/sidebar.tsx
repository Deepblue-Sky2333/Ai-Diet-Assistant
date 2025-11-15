'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { apiClient } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { LayoutDashboard, ShoppingCart, Utensils, Calendar, BarChart3, MessageSquare, Settings, LogOut } from '@/components/icons';
import { useEffect, useRef, useState } from 'react';

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Supermarket', href: '/dashboard/supermarket', icon: ShoppingCart },
  { name: 'Meal Logs', href: '/dashboard/meals', icon: Utensils },
  { name: 'Diet Plans', href: '/dashboard/plans', icon: Calendar },
  { name: 'Nutrition', href: '/dashboard/nutrition', icon: BarChart3 },
  { name: 'AI Chat', href: '/dashboard/chat', icon: MessageSquare },
  { name: 'Settings', href: '/dashboard/settings', icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const { toast } = useToast();
  const [indicatorStyle, setIndicatorStyle] = useState({ top: 0, height: 0 });
  const navRefs = useRef<(HTMLAnchorElement | null)[]>([]);
  const [mounted, setMounted] = useState(false);

  const handleLogout = async () => {
    try {
      await apiClient.logout();
      toast({
        title: 'Logged out',
        description: 'See you next time!',
      });
      router.push('/login');
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to logout',
        variant: 'destructive',
      });
    }
  };

  const activeIndex = navigation.findIndex(item => pathname === item.href);

  useEffect(() => {
    setMounted(true);
  }, []);


  useEffect(() => {
    if (!mounted) return;
    
    const timer = setTimeout(() => {
      if (activeIndex !== -1 && navRefs.current[activeIndex]) {
        const activeElement = navRefs.current[activeIndex];
        if (activeElement) {
          setIndicatorStyle({
            top: activeElement.offsetTop,
            height: activeElement.offsetHeight,
          });
        }
      }
    }, 0);

    return () => clearTimeout(timer);
  }, [activeIndex, pathname, mounted]);


  return (
    <div className="flex h-full w-64 flex-col bg-white dark:bg-gray-950 border-r border-gray-200 dark:border-gray-800">
      <div className="flex h-16 items-center border-b border-gray-200 dark:border-gray-800 px-6">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center">
            <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <span className="text-lg font-semibold bg-gradient-to-r from-emerald-600 to-teal-600 bg-clip-text text-transparent">
            AI Diet
          </span>
        </div>
      </div>
      <nav className="flex-1 space-y-1 px-3 py-4 relative">
        {activeIndex !== -1 && indicatorStyle.height > 0 && (
          <div
            className="absolute left-3 w-[calc(100%-1.5rem)] bg-emerald-100 dark:bg-emerald-950 rounded-lg transition-all duration-300 ease-out"
            style={{
              top: `${indicatorStyle.top}px`,
              height: `${indicatorStyle.height}px`,
            }}
          />
        )}
        {navigation.map((item, index) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.name}
              href={item.href}
              ref={(el) => {
                navRefs.current[index] = el;
              }}
              className={cn(
                'relative flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors z-10',
                isActive
                  ? 'text-emerald-900 dark:text-emerald-100'
                  : 'text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-gray-100'
              )}
            >
              <item.icon className="h-5 w-5" />
              {item.name}
            </Link>
          );
        })}
      </nav>
      <div className="border-t border-gray-200 dark:border-gray-800 p-4">
        <Button
          variant="outline"
          className="w-full justify-start gap-3 border-gray-300 dark:border-gray-700"
          onClick={handleLogout}
        >
          <LogOut className="h-5 w-5" />
          Logout
        </Button>
      </div>
    </div>
  );
}
