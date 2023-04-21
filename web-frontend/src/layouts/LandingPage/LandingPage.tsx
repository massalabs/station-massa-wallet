// react function component with typescript that is a layout for the landing page
// and have the background image svg

import React from 'react';

interface LandingPageProps {
  children: React.ReactNode;
}

export default function LandingPage(props: LandingPageProps) {
  return (
    <div className="bg-landing-page bg-no-repeat bg-cover bg-center min-h-screen">
      {props.children}
    </div>
  );
}
