import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function SignLayout(props: LayoutProps) {
  return (
    <div className="bg-primary h-screen flex justify-center p-24">
      <div
        className="flex flex-col justify-center
        text-f-primary w-fit h-fit"
      >
        {props.children}
      </div>
    </div>
  );
}
