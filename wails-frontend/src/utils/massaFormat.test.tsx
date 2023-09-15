import { Unit, formatStandard } from '.';

//test remove trailing zeros

// Test case 1: Pass a BigInt of 10000000000000
test('Remove trailing zeros from BigInt', () => {
  const input = 10000000000000n;
  const result = formatStandard(input, Unit.MAS); // Convert BigInt to string
  expect(result).toBe('10,000');
});

// Test case 2: Pass a string "10000000000000"
test('Remove trailing zeros from string', () => {
  const input = '100000000000000';
  const result = formatStandard(input, Unit.MAS);
  expect(result).toBe('10,000');
});
