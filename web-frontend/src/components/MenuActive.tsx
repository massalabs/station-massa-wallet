import React from 'react';

interface MenuActiveProps {
  children: React.ReactNode;
}

// TODO: Urbane font family

export default function MenuActive(props: MenuActiveProps) {
  return (
    <p className="text-base font-semibold text-black leading-5 tracking-massa">
      {props.children}
    </p>
  );
}
