import { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function SignLayout(props: LayoutProps) {
  return (
    <div className="bg-primary h-screen flex justify-center items-center p-14">
      <div
        className="flex flex-col justify-center items-center
        text-f-primary w-fit h-full"
      >
        {props.children}
      </div>
    </div>
  );
}
