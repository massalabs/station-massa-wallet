import React from 'react';

export interface BannerProps {
  children: React.ReactNode;
}

// TODO: urban font family

export default function Banner(props: BannerProps) {
  return (
    <h1 className="text-center text-4xl font-semibold text-black leading-10 tracking-massa">
      {props.children}
    </h1>
  );
}
