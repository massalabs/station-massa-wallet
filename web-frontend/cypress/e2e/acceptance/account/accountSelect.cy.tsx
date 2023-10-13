import { compareSnapshot } from '../../../compareSnapshot';

describe('E2E | Acceptance | Account', function () {
  describe('/account-select', () => {
    it('should match snapshot', () => {
      cy.visit('/');

      compareSnapshot(cy, 'account-select');
    });
  });
});
