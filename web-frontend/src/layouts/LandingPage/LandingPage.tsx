// react function component with typescript that is a layout for the landing page
// and have the background image svg

import React from 'react';

import { FiX } from 'react-icons/fi';
import { useLocation, Link } from 'react-router-dom';

import { useQuery } from '@/custom/api/useQuery';
import { routeFor } from '@/utils';

interface LandingPageProps {
  children: React.ReactNode;
}

function CloseElement() {
  const { pathname } = useLocation();
  const query = useQuery();
  const isRoot = pathname === '/account-select';
  const rootValue = query.get('from') ? `${query.get('from')}/home` : undefined;

  return !isRoot ? (
    <Link
      className="fixed flex w-full justify-end text-f-primary pt-8 pr-8"
      to={routeFor(rootValue ?? 'index')}
    >
      <button
        data-testid="popup-modal-header-close"
        className="text-neutral bg-primary rounded-lg text-sm p-1.5 ml-auto inline-flex items-center
                        hover:bg-tertiary hover:text-c-primary"
        type="button"
      >
        <FiX className="w-7 h-7" />
      </button>
    </Link>
  ) : null;
}

export default function LandingPage(props: LandingPageProps) {
  return (
    <div
      className={`bg-primary bg-landing-page bg-no-repeat bg-cover
      bg-center min-h-screen`}
    >
      <CloseElement />
      {props.children}
    </div>
  );
}
