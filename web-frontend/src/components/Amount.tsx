import MenuActive from './MenuActive';

interface AmountProps {
  amount: number;
  fractionDigits?: number;
}

export default function Amount(props: AmountProps) {
  const formattedAmount = props.amount.toFixed(props.fractionDigits || 2);

  return <MenuActive>{formattedAmount}</MenuActive>;
}
