import { compareSnapshot } from '../../compareSnapshot';

describe('Component | Integration ', function () {
  it('pass', () => {
    compareSnapshot(cy, 'example');
  });
});
