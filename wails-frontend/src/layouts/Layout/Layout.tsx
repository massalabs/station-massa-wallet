import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function Layout(props: LayoutProps) {
  return (
    <div className="bg-primary min-h-screen">
      <div
        className="flex flex-col justify-center h-screen
        max-w-sm text-f-primary m-auto"
      >
        {props.children}
      </div>
    </div>
  );
}
