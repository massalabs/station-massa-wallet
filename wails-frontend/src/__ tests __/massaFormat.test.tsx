import { Unit, formatStandard } from '../utils/massaFormat';

describe('formatStandard', () => {
  it('formatStandard with min  string value', () => {
    const value = '0000000000';

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('0');
  });

  it('formatStandard with min bigint value', () => {
    const value = 0n;

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('0');
  });

  it('formatStandard with mid range string value', () => {
    const value = '10000000000000';

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('10,000');
  });

  it('formatStandard with mid range bigint value', () => {
    const value = 10000000000000n;

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('10,000');
  });

  it('formatStandard with max string value', () => {
    const value = '922337203600000000000';

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('922,337,203,600');
  });

  it('formatStandard with max bigint value', () => {
    const value = 9223372036854775807n;

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('9,223,372,036.854776');
  });
});
