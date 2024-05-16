import { useEffect, useState } from 'react';

import { Client } from '@massalabs/massa-web3';
import { useParams } from 'react-router-dom';

import { prepareSCCall } from '@/utils/prepareSCCall';

export function usePrepareScCall() {
  const { nickname } = useParams();
  const [client, setClient] = useState<Client>();

  useEffect(() => {
    if (!nickname) {
      throw new Error('Nickname not found');
    }
    prepareSCCall(nickname).then((result) => {
      setClient(result?.client);
    });
  }, [nickname, setClient]);

  return { client };
}
