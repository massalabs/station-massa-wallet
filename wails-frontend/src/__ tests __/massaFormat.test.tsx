import { Unit, formatStandard } from '../utils/massaFormat';

describe('formatStandard with min  string value', () => {
  it('should format a value of 10,000', () => {
    const value = '0000000000';

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('0');
  });
});

describe('formatStandard with min bigint value', () => {
  it('should return 10,000', () => {
    const value = 0n;

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('0');
  });
});

describe('formatStandard with mid range string value', () => {
  it('should return 10,000', () => {
    const value = '10000000000000';

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('10,000');
  });
});

describe('formatStandard with mid range bigint value', () => {
  it('should return 10,000', () => {
    const value = 10000000000000n;

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('10,000');
  });
});

describe('formatStandard with max  string value', () => {
  it('should format a value of 10,000', () => {
    const value = '922337203600000000000';

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('922,337,203,600');
  });
});

describe('formatStandard with max bigint value', () => {
  it('should return 10,000', () => {
    const value = 9223372036854775807n;

    const result = formatStandard(value, Unit.NanoMAS);

    expect(result).toBe('9,223,372,036.854776');
  });
});
