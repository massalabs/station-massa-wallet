import { ReactNode } from 'react';

interface DashboardLayoutProps {
  children: ReactNode;
}

export function DashboardLayout(props: DashboardLayoutProps) {
  return (
    <div className="bg-primary min-h-screen">
      <div
        className="flex flex-col justify-center h-screen
        max-w-xs min-w-fit text-f-primary m-auto"
      >
        {props.children}
      </div>
    </div>
  );
}
