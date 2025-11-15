import { Loader2 } from '@/components/icons';

export default function Loading() {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <Loader2 className="h-12 w-12 animate-spin text-emerald-600 mx-auto" />
        <p className="mt-4 text-gray-600 dark:text-gray-400">Loading AI chat...</p>
      </div>
    </div>
  );
}
