import { faker } from '@faker-js/faker';
import { Response, Server } from 'miragejs';

import { RuleType, SignRule } from '../../models/ConfigModel';
import { AppSchema } from '../types';

export function routesForSignRules(server: Server) {
  server.post('accounts/:nickname/signrules', (schema: AppSchema, request) => {
    const { nickname } = request.params;
    const account = schema.findBy('account', { nickname });

    if (!account) {
      return new Response(
        404,
        {},
        {
          code: '404',
          message: 'Account not found',
        },
      );
    }

    const body = JSON.parse(request.requestBody);
    const { ruleType, contract, enabled, name } = body;

    if (!ruleType || !contract || enabled === undefined) {
      return new Response(
        400,
        {},
        {
          code: '400',
          message: 'Missing required fields',
        },
      );
    }

    if (
      ![RuleType.AutoSign, RuleType.DisablePasswordPrompt].includes(ruleType)
    ) {
      return new Response(
        422,
        {},
        {
          code: '422',
          message: 'Invalid rule type',
        },
      );
    }

    const id = faker.string.uuid();
    const rule: SignRule = {
      id,
      ruleType,
      contract,
      enabled,
      name,
    };

    schema.create('signRule', { ...rule, accountNickname: nickname });

    return new Response(200, {}, { id });
  });

  server.delete(
    'accounts/:nickname/signrules/:ruleId',
    (schema: AppSchema, request) => {
      const { nickname, ruleId } = request.params;
      const account = schema.findBy('account', { nickname });

      if (!account) {
        return new Response(
          404,
          {},
          {
            code: '404',
            message: 'Account not found',
          },
        );
      }

      const rule = schema.findBy('signRule', {
        id: ruleId,
        accountNickname: nickname,
      });

      if (!rule) {
        return new Response(
          404,
          {},
          {
            code: '404',
            message: 'Rule not found',
          },
        );
      }

      rule.destroy();
      return new Response(200);
    },
  );

  server.put(
    'accounts/:nickname/signrules/:ruleId',
    (schema: AppSchema, request) => {
      const { nickname, ruleId } = request.params;
      const account = schema.findBy('account', { nickname });

      if (!account) {
        return new Response(
          404,
          {},
          {
            code: '404',
            message: 'Account not found',
          },
        );
      }

      const rule = schema.findBy('signRule', {
        id: ruleId,
        accountNickname: nickname,
      });

      if (!rule) {
        return new Response(
          404,
          {},
          {
            code: '404',
            message: 'Rule not found',
          },
        );
      }

      const body = JSON.parse(request.requestBody);
      const { ruleType, contract, enabled } = body;

      if (!ruleType || !contract || enabled === undefined) {
        return new Response(
          400,
          {},
          {
            code: '400',
            message: 'Missing required fields',
          },
        );
      }

      if (
        ![RuleType.AutoSign, RuleType.DisablePasswordPrompt].includes(ruleType)
      ) {
        return new Response(
          422,
          {},
          {
            code: '422',
            message: 'Invalid rule type',
          },
        );
      }

      rule.update({
        ruleType,
        contract,
        enabled,
      });

      return new Response(200, {}, { id: ruleId });
    },
  );
}
