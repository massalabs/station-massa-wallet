import React from 'react';

export interface BodyProps {
  children: React.ReactNode;
}

// TODO: Poppins font family

export default function Body(props: BodyProps) {
  return (
    <p className="text-base font-medium text-black leading-6 tracking-massa">
      {props.children}
    </p>
  );
}
