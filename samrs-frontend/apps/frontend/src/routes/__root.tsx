import { Spinner } from '@/components/ui/spinner';
import { createRootRoute, Outlet } from '@tanstack/react-router';
import { Suspense } from 'react';

export const Loading = () => (
  <div className="flex items-center justify-center">
    <Spinner className='w-8 h-8' />
  </div>
);

export const rootRoute = createRootRoute({
  component: () => (
    <>
      <Suspense fallback={<Loading />}>
        <Outlet />
      </Suspense>
    </>
  ),
});